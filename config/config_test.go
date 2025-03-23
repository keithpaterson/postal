package config

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// I couldn't get the matchers to work properly so ...
func compareConfig(actual *Config, expected *Config) {
	Expect(actual.Request).To(Equal(expected.Request))

	compareProperties(actual.Properties, expected.Properties)
	compareJWTPayloads(actual.JWT.Claims, expected.JWT.Claims)
}

func compareProperties(actual Properties, expected Properties) {
	if expected == nil {
		Expect(actual).To(BeNil())
	} else {
		Expect(len(actual)).To(Equal(len(expected)))
		for k, v := range expected {
			Expect(k).To(BeKeyOf(actual))
			Expect(actual[k]).To(Equal(v))
		}
	}
}

func compareJWTPayloads(actual JWTClaims, expected JWTClaims) {
	Expect(actual).To(Equal(expected))
	if expected == nil {
		Expect(actual).To(BeNil())
	} else {
		Expect(len(actual)).To(Equal(len(expected)))
		for k, v := range expected {
			Expect(k).To(BeKeyOf(actual))
			Expect(actual[k]).To(Equal(v))
		}
	}
}

var _ = Describe("Config", func() {
	type expectation struct {
		cfg *Config
		err error
	}

	var (
		cfg *Config

		expectError = func(err error) expectation {
			return expectation{cfg: nil, err: err}
		}
		expectConfig = func(cfg *Config) expectation {
			return expectation{cfg: cfg, err: nil}
		}
	)

	BeforeEach(func() {
		cfg = NewConfig()
	})

	DescribeTable("Load",
		func(reader io.Reader, expect expectation) {
			// Act
			err := cfg.Load(reader)

			// Assert
			if expect.err != nil {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
				compareConfig(cfg, expect.cfg)
			}
		},
		Entry("nil reader returns error", nil, expectError(ErrNilReader)),
		Entry("invalid data returns error", bytes.NewReader([]byte("invalid data")), expectError(toml.ParseError{})),
		Entry("valid request returns config", bytes.NewReader(validRequestData), expectConfig(&Config{Request: validReq})),
		Entry("valid properties returns config", bytes.NewReader(validPropertiesData), expectConfig(&Config{Request: newRequestConfig(), Properties: validProps})),
		Entry("valid jwt returns config", bytes.NewReader(validJWTData), expectConfig(&Config{Request: newRequestConfig(), JWT: validJWT})),
	)
})
