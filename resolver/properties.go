package resolver

import (
	"fmt"
)

type propertiesResolver struct {
	resolverImpl
}

func newPropertiesResolver(root *rootResolver) *propertiesResolver {
	return &propertiesResolver{resolverImpl{root: root, log: root.log.Named("properties")}}
}

func (r *propertiesResolver) resolve(name string, token string) (string, bool) {
	if name != "prop" {
		return token, false
	}

	value, ok := r.root.config.Properties[token]
	if !ok {
		return token, false
	}
	strVal := fmt.Sprintf("%v", value) //value.(string)

	return r.resolveValue(strVal), true
}
