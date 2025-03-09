// package sender declares the interface for a request sender and provides a factory
// for instantiating whichever sender you need.
package sender

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/logging"
	"github.com/keithpaterson/postal/sender/curl"
	"github.com/keithpaterson/postal/sender/native"
	"github.com/keithpaterson/postal/validate"

	"go.uber.org/zap"
)

// enum to identify which sender you want
const (
	NativeSender SenderType = iota
	CurlSender
)

const (
	NativeSenderName = "native"
	CurlSenderName   = "curl"
)

const (
	errMsgInvalid = "invalid sender"
)

type SenderType int

type Sender interface {
	Send(config *config.Config) error
}

type rootSender struct {
	id  SenderType
	log *zap.SugaredLogger
}

func NewSender(id SenderType) (*rootSender, error) {
	return &rootSender{id: id, log: logging.NamedLogger("sender")}, nil
}

func NewNamedSender(name string) (*rootSender, error) {
	var err error
	var id SenderType
	if id, err = toSenderType(name); err != nil {
		return nil, err
	}
	return NewSender(id)
}

func (s *rootSender) Send(cfg *config.Config) error {
	var err error
	var actualCfg *config.Config
	if actualCfg, err = validate.ValidateConfig(cfg); err != nil {
		return err
	}
	s.log.Debugw("Send", "validated config", fmt.Sprintf("%#v", actualCfg))

	switch s.id {
	case NativeSender:
		e := native.NewSender(s.log)
		return e.Send(actualCfg)
	case CurlSender:
		e := curl.NewSender(s.log)
		return e.Send(actualCfg)
	default:
		return errors.New(errMsgInvalid)
	}
}

func toSenderType(name string) (SenderType, error) {
	switch strings.ToLower(name) {
	case NativeSenderName:
		return NativeSender, nil
	case CurlSenderName:
		return CurlSender, nil
	default:
		return -1, fmt.Errorf("%s name '%s'", errMsgInvalid, name)
	}
}
