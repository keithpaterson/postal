package config

import "slices"

const (
	// default JWT algorithm if it is not otherwise specified
	DefaultAlgorithm = "HS256"
)

// algorithm IDs
const (
	AlgNone JWTAlgorithm = iota
	AlgHS256
	AlgHS384
	AlgHS512
	AlgRS256
	AlgRS384
	AlgRS512
	AlgES256
	AlgES384
	AlgES512
	AlgPS256
	AlgPS384
	AlgPS512
	AlgEdDSA

	// must always be the last
	algMax
)

type JWTAlgorithm int

// JWT Algorithm names used in config files
var algorithmNames = []string{
	"None",
	"HS256", "HS384", "HS512",
	"RS256", "RS384", "RS512",
	"ES256", "ES384", "ES512",
	"PS256", "PS384", "P5512",
	//unsupported: "EDDSA",
}

// JWTConfig holds the the data used to generate a JSON Web Token.
//
// Claims are all the JWT clams as name=value pairs
//
// SigningKey is the key used to sign the token; it is not recommended to store this value in your configuration
// files, but you can supply it on the command-line
type JWTConfig struct {
	Header     JWTHeader `toml:"header,omitempty"      validate:"required"`
	Claims     JWTClaims `toml:"claims,omitempty"      validate:"required,dive,gt=0"`
	SigningKey string    `toml:"signing-key,omitempty" validate:"omitempty,gt=0"`
}

type JWTHeader struct {
	Alg string `toml:"alg,required" validate:"required,oneof=None HS256 HS384 HS512 RS256 RS384 RS512 ES256 ES384 ES512 PS256 PS384 P5512 EDDSA"`
}

type JWTClaims map[string]string

func (a JWTAlgorithm) String() string {
	return algorithmNames[a]
}

func (h JWTHeader) Algorithm() JWTAlgorithm {
	index := slices.Index(algorithmNames, h.Alg)
	if index < 0 || index >= int(algMax) {
		index = int(AlgNone)
	}
	return JWTAlgorithm(index)
}
