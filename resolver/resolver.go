package resolver

import (
	"github.com/keithpaterson/postal/config"

	"github.com/keithpaterson/go-tools/resolver"
	"github.com/keithpaterson/postal/logging"
	"go.uber.org/zap"
)

type wrapResolver struct {
	log *zap.SugaredLogger
	cfg *config.Config
}

func NewResolver(cfg *config.Config) *wrapResolver {
	return &wrapResolver{log: logging.NamedLogger("resolver"), cfg: cfg}
}

func (r *wrapResolver) Resolve(input string) string {
	return resolver.NewResolver(&resolver.ResolverConfig{Properties: resolver.Properties(r.cfg.Properties)}).
		WithStandardResolvers().
		WithResolver("jwt", newJWTResolver(r.log, &r.cfg.JWT)).
		Resolve(input)
}
