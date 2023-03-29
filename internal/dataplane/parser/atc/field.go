package atc

import (
	"fmt"
	"strings"
)

// This file defines the field types, field transforms, and field constants used in Kong's expression router.
// https://docs.konghq.com/gateway/latest/reference/router-expressions-language/ is the upstream reference that
// describes these fields.

// TransformLower instructs Kong to transform a field (for example, http.path) to lowercase before comparing it to
// a value in a predicate expression. It can only be applied to the left side of a predicate expression.
type TransformLower struct {
	inner LHS
}

func (t TransformLower) FieldType() FieldType {
	return FieldTypeString
}

func (t TransformLower) String() string {
	return fmt.Sprintf("lower(%s)", t.inner.String())
}

func NewTransformerLower(inner LHS) TransformLower {
	return TransformLower{inner: inner}
}

// StringField is defined for fields with constant name and having string type.
// The inner string value is the name of the field.
type StringField string

func (f StringField) FieldType() FieldType {
	return FieldTypeString
}

func (f StringField) String() string {
	return string(f)
}

// https://docs.konghq.com/gateway/latest/reference/router-expressions-language/#available-fields

const (
	FieldNetProtocol StringField = "net.protocol"
	FieldTLSSNI      StringField = "tls.sni"
	FieldHTTPMethod  StringField = "http.method"
	FieldHTTPHost    StringField = "http.host"
	FieldHTTPPath    StringField = "http.path"
)

// IntField is defined for fields with constant name and having integer type.
// The inner string value is the name of the field.
type IntField string

func (f IntField) FieldType() FieldType {
	return FieldTypeInt
}

func (f IntField) String() string {
	return string(f)
}

// https://docs.konghq.com/gateway/latest/reference/router-expressions-language/#available-fields

const (
	FieldNetPort IntField = "net.port"
)

// HTTPHeaderField extracts the value of an HTTP header from the request.
type HTTPHeaderField struct {
	HeaderName string
}

func (f HTTPHeaderField) FieldType() FieldType {
	return FieldTypeString
}

func (f HTTPHeaderField) String() string {
	return "http.headers." + strings.ToLower(strings.ReplaceAll(f.HeaderName, "-", "_"))
}
