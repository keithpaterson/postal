package cacert

import (
	"github.com/keithpaterson/postal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func certConfig(values ...string) config.CacertConfig {
	cfg := config.NewConfig()
	cfg.Cacert.Certificates = values
	return cfg.Cacert
}

func certConfigWithExt(ext [2]string, values ...string) config.CacertConfig {
	cfg := config.NewConfig()
	cfg.Cacert.Certificates = values
	cfg.Cacert.CertificateFileExtensions = ext
	return cfg.Cacert
}

var _ = Describe("Certificates", func() {
	Describe("Builder", func() {
		Context("WithCertificate", func() {
			DescribeTable("cert value parsing",
				func(value string, expect error) {
					// Arrange
					builder := Certificates()

					// Act
					actual, err := builder.WithCertificate(value).Build()

					// Assert
					if expect != nil {
						Expect(err).To(MatchError(expect))
						Expect(actual).To(BeEmpty())
					} else {
						Expect(err).ToNot(HaveOccurred())
						Expect(actual).ToNot(BeEmpty())
					}
				},
				Entry("valid certificate data succeeds", testValidCertificate+","+testValidKey, nil),
				Entry("valid certificate files succeeds", "file:testdata/unittest.pem,testdata/unittest.key", nil),
				Entry("valid certificate file basename succeeds", "file:testdata/unittest", nil),
				Entry("invalid certificate data missing comma returns error", "no comma in data", ErrInvalidCertificateDataFormat),
				Entry("invalid certificate data extra comma returns error", "string:too,many,commas", ErrInvalidCertificateDataFormat),
				Entry("invalid certificate file extra comma returns error", "file:too,many,commas", ErrInvalidCertificateFileFormat),
				Entry("invalid certificate data returns error", testValidCertificateRequest+","+testValidKey, ErrInvalidCertificate),
				Entry("invalid certificate file returns error", "file:testdata/unittest.crt,testdata/unittest.key", ErrInvalidCertificate),
			)
			It("will succeed with no certificates", func() {
				// Arrange
				builder := Certificates()

				// Act
				actual, err := builder.Build()

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(BeEmpty())
			})
		})

		Context("WithConfig", func() {

			DescribeTable("cert value parsing",
				func(cfg config.CacertConfig, expect error) {
					// Arrange
					builder := Certificates()

					// Act
					actual, err := builder.FromConfig(cfg).Build()

					// Assert
					if expect != nil {
						Expect(err).To(MatchError(expect))
						Expect(actual).To(BeEmpty())
					} else {
						Expect(err).ToNot(HaveOccurred())
						Expect(actual).ToNot(BeEmpty())
					}
				},
				Entry("valid certificate data succeeds", certConfig(testValidCertificate+","+testValidKey), nil),
				Entry("valid certificate files succeeds", certConfig("file:testdata/unittest.pem,testdata/unittest.key"), nil),
				Entry("valid certificate file basename succeeds", certConfig("file:testdata/unittest"), nil),
				Entry("invalid certificate data missing comma returns error", certConfig("no comma in data"), ErrInvalidCertificateDataFormat),
				Entry("invalid certificate data extra comma returns error", certConfig("string:too,many,commas"), ErrInvalidCertificateDataFormat),
				Entry("invalid certificate file extra comma returns error", certConfig("file:too,many,commas"), ErrInvalidCertificateFileFormat),
				Entry("invalid certificate data returns error", certConfig(testValidCertificateRequest+","+testValidKey), ErrInvalidCertificate),
				Entry("invalid certificate file returns error", certConfig("file:testdata/unittest.crt,testdata/unittest.key"), ErrInvalidCertificate),
			)
			It("will succeed with no certificates", func() {
				// Arrange
				builder := Certificates()

				// Act
				actual, err := builder.FromConfig(certConfig()).Build()

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(BeEmpty())
			})
		})

		DescribeTable("custom file extensions",
			func(ext [2]string, value string, expect error) {
				// Arrange
				builder := Certificates()

				// Act
				actual, err := builder.FromConfig(certConfigWithExt(ext, value)).Build()

				// Assert
				if expect != nil {
					Expect(err).To(MatchError(expect))
					Expect(actual).To(BeEmpty())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(actual).ToNot(BeEmpty())
				}
			},
			Entry("missing files", [2]string{".foo", ".bar"}, "file:testdata/unittest", ErrInvalidCertificate),
			Entry("bad file data", [2]string{".csr", ".key"}, "file:testdata/unittest", ErrInvalidCertificate),
		)
	})
})
