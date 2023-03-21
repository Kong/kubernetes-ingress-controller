package atc

import (
	"net/http"
	"strconv"
	"strings"
)

type FieldType int

const (
	FieldTypeInt = iota
	FieldTypeString
	FieldTypeSingleIP
	FieldTypeIPCIDR
)

type LHS interface {
	FieldType() FieldType
	// TODO(naming): use a better name for this method? "String" is too gerneral
	String() string
	ExtractValue(*http.Request) Literal
}

type BinaryOperator string

var (
	OpEqual        BinaryOperator = "=="
	OpNotEqual     BinaryOperator = "!="
	OpRegexMatch   BinaryOperator = "~"
	OpPrefixMatch  BinaryOperator = "^="
	OpSuffixMatch  BinaryOperator = "=^"
	OpIn           BinaryOperator = "in"
	OpNotIn        BinaryOperator = "not in"
	OpContains     BinaryOperator = "contains"
	OpLessThan     BinaryOperator = "<"
	OpLessEqual    BinaryOperator = "<="
	OpGreaterThan  BinaryOperator = ">"
	OpGreaterEqual BinaryOperator = ">="
)

type LiteralType int

const (
	LiteralTypeInt LiteralType = iota
	LiteralTypeString
	// TODO: define subtypes of IP literals(IPv4/IPv6;single IP/IP CIDR).
	LiteralTypeIP
)

type Literal interface {
	Type() LiteralType
	String() string
}

var _ Literal = StringLiteral("")

type StringLiteral string

func (l StringLiteral) Type() LiteralType {
	return LiteralTypeString
}

func (l StringLiteral) String() string {
	str := string(l)
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "\"", "\\\"")
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\t", "\\t")

	return "\"" + str + "\""
}

var _ Literal = IntLiteral(0)

type IntLiteral int

func (l IntLiteral) Type() LiteralType {
	return LiteralTypeInt
}

func (l IntLiteral) String() string {
	return strconv.Itoa(int(l))
}

type Predicate struct {
	field LHS
	op    BinaryOperator
	value Literal
}

func (p Predicate) Matches(req *http.Request) bool {
	// TODO: add logics to the matches
	return true
}

func (p Predicate) Expression() string {
	lhs := p.field.String()
	op := string(p.op)
	rhs := p.value.String()
	return lhs + " " + op + " " + rhs
}

// NewPredicate generates a single predicate.
// TODO: check validity of LHS, op and RHS.
func NewPredicate(lhs LHS, op BinaryOperator, rhs Literal) Predicate {
	return Predicate{
		field: lhs,
		op:    op,
		value: rhs,
	}
}

func NewPredicateNetProtocol(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldNetProtocol{},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPPath(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldHTTPPath{},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPrediacteHTTPHost(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldHTTPHost{},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPMethod(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldHTTPMethod{},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPHeader(key string, op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldHTTPHeader{
			headerName: key,
		},
		op:    op,
		value: StringLiteral(value),
	}
}
