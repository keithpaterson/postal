package config

// RuntimeConfig holds any values that may change runtime behaviour, such as the DryRun flag.
//
// Runtime configuration is not persisted or loaded (e.g. from configuration files)
//
// Note that Runtime configuration is not affected by token replacing; in other words these properties
// should never contain tokens, except where explicitly noted.
type RuntimeConfig struct {
	// DryRun is true when the program should validate inputs but not perform any actual action.
	DryRun bool
}
