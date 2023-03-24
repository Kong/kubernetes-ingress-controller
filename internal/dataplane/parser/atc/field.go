package atc

import (
	"fmt"
	"strings"
)

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

const (
	FieldNetPort IntField = "net.port"
)

type FieldHTTPHeader struct {
	HeaderName string
}

func (f FieldHTTPHeader) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPHeader) String() string {
	return "http.header." + strings.ToLower(strings.ReplaceAll(f.HeaderName, "-", "_"))
}
