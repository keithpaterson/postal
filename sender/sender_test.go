package sender

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/resweave-utils/utility/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sender", func() {
	DescribeTable("toSenderType",
		func(name string, expect SenderType, withError bool) {
			// Act
			actual, err := toSenderType(name)

			// Assert
			if withError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(actual).To(Equal(expect))
		},
		Entry(nil, NativeSenderName, NativeSender, false),
		Entry(nil, CurlSenderName, CurlSender, false),
		Entry(nil, "foo", SenderType(-1), true),
		Entry(nil, "", SenderType(-1), true),
		Entry(nil, "anything", SenderType(-1), true),
	)

	DescribeTable("NewNamedSender",
		func(name string, expect error) {
			// Arrange
			// Act
			actual, err := NewNamedSender(name)

			// Assert
			if expect != nil {
				Expect(err).To(HaveOccurred())
				Expect(actual).To(BeNil())
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).ToNot(BeNil())
			}
		},
		Entry(NativeSenderName, NativeSenderName, nil),
		Entry(CurlSenderName, CurlSenderName, nil),
		Entry("empty name", "", ErrInvalidSender),
		Entry("unsupported name", "foo", ErrInvalidSender),
	)

	DescribeTable("Send",
		func(id SenderType, expect error) {
			// Arrange
			host, tearDown := test.HttpService().
				WithMethod("GET").
				WithPath("/foo").
				ReturnStatusCode(http.StatusOK).
				Start()
			defer tearDown()

			cfg := config.NewConfig()
			cfg.Request = config.RequestConfig{
				URL:    fmt.Sprintf("%s/foo", host),
				Method: "GET",
			}

			// Act
			sender, err := NewSender(id)
			Expect(err).ToNot(HaveOccurred())
			err = sender.Send(cfg)

			// Assert
			if expect != nil {
				Expect(err).To(MatchError(expect))
			} else {
				Expect(err).ToNot(HaveOccurred())
			}
		},
		Entry("native", NativeSender, nil),
		Entry("curl", CurlSender, errors.New("not implemented")),
		Entry("invalid id (100)", SenderType(100), ErrInvalidSender),
	)

	Context("Send With Verify", func() {
		It("will fail if the config can't be verified", func() {
			// Arrange
			cfg := config.NewConfig()
			cfg.Request = config.RequestConfig{
				URL:    "http://test.io/foo",
				Method: "CAT", // this is invalid
			}

			// Act
			sender, err := NewSender(NativeSender)
			Expect(err).ToNot(HaveOccurred())
			err = sender.Send(cfg)

			// Assert
			Expect(err).To(HaveOccurred())
		})
	})
})
