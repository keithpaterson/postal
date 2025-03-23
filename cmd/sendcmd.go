package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/logging"
	"github.com/keithpaterson/postal/sender"

	"github.com/spf13/cobra"
)

const (
	algFlag, aFlag      = "alg", "a"
	bodyFlag, bFlag     = "body", "b"
	cacertFlag          = "cacert"
	configFlag, cFlag   = "config", "c"
	headerFlag, hFlag   = "header", "H"
	jwtFlag             = "jwt"
	methodFlag, mFlag   = "method", "m"
	outFileFlag, fFlag  = "out-file", "f"
	outFmtFlag, oFlag   = "out-format", "o"
	propFlag, pFlag     = "prop", "p"
	signingKeyFlag      = "signing-key"
	templateFlag, tFlag = "template", "t"
	urlFlag, uFlag      = "url", "u"
	usingFlag           = "using"
)

var (
	ErrInvalidConfigFile    = errors.New("config file error")
	ErrInvalidPropertyValue = errors.New("invalid property value")
	ErrInvalidHeader        = errors.New("invalid header value")
	ErrInvalidJWTClaim      = errors.New("invalid JWT claim")
)

func NewSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send -c filename [-c filename...] [flags]",
		Short: "send a message",
		Run:   sendMessage,
	}

	cmd.Flags().StringP(algFlag, aFlag, config.DefaultAlgorithm, "JWT algorithm")
	cmd.Flags().StringP(bodyFlag, bFlag, "", "body specification")
	cmd.Flags().String(cacertFlag, "", "CA certification specification")
	cmd.Flags().StringArrayP(configFlag, cFlag, []string{}, "one or more config file names")
	cmd.Flags().StringArrayP(headerFlag, hFlag, []string{}, "one or more HTTP headers (key=value)")
	cmd.Flags().StringArray(jwtFlag, []string{}, "one or more JWT claims (key=value)")
	cmd.Flags().StringP(methodFlag, mFlag, "", "HTTP method")
	cmd.Flags().StringP(outFileFlag, fFlag, "stdout", "specify a filename to write the result into")
	cmd.Flags().StringP(outFmtFlag, oFlag, "text", fmt.Sprintf("output format, one of [%s]", config.OutFmtNames))
	cmd.Flags().StringArrayP(propFlag, pFlag, []string{}, "one or more properties (key=value)")
	cmd.Flags().String(signingKeyFlag, "", "your signing key; used to sign the JWT token")
	cmd.Flags().StringP(templateFlag, tFlag, "${response:body}", "template for writing text response output")
	cmd.Flags().StringP(urlFlag, uFlag, "", "URL")
	cmd.Flags().String(usingFlag, sender.NativeSenderName, fmt.Sprintf("Identifies which sender to use: one of [%s]", sender.Names))

	cmd.MarkFlagRequired(configFlag)

	return cmd
}

func sendMessage(cmd *cobra.Command, args []string) {
	if err := sendMessageE(cmd, args); err != nil {
		fmt.Println("ERROR: failed to send request:", err)
	}
}

func sendMessageE(cmd *cobra.Command, _ []string) error {
	setupLogging(cmd)
	log := logging.NamedLogger("sendcmd")

	parser := &sendCmdParser{}
	cfg, err := parser.parseConfig(cmd)
	if err != nil {
		return err
	}

	name, _ := cmd.Flags().GetString(usingFlag)

	log.Debug("dryRun:", parser.dryRun)
	log.Debugf("%#v", *cfg)

	sender, err := sender.NewNamedSender(name)
	if err != nil {
		return err
	}
	return sender.Send(cfg)
}

type sendCmdParser struct {
	cmd *cobra.Command
	cfg *config.Config

	dryRun bool
}

func (p *sendCmdParser) parseConfig(cmd *cobra.Command) (*config.Config, error) {
	p.cmd = cmd
	defer func() { p.cmd = nil; p.cfg = nil }()

	p.dryRun, _ = cmd.Flags().GetBool(dryRunFlag)

	var err error
	p.cfg = config.NewConfig()

	// order is important here:
	// - config files are lowest-order data sources; bring all of them in first
	if err = p.loadConfig(); err != nil {
		return nil, err
	}

	// command-line arguments are higher-order data sources; these supercede anything
	// spcified in the config files
	if err = p.processProperties(); err != nil {
		return nil, err
	}
	if err = p.processRequestArgs(); err != nil {
		return nil, err
	}
	if err = p.processJWT(); err != nil {
		return nil, err
	}
	if err = p.processOutput(); err != nil {
		return nil, err
	}
	return p.cfg, nil
}

