package curl

import (
	"errors"
	"postal/config"
)

type sender struct {
}

func NewSender() *sender {
	return &sender{}
}

func (s *sender) Send(cfg *config.Config) error {
	return errors.New("not implemented")
}
