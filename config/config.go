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
// Any Config field can contain ${tokens} which will be resolved
// prior to handling (see package resolver for more info).
//
// Validation only occurs after the configuration file(s) have been loaded,
// the command-line arguments have been processed, and the tokens have been
// resolved.
//
// If the Config does not pass validation an error is raised
// and processing terminates with an error.
//
// This allows you to specify data in your configuration files that may be invalid due
// to the presence of the tokens. After the tokens are resolved, the configuration must
// be valid.
//
// For Example:
//
// After reading config files you may have:
//
//	Config{Request:{Method: "${prop:method}", URL: "${prop.url}"}}
//
// This would be considered invalid because the URL does not conform to the usual URL format.
//
// After resolving, your config may change to:
//
//	Config{Request{Method: "GET", URL: "http://test.io/foo"}}
//
// This would be considered valid.
//
// In the case where some tokens cannot be resolved, you might end up with something like:
//
//	Config{Request{Method: "GET", URL: "${prop.url}"}}
//
// This will fail validation (improper URL format) and result in an error message
type Config struct {
	Request    RequestConfig `toml:"request,omitempty"    validate:"required"`
	JWT        JWTConfig     `toml:"jwt,omitempty"        validate:"omitempty"`
	Cacert     CacertConfig  `toml:"cacert,omitempty"     validate:"omitempty"`
	Properties Properties    `toml:"properties,omitempty" validate:"omitempty,dive,gt=0"`
	Output     OutputConfig  `toml:"output,omitempty"     validate:"omitempty"`
}

type Properties map[string]any

func NewConfig() *Config {
	cfg := &Config{Request: newRequestConfig(), Cacert: newCacertConfig(), Output: newOutputConfig()}
	return cfg
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
