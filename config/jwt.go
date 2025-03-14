package config

import (
	"slices"
	"strconv"
	"strings"
)

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
	//unsupported: AlgEdDSA

	// must always be the last
	algMax
)

type JWTAlgorithm int

// JWT Algorithm names used in config files
var algorithmNames = []string{
	"none",
	"hs256", "hs384", "hs512",
	"rs256", "rs384", "rs512",
	"es256", "es384", "es512",
	"ps256", "ps384", "ps512",
	//unsupported: "eddsa",
}

// JWTConfig holds the the data used to generate a JSON Web Token.
//
// Claims are all the JWT clams as name=value pairs
//
// SigningKey is the key used to sign the token; it is not recommended to store the actual value in your configuration
// files, but you can supply it via the command-line.
type JWTConfig struct {
	Header JWTHeader `toml:"header,omitempty"      validate:"required"`

	// Claims are all the JWT claims as name=value pairs
	Claims JWTClaims `toml:"claims,omitempty"      validate:"required,dive,gt=0"`

	// SigningKey is used to sign the token.
	// Accepted formats are:
	//  "string:text"      : generally used for HMAC signatures, this is your signature key in raw (readable) text format
	//  "hex:01 02 03 ..." : hexadecimal representation of your signature (space-separated values), this is your signature key in readable format
	//  "file:filename"    : locates a file containing your signature in PEM format
	//  "pemdata:data"     : provides the PEM formatted signature as a text block.
	//
	// Generally is isn't recommended that you store your signing key in the configuration.
	// Storing the filename is generally considered safer.
	// You can provide the signing key via command-line using one of these format strings.
	SigningKey string `toml:"signing-key,omitempty" validate:"omitempty,gt=0"`
}

type JWTHeader struct {
	Alg string `toml:"alg,required" validate:"required,oneof=none hs256 hs384 hs512 rs256 rs384 rs512 es256 es384 es512 ps256 ps384 ps512"`
}

type JWTClaims map[string]string

func (a JWTAlgorithm) String() string {
	if a < 0 || a >= algMax {
		return strconv.Itoa(int(a))
	}
	return algorithmNames[a]
}

func (h JWTHeader) Algorithm() JWTAlgorithm {
	index := slices.Index(algorithmNames, strings.ToLower(h.Alg))
	if index < 0 {
		index = int(AlgNone)
	}
	return JWTAlgorithm(index)
}
