package cmd

import (
	"fmt"
	"os"
	"postal/config"
	"postal/executor"
	"strings"

	"github.com/spf13/cobra"
)

const (
	fileFlag, fFlag   = "file", "f"
	propFlag, pFlag   = "prop", "p"
	methodFlag, mFlag = "method", "m"
	urlFlag, uFlag    = "url", "u"
	bodyFlag, bFlag   = "body", "b"
	mimeTypeFlag      = "mime-type"
	jwtFlag           = "jwt"
	algFlag, aFlag    = "alg", "a"
	execFlag          = "using"
)

var (
	executorNames = strings.Join([]string{executor.StrNativeExecutor, executor.StrCurlExecutor}, ", ")
)

func NewSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "send a message",
		RunE:  sendMessage,
	}

	cmd.Flags().String(execFlag, executor.StrNativeExecutor, fmt.Sprintf("Identifies which executor to use: one of [%s]", executorNames))
	cmd.Flags().StringArrayP(fileFlag, fFlag, []string{}, "config file")
	cmd.Flags().StringArrayP(propFlag, pFlag, []string{}, "one or more properties (key=value)")
	cmd.Flags().StringP(methodFlag, mFlag, "", "HTTP method")
	cmd.Flags().StringP(urlFlag, uFlag, "", "URL")
	cmd.Flags().StringP(bodyFlag, bFlag, "", "body specification")
	cmd.Flags().String(mimeTypeFlag, "", "Body Mime Type")
	cmd.Flags().StringArray(jwtFlag, []string{}, "one or more JWT claims (key=value)")
	cmd.Flags().StringP(algFlag, aFlag, "", "JWT algorithm")

	cmd.MarkFlagRequired("file")

	return cmd
}

func sendMessage(cmd *cobra.Command, _ []string) error {
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

	fmt.Println("dryRun:", dryRun)
	fmt.Printf("%#v\n", *cfg)

	return executor.RunNamed(name, cfg)
}

func loadConfig(cmd *cobra.Command, cfg *config.Config) error {
	var err error
	var filenames []string
	if filenames, err = cmd.Flags().GetStringArray(fileFlag); err != nil {
		return fmt.Errorf("failed to process file flag: %w", err)
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
		return fmt.Errorf("failed to process properties: %w", err)
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
			return fmt.Errorf("failed to process method flag: %w", err)
		}
		cfg.Request.Method = method
	}

	if cmd.Flags().Changed(urlFlag) {
		var url string
		if url, err = cmd.Flags().GetString(urlFlag); err != nil {
			return fmt.Errorf("failed to process url flag: %w", err)
		}
		cfg.Request.URL = url
	}

	if cmd.Flags().Changed(bodyFlag) {
		var body string
		if body, err = cmd.Flags().GetString(bodyFlag); err != nil {
			return fmt.Errorf("failed to process body flag: %w", err)
		}
		cfg.Request.Body = body
		// TODO(keithpaterson): we use the body information to determine Mime Type
	}

	// .. and this explicitly overrides any implied mime type
	if cmd.Flags().Changed(mimeTypeFlag) {
		var mimeType string
		if mimeType, err = cmd.Flags().GetString(mimeTypeFlag); err != nil {
			return fmt.Errorf("failed to process mime-type flag: %w", err)
		}
		cfg.Request.MimeType = mimeType
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
		cfg.JWT.Payload[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return nil
}
