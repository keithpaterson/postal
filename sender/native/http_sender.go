package native

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/logging"

	"github.com/keithpaterson/resweave-utils/client"
	"github.com/keithpaterson/resweave-utils/header"
	"go.uber.org/zap"
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
	c := client.DefaultHTTPClient().WithLogger(logging.Logger())
	if resp, err = c.Execute(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	// how to parse the response?
	var data []byte
	if data, err = io.ReadAll(resp.Body); err != nil {
		s.log.Debugw("execute", "response", data)
	}
	fmt.Println("\nresponse:\n>>>\n", string(data), "\n<<<")

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
		return nil, errors.New("invalid Request.Body spec: expect 'name:data'")
	}

	var err error
	var body []byte

	switch name {
	case "json":
		// data is raw json
		body = []byte(data)
		if _, ok := s.cfg.Request.Headers["content-type"]; !ok {
			s.cfg.Request.Headers["content-type"] = header.MimeTypeJson
		}
	case "file":
		var file *os.File
		if file, err = os.Open(data); err != nil {
			return nil, err
		}
		if body, err = io.ReadAll(file); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported Request.Body spec '%s'", name)
	}

	return body, nil
}
