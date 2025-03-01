package config

import (
	"bytes"

	"github.com/BurntSushi/toml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestConfig", func() {
	type expectation struct {
		req *RequestConfig
		err error
	}
	var (
		validData = []byte(`method = "POST"
url = "http://test.io/fing/fang/fong"
body = 'json:{"this":"that","then":123}'`)
	)
	var (
		validReq = RequestConfig{Method: "POST", URL: "http://test.io/fing/fang/fong", Body: `json:{"this":"that","then":123}`}
	)
	DescribeTable("from TOML",
		func(input []byte, expect expectation) {
			// Act
			var actual RequestConfig
			_, err := toml.NewDecoder(bytes.NewReader(input)).Decode(&actual)

			// Assert
			if expect.err != nil {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(*expect.req))
			}
		},
		Entry("invalid data fails", []byte("invalid data"), expectation{req: nil, err: toml.ParseError{}}),
		Entry("valid data succeeds", validData, expectation{req: &validReq, err: nil}),
	)
})
