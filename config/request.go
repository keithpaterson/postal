package config

type RequestConfig struct {
	Method   string `toml:"method,omitempty"`
	URL      string `toml:"url,omitempty"`
	Body     string `toml:"body,omitempty"`
	MimeType string `toml:"mime-type,omitempty"`
}
