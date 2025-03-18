package cacert_test

import (
	"testing"

	"github.com/keithpaterson/postal/cacert"
	"github.com/keithpaterson/postal/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModels(t *testing.T) {
	logging.Disable()
	RegisterFailHandler(Fail)
	cacert.LoadTestPrivateKeys()
	RunSpecs(t, "Cacert Suite")
}
