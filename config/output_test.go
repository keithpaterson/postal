package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OutputConfig", func() {
	Context("OutFormat", func() {
		DescribeTable("String",
			func(out OutFormat, expect string) {
				// Act
				actual := out.String()

				// Assert
				Expect(actual).To(Equal(expect))
			},
			Entry(nil, OutFmtRaw, "raw"),
			Entry(nil, OutFmtText, "text"),
			Entry(nil, OutFormat(-100), "undefined"),
			Entry(nil, OutFormat(100), "undefined"),
		)
	})

	DescribeTable("OutFormat()",
		func(value string, expect OutFormat) {
			// Arrange
			cfg := OutputConfig{Format: value}

			// Act
			actual := cfg.OutFormat()

			// Assert
			Expect(actual).To(Equal(expect))
		},
		Entry(nil, "raw", OutFmtRaw),
		Entry(nil, "text", OutFmtText),
		Entry(nil, "", OutFmtText),
		Entry(nil, "something", OutFmtText),
		Entry(nil, "anything", OutFmtText),
		Entry(nil, "undefined", OutFmtText),
	)
})
