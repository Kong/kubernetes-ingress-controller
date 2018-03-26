package apis

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var GroupName = "configuration.konghq.com"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
