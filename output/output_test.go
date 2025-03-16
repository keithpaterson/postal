package output

import (
	"github.com/keithpaterson/postal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Output", func() {
	DescribeTable("NewOutputter",
		func(name string, expected interface{}) {
			// Arrange
			cfg := config.NewConfig()
			out := &cfg.Output
			out.Format = name

			// Act
			outputter := NewOutputter(cfg)

			// Assert
			Expect(outputter).To(BeAssignableToTypeOf(expected))
		},
		Entry("valid raw", "raw", &rawOutputter{}),
		Entry("valid text", "text", &textOutputter{}),
		Entry("invalid foo", "foo", &textOutputter{}),
		Entry("invalid bar", "bar", &textOutputter{}),
		Entry("invalid empty", "", &textOutputter{}),
	)
})
