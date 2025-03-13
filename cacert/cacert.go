package cacert

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/keithpaterson/postal/logging"
	ulog "github.com/keithpaterson/resweave-utils/logging"
)

var (
	ErrCacertInvalidPool = errors.New("invalid cert pool")
	ErrCacertInvalidCert = errors.New("invalid cert")

	errCacertNotInitialized = errors.New("not initialized")
	errCacertInvalidPEMData = errors.New("invalid PEM data")
)

const (
	CertPoolCustom certPoolType = iota
	CertPoolSystem
)

type certPoolType int

type cacertBuilder struct {
	log  *zap.SugaredLogger
	pool *x509.CertPool

	err error
}

func Builder(poolType certPoolType) *cacertBuilder {
	builder := &cacertBuilder{log: logging.NamedLogger("cacert builder")}
	return builder.WithPool(poolType)
}

func (b *cacertBuilder) WithPool(poolType certPoolType) *cacertBuilder {
	var err error
	switch poolType {
	case CertPoolCustom:
		b.pool = x509.NewCertPool()
	case CertPoolSystem:
		var pool *x509.CertPool
		if pool, err = x509.SystemCertPool(); err != nil {
			// in the error case keep whatever the previous pool was
			b.err = err
			fmt.Println("ERROR: could not load system cert pool:", err)
			b.log.Errorw("WithPool", ulog.LogKeyError, err)
		} else {
			b.err = nil
			b.pool = pool
		}
	default:
		// in the default case keep whatever the previous pool was
		fmt.Println("WARNING: unsupported pool type:", poolType)
	}

	return b
}

func (b *cacertBuilder) WithCert(certValue string) *cacertBuilder {
	if b.pool == nil {
		if b.err != nil {
			b.err = errCacertNotInitialized
		}
		fmt.Println("ERROR: you must choose a pool before adding certs")
		return b
	}

	dataType, data, ok := strings.Cut(certValue, ":")
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
			b.err = err
			return b
		}
	}

	if !b.pool.AppendCertsFromPEM(pemData) {
		b.err = errCacertInvalidPEMData
	}

	return b
}

func (b *cacertBuilder) Pool() (*x509.CertPool, error) {
	if b.pool == nil {
		return nil, fmt.Errorf("%w: %w", ErrCacertInvalidPool, b.err)
	}
	if b.err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCacertInvalidCert, b.err)
	}

	return b.pool, nil
}