func (p *sendCmdParser) loadConfig() error {
	var err error
	var filenames []string
	if filenames, err = p.cmd.Flags().GetStringArray(configFlag); err != nil {
		return p.flagError(configFlag, err)
	}

	for _, filename := range filenames {
		var file *os.File
		var err error
		if file, err = os.Open(filename); err != nil {
			return fmt.Errorf("%w: failed to open file: %w", ErrInvalidConfigFile, err)
		}
		if err = p.cfg.Load(file); err != nil {
			file.Close()
			return fmt.Errorf("%w: failed to load file: %w", ErrInvalidConfigFile, err)
		}
		file.Close()
	}
	return nil
}

func (p *sendCmdParser) processProperties() error {
	var err error
	var props []string
	if props, err = p.cmd.Flags().GetStringArray(propFlag); err != nil {
		return p.flagError(propFlag, err)
	}

	for _, prop := range props {
		var key, value string
		var ok bool
		if key, value, ok = strings.Cut(prop, "="); !ok {
			return fmt.Errorf("%w: failed to parse value '%s'", ErrInvalidPropertyValue, prop)
		}
		p.cfg.Properties[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return nil
}

func (p *sendCmdParser) processRequestArgs() error {
	var err error

	if p.cmd.Flags().Changed(methodFlag) {
		var method string
		if method, err = p.cmd.Flags().GetString(methodFlag); err != nil {
			return p.flagError(methodFlag, err)
		}
		p.cfg.Request.Method = method
	}

	if p.cmd.Flags().Changed(urlFlag) {
		var url string
		if url, err = p.cmd.Flags().GetString(urlFlag); err != nil {
			return p.flagError(urlFlag, err)
		}
		p.cfg.Request.URL = url
	}

	if p.cmd.Flags().Changed(bodyFlag) {
		var body string
		if body, err = p.cmd.Flags().GetString(bodyFlag); err != nil {
			return p.flagError(bodyFlag, err)
		}
		p.cfg.Request.Body = body
		// TODO(keithpaterson): we use the body information to determine Mime Type
	}

	return p.processRequestHeaders()
}

func (p *sendCmdParser) processRequestHeaders() error {
	var err error
	var headers []string
	if headers, err = p.cmd.Flags().GetStringArray(headerFlag); err != nil {
		return p.flagError(headerFlag, err)
	}

	for _, header := range headers {
		key, value, ok := strings.Cut(header, "=")
		if !ok {
			return fmt.Errorf("%w: expect key=value, got '%s'", ErrInvalidHeader, header)
		}
		p.cfg.Request.Headers[key] = value
	}
	return nil
}

func (p *sendCmdParser) processJWT() error {

	var err error
	if p.cmd.Flags().Changed(signingKeyFlag) {
		var key string
		if key, err = p.cmd.Flags().GetString(signingKeyFlag); err != nil {
			return fmt.Errorf("failed to process JWT signing key: %w", err)
		}
		p.cfg.JWT.SigningKey = key
	}

	if p.cmd.Flags().Changed(algFlag) {
		var algorithm string
		if algorithm, err = p.cmd.Flags().GetString(algFlag); err != nil {
			return fmt.Errorf("failed to process JWT algorithm: %w", err)
		}
		p.cfg.JWT.Header.Alg = algorithm
	}

	return p.processJWTClaims()
}

func (p *sendCmdParser) processJWTClaims() error {
	var err error

	var claims []string
	if claims, err = p.cmd.Flags().GetStringArray(jwtFlag); err != nil {
		return fmt.Errorf("failed to process JWT claims: %w", err)
	}

	for _, claim := range claims {
		var key, value string
		var ok bool
		if key, value, ok = strings.Cut(claim, "="); !ok {
			return fmt.Errorf("%w: expect key=value, got '%s'", ErrInvalidJWTClaim, claim)
		}
		p.cfg.JWT.Claims[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return nil
}

func (p *sendCmdParser) processOutput() error {
	var err error

	var template string
	if template, err = p.cmd.Flags().GetString(templateFlag); err != nil {
		return p.flagError(templateFlag, err)
	}
	p.cfg.Output.Template = template

	var outFormat string
	if outFormat, err = p.cmd.Flags().GetString(outFmtFlag); err != nil {
		return p.flagError(outFmtFlag, err)
	}
	p.cfg.Output.Format = outFormat

	var outFile string
	if outFile, err = p.cmd.Flags().GetString(outFileFlag); err != nil {
		return p.flagError(outFileFlag, err)
	}
	p.cfg.Output.Filename = outFile

	return nil
}

func (p *sendCmdParser) flagError(name string, err error) error {
	if err != nil {
		return fmt.Errorf("failed to process %s flag: %w", name, err)
	} else {
		return fmt.Errorf("failed to process %s flag", name)
	}
}
