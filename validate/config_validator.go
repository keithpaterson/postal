package validate

import (
	"encoding/json"
	"strings"

	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/resolver"

	"github.com/keithpaterson/resweave-utils/utility/rw"
)

// ValidateConfig resolves tokens in the config data, validates the result and returns a new (valid) config object
//
// If validation fails, returns the original config object with the error.
func ValidateConfig(cfg *config.Config) (*config.Config, error) {
	var raw []byte
	var err error
	if raw, err = json.Marshal(cfg); err != nil {
		return cfg, err
	}

	res := resolver.NewResolver(cfg)
	resolvedStr := res.Resolve(string(raw))

	var resolved config.Config
	if err = rw.UnmarshalJson(strings.NewReader(resolvedStr), &resolved); err != nil {
		return cfg, err
	}

	if err = ValidateStruct(resolved); err != nil {
		return cfg, err
	}

	return &resolved, nil
}
