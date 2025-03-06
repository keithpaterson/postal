package resolver

import (
	"postal/config"
	"postal/logging"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var (
	regToken = regexp.MustCompile(`(?U)\${(.*)}`) // matches a "${name:value}" token, captures "name:token"
)

// A Resolver is used to convert property tokens in the form '${token-name} into actual data.
type Resolver interface {
	// resolves all tokens in the string and returns the result
	Resolve(input string) (string, error)
}

type rootResolver struct {
	log       *zap.SugaredLogger
	config    *config.Config
	resolvers resolversMap
}

func NewResolver(config *config.Config) *rootResolver {
	root := &rootResolver{config: config, log: logging.NamedLogger("resolver")}

	dtr := newDateTimeResolver(root)
	root.resolvers = resolversMap{
		"env":      newEnvResolver(root),
		"prop":     newPropertiesResolver(root),
		"datetime": dtr,
		"date":     dtr,
		"time":     dtr,
		"epoch":    dtr,
	}
	return root
}

type resolver interface {
	resolve(name string, token string) (string, bool)
}

type resolverImpl struct {
	root *rootResolver
	log  *zap.SugaredLogger
}
type resolversMap map[string]resolver

func (ri *resolverImpl) resolveValue(value string) string {
	return ri.root.Resolve(value)
}

//
// rootResolver
//

func (r *rootResolver) Resolve(input string) string {
	tokens := regToken.FindAllString(input, -1)
	if tokens == nil {
		return input
	}

	replacements := r.resolveTokenValues(tokens)
	replacer := strings.NewReplacer(replacements...)
	return replacer.Replace(input)
}

func (r *rootResolver) resolveTokenValues(tokens []string) []string {
	result := make([]string, len(tokens)*2)
	for index, token := range tokens {
		matches := regToken.FindStringSubmatch(token)
		if len(matches) < 2 {
			continue // shouldn't happen but this appears to be an invalid token
		}
		if resolved, ok := r.resolveToken(matches[1]); ok {
			offset := index * 2
			result[offset] = token
			result[offset+1] = resolved
		}
	}
	return result
}

func (r *rootResolver) resolveToken(token string) (string, bool) {
	// expect "name:value"
	// e.g. "prop:foo", "env:MY_ENV_VAR", "datetime:now.(RSS3339)", "time:now.(TimeOnly) + 30s"
	var name, value string
	var ok bool
	if name, value, ok = strings.Cut(token, ":"); !ok {
		value = name
		name = "prop"
	}
	name = strings.ToLower(name)

	resolver, ok := r.resolvers[name]
	if !ok {
		// log a warning?
		return token, false
	}
	return resolver.resolve(name, value)
}
