package cmd

import (
	"maps"

	"github.com/keithpaterson/postal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	noArgs = testData{[]string{}, makeParsedConfig(nil, nil, nil, nil, nil)}
)

type testData struct {
	args []string
	cfg  *config.Config
}

func makeParsedConfig(request *config.RequestConfig, jwt *config.JWTConfig, cacert *config.CacertConfig, properties config.Properties, output *config.OutputConfig) *config.Config {
	// start with the base config
	cfg := config.NewConfig()

	if request != nil {
		cfg.Request = *request
	}
	if jwt != nil {
		cfg.JWT = *jwt
	}
	if cacert != nil {
		cfg.Cacert = *cacert
	}
	maps.Copy(cfg.Properties, properties)
	if output != nil {
		cfg.Output = *output
	}

	return cfg
}

var _ = Describe("SendCmd", func() {
	// TODO(keithpaterson): consider inducing flag errors somehow and testing those error paths
	DescribeTable("parseConfig",
		func(data testData, expect error) {
			// Arrange
			cmd := NewSendCommand()
			cmd.ParseFlags(data.args)

			parser := &sendCmdParser{}

			// Act
			cfg, err := parser.parseConfig(cmd)

			// Assert
			if expect != nil {
				Expect(err).To(MatchError(expect))
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).To(Equal(data.cfg))
			}
		},
		Entry("no args", noArgs, nil),
		// config file tests
		Entry("missing config file returns error", testData{[]string{"-c", "not-found"}, noArgs.cfg}, ErrInvalidConfigFile),
		Entry("invalid config file returns error", testData{[]string{"-c", "testdata/invalid.cfg"}, noArgs.cfg}, ErrInvalidConfigFile),
		Entry("valid config file succeeds", testData{[]string{"-c", "testdata/valid.cfg"}, makeParsedConfig(&config.RequestConfig{Method: "GET", URL: "https://test.io/test", Headers: config.HeadersConfig{}}, nil, nil, nil, nil)}, nil),
		// property tests
		Entry("invalid properties returns error", testData{[]string{"-p", "not valid"}, noArgs.cfg}, ErrInvalidPropertyValue),
		Entry("valid properties returns error", testData{[]string{"-p", "foo=bar"}, makeParsedConfig(nil, nil, nil, config.Properties{"foo": "bar"}, nil)}, nil),
		// request tests
		Entry("valid request method succeeds", testData{[]string{"-m", "PUT"}, makeParsedConfig(&config.RequestConfig{Method: "PUT", Headers: config.HeadersConfig{}}, nil, nil, nil, nil)}, nil),
		Entry("valid request url succeeds", testData{[]string{"-u", "http://test.io"}, makeParsedConfig(&config.RequestConfig{URL: "http://test.io", Headers: config.HeadersConfig{}}, nil, nil, nil, nil)}, nil),
		Entry("valid request body succeeds", testData{[]string{"-b", "json:{}"}, makeParsedConfig(&config.RequestConfig{Body: "json:{}", Headers: config.HeadersConfig{}}, nil, nil, nil, nil)}, nil),
		// headers tests
		Entry("invalid request header returns error", testData{[]string{"-H", "missing equals sign"}, noArgs.cfg}, ErrInvalidHeader),
		Entry("one request header succeeds", testData{[]string{"-H", "this=that,them"}, makeParsedConfig(&config.RequestConfig{Headers: config.HeadersConfig{"this": "that,them"}}, nil, nil, nil, nil)}, nil),
		Entry("two request headers succeeds", testData{[]string{"-H", "this=that,them", "-H", "another=header"}, makeParsedConfig(&config.RequestConfig{Headers: config.HeadersConfig{"this": "that,them", "another": "header"}}, nil, nil, nil, nil)}, nil),
		// jwt tests
		Entry("valid jwt algorithm is stored", testData{[]string{"-a", "ES256"}, makeParsedConfig(nil, &config.JWTConfig{Header: config.JWTHeader{Alg: "ES256"}, Claims: make(config.JWTClaims)}, nil, nil, nil)}, nil),
		Entry("invalid jwt algorithm is stored", testData{[]string{"-a", "anything"}, makeParsedConfig(nil, &config.JWTConfig{Header: config.JWTHeader{Alg: "anything"}, Claims: make(config.JWTClaims)}, nil, nil, nil)}, nil),
		Entry("jwt signing key is stored", testData{[]string{"--signing-key", "validated elsewhere"}, makeParsedConfig(nil, &config.JWTConfig{Header: config.JWTHeader{Alg: "hs256"}, SigningKey: "validated elsewhere", Claims: make(config.JWTClaims)}, nil, nil, nil)}, nil),
		Entry("invalid jwt claim returns error", testData{[]string{"--jwt", "not valid"}, noArgs.cfg}, ErrInvalidJWTClaim),
		Entry("one jwt claim succeeds", testData{[]string{"--jwt", "foo=bar"}, makeParsedConfig(nil, &config.JWTConfig{Header: config.JWTHeader{Alg: "hs256"}, Claims: config.JWTClaims{"foo": "bar"}}, nil, nil, nil)}, nil),
		Entry("two jwt claim succeeds", testData{[]string{"--jwt", "foo=bar", "--jwt", "this=that,those"}, makeParsedConfig(nil, &config.JWTConfig{Header: config.JWTHeader{Alg: "hs256"}, Claims: config.JWTClaims{"foo": "bar", "this": "that,those"}}, nil, nil, nil)}, nil),
	)
})
