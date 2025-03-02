package native

import (
	"fmt"
	"net/url"
	"postal/config"
)

type sender struct {
	cfg *config.Config
}

func NewSender() *sender {
	return &sender{}
}

func (s *sender) Send(cfg *config.Config) error {
	s.cfg = cfg

	var err error

	// parse the URL here so that we can determine the scheme; from that we can call an appropriate function to handle that scheme
	var target *url.URL
	if target, err = url.Parse(s.cfg.Request.URL); err != nil {
		return err
	}
	switch target.Scheme {
	case "http", "https":
		err = sendHttp(s.cfg)
	default:
		err = fmt.Errorf("unsupported scheme '%s'", target.Scheme)
	}
	return err
}
