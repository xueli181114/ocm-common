package state_machine_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStateMachine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateMachine Suite")
}
