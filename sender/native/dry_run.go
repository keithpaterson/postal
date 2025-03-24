package native

import (
	"fmt"
	"net/http"

	"github.com/keithpaterson/postal/config"
)

func (s *httpSender) dryRun(req *http.Request) error {
	// for now, dry run always outputs to the console
	fmt.Println("DRY RUN")
	fmt.Println("-------")
	fmt.Println("\nConfiguration:")
	s.dryCfgRequest(s.cfg.Request)
	s.dryCacert(s.cfg.Cacert)
	s.dryProperties(s.cfg.Properties)
	s.dryOutput(s.cfg.Output)

	s.dryHttpRequest(req)

	fmt.Println()

	return nil
}

func (s *httpSender) dryCfgRequest(req config.RequestConfig) {
	fmt.Println("  Request:")
	fmt.Println("   ", req.Method, req.URL)

	if req.Body != "" {
		fmt.Printf("    with body:\n      %s", req.Body)
	}
	fmt.Println()

	if len(req.Headers) > 0 {
		fmt.Print("    with headers:\n      ")
		for key, value := range s.cfg.Request.Headers {
			fmt.Printf("%s=%s; ", key, value)
		}
	}
	fmt.Println()
}

func (s *httpSender) dryCacert(cacert config.CacertConfig) {
	if cacert.Pool() == config.CertPoolNone {
		return
	}

	fmt.Println("  CA Certificates:")
	fmt.Println("    Pool:", cacert.PoolName)

	if cacert.CaCrt != "" {
		fmt.Println("    CA CRT:", cacert.CaCrt)
	}
	if len(cacert.Certificates) > 0 {
		fmt.Println("    Certificates:")
		for _, c := range cacert.Certificates {
			fmt.Println("     ", c)
		}
		if cacert.CertificateFileExtensions != [2]string{"", ""} {
			fmt.Printf("    ('file:' uses extensions: %v)", cacert.CertificateFileExtensions)
		}
	}
}

func (s *httpSender) dryProperties(props config.Properties) {
	if len(props) == 0 {
		return
	}
	fmt.Println("  Properties:")
	for key, value := range props {
		fmt.Println("   ", key, "=", value)
	}
}

func (s *httpSender) dryOutput(output config.OutputConfig) {
	fmt.Println("  Output:")
	fmt.Println("    Format:", output.Format)
	fmt.Println("    Filename:", output.Filename)
	fmt.Println("    using Template:")
	fmt.Println("      >>>")
	fmt.Println(output.Template)
	fmt.Println("      <<<")
}

func (s *httpSender) dryHttpRequest(req *http.Request) {
	if nil == req {
		return
	}

	fmt.Println("\nHTTP Request:")
	fmt.Println(" ", req.Method, req.URL.String())
	if len(req.Header) > 0 {
		fmt.Print("  with Header:\n    ")
		for key, value := range req.Header {
			fmt.Printf("%s=%s; ", key, value)
		}
		fmt.Println()
	}
	fmt.Println("  Content Length:", req.ContentLength)
}
