// package sender declares the interface for a request sender and provides a factory
// for instantiating whichever sender you need.
package sender

import (
	"errors"
	"fmt"
	"slices"

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

var (
	ErrInvalidSender = errors.New("invalid sender")
)

var (
	Names = []string{NativeSenderName, CurlSenderName}
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

	// Runtime info doesn't get persisted, so copy the original information
	actualCfg.Runtime = cfg.Runtime

	switch s.id {
	case NativeSender:
		e := native.NewSender(s.log)
		return e.Send(actualCfg)
	case CurlSender:
		e := curl.NewSender(s.log)
		return e.Send(actualCfg)
	default:
		return ErrInvalidSender
	}
}

func toSenderType(name string) (SenderType, error) {
	index := slices.Index(Names, name)
	if index < 0 {
		return -1, fmt.Errorf("%w: name '%s'", ErrInvalidSender, name)
	}
	return SenderType(index), nil
}
