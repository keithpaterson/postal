package cacert

import "github.com/keithpaterson/postal/internal/util/test"

var (
	testValidCaCrt              string
	testValidKey                string
	testValidCertificateRequest string
	testValidCertificate        string
)

func LoadTestPrivateKeys() {
	testValidCaCrt = test.MustLoadTestString("testCA.pem")
	testValidKey = test.MustLoadTestString("unittest.key")
	testValidCertificateRequest = test.MustLoadTestString("unittest.csr")
	testValidCertificate = test.MustLoadTestString("unittest.pem")
}
