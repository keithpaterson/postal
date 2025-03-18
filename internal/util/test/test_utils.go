// internal utilities only used during test
package test

import (
	"os"

	"github.com/onsi/gomega"
)

func MustLoadTestData(filename string) []byte {
	data, err := os.ReadFile("testdata/" + filename)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	return data
}

func MustLoadTestString(filename string) string {
	return string(MustLoadTestData(filename))
}
