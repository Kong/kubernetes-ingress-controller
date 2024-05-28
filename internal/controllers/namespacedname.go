package controllers

import (
	"github.com/samber/mo"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OptionalNamespacedName is a wrapper around mo.Option[k8stypes.NamespacedName] that provides
// additional Matches and MatchesNN methods for matching against client.Object and
// k8stypes.NamespacedName.
type OptionalNamespacedName struct {
	mo.Option[k8stypes.NamespacedName]
}

// NewOptionalNamespacedName creates a new OptionalNamespacedName with the provided value.
func NewOptionalNamespacedName(onn mo.Option[k8stypes.NamespacedName]) OptionalNamespacedName {
	return OptionalNamespacedName{onn}
}

// Get calls the underlying mo.Option.Get.
func (onn OptionalNamespacedName) Get() (k8stypes.NamespacedName, bool) {
	return onn.Option.Get()
}

// IsPresent calls the underlying mo.Option.IsPresent.
func (onn OptionalNamespacedName) IsPresent() bool {
	return onn.Option.IsPresent()
}

// Matches returns true if the OptionalNamespacedName is present and the provided object's
// namespace and name match the OptionalNamespacedName's namespace and name.
// It also returns true if the OptionalNamespacedName is not present as it is considered
// to match everything.
func (onn OptionalNamespacedName) Matches(obj client.Object) bool {
	n, ok := onn.Option.Get()
	if !ok {
		return true
	}

	return n.Namespace == obj.GetNamespace() && n.Name == obj.GetName()
}

// MatchesNN returns true if the OptionalNamespacedName is present and the provided
// k8stypes.NamespacedName matches the OptionalNamespacedName's namespace and name.
// It also returns true if the OptionalNamespacedName is not present as it is considered
// to match everything.
func (onn OptionalNamespacedName) MatchesNN(nn k8stypes.NamespacedName) bool {
	n, ok := onn.Option.Get()
	if !ok {
		return true
	}

	return n.Namespace == nn.Namespace && n.Name == nn.Name
}
