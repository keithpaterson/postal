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
// The signing key value specified in configs or on the command-line must conform to the format:
//
//		 "type:value", where type can be one of:
//		   "string"  : the value is a text string.  Note that using "value" implies a string key
//		   "hex"     : the value is a space-separated list of hexadecimal values,
//	                e.g. "hex:01 02 03 04 05" -> "[0x01, 0x02, 0x03, 0x04, 0x05]"
//		   "pemfile" : the value is a PEM filename
//		   "pemdata" : the value is a PEM signature (e.g. the value you would find in the PEM file)
//
// Assuming no errors, 'token' will contain the complete, signed JWT token.
package jwt
