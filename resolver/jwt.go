package resolver

import (
	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/jwt"

	"github.com/keithpaterson/go-tools/resolver"
	"github.com/keithpaterson/resweave-utils/logging"
	"go.uber.org/zap"
)

type jwtResolver struct {
	resolver.ResolverImpl

	log *zap.SugaredLogger
	cfg *config.JWTConfig
}

func newJWTResolver(log *zap.SugaredLogger, cfg *config.JWTConfig) *jwtResolver {
	return &jwtResolver{log: log.Named("jwt"), cfg: cfg}
}

// only accepts "jwt:token"
func (r *jwtResolver) Resolve(name string, token string) (string, bool) {
	if name != "jwt" || token != "token" {
		return token, false
	}

	builder := jwt.NewBuilder()
	value, err := builder.MakeToken(*r.cfg)
	if err != nil {
		r.log.Errorw("resolve", logging.LogKeyError, err)
		return token, false
	}
	return value, true
}
