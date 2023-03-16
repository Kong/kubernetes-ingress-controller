package atc

import (
	"errors"
	"net"
	"net/http"
	"strconv"
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

func (t TransformLower) ExtractValue(req *http.Request) Literal {
	innerVal := t.inner.ExtractValue(req)
	str, ok := innerVal.(StringLiteral)
	if !ok {
		return StringLiteral("")
	}
	return StringLiteral(string(str))
}

type FieldNetProtocol struct{}

func (f FieldNetProtocol) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldNetProtocol) String() string {
	return "net.protocol"
}

func (f FieldNetProtocol) ExtractValue(req *http.Request) Literal {
	if req.TLS != nil {
		return StringLiteral("https")
	}
	return StringLiteral("http")
}

type FieldNetPort struct{}

func (f FieldNetPort) FieldType() FieldType {
	return FieldTypeInt
}

func (f FieldNetPort) String() string {
	return "net.port"
}

func (f FieldNetPort) ExtractValue(req *http.Request) Literal {
	_, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		if errors.Is(err, &net.AddrError{}) && strings.Contains(err.Error(), "missing ports") {
			if req.TLS != nil {
				return IntLiteral(443)
			}
			return IntLiteral(80)
		}
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return IntLiteral(0)
	}
	return IntLiteral(intPort)
}

type FieldTLSSNI struct{}

func (f FieldTLSSNI) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldTLSSNI) String() string {
	return "tls.sni"
}

func (f FieldTLSSNI) ExtractValue(req *http.Request) Literal {
	host, _, _ := net.SplitHostPort(req.Host)
	return StringLiteral(host)
}

type FieldHTTPMethod struct{}

func (f FieldHTTPMethod) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPMethod) String() string {
	return "http.method"
}

func (f FieldHTTPMethod) ExtractValue(req *http.Request) Literal {
	method := req.Method
	if method == "" {
		method = "GET"
	}
	return StringLiteral(method)
}

type FieldHTTPHost struct{}

func (f FieldHTTPHost) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPHost) String() string {
	return "http.host"
}

func (f FieldHTTPHost) ExtractValue(req *http.Request) Literal {
	host, _, _ := net.SplitHostPort(req.Host)
	return StringLiteral(host)
}

type FieldHTTPPath struct{}

func (f FieldHTTPPath) FieldType() FieldType {
	return FieldTypeString
}

func (f FieldHTTPPath) String() string {
	return "http.path"
}

func (f FieldHTTPPath) ExtractValue(req *http.Request) Literal {
	return StringLiteral(req.URL.RawPath)
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

func (f FieldHTTPHeader) ExtractValue(req *http.Request) Literal {
	return (StringLiteral(req.Header.Get(f.headerName)))

}
