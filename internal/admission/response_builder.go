package admission

import (
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
)

type ResponseBuilder struct {
	uid      k8stypes.UID
	message  string
	allowed  bool
	warnings []string
}

func NewResponseBuilder(uid k8stypes.UID) *ResponseBuilder {
	return &ResponseBuilder{uid: uid}
}

func (r *ResponseBuilder) Allowed(allowed bool) *ResponseBuilder {
	r.allowed = allowed
	return r
}

func (r *ResponseBuilder) WithMessage(msg string) *ResponseBuilder {
	r.message = msg
	return r
}

func (r *ResponseBuilder) WithWarning(warning string) *ResponseBuilder {
	r.warnings = append(r.warnings, warning)
	return r
}

func (r *ResponseBuilder) Build() *admissionv1.AdmissionResponse {
	var code int32
	if !r.allowed {
		code = 400
	}

	return &admissionv1.AdmissionResponse{
		UID:     r.uid,
		Allowed: r.allowed,
		Result: &metav1.Status{
			Message: r.message,
			Code:    code,
		},
		Warnings: r.warnings,
	}
}
