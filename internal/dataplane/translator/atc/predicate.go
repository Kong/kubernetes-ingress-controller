package atc

import (
	"errors"
	"fmt"
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

// LHS is the left hand side (the field) of a predicate expression.
type LHS interface {
	// FieldType returns the FieldType iota indicating the LHS type.
	FieldType() FieldType

	// String returns a string representation of the LHS.
	String() string
}

// BinaryOperator is an operator that accepts two arguments within a predicate expression.
type BinaryOperator string

const (
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
	// LiteralTypeIP is a type that represents a literal of a single IP address or an IP CIDR.
	//
	// TODO: define subtypes of IP literals(IPv4/IPv6;single IP/IP CIDR).
	LiteralTypeIP
)

// Literal is the right hand side (the value) of a predicate expression.
type Literal interface {
	// Type returns the LiteralType iota indicating the Literal type.
	Type() LiteralType

	// String returns a string representation of the Literal.
	String() string
}

var _ Literal = StringLiteral("")

// StringLiteral is a string Literal.
type StringLiteral string

func (l StringLiteral) Type() LiteralType {
	return LiteralTypeString
}

func (l StringLiteral) String() string {
	str := string(l)
	// replace the escape characters: '\', '\n', '\t', '\r', '\"'
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "\"", "\\\"")
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\t", "\\t")

	return fmt.Sprintf("\"%s\"", str)
}

var _ Literal = IntLiteral(0)

// IntLiteral is an integer Literal.
type IntLiteral int

func (l IntLiteral) Type() LiteralType {
	return LiteralTypeInt
}

func (l IntLiteral) String() string {
	return strconv.Itoa(int(l))
}

// Predicate is an expression consisting of two arguments and a comparison operator. Kong's expression router evaluates
// these to true or false.
type Predicate struct {
	field LHS
	op    BinaryOperator
	value Literal
}

// Expression returns a string representation of a Predicate.
func (p Predicate) Expression() string {
	lhs := p.field.String()
	op := string(p.op)
	rhs := p.value.String()
	return fmt.Sprintf("%s %s %s", lhs, op, rhs)
}

// IsEmpty returns true if a Predicate has no value to compare against.
func (p Predicate) IsEmpty() bool {
	return p.value == nil
}

var (
	ErrTypeNotMatch    = errors.New("type does not match on sides of predicate")
	ErrOperatorInvalid = errors.New("operator is not valid for the types of sides of predicate")
)

// NewPredicate generates a single predicate.
func NewPredicate(lhs LHS, op BinaryOperator, rhs Literal) (Predicate, error) {
	// Check for predicates on string fields.
	if lhs.FieldType() == FieldTypeString {
		if rhs.Type() != LiteralTypeString {
			return Predicate{}, ErrTypeNotMatch
		}
		if op == OpGreaterThan || op == OpGreaterEqual || op == OpLessThan || op == OpLessEqual {
			return Predicate{}, ErrOperatorInvalid
		}
	}
	// Check for predicates on integer fields.
	if lhs.FieldType() == FieldTypeInt {
		if rhs.Type() != LiteralTypeInt {
			return Predicate{}, ErrTypeNotMatch
		}
		if op == OpContains || op == OpPrefixMatch || op == OpSuffixMatch || op == OpRegexMatch {
			return Predicate{}, ErrOperatorInvalid
		}
	}
	return Predicate{
		field: lhs,
		op:    op,
		value: rhs,
	}, nil
}

func NewPredicateNetProtocol(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldNetProtocol,
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPPath(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldHTTPPath,
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPrediacteHTTPHost(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldHTTPHost,
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPMethod(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldHTTPMethod,
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPHeader(key string, op BinaryOperator, value string) Predicate {
	return Predicate{
		field: HTTPHeaderField{
			HeaderName: key,
		},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateTLSSNI(op BinaryOperator, value string) Predicate {
	return Predicate{
		field: FieldTLSSNI,
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPQuery(key string, op BinaryOperator, value string) Predicate {
	return Predicate{
		field: HTTPQueryField{
			QueryParamName: key,
		},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPPathSingleSegment(index int, op BinaryOperator, value string) Predicate {
	return Predicate{
		field: HTTPPathSingleSegmentField{
			Index: index,
		},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPPathSegmentInterval(start, end int, op BinaryOperator, value string) Predicate {
	return Predicate{
		field: HTTPPathSegmentIntervalField{
			Start: start,
			End:   end,
		},
		op:    op,
		value: StringLiteral(value),
	}
}

func NewPredicateHTTPPathSegmentLength(op BinaryOperator, value int) Predicate {
	return Predicate{
		field: FieldHTTPPathSegmentsLen,
		op:    op,
		value: IntLiteral(value),
	}
}
