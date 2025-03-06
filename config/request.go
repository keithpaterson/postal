package config

// RequestConfig stores the properties of the request that will be sent.
//
// # Method indicates the HTTP method to use
//
// # URL specifies the fully-formed URL that we will send to
//
// Body supplies any body data (e.g. for POST, PUT, PUSH methods)
//
// Headers is a list of name-value pairs for any additional request headers you need.
type RequestConfig struct {
	Method  string        `toml:"method,omitempty"    validate:"method,required"`
	URL     string        `toml:"url,omitempty"       validate:"url,required"`
	Body    string        `toml:"body,omitempty"      validate:"required"`
	Headers HeadersConfig `toml:"headers,omitempty"   validate:"omitempty"`
}

type HeadersConfig map[string]string

// Notes:
// URL must be string here since we could craft an invalid url if it contains any tokens.
