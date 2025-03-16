package output

import (
	"net/http"
	"strings"

	"github.com/keithpaterson/go-tools/resolver"
	"github.com/keithpaterson/postal/config"
	"go.uber.org/zap"
)

type template struct {
	cfg *config.Config
	log *zap.SugaredLogger
}

func newTemplate(cfg *config.Config, log *zap.SugaredLogger) *template {
	return &template{cfg: cfg, log: log.Named("template")}
}

func (t *template) Apply(resp *http.Response) string {
	template := t.cfg.Output.Template
	if template == "" {
		template = "${response:body}"
	}

	return t.resolve(t.fixTemplate(template), resp)
}

func (t *template) fixTemplate(input string) string {
	replacer := strings.NewReplacer(templateReplacements...)
	return replacer.Replace(input)
}

func (t *template) resolve(input string, resp *http.Response) string {
	return resolver.NewResolver(&resolver.ResolverConfig{Properties: resolver.Properties(t.cfg.Properties)}).
		WithStandardResolvers().
		WithResolver("response", newResponseResolver(resp, t.log)).
		Resolve(input)
}

var (
	// common replacements for non-token text we might find in the template specification
	// TODO(keithpaterson): make tab expansion configurable (via config.OutputConfig.Options?)
	templateReplacements = []string{
		"\\n", "\n", // newlines become new lines
		"\\t", "    ", // tabs expand to 4 spaces
	}
)
