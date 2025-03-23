package config

// RuntimeConfig holds any values that may change runtime behaviour, such as the DryRun flag.
type RuntimeConfig struct {
	// DryRun is true when the program should validate inputs but not perform any actual action.
	DryRun bool
}
