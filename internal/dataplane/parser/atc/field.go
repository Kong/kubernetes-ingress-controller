package atc

import (
	"strings"
)

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
