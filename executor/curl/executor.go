package curl

import (
	"errors"
	"postal/config"
)

type executor struct {
}

func NewExecutor() *executor {
	return &executor{}
}

func (e *executor) Execute(cfg *config.Config) error {
	return errors.New("not implemented")
}
