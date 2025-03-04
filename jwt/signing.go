package jwt

import (
	"errors"
	"fmt"
	"os"
	"postal/config"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidKeyType       = errors.New("invalid key type")
	ErrNoSigningKey         = errors.New("signing key not provided")
	ErrInvalidSigningValue  = errors.New("invalid sigining value")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidPemData       = errors.New("invalid PEM data")
	ErrParsePEMFailed       = errors.New("failed to parse PEM data")
)

func (b *jwtBuilder) getSigningKey() (any, error) {
	if b.jwt.SigningKey == "" {
		return nil, ErrNoSigningKey
	}

	keyType, value, ok := strings.Cut(b.jwt.SigningKey, ":")
	if !ok {
		value = keyType
		keyType = "string"
	}
	keyType = strings.TrimSpace(keyType)
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, ErrNoSigningKey
	}
	return b.parseSigningKey(keyType, value)
}

func (b *jwtBuilder) getSigningMethod() (jwt.SigningMethod, error) {
	switch b.jwt.Header.Algorithm() {
	case config.AlgHS256:
		return jwt.SigningMethodHS256, nil
	case config.AlgHS384:
		return jwt.SigningMethodHS384, nil
	case config.AlgHS512:
		return jwt.SigningMethodHS512, nil
	case config.AlgRS256:
		return jwt.SigningMethodRS256, nil
	case config.AlgRS384:
		return jwt.SigningMethodRS384, nil
	case config.AlgRS512:
		return jwt.SigningMethodRS512, nil
	case config.AlgES256:
		return jwt.SigningMethodES256, nil
	case config.AlgES384:
		return jwt.SigningMethodES384, nil
	case config.AlgES512:
		return jwt.SigningMethodES512, nil
	case config.AlgPS256:
		return jwt.SigningMethodPS256, nil
	case config.AlgPS384:
		return jwt.SigningMethodPS384, nil
	case config.AlgPS512:
		return jwt.SigningMethodPS512, nil
	// unsupported:
	case config.AlgEdDSA:
		return jwt.SigningMethodEdDSA, nil
	default:
		return nil, fmt.Errorf("%w (%d: %s)", ErrInvalidSigningMethod, b.jwt.Header.Algorithm(), b.jwt.Header.Alg)
	}
}

func (b *jwtBuilder) parseSigningKey(keyType string, value string) (any, error) {
	var err error
	var rawData []byte
	switch keyType {
	case "string":
		rawData, err = b.fromString(value)
	case "hex":
		rawData, err = b.fromHexArray(value)
	case "pemfile":
		rawData, err = b.fromFile(value)
	case "pemdata":
		rawData = []byte(value)
	default:
		return nil, fmt.Errorf("%w '%s'", ErrInvalidKeyType, keyType)
	}

	if err != nil {
		return nil, err
	}
	return b.decodePemData(rawData)
}

func (b *jwtBuilder) fromString(input string) ([]byte, error) {
	return []byte(input), nil
}

func (b *jwtBuilder) fromHexArray(input string) ([]byte, error) {
	var err error
	var tmp int64
	values := strings.Split(input, " ")
	result := make([]byte, len(values))
	for index, value := range values {
		if tmp, err = strconv.ParseInt(value, 16, 8); err != nil {
			return nil, fmt.Errorf("%w '%s' at byte %d", ErrInvalidSigningValue, value, index)
		}
		result[index] = byte(tmp)
	}
	return result, nil
}

func (b *jwtBuilder) fromFile(filename string) ([]byte, error) {
	var err error
	var raw []byte
	if raw, err = os.ReadFile(filename); err != nil {
		return nil, err
	}
	return raw, nil
}

func (b *jwtBuilder) decodePemData(data []byte) (any, error) {
	switch b.jwt.Header.Algorithm() {
	case config.AlgHS256, config.AlgHS384, config.AlgHS512:
		return data, nil
	case config.AlgRS256, config.AlgRS384, config.AlgRS512:
		return jwt.ParseRSAPrivateKeyFromPEM(data)
	case config.AlgES256, config.AlgES384, config.AlgES512:
		return jwt.ParseECPrivateKeyFromPEM(data)
	case config.AlgPS256, config.AlgPS384, config.AlgPS512:
		return jwt.ParseRSAPrivateKeyFromPEM(data)
		// unsupported:
	case config.AlgEdDSA:
		return jwt.ParseEdPrivateKeyFromPEM(data)
	default:
		// should not happen
		return nil, ErrParsePEMFailed
	}
}
