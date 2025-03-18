package jwt

import "github.com/keithpaterson/postal/internal/util/test"

const (
	testHmacPrivateKey = `This is a terrible HMAC key.`
)

var (
	testRsaPrivateKey      string
	testEcdsa256PrivateKey string
	testEcdsa384PrivateKey string
	testEcdsa512PrivateKey string
	//testEddsaPrivateKey    string
)

func LoadTestPrivateKeys() {
	testRsaPrivateKey = test.MustLoadTestString("test_rsa")
	testEcdsa256PrivateKey = test.MustLoadTestString("test_ecdsa_256")
	testEcdsa384PrivateKey = test.MustLoadTestString("test_ecdsa_384")
	testEcdsa512PrivateKey = test.MustLoadTestString("test_ecdsa_512")
	//testEddsaPrivateKey    = test.MustLoadTestString("test_eddsa")
}
