package resolver

import (
	"postal/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWT Resolver", func() {
	type expectation struct {
		value string
		ok    bool
	}
	var (
		cfg      *config.Config
		root     *rootResolver
		resolver *jwtResolver
	)
	BeforeEach(func() {
		cfg = config.NewConfig()
		root = NewResolver(cfg)
		resolver = newJWTResolver(root)
	})

	DescribeTable("resolve",
		func(tokens []string, expect []expectation) {
			// Arrange
			cfg.JWT.Header.Alg = config.AlgHS256.String()
			cfg.JWT.Claims = config.JWTClaims{"iss": "test", "test": "yup,this-is;a-test"}
			cfg.JWT.SigningKey = "booga booga"

			// Act & Assert
			for index, token := range tokens {
				actual, ok := resolver.resolve("jwt", token)
				Expect(actual).To(Equal(expect[index].value))
				Expect(ok).To(Equal(expect[index].ok))
			}
		},
		Entry("not 'token', no change", []string{"input"}, []expectation{{"input", false}}),
		Entry("is 'token', generates token", []string{"token"}, []expectation{{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0IiwidGVzdCI6Inl1cCx0aGlzLWlzO2EtdGVzdCJ9.WMToCOCoctMMd-iTzPO2WZPJJj_xbh2NlfMjttY0SkE", true}}),
	)

	It("will ignore requests for invalid token name", func() {
		// Arrange

		// Act
		actual, ok := resolver.resolve("crumb", "fling")

		// Assert
		Expect(ok).To(BeFalse())
		Expect(actual).To(Equal("fling"))
	})
})
