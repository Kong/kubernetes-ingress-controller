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
}

type TransformLower struct {
	inner LHS
}

func (t TransformLower) FieldType() FieldType {
	return FieldTypeString
}

func (t TransformLower) String() string {
	return "lower(" + t.inner.String() + ")"
}

type FieldNetProtocol struct{}

func (f FieldNetProtocol) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldNetProtocol) String() string {
	return "net.protocol"
}

type FieldNetPort struct{}

func (f FieldNetPort) FieldType() FieldType {
	return FieldTypeInt
}

func (f FieldNetPort) String() string {
	return "net.port"
}

type FieldTLSSNI struct{}

func (f FieldTLSSNI) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldTLSSNI) String() string {
	return "tls.sni"
}

type FieldHTTPMethod struct{}

func (f FieldHTTPMethod) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPMethod) String() string {
	return "http.method"
}

type FieldHTTPHost struct{}

func (f FieldHTTPHost) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPHost) String() string {
	return "http.host"
}

type FieldHTTPPath struct{}

func (f FieldHTTPPath) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPPath) String() string {
	return "http.path"
}

type FieldHTTPHeader struct {
	headerName string
}

func (f FieldHTTPHeader) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPHeader) String() string {
	return "http.header." + strings.ToLower(strings.ReplaceAll(f.headerName, "-", "_"))
}

type BinaryOperator string

var (
	OpEqual        BinaryOperator = "=="
	OpNotEqual     BinaryOperator = "!="
	OpRegexMatch   BinaryOperator = "~"
	OpPrefixMatch  BinaryOperator = "^="
	OpSuffixMatch  BinaryOperator = "=^"
	OpContains     BinaryOperator = "in"
	OpNotContains  BinaryOperator = "not in"
	OpLessThan     BinaryOperator = "<"
	OpLessEqual    BinaryOperator = "<="
	OpGreaterThan  BinaryOperator = ">"
	OpGreaterEqual BinaryOperator = ">="
)

type LiteralType int

const (
	LiteralTypeInt LiteralType = iota
	LiteralTypeString
	// TODO: define subtypes of IP literals(IPv4/IPv6;single IP/IP CIDR)
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

// TODO: define more concrete function to generate predicates with specified fields
// like NewPredicateHTTPPath(path string, op BinaryOperator, value string) Predicate
