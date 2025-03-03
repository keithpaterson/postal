package config

type RequestConfig struct {
	Method  string        `toml:"method,omitempty"    validate:"method,required"`
	URL     string        `toml:"url,omitempty"       validate:"url,required"`
	Body    string        `toml:"body,omitempty"      validate:"required"`
	Headers HeadersConfig `toml:"headers,omitempty"   validate:"omitempty"`
}

type HeadersConfig map[string]string

// Notes:
// URL must be string here since we could craft an invalid url if it contains any tokens.
