package sql_parser

import (
	. "github.com/openshift-online/ocm-common/pkg/utils/parser/state_machine"
	. "github.com/openshift-online/ocm-common/pkg/utils/parser/string_parser"
)

const (
	braceTokenFamily     = "BRACE"
	opTokenFamily        = "OP"
	logicalOpTokenFamily = "LOGICAL"
	columnTokenFamily    = "COLUMN"

	othersTokenFamily      = "OTHERS"
	valueTokenFamily       = "VALUE"
	quotedValueTokenFamily = "QUOTED"
	openBrace              = "OPEN_BRACE"
	closedBrace            = "CLOSED_BRACE"
	comma                  = "COMMA"
	column                 = "COLUMN"
	value                  = "VALUE"
	quotedValue            = "QUOTED_VALUE"
	eq                     = "EQ"
	notEq                  = "NOT_EQ"
	gt                     = "GREATER_THAN"
	lt                     = "LESS_THAN"
	gte                    = "GREATER_THAN_OR_EQUAL"
	lte                    = "LESS_THAN_OR_EQUAL"
	like                   = "LIKE"
	ilike                  = "ILIKE"
	in                     = "IN"
	listOpenBrace          = "LIST_OPEN_BRACE"
	quotedValueInList      = "QUOTED_VALUE_IN_LIST"
	valueInList            = "VALUE_IN_LIST"
	and                    = "AND"
	or                     = "OR"
	not                    = "NOT"

	// Define the names of the tokens to be parsed

	jsonbFamily           = "JSONB"                    // Each JSONB token will be associated to the JSONB family
	jsonbField            = "JSON_FIELD"               // Each JSONB field
	jsonbArrow            = "JSONB_ARROW"              // The JSONB arrow token (->)
	jsonbToString         = "JSONB_TOSTRING"           // The JSONB to-string token (->>)
	jsonbContains         = "@>"                       // The JSONB @> token
	jsonbFieldToStringify = "JSONB_FIELD_TO_STRINGIFY" // The field that will contain the `string` value, ie: ->> FIELD
)

