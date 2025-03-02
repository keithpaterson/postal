package native

import (
	"fmt"
	"net/url"
	"postal/config"

	"go.uber.org/zap"
)

type sender struct {
	cfg *config.Config
	log *zap.SugaredLogger
}

func NewSender(log *zap.SugaredLogger) *sender {
	return &sender{log: log.Named("native")}
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
		err = sendHttp(s.cfg, s.log)
	default:
		err = fmt.Errorf("unsupported scheme '%s'", target.Scheme)
	}
	return err
}
