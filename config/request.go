package config

type RequestConfig struct {
	Method   string `toml:"method,omitempty"    validate:"required,oneof=POST GET PUT PATCH DELETE"`
	URL      string `toml:"url,omitempty"       validate:"url,required"`
	Body     string `toml:"body,omitempty"      validate:"required"`
	MimeType string `toml:"mime-type,omitempty" validate:"omitempty"`
}
