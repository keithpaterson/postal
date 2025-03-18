package jwt_test

import (
	"testing"

	"github.com/keithpaterson/postal/jwt"
	"github.com/keithpaterson/postal/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModels(t *testing.T) {
	logging.Disable()
	RegisterFailHandler(Fail)
	jwt.LoadTestPrivateKeys()
	RunSpecs(t, "JWT Suite")
}
