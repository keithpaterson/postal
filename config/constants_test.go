package config

// used during config_test
var (
	// raw "file" data
	validRequestData = []byte(`
[request]
	method = "POST"
	url = "http://test.io"
	body = 'json:{"this":"that","then":123}'`)

	validPropertiesData = []byte(`
[properties]
	one = "1"
	two = 2
	three = "three"`)

	validJWTData = []byte(`
[jwt]
	[jwt.header]
		alg = "HS256"
	
	[jwt.claims]
		iss = "foo"
		sub = "this=x,that=y,those=z"
		aud = "urn:testything"
		exp = "987654321"
		foo = "bar"
		bar = "foo"`)

	// what the "file" data should decode into:
	validReq   = RequestConfig{Method: "POST", URL: "http://test.io", Body: `json:{"this":"that","then":123}`, Headers: make(HeadersConfig)}
	validProps = Properties{"one": "1", "two": int64(2), "three": "three"}
	validJWT   = JWTConfig{
		Header: JWTHeader{Alg: "HS256"},
		Claims: JWTClaims{"iss": "foo", "sub": "this=x,that=y,those=z", "aud": "urn:testything", "exp": "987654321", "foo": "bar", "bar": "foo"}}
)
