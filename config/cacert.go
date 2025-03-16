package config

import (
	"slices"
	"strconv"
)

const (
	CertPoolNone CertPoolType = iota
	CertPoolEmpty
	CertPoolSystem

	// must be last
	certPoolMax
)

const (
	CertPoolInvalid = CertPoolType(-1)
)

type CertPoolType int

func (t CertPoolType) String() string {
	if t < 0 || t >= certPoolMax {
		return strconv.Itoa(int(t))
	}
	return certPoolNames[t]
}

var certPoolNames = []string{"none", "empty", "system"}

// CacertConfig holds the configuration for Ca Certificates; If specified, this information is used in HTTP requests
type CacertConfig struct {
	// PoolName identifies which CACert pool to instantiate.  may be "none", in which case this struct is not processsed.
	// by default, PoolName is set to "none"
	PoolName string `toml:"pool,omitempty"                            validate:"omitempty,oneof=none empty system"`

	// CaCrt identifies the Root CA certificate data to use.
	// This can be specified as one of:
	//  "string:<data>" wehere <data> is the actual certificate in PEM format including the header and footer
	//  "file:<filename>"" where <filename> locates a file containing the certificate in PEM format.
	CaCrt string `toml:"ca-crt"                                       validate:"omitempty,gt=0"`

	// Certificates identifies any certificates you have registered with the Root CA
	// each Certificate is configured as <certificate,key> pairs using one of the formats:
	//  "string:<cert>,<key>"        : <cert> is the actual certificate in PEM format.
	//                                 <key> is the private key in PEM format.
	//                                 ensure that you include the header and footer for each.
	//  "file:<certfile>,<keyfile>"  : <certfile> locates a file containing the certificate.
	//                                 <keyfile> is the private key.
	//                                 both files are expected to be in PEM format
	//  "file:<basename>"            : short-hand notation for identifying cert,key filename pairs
	//                                 <basename> locates the certificate and key files. See CertificateFileExtensions for detail.
	Certificates []string `toml:"certificates,omitempty"              validate:"omitempty"`

	// CertificateFileExtensions allows you to override the default file extensions used to load the certificate and key files
	// when processing the short-hand notation for Certificates.
	// by default, this is set to [".pem", ".key"].
	//   For example: "file:foo" is equivalent to "file:foo.pem,foo.key"
	CertificateFileExtensions [2]string `toml:"file-ext,omitempty"    validate:"omitempty,len=2,dive,gt=1,startswith=."`
}

func newCacertConfig() CacertConfig {
	return CacertConfig{PoolName: "none"}
}

func (c CacertConfig) Pool() CertPoolType {
	index := slices.Index(certPoolNames, c.PoolName)
	if index < 0 || index >= int(certPoolMax) {
		return CertPoolInvalid
	}
	return CertPoolType(index)
}
