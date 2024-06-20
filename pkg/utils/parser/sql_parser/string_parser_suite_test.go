package sql_parser_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStringParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StringParser Suite")
}
