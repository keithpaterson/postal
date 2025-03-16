package output

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/keithpaterson/go-tools/resolver"
	"github.com/keithpaterson/resweave-utils/logging"
	"go.uber.org/zap"
)

type responseResolver struct {
	resolver.ResolverImpl

	log  *zap.SugaredLogger
	resp *http.Response

	// store and cache (some) values
	body *cached
}

type cached struct {
	cached bool
	value  any
}

func newResponseResolver(resp *http.Response, log *zap.SugaredLogger) *responseResolver {
	return &responseResolver{log: log.Named("response resolver"), resp: resp, body: &cached{cached: false}}
}

func (r *responseResolver) Resolve(name string, token string) (string, bool) {
	if name != "response" {
		return token, false
	}

	return r.resolveToken(token)
}

func (r *responseResolver) resolveToken(token string) (string, bool) {
	var err error
	var value string
	switch strings.ToLower(token) {
	case "body":
		var body []byte
		if body, err = r.getBody(); err != nil {
			r.log.Errorw("resolveToken", "token", "body", logging.LogKeyError, err)
			return token, false
		}
		value = string(body)
	case "headers":
		value = r.getAllHeaders(r.resp.Header)
	case "status":
		value = r.resp.Status
	case "status-code", "statuscode":
		value = strconv.Itoa(r.resp.StatusCode)
	case "content-length", "contentlength":
		value = strconv.FormatInt(r.resp.ContentLength, 10)
	}

	// header is a special case..
	if strings.HasPrefix(token, "header=") {
		var ok bool
		var name string
		if _, name, ok = strings.Cut(token, "="); !ok || name == "" {
			return token, false // invalid token
		}
		value = r.getHeader(name)
	}

	return value, true
}

func (r *responseResolver) getBody() ([]byte, error) {
	if r.body.hasValue() {
		return r.body.getBytes(), nil
	}

	body, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return nil, err
	}
	r.body.set(body)
	return body, nil
}

func (r *responseResolver) getHeader(name string) string {
	var ok bool
	var hdrValues []string
	value := ""
	if hdrValues, ok = r.resp.Header[name]; ok {
		value = strings.Join(hdrValues, ",")
	}
	return value
}

func (r *responseResolver) getAllHeaders(headers http.Header) string {
	result := make([]string, len(headers))
	index := 0
	for name, value := range headers {
		result[index] = fmt.Sprintf("%s=%s", name, strings.Join(value, ","))
		index++
	}
	return strings.Join(result, "; ")
}

// cached struct functions

func (c *cached) hasValue() bool {
	return c.cached
}

func (c *cached) getBytes() []byte {
	if c.cached {
		return c.value.([]byte)
	}
	return nil
}

func (c *cached) set(value any) {
	if c.cached {
		panic("ERROR: invalid to set an already cached value")
	}
	c.cached = true
	c.value = value
}
