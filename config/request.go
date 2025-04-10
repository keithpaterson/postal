package config

// RequestConfig stores the properties of the request.
type RequestConfig struct {
	// Method is the HTTP Method (GET, POST, PUT, PATCH, DELETE)
	Method string `toml:"method,omitempty"            validate:"method,required"`
	URL    string `toml:"url,omitempty"               validate:"url,required"`
	// Body specifies any data to put in the request body, may be empty.
	// Supported formats:
	//  "json:{json-data}" : must be a valid json blob after resolving.
	//                       sets "Content-Type" header to "application/json"
	//  "file:<filename>"  : load body data from a file.  <filename> must be valid after resolving.
	//                       you should include the "Content-Type" header (it is not inferred)
	Body string `toml:"body,omitempty"                validate:"omitempty,gt=0"`
	// Headers are name=value pairs for any request headers you need.
	Headers HeadersConfig `toml:"headers,omitempty"   validate:"omitempty,dive,gt=0"`
}

// TODO(keithpaterson): add timeout, backoff info either here or in another struct (or a child struct?)

type HeadersConfig map[string]string

func newRequestConfig() RequestConfig {
	return RequestConfig{
		Headers: make(HeadersConfig),
	}
}
