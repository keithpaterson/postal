package resolver

import (
	"os"
)

type envResolver struct {
	resolverImpl
}

func newEnvResolver(root *rootResolver) *envResolver {
	return &envResolver{resolverImpl{root: root}}
}

func (r *envResolver) resolve(name string, token string) (string, bool) {
	if name != "env" {
		return token, false
	}

	envValue, ok := os.LookupEnv(token)
	if !ok {
		return token, false
	}

	return r.resolveValue(envValue), true
}
