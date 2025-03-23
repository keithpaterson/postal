package cmd_test

import (
	"testing"

	"github.com/keithpaterson/postal/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	logging.Disable()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}
