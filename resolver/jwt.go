package resolver

import (
	"github.com/keithpaterson/postal/jwt"

	"github.com/keithpaterson/resweave-utils/logging"
)

type jwtResolver struct {
	resolverImpl
}

func newJWTResolver(root *rootResolver) *jwtResolver {
	return &jwtResolver{resolverImpl{root: root, log: root.log.Named("jwt")}}
}

// only accepts "jwt:token"
func (r *jwtResolver) resolve(name string, token string) (string, bool) {
	if name != "jwt" || token != "token" {
		return token, false
	}

	builder := jwt.NewBuilder()
	value, err := builder.MakeToken(r.root.config.JWT)
	if err != nil {
		r.log.Errorw("resolve", logging.LogKeyError, err)
		return token, false
	}
	return value, true
}
