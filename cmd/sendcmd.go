package cmd

import (
	"fmt"
	"os"
	"postal/config"
	"postal/logging"
	"postal/sender"
	"strings"

	"github.com/spf13/cobra"
)

const (
	fileFlag, fFlag   = "file", "f"
	propFlag, pFlag   = "prop", "p"
	methodFlag, mFlag = "method", "m"
	urlFlag, uFlag    = "url", "u"
	bodyFlag, bFlag   = "body", "b"
	jwtFlag           = "jwt"
	algFlag, aFlag    = "alg", "a"
	headerFlag, hFlag = "header", "H"
	execFlag          = "using"
)

var (
	senderNames = strings.Join([]string{sender.NativeSenderName, sender.CurlSenderName}, ", ")
)

func NewSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "send a message",
		Run:   sendMessage,
	}

	cmd.Flags().String(execFlag, sender.NativeSenderName, fmt.Sprintf("Identifies which executor to use: one of [%s]", senderNames))
	cmd.Flags().StringArrayP(fileFlag, fFlag, []string{}, "config file")
	cmd.Flags().StringArrayP(propFlag, pFlag, []string{}, "one or more properties (key=value)")
	cmd.Flags().StringP(methodFlag, mFlag, "", "HTTP method")
	cmd.Flags().StringP(urlFlag, uFlag, "", "URL")
	cmd.Flags().StringP(bodyFlag, bFlag, "", "body specification")
	cmd.Flags().StringArray(jwtFlag, []string{}, "one or more JWT claims (key=value)")
	cmd.Flags().StringArrayP(headerFlag, hFlag, []string{}, "one or HTTP headers (key=value)")
	cmd.Flags().StringP(algFlag, aFlag, config.DefaultAlgorithm, "JWT algorithm")

	cmd.MarkFlagRequired("file")

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

	dryRun, _ := cmd.Flags().GetBool(dryRunFlag)

	var err error
	cfg := config.NewConfig()

	// order is important here:
	// - config files are lowest-order data sources; bring all of them in first
	if err = loadConfig(cmd, cfg); err != nil {
		return err
	}

	// command-line arguments are higher-order data sources; these supercede anything
	// spcified in the config files
	if err = processProperties(cmd, cfg); err != nil {
		return err
	}
	if err = processRequestArgs(cmd, cfg); err != nil {
		return err
	}
	if err = processJWTClaims(cmd, cfg); err != nil {
		return err
	}

	name, _ := cmd.Flags().GetString(execFlag)

	log.Debug("dryRun:", dryRun)
	log.Debugf("%#v", *cfg)

	sender, err := sender.NewNamedSender(name)
	if err != nil {
		return err
	}
	return sender.Send(cfg)
}

func loadConfig(cmd *cobra.Command, cfg *config.Config) error {
	var err error
	var filenames []string
	if filenames, err = cmd.Flags().GetStringArray(fileFlag); err != nil {
		return flagError(fileFlag, err)
	}

	for _, filename := range filenames {
		var file *os.File
		var err error
		if file, err = os.Open(filename); err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		if err = cfg.Load(file); err != nil {
			file.Close()
			return fmt.Errorf("failed to load file: %w", err)
		}
		file.Close()
	}
	return nil
}

func processProperties(cmd *cobra.Command, cfg *config.Config) error {
	var err error
	var props []string
	if props, err = cmd.Flags().GetStringArray(propFlag); err != nil {
		return flagError(propFlag, err)
	}

	for _, prop := range props {
		var key, value string
		var ok bool
		if key, value, ok = strings.Cut(prop, "="); !ok {
			return fmt.Errorf("failed to parse property value '%s'", prop)
		}
		cfg.Properties[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return nil
}

func processRequestArgs(cmd *cobra.Command, cfg *config.Config) error {
	var err error

	if cmd.Flags().Changed(methodFlag) {
		var method string
		if method, err = cmd.Flags().GetString(methodFlag); err != nil {
			return flagError(methodFlag, err)
		}
		cfg.Request.Method = method
	}

	if cmd.Flags().Changed(urlFlag) {
		var url string
		if url, err = cmd.Flags().GetString(urlFlag); err != nil {
			return flagError(urlFlag, err)
		}
		cfg.Request.URL = url
	}

	if cmd.Flags().Changed(bodyFlag) {
		var body string
		if body, err = cmd.Flags().GetString(bodyFlag); err != nil {
			return flagError(bodyFlag, err)
		}
		cfg.Request.Body = body
		// TODO(keithpaterson): we use the body information to determine Mime Type
	}

	return processRequestHeaders(cmd, cfg)
}

func processRequestHeaders(cmd *cobra.Command, cfg *config.Config) error {
	var err error
	var headers []string
	if headers, err = cmd.Flags().GetStringArray(headerFlag); err != nil {
		return flagError(headerFlag, err)
	}

	for _, header := range headers {
		key, value, ok := strings.Cut(header, "=")
		if !ok {
			return fmt.Errorf("invalid header specification '%s' (expect key=value)", header)
		}
		cfg.Request.Headers[key] = value
	}
	return nil
}

func processJWTClaims(cmd *cobra.Command, cfg *config.Config) error {
	var err error

	if cmd.Flags().Changed(algFlag) {
		var algorithm string
		if algorithm, err = cmd.Flags().GetString(algFlag); err != nil {
			return fmt.Errorf("failed to process JWT algorithm: %w", err)
		}
		cfg.JWT.Header.Alg = algorithm
	}

	var claims []string
	if claims, err = cmd.Flags().GetStringArray(jwtFlag); err != nil {
		return fmt.Errorf("failed to process JWT claims: %w", err)
	}

	for _, claim := range claims {
		var key, value string
		var ok bool
		if key, value, ok = strings.Cut(claim, "="); !ok {
			return fmt.Errorf("failed to parse JWT claim value '%s'", claim)
		}
		cfg.JWT.Claims[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return nil
}

func flagError(name string, err error) error {
	if err != nil {
		return fmt.Errorf("failed to process %s flag: %w", name, err)
	} else {
		return fmt.Errorf("failed to process %s flag", name)
	}
}
