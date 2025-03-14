package cacert

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/keithpaterson/postal/config"
)

var (
	ErrInvalidCertificateDataFormat = errors.New("invalid certificate format: expected 'type:value,value'")
	ErrInvalidCertificateFileFormat = errors.New("invalid certificate format: expected 'type:value[,value]'")
	ErrInvalidCertificate           = errors.New("invalid certificate")
)

type certsBuilder struct {
	certs      []string
	extensions [2]string
}

func Certificates() *certsBuilder {
	return &certsBuilder{extensions: [2]string{".pem", ".key"}}
}

func (b *certsBuilder) FromConfig(cfg config.CacertConfig) *certsBuilder {
	b.certs = cfg.Certificates

	if !b.allEmpty(cfg.CertificateFileExtensions[:]...) {
		b.extensions = cfg.CertificateFileExtensions
	}
	return b
}

func (b *certsBuilder) WithCertificate(value string) *certsBuilder {
	b.certs = append(b.certs, value)
	return b
}

func (b *certsBuilder) Build() ([]tls.Certificate, error) {
	if len(b.certs) < 1 {
		return nil, nil
	}

	var err error
	certificates := make([]tls.Certificate, len(b.certs))
	for index, certValue := range b.certs {
		dataType, dataValue, ok := strings.Cut(certValue, ":")
		if !ok {
			dataValue = dataType
			dataType = "string"
		}

		if certificates[index], err = b.resolveCertificate(dataType, dataValue); err != nil {
			return nil, err
		}
	}
	return certificates, nil
}

func (b *certsBuilder) resolveCertificate(dataType string, dataValue string) (cert tls.Certificate, err error) {
	switch dataType {
	case "string", "pemdata":
		var crtData, keyData []byte
		if crtData, keyData, err = b.splitData(dataValue); err != nil {
			return
		}
		if cert, err = tls.X509KeyPair(crtData, keyData); err != nil {
			err = fmt.Errorf("%w: %w", ErrInvalidCertificate, err)
			return
		}
	case "file", "pemfile":
		var crtFile, keyFile string
		if crtFile, keyFile, err = b.splitFilenames(dataValue); err != nil {
			return
		}
		if cert, err = tls.LoadX509KeyPair(crtFile, keyFile); err != nil {
			err = fmt.Errorf("%w: %w", ErrInvalidCertificate, err)
			return
		}
	}
	return
}

func (b *certsBuilder) splitData(dataValue string) ([]byte, []byte, error) {
	values := strings.Split(dataValue, ",")
	if len(values) != 2 {
		return nil, nil, ErrInvalidCertificateDataFormat
	}
	return []byte(values[0]), []byte(values[1]), nil
}

func (b *certsBuilder) splitFilenames(dataValue string) (string, string, error) {
	values := strings.Split(dataValue, ",")
	switch len(values) {
	case 1:
		basename := strings.TrimSpace(values[0])
		return basename + b.extensions[0], basename + b.extensions[1], nil
	case 2:
		return strings.TrimSpace(values[0]), strings.TrimSpace(values[1]), nil
	default:
		return "", "", ErrInvalidCertificateFileFormat
	}
}

// could be a utility function but leave it here for now
func (b *certsBuilder) allEmpty(value ...string) bool {
	for _, v := range value {
		if v != "" {
			return false
		}
	}
	return true
}
