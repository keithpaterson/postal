package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWTConfig", func() {
	Context("JWTAlgorithm", func() {
		DescribeTable("String",
			func(alg JWTAlgorithm, expect string) {
				// Act
				actual := alg.String()

				// Assert
				Expect(actual).To(Equal(expect))
			},
			Entry(nil, AlgHS256, "hs256"),
			Entry(nil, AlgHS384, "hs384"),
			Entry(nil, AlgHS512, "hs512"),
			Entry(nil, AlgRS256, "rs256"),
			Entry(nil, AlgRS384, "rs384"),
			Entry(nil, AlgRS512, "rs512"),
			Entry(nil, AlgES256, "es256"),
			Entry(nil, AlgES384, "es384"),
			Entry(nil, AlgES512, "es512"),
			Entry(nil, AlgPS256, "ps256"),
			Entry(nil, AlgPS384, "ps384"),
			Entry(nil, AlgPS512, "ps512"),
			Entry(nil, JWTAlgorithm(-100), "-100"),
			Entry(nil, JWTAlgorithm(100), "100"),
		)
	})

	Context("JWTHeader", func() {
		DescribeTable("Algorithm",
			func(alg string, expect JWTAlgorithm) {
				// Arrange
				hdr := JWTHeader{Alg: alg}

				// Act
				actual := hdr.Algorithm()

				// Assert
				Expect(actual).To(Equal(expect))
			},
			Entry(nil, "hs256", AlgHS256),
			Entry(nil, "hs384", AlgHS384),
			Entry(nil, "hs512", AlgHS512),
			Entry(nil, "rs256", AlgRS256),
			Entry(nil, "rs384", AlgRS384),
			Entry(nil, "rs512", AlgRS512),
			Entry(nil, "es256", AlgES256),
			Entry(nil, "es384", AlgES384),
			Entry(nil, "es512", AlgES512),
			Entry(nil, "ps256", AlgPS256),
			Entry(nil, "ps384", AlgPS384),
			Entry(nil, "ps512", AlgPS512),
			Entry(nil, "eddsa", AlgNone),
			Entry(nil, "", AlgNone),
			Entry(nil, "something", AlgNone),
			Entry(nil, "anything", AlgNone),
		)
	})
})
