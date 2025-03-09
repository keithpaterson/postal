package curl

import (
	"errors"

	"github.com/keithpaterson/postal/config"

	"go.uber.org/zap"
)

type sender struct {
	log *zap.SugaredLogger
}

func NewSender(log *zap.SugaredLogger) *sender {
	return &sender{log: log.Named("curl")}
}

func (s *sender) Send(cfg *config.Config) error {
	return errors.New("not implemented")
}