func BasicSQLGrammar() Grammar {
	grammar := Grammar{
		Tokens: []TokenDefinition{
			{Name: openBrace, StateData: braceTokenFamily, Acceptor: StringAcceptor(`(`)},
			{Name: closedBrace, StateData: braceTokenFamily, Acceptor: StringAcceptor(`)`)},
			{Name: column, StateData: columnTokenFamily, Acceptor: RegexpAcceptor(`(?i)[A-Z][A-Z0-9_.]*`)},
			{Name: value, StateData: valueTokenFamily, Acceptor: RegexpAcceptor(`[^'() ]*`)},
			{Name: quotedValue, StateData: quotedValueTokenFamily, Acceptor: RegexpAcceptor(`'([^']|\\')*'`)},
			{Name: eq, StateData: opTokenFamily, Acceptor: StringAcceptor(`=`)},
			{Name: gt, StateData: opTokenFamily, Acceptor: StringAcceptor(`>`)},
			{Name: lt, StateData: opTokenFamily, Acceptor: StringAcceptor(`<`)},
			{Name: gte, StateData: opTokenFamily, Acceptor: StringAcceptor(`>=`)},
			{Name: lte, StateData: opTokenFamily, Acceptor: StringAcceptor(`<=`)},
			{Name: comma, Acceptor: StringAcceptor(`,`)},
			{Name: notEq, StateData: opTokenFamily, Acceptor: StringAcceptor(`<>`)},
			{Name: like, StateData: opTokenFamily, Acceptor: RegexpAcceptor(`(?i)LIKE`)},
			{Name: ilike, StateData: opTokenFamily, Acceptor: RegexpAcceptor(`(?i)ILIKE`)},
			{Name: in, StateData: opTokenFamily, Acceptor: RegexpAcceptor(`(?i)IN`)},
			{Name: listOpenBrace, StateData: braceTokenFamily, Acceptor: StringAcceptor(`(`)},
			{Name: quotedValueInList, StateData: quotedValueTokenFamily, Acceptor: RegexpAcceptor(`'([^']|\\')*'`)},
			{Name: valueInList, StateData: valueTokenFamily, Acceptor: RegexpAcceptor(`[^'() ]*`)},
			{Name: and, StateData: logicalOpTokenFamily, Acceptor: RegexpAcceptor(`(?i)AND`)},
			{Name: or, StateData: logicalOpTokenFamily, Acceptor: RegexpAcceptor(`(?i)OR`)},
			{Name: not, StateData: logicalOpTokenFamily, Acceptor: RegexpAcceptor(`(?i)NOT`)},
			{Name: jsonbArrow, StateData: jsonbFamily, Acceptor: StringAcceptor(`->`)},
			{Name: jsonbField, StateData: jsonbFamily, Acceptor: RegexpAcceptor(`'([^']|\\')*'`)},
			{Name: jsonbToString, StateData: jsonbFamily, Acceptor: StringAcceptor(`->>`)},
			{Name: jsonbContains, StateData: jsonbFamily, Acceptor: StringAcceptor(`@>`)},
			{Name: jsonbFieldToStringify, StateData: jsonbFamily, Acceptor: RegexpAcceptor(`'([^']|\\')*'`)},
		},
		Transitions: []TokenTransitions{
			{TokenName: StartState, ValidTransitions: []string{column, openBrace}},
			{TokenName: openBrace, ValidTransitions: []string{column, openBrace}},
			{TokenName: column, ValidTransitions: []string{gt, lt, gte, lte, eq, notEq, like, ilike, in, not, jsonbArrow}},
			{TokenName: eq, ValidTransitions: []string{quotedValue, value}},
			{TokenName: notEq, ValidTransitions: []string{quotedValue, value}},
			{TokenName: gt, ValidTransitions: []string{quotedValue, value}},
			{TokenName: lt, ValidTransitions: []string{quotedValue, value}},
			{TokenName: lte, ValidTransitions: []string{quotedValue, value}},
			{TokenName: gte, ValidTransitions: []string{quotedValue, value}},
			{TokenName: like, ValidTransitions: []string{quotedValue, value}},
			{TokenName: ilike, ValidTransitions: []string{quotedValue, value}},
			{TokenName: quotedValue, ValidTransitions: []string{or, and, closedBrace, EndState}},
			{TokenName: value, ValidTransitions: []string{or, and, closedBrace, EndState}},
			{TokenName: closedBrace, ValidTransitions: []string{or, and, closedBrace, EndState}},
			{TokenName: and, ValidTransitions: []string{column, openBrace}},
			{TokenName: or, ValidTransitions: []string{column, openBrace}},
			{TokenName: not, ValidTransitions: []string{in}},
			{TokenName: in, ValidTransitions: []string{listOpenBrace}},
			{TokenName: listOpenBrace, ValidTransitions: []string{quotedValueInList, valueInList}},
			{TokenName: quotedValueInList, ValidTransitions: []string{comma, closedBrace}},
			{TokenName: valueInList, ValidTransitions: []string{comma, closedBrace}},
			{TokenName: comma, ValidTransitions: []string{quotedValueInList, valueInList}},
			{TokenName: jsonbArrow, ValidTransitions: []string{jsonbField}},
			{TokenName: jsonbField, ValidTransitions: []string{jsonbArrow, jsonbToString, jsonbContains}},
			{TokenName: jsonbToString, ValidTransitions: []string{jsonbFieldToStringify}},
			{TokenName: jsonbFieldToStringify, ValidTransitions: []string{eq, notEq, like, ilike, in, not}},
			{TokenName: jsonbContains, ValidTransitions: []string{quotedValue}},
		},
	}

	return grammar
}
