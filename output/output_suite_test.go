package output_test

import (
	"testing"

	"github.com/keithpaterson/postal/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOutput(t *testing.T) {
	logging.Disable()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Output Suite")
}
