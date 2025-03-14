package cacert

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/keithpaterson/postal/config"
)

type cacertParser struct {
	cfg config.CacertConfig
}

func FromConfig(cfg config.CacertConfig) cacertParser {
	return cacertParser{cfg: cfg}
}

func (p cacertParser) GetCertificatePool() (*x509.CertPool, error) {
	return Pool().FromConfig(p.cfg).Build()
}

func (p cacertParser) GetCertificates() ([]tls.Certificate, error) {
	return Certificates().FromConfig(p.cfg).Build()
}
