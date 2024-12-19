package builder

import (
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// ParentReferenceBuilder is a builder for constructing gatewayapi.ParentReference objects.
// Primarily used for testing.
type ParentReferenceBuilder struct {
	pr gatewayapi.ParentReference
}

// NewParentReference creates a new instance of ParentReferenceBuilder.
func NewParentReference() *ParentReferenceBuilder {
	return &ParentReferenceBuilder{
		pr: gatewayapi.ParentReference{},
	}
}

// Group sets the Group field of the ParentReference.
func (b *ParentReferenceBuilder) Group(group gatewayapi.Group) *ParentReferenceBuilder {
	b.pr.Group = &group
	return b
}

// Kind sets the Kind field of the ParentReference.
func (b *ParentReferenceBuilder) Kind(kind gatewayapi.Kind) *ParentReferenceBuilder {
	b.pr.Kind = &kind
	return b
}

// Namespace sets the Namespace field of the ParentReference.
func (b *ParentReferenceBuilder) Namespace(namespace gatewayapi.Namespace) *ParentReferenceBuilder {
	b.pr.Namespace = &namespace
	return b
}

// Name sets the Name field of the ParentReference.
func (b *ParentReferenceBuilder) Name(name gatewayapi.ObjectName) *ParentReferenceBuilder {
	b.pr.Name = name
	return b
}

// SectionName sets the SectionName field of the ParentReference.
func (b *ParentReferenceBuilder) SectionName(sectionName gatewayapi.SectionName) *ParentReferenceBuilder {
	b.pr.SectionName = &sectionName
	return b
}

// Build returns the configured ParentReference.
func (b *ParentReferenceBuilder) Build() gatewayapi.ParentReference {
	return b.pr
}
