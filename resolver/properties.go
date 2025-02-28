package resolver

type propertiesResolver struct {
	resolverImpl
}

func newPropertiesResolver(root *rootResolver) *propertiesResolver {
	return &propertiesResolver{resolverImpl{root: root}}
}

func (r *propertiesResolver) resolve(name string, token string) (string, bool) {
	if name != "prop" {
		return token, false
	}

	value, ok := r.root.config.Properties[token]
	if !ok {
		return token, false
	}
	strVal := value.(string)

	return r.resolveValue(strVal), true
}
