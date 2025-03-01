package config

const (
	defaultAlgorithm = "HS256"
)

type JWTConfig struct {
	Header  JWTHeader  `toml:"header,omitempty"`
	Payload JWTPayload `toml:"payload,omitempty"`
}

type JWTHeader struct {
	Alg string `toml:"alg,required"`
}

type JWTPayload map[string]string

// TODO(keithpaterson): add some validation:
// see https://datatracker.ietf.org/doc/html/rfc7518#section-3 for the list of supported algorithms
// - header.alg must be one of [
//     HS256, HS384, HS512,
//     RS256, RS384, RS512,
//     ES256, ES384, ES512,
//     PS256, PS384, PS512,
//     none]
// - payload must have at least one key
// - payload keys must be strings
// - payload values must be strings (?) or at least can be converted to string
// - are there any truly REQUIRED keys?
