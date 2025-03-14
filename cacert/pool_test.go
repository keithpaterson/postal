package cacert

import (
	"github.com/keithpaterson/postal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func poolConfig(poolType config.CertPoolType, cacrt string) config.CacertConfig {
	return config.CacertConfig{PoolName: poolType.String(), CaCrt: cacrt}
}

var _ = Describe("Pool", func() {
	Describe("Builder", func() {
		DescribeTable("WithPool",
			func(poolType config.CertPoolType, expect error) {
				// Arrange
				builder := Pool()

				// Act
				pool, err := builder.WithPool(poolType).Build()

				// Assert
				if expect != nil {
					Expect(pool).To(BeNil())
					Expect(err).To(MatchError(expect))
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(pool).ToNot(BeNil())
				}
			},
			Entry("invalid pool returns error", config.CertPoolType(100), ErrInvalidPool),
			Entry("empty pool returns pool", config.CertPoolEmpty, nil),
			Entry("system pool returns pool", config.CertPoolSystem, nil),
		)
		It("will return nil, nil for the none pool", func() {
			// Arrange
			builder := Pool()

			// Act
			pool, err := builder.WithPool(config.CertPoolNone).Build()

			// Assert
			Expect(pool).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})

		DescribeTable("WithCaCrt",
			func(poolType config.CertPoolType, data string, expect error) {
				// Arrange
				builder := Pool()

				// Act
				pool, err := builder.WithPool(poolType).WithCACrt(data).Build()

				// Assert
				if expect != nil {
					Expect(err).To(MatchError(expect))
					Expect(pool).To(BeNil())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(pool).ToNot(BeNil())
				}
			},
			Entry("invalid pool type returns error", config.CertPoolType(100), "not valid data", ErrInvalidPool),
			Entry("invalid cert data returns error", config.CertPoolEmpty, "not valid data", ErrInvalidCert),
			Entry("invalid cert string data returns error", config.CertPoolEmpty, "string:not valid data", ErrInvalidCert),
			Entry("invalid cert file returns error", config.CertPoolEmpty, "file:testdata/does-not-exist.baha", ErrInvalidCert),
			Entry("valid cert from string returns pool", config.CertPoolEmpty, "string:"+testValidCaCrt, nil),
			Entry("valid cert from file returns pool", config.CertPoolEmpty, "file:testdata/unittest.pem", nil),
		)

		DescribeTable("WithConfig",
			func(cfg config.CacertConfig, expect error) {
				// Arrange
				builder := Pool()

				// Act
				pool, err := builder.FromConfig(cfg).Build()

				// Assert
				if cfg.Pool() == config.CertPoolNone {
					if expect != nil {
						Expect(err).To(MatchError(expect))
					} else {
						Expect(err).ToNot(HaveOccurred())

					}
					Expect(pool).To(BeNil())
					return
				}

				if expect != nil {
					Expect(err).To(MatchError(expect))
					Expect(pool).To(BeNil())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(pool).ToNot(BeNil())
				}
			},
			Entry("empty config returns error", config.CacertConfig{}, ErrInvalidPool),
			Entry("none pool returns nil, nil", poolConfig(config.CertPoolNone, ""), nil),
			Entry("none pool with valid cacrt returns nil, nil", poolConfig(config.CertPoolNone, testValidCaCrt), nil),
			Entry("none pool with invalid cacrt returns nil, nil", poolConfig(config.CertPoolNone, "not a cacrt"), nil),
			Entry("empty pool returns pool", poolConfig(config.CertPoolEmpty, ""), nil),
			Entry("empty pool with valid cacrt returns pool", poolConfig(config.CertPoolEmpty, testValidCaCrt), nil),
			Entry("empty pool with invalid cacrt returns error", poolConfig(config.CertPoolEmpty, "not a cacrt"), ErrInvalidCert),
			Entry("system pool returns pool", poolConfig(config.CertPoolSystem, ""), nil),
			Entry("system pool with valid cacrt returns pool", poolConfig(config.CertPoolSystem, testValidCaCrt), nil),
			Entry("system pool with invalid cacrt returns error", poolConfig(config.CertPoolSystem, "not a cacrt"), ErrInvalidCert),
		)
	})
})
