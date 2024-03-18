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
	FieldNetDstPort          IntField = "net.dst.port"
	FieldHTTPPathSegmentsLen IntField = "http.path.segments.len"
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

// HTTPQueryField extracts the value of an HTTP query parameter from the query string of the request.
type HTTPQueryField struct {
	QueryParamName string
}

func (f HTTPQueryField) FieldType() FieldType {
	return FieldTypeString
}

func (f HTTPQueryField) String() string {
	return "http.queries." + f.QueryParamName
}

// HTTPPathSingleSegmentField represensts a single segment of HTTP path with 0 based index.
type HTTPPathSingleSegmentField struct {
	Index int
}

func (f HTTPPathSingleSegmentField) FieldType() FieldType {
	return FieldTypeString
}

func (f HTTPPathSingleSegmentField) String() string {
	return fmt.Sprintf("http.path.segments.%d", f.Index)
}

// HTTPPathSegmentIntervalField represents a closed interval of segments in HTTP path.
type HTTPPathSegmentIntervalField struct {
	Start int
	End   int
}

func (f HTTPPathSegmentIntervalField) FieldType() FieldType {
	return FieldTypeString
}

func (f HTTPPathSegmentIntervalField) String() string {
	return fmt.Sprintf("http.path.segments.%d_%d", f.Start, f.End)
}
