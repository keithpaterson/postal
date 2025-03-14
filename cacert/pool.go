package cacert

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/keithpaterson/postal/config"
)

var (
	ErrInvalidPool = errors.New("invalid cert pool")
	ErrInvalidCert = errors.New("invalid cert")

	errInvalidPoolType = errors.New("invalid pool type")
	errInvalidPEMData  = errors.New("invalid PEM data")
)

type poolBuilder struct {
	cfg config.CacertConfig
}

func Pool() *poolBuilder {
	builder := &poolBuilder{}
	return builder
}

func (b *poolBuilder) FromConfig(cfg config.CacertConfig) *poolBuilder {
	b.cfg = cfg
	return b
}

func (b *poolBuilder) WithPool(poolType config.CertPoolType) *poolBuilder {
	b.cfg.PoolName = poolType.String()
	return b
}

func (b *poolBuilder) WithPoolName(name string) *poolBuilder {
	b.cfg.PoolName = name
	return b
}

func (b *poolBuilder) WithCACrt(crtValue string) *poolBuilder {
	b.cfg.CaCrt = crtValue
	return b
}

func (b *poolBuilder) Build() (*x509.CertPool, error) {
	pool, err := b.resolvePool()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidPool, err)
	}

	if pool == nil {
		// this is valid when pool type is 'none'
		return nil, nil
	}

	if pool, err = b.resolveCaCrt(pool); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidCert, err)
	}
	return pool, nil
}

func (b *poolBuilder) resolvePool() (*x509.CertPool, error) {
	switch b.cfg.Pool() {
	case config.CertPoolNone:
		return nil, nil
	case config.CertPoolEmpty:
		return x509.NewCertPool(), nil
	case config.CertPoolSystem:
		var pool *x509.CertPool
		var err error
		if pool, err = x509.SystemCertPool(); err != nil {
			return nil, fmt.Errorf("failed to load system certificate pool: %w", err)
		}
		return pool, nil
	default:
		return nil, errInvalidPoolType
	}
}

func (b *poolBuilder) resolveCaCrt(pool *x509.CertPool) (*x509.CertPool, error) {
	if b.cfg.CaCrt == "" {
		return pool, nil
	}

	dataType, data, ok := strings.Cut(b.cfg.CaCrt, ":")
	if !ok {
		data = dataType
		dataType = "string"
	}

	var err error
	var pemData []byte
	switch dataType {
	case "string", "pemdata":
		pemData = []byte(data)
	case "file", "pemfile":
		if pemData, err = os.ReadFile(data); err != nil {
			return nil, fmt.Errorf("failed to read ca-crt file: %w", err)
		}
	}

	if !pool.AppendCertsFromPEM(pemData) {
		return nil, errInvalidPEMData
	}

	return pool, nil

}
