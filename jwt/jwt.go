package jwt

import (
	"postal/config"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type jwtBuilder struct {
	log *zap.SugaredLogger
	jwt config.JWTConfig
}

func NewBuilder(log *zap.SugaredLogger) *jwtBuilder {
	return &jwtBuilder{log: log.Named("jwt")}
}

func (b *jwtBuilder) MakeToken(cfg config.JWTConfig) (string, error) {
	b.jwt = cfg
	claims := make(jwt.MapClaims)
	for key, value := range b.jwt.Claims {
		claims[key] = value
	}

	var err error
	var method jwt.SigningMethod
	if method, err = b.getSigningMethod(); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(method, claims)
	return b.signToken(token)
}

func (b *jwtBuilder) signToken(token *jwt.Token) (string, error) {
	var err error
	var key any
	if key, err = b.getSigningKey(); err != nil {
		return "", err
	}

	var result string
	if result, err = token.SignedString(key); err != nil {
		return "", err
	}
	return result, nil
}
