package string_scanner_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStringscanner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stringscanner Suite")
}
