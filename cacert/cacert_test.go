package cacert

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cacert", func() {
	Describe("Builder", func() {
		DescribeTable("WithPool",
			func(poolType certPoolType, expect error) {
				// Arrange
				builder := Builder(poolType)

				// Act
				pool, err := builder.Pool()

				// Assert
				if expect != nil {
					Expect(pool).To(BeNil())
					Expect(err).To(MatchError(expect))
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(pool).ToNot(BeNil())
				}
			},
			Entry("invalid pool returns error", certPoolType(100), ErrCacertInvalidPool),
			Entry("custom pool returns pool", CertPoolCustom, nil),
			Entry("system pool returns pool", CertPoolSystem, nil),
		)

		DescribeTable("WithCert",
			func(poolType certPoolType, data string, expect error) {
				// Arrange
				builder := Builder(poolType)

				// Act
				pool, err := builder.WithCert(data).Pool()

				// Assert
				if expect != nil {
					Expect(pool).To(BeNil())
					Expect(err).To(MatchError(expect))
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(pool).ToNot(BeNil())
				}
			},
			Entry("invalid pool type returns error", certPoolType(100), "not valid data", ErrCacertInvalidPool),
			Entry("invalid cert data returns error", CertPoolCustom, "not valid data", ErrCacertInvalidCert),
			Entry("invalid cert string data returns error", CertPoolCustom, "string:not valid data", ErrCacertInvalidCert),
			Entry("invalid cert file returns error", CertPoolCustom, "file:testdata/does-not-exist.baha", ErrCacertInvalidCert),
			Entry("valid cert from string returns pool", CertPoolCustom, "string:"+testValidCert, nil),
			Entry("valid cert from file returns pool", CertPoolCustom, "file:testdata/unittest.pem", nil),
		)
	})
})
