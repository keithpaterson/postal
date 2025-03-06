package jwt

import (
	"fmt"
	"postal/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func makeConfig(alg config.JWTAlgorithm, key string) config.JWTConfig {
	return config.JWTConfig{
		Header:     config.JWTHeader{Alg: alg.String()},
		SigningKey: key,
		Claims:     config.JWTClaims{"iss": "test", "sub": "test", "aud": "test", "exp": "123456788", "nbf": "23456879", "iat": "34567890", "jti": "test"},
	}
}

var _ = Describe("JWT", func() {
	var (
		emptyJWT               = config.JWTConfig{}
		headerOnlyJWT          = config.JWTConfig{Header: config.JWTHeader{Alg: "HS256"}}
		invalidSignkingKeyJWT1 = config.JWTConfig{Header: config.JWTHeader{Alg: "HS256"}, SigningKey: "string:"}
		invalidSignkingKeyJWT2 = config.JWTConfig{Header: config.JWTHeader{Alg: "HS256"}, SigningKey: "hex:"}
		invalidSignkingKeyJWT3 = config.JWTConfig{Header: config.JWTHeader{Alg: "HS256"}, SigningKey: "hex:this is not hex data"}
		noClaimsJWT            = config.JWTConfig{Header: config.JWTHeader{Alg: "HS256"}, SigningKey: "hex:01 02 03 04 05 06"}
		withClaimsJWT1         = config.JWTConfig{Header: config.JWTHeader{Alg: "HS256"}, SigningKey: "hex:01 02 03 04 05 06",
			Claims: config.JWTClaims{"iss": "foo", "abc": "bar", "xxx": "YYY", "and": "so-on"}}
		withClaimsJWT2 = config.JWTConfig{Header: config.JWTHeader{Alg: "HS256"}, SigningKey: "value without prefix",
			Claims: config.JWTClaims{"iss": "foo", "abc": "bar", "xxx": "YYY", "and": "so-on"}}

		noClaimsToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.nLjmntCckfwgWjjKNpKAilc_UwRffT0sOJTIUlxK9XM"
		withClaimsToken1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYmMiOiJiYXIiLCJhbmQiOiJzby1vbiIsImlzcyI6ImZvbyIsInh4eCI6IllZWSJ9.PBmywc18RuDgL4FgvVZ363vTBjl5unHHJsOPvuRN1r0"
		withClaimsToken2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYmMiOiJiYXIiLCJhbmQiOiJzby1vbiIsImlzcyI6ImZvbyIsInh4eCI6IllZWSJ9.dvsUkXV_ogaDSoxywque9Y5AFAr3dSlGtWQey6PDPts"
	)

	type expectation struct {
		jwtString string
		err       error
	}

	var (
		builder *jwtBuilder
	)

	BeforeEach(func() {
		builder = NewBuilder()
	})

	DescribeTable("MakeToken",
		func(jwtCfg config.JWTConfig, expect expectation) {
			// Act
			actual, err := builder.MakeToken(jwtCfg)

			// Assert
			if expect.err != nil {
				Expect(err).To(MatchError(expect.err))
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(expect.jwtString))
			}
		},
		Entry("empty config returns error", emptyJWT, expectation{"", ErrInvalidSigningMethod}),
		Entry("no signing key returns error", headerOnlyJWT, expectation{"", ErrNoSigningKey}),
		Entry("invalid signing key (no string value) returns error", invalidSignkingKeyJWT1, expectation{"", ErrNoSigningKey}),
		Entry("invalid signing key (no hex value) returns error (2)", invalidSignkingKeyJWT2, expectation{"", ErrNoSigningKey}),
		Entry("invalid signing key (invalid hex value) returns error", invalidSignkingKeyJWT3, expectation{"", ErrInvalidSigningValue}),
		Entry("with no claims generates token", noClaimsJWT, expectation{noClaimsToken, nil}),
		Entry("with claims generates token", noClaimsJWT, expectation{noClaimsToken, nil}),
		Entry("with claims generates token", withClaimsJWT1, expectation{withClaimsToken1, nil}),
		Entry("with claims generates token", withClaimsJWT2, expectation{withClaimsToken2, nil}),
	)

	Context("Algorithms", func() {
		DescribeTable("HMAC",
			func(alg config.JWTAlgorithm, keyval string) {
				// Arrange
				cfg := makeConfig(alg, fmt.Sprintf("string:%s", keyval))

				// Act
				actual, err := builder.MakeToken(cfg)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).ToNot(BeEmpty())
			},
			Entry(config.AlgHS256.String(), config.AlgHS256, testHmacPrivateKey),
			Entry(config.AlgHS384.String(), config.AlgHS384, testHmacPrivateKey),
			Entry(config.AlgHS512.String(), config.AlgHS512, testHmacPrivateKey),
		)

		DescribeTable("ECDSA/EdDSA",
			func(alg config.JWTAlgorithm, keyval string) {
				// Arrange
				cfg := makeConfig(alg, fmt.Sprintf("pemdata:%s", keyval))

				// Act
				actual, err := builder.MakeToken(cfg)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).ToNot(BeEmpty())
			},
			Entry(config.AlgES256.String(), config.AlgES256, testEcdsa256PrivateKey),
			Entry(config.AlgES384.String(), config.AlgES384, testEcdsa384PrivateKey),
			Entry(config.AlgES512.String(), config.AlgES512, testEcdsa512PrivateKey),
			//Entry(config.AlgES512.String(), config.AlgEdDSA, testEdDSAPrivateKey),
		)

		DescribeTable("RSS",
			func(alg config.JWTAlgorithm, keyval string) {
				// Arrange
				cfg := makeConfig(alg, fmt.Sprintf("pemdata:%s", keyval))

				// Act
				actual, err := builder.MakeToken(cfg)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).ToNot(BeEmpty())
			},
			Entry(config.AlgRS256.String(), config.AlgRS256, testRsaPrivateKey),
			Entry(config.AlgRS384.String(), config.AlgRS384, testRsaPrivateKey),
			Entry(config.AlgRS512.String(), config.AlgRS512, testRsaPrivateKey),
		)

		DescribeTable("PSS",
			func(alg config.JWTAlgorithm, keyval string) {
				// Arrange
				cfg := makeConfig(alg, fmt.Sprintf("pemdata:%s", keyval))

				// Act
				actual, err := builder.MakeToken(cfg)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).ToNot(BeEmpty())
			},
			Entry(config.AlgPS256.String(), config.AlgPS256, testRsaPrivateKey),
			Entry(config.AlgPS384.String(), config.AlgPS384, testRsaPrivateKey),
			Entry(config.AlgPS512.String(), config.AlgPS512, testRsaPrivateKey),
		)
	})
})
