package config

import (
	"errors"
	"io"

	"github.com/BurntSushi/toml"
)

var (
	ErrNilReader = errors.New("invalid nil reader")
)

// Config holds the full configuration for the request being sent.
//
// All Config fields can contain property tokens which will be resolved
// prior to sending the request (see package resolver for more info)
//
// If the resolved configuration does not pass validation an error is raised
// and the request will not be sent.
//
// This means that the data in your configuration files does not need to conform
// to valid formats.  The resolved values, on the other hand, are required to conform.
type Config struct {
	Request    RequestConfig `toml:"request,omitempty"    validate:"required"`
	JWT        JWTConfig     `toml:"jwt,omitempty"        validate:"omitempty"`
	Properties Properties    `toml:"properties,omitempty" validate:"omitempty,dive,gt=0"`
}

type Properties map[string]any

func NewConfig() *Config {
	return &Config{}
}

// Load populates the configuration using the toml data stored in the reader.
func (c *Config) Load(reader io.Reader) error {
	if reader == nil {
		return ErrNilReader
	}
	if _, err := toml.NewDecoder(reader).Decode(c); err != nil {
		return err
	}
	return nil
}
