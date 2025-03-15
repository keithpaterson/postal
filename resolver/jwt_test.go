package resolver

import (
	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWT Resolver", func() {
	type expectation struct {
		value string
		ok    bool
	}
	var (
		resolver *jwtResolver
	)
	BeforeEach(func() {
		resolver = newJWTResolver(logging.NamedLogger("test"), config.NewConfig().JWT)
	})

	DescribeTable("resolve",
		func(tokens []string, expect []expectation) {
			// Arrange
			resolver.cfg.Header.Alg = config.AlgHS256.String()
			resolver.cfg.Claims = config.JWTClaims{"iss": "test", "test": "yup,this-is;a-test"}
			resolver.cfg.SigningKey = "booga booga"

			// Act & Assert
			for index, token := range tokens {
				actual, ok := resolver.Resolve("jwt", token)
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
		actual, ok := resolver.Resolve("crumb", "fling")

		// Assert
		Expect(ok).To(BeFalse())
		Expect(actual).To(Equal("fling"))
	})
})
