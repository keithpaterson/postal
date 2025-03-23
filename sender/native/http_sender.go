package native

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/keithpaterson/postal/cacert"
	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/output"
	"github.com/keithpaterson/postal/validate"

	"github.com/keithpaterson/resweave-utils/client"
	"github.com/keithpaterson/resweave-utils/header"
	ulog "github.com/keithpaterson/resweave-utils/logging"
	"go.uber.org/zap"
)

var (
	ErrInvalidBodySpec     = errors.New("invalid Request.Body spec: expect 'name:data'")
	ErrUnsupportedBodySpec = errors.New("unsupported Request.Body spec")
	ErrInvalidBody         = errors.New("invalid body")
	ErrInvalidCert         = errors.New("invalid TLS certificate")
)

type httpSender struct {
	cfg *config.Config
	log *zap.SugaredLogger
}

func sendHttp(cfg *config.Config, log *zap.SugaredLogger) error {
	s := &httpSender{cfg: cfg, log: log.Named("http")}
	return s.execute()
}

func (s *httpSender) execute() error {
	s.log.Debugw("execute", "status", "started")
	defer s.log.Debugw("execute", "status", "completed")

	var err error
	var body []byte
	if body, err = s.getBodyData(); err != nil {
		return err
	}

	var req *http.Request
	if req, err = s.newRequest(body); err != nil {
		return err
	}
	for key, value := range s.cfg.Request.Headers {
		req.Header.Add(key, value)
	}

	var resp *http.Response
	c := client.NewHTTPClient("test").WithRetryHandler(client.NewRetryCounter(0))
	if err = s.configureTLS(c.Client); err != nil {
		return err
	}

	if resp, err = c.Execute(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	writer := output.NewOutputter(s.cfg)
	if err = writer.Write(resp); err != nil {
		return err
	}

	return nil
}

func (s *httpSender) configureTLS(client *http.Client) error {
	if s.cfg.Cacert.Pool() != config.CertPoolNone {
		parser := cacert.FromConfig(s.cfg.Cacert)
		pool, err := parser.GetCertificatePool()
		if err != nil {
			s.log.Errorw("execute", ulog.LogKeyStatus, "failed to build cert pool", ulog.LogKeyError, err)
			return fmt.Errorf("%w: invalid pool: %w", ErrInvalidCert, err)
		}
		certificates, err := parser.GetCertificates()
		if err != nil {
			s.log.Errorw("execute", ulog.LogKeyStatus, "failed to build certificate list", ulog.LogKeyError, err)
			return fmt.Errorf("failed to configure TLS: %w", err)
		}
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: certificates,
			},
		}
	}
	return nil
}

func (s *httpSender) newRequest(body []byte) (*http.Request, error) {
	return http.NewRequest(s.cfg.Request.Method, s.cfg.Request.URL, bytes.NewBuffer(body))
}

func (s *httpSender) getBodyData() ([]byte, error) {
	// body specification is one of:
	//   "json:{json-data}"
	//   "file:file-name"
	//   TODO(keithpaterson): others
	if s.cfg.Request.Body == "" {
		return nil, nil
	}

	name, data, ok := strings.Cut(s.cfg.Request.Body, ":")
	if !ok {
		return nil, ErrInvalidBodySpec
	}

	var err error
	var body []byte

	switch name {
	case "json":
		// data is raw json
		body = []byte(data)
		if err = validate.ValidateJson(body); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidBody, err)
		}
		if _, ok := s.cfg.Request.Headers["content-type"]; !ok {
			s.cfg.Request.Headers["content-type"] = header.MimeTypeJson
		}
	case "file":
		if body, err = os.ReadFile(data); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidBody, err)
		}
	default:
		return nil, fmt.Errorf("%w: '%s'", ErrUnsupportedBodySpec, name)
	}

	return body, nil
}
