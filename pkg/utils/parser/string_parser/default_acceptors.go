package string_parser

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/parser/state_machine"
	"regexp"
	"strings"
)

func StringAcceptor(valueToAccept string) state_machine.Acceptor[string] {
	return func(currentValue string) bool { return currentValue == valueToAccept }
}

func RegexpAcceptor(regexpToAccept string) state_machine.Acceptor[string] {
	if !strings.HasPrefix(regexpToAccept, "^") {
		regexpToAccept = fmt.Sprintf(`^%s`, regexpToAccept)
	}

	if !strings.HasSuffix(regexpToAccept, `$`) {
		regexpToAccept = fmt.Sprintf(`%s$`, regexpToAccept)
	}

	return func(currentValue string) bool {
		matched, _ := regexp.Match(regexpToAccept, []byte(currentValue))
		return matched
	}
}
