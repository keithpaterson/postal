// package jwt builds JWT token strings for requests
//
// Token data is pulled from the config.JWTConfig structure, e.g.
//
// Here is a contrived example for your amusement:
//
//	cfg := config.JwtConfig{Header: {Alg: "HS256"}, Claims: {"iss": "blah", "exp": "1234567", ...}, SigningKey: []byte{your signature key}}
//	builder := jwt.NewBuilder()
//	token, err := builder.MakeToken(cfg)
//
// SigningKey is your unique signature used to sign the token and is required at runtime in order to generate a valid token.
// The config structure does support storing the key in your config files; note that doing so is for convenience only and is
// normally considered poor security practice.
//
// You should prefer passing the key in at the command-line in order to better protect it.
//
// Assuming no errors, 'token' will contain the complete, signed JWT token.
package jwt
