package native_test

import (
	"testing"

	"github.com/keithpaterson/postal/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSenderNative(t *testing.T) {
	logging.Disable()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sender Native Suite")
}
