package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CacertConfig", func() {
	Context("CacertPoolType", func() {
		DescribeTable("String",
			func(t CertPoolType, expect string) {
				// Act
				actual := t.String()

				// Assert
				Expect(actual).To(Equal(expect))
			},
			Entry(nil, CertPoolNone, "none"),
			Entry(nil, CertPoolEmpty, "empty"),
			Entry(nil, CertPoolSystem, "system"),
			Entry(nil, CertPoolType(-100), "-100"),
			Entry(nil, CertPoolType(100), "100"),
		)
	})

	DescribeTable("Pool",
		func(name string, expect CertPoolType) {
			// Arrange
			cfg := CacertConfig{PoolName: name}

			// Act
			actual := cfg.Pool()

			// Assert
			Expect(actual).To(Equal(expect))
		},
		Entry(nil, "none", CertPoolNone),
		Entry(nil, "empty", CertPoolEmpty),
		Entry(nil, "system", CertPoolSystem),
		Entry(nil, "", CertPoolInvalid),
		Entry(nil, "caterpillar", CertPoolInvalid),
		Entry(nil, "anything else", CertPoolInvalid),
	)
})
