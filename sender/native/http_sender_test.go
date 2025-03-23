package native

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/keithpaterson/postal/cacert"
	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/logging"
	"github.com/keithpaterson/resweave-utils/utility/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Native Sender", func() {
	DescribeTable("configureTLS()",
		func(pool config.CertPoolType, expect error) {
			// Arrange
			cfg := config.NewConfig()
			cfg.Cacert.PoolName = pool.String()

			sender := &httpSender{cfg: cfg, log: logging.NamedLogger("test")}
			client := http.DefaultClient

			// Act
			err := sender.configureTLS(client)

			// Assert
			if expect != nil {
				Expect(err).To(MatchError(expect))
			} else {
				Expect(err).ToNot(HaveOccurred())
			}
		},
		Entry(config.CertPoolNone.String(), config.CertPoolNone, nil),
		Entry(config.CertPoolEmpty.String(), config.CertPoolEmpty, nil),
		Entry(config.CertPoolSystem.String(), config.CertPoolSystem, nil),
		Entry("invalid", config.CertPoolType(-1), ErrInvalidCert),
	)
	DescribeTable("getBodyData()",
		func(body string, expect error) {
			// Arrange
			cfg := config.NewConfig()
			cfg.Request.Body = body

			sender := &httpSender{cfg: cfg, log: logging.NamedLogger("test")}

			// Act
			actual, err := sender.getBodyData()

			// Assert
			if expect != nil {
				Expect(err).To(MatchError(expect))
				Expect(actual).To(BeNil())
			} else {
				Expect(err).ToNot(HaveOccurred())
				if body == "" {
					Expect(actual).To(BeNil())
				} else {
					Expect(actual).ToNot(BeNil())
				}
			}
		},
		Entry("empty", "", nil),
		Entry("invalid format", "invalid format", ErrInvalidBodySpec),
		Entry("invalid json", "json:not json", ErrInvalidBody),
		Entry("valid json", `json:{"valid": "json"}`, nil),
		Entry("file not found", "file:file-not-found", ErrInvalidBody),
		Entry("invalid type", "unsupported:blah blah", ErrUnsupportedBodySpec),
	)

	Context("execute()", func() {
		DescribeTable("error conditions",
			func(cfg *config.Config, expect error) {
				// Act
				err := sendHttp(cfg, logging.NamedLogger("test"))

				// Assert
				if expect != nil {
					Expect(err).To(MatchError(expect))
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
			},
			Entry("bad method", requestCfg("not valid", "none", `json:{"name":"test"}`), errors.New(`net/http: invalid method "NOT VALID"`)),
			Entry("bad json", requestCfg("get", "none", `json:this is not json`), ErrInvalidBody),
			Entry("invalid cert", requestCfg("get", "invalid pool", `json:{"name":"test"}`), cacert.ErrInvalidPool),
		)

		It("succeeds with valid data", func() {
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
			err := sendHttp(cfg, logging.NamedLogger("test"))

			// Assert
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func requestCfg(method string, certPool string, body string) *config.Config {
	cfg := config.NewConfig()
	cfg.Request.Method = strings.ToUpper(method)
	cfg.Request.Body = body
	cfg.Request.URL = "http://test.io/foo"
	cfg.Request.Headers["content-type"] = "application/json"
	cfg.Cacert.PoolName = certPool
	return cfg
}
