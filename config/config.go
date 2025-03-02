package config

import (
	"errors"
	"io"

	"github.com/BurntSushi/toml"
)

var (
	ErrNilReader = errors.New("invalid nil reader")
)

type Config struct {
	Request    RequestConfig `toml:"request,omitempty"    validate:"required"`
	JWT        JWTConfig     `toml:"jwt,omitempty"        validate:"omitempty"`
	Properties Properties    `toml:"properties,omitempty" validate:"omitempty,dive,gt=0"`
}

type Properties map[string]any

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(reader io.Reader) error {
	if reader == nil {
		return ErrNilReader
	}
	if _, err := toml.NewDecoder(reader).Decode(c); err != nil {
		return err
	}
	return nil
}
