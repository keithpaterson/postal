package sender_test

import (
	"testing"

	"github.com/keithpaterson/postal/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSender(t *testing.T) {
	logging.Disable()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sender Suite")
}
