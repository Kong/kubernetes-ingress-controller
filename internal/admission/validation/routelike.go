package validation

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func ValidateRouteSourceAnnotations(obj client.Object) error {
	protocols := annotations.ExtractProtocolNames(obj.GetAnnotations())
	for _, protocol := range protocols {
		if !util.ValidateProtocol(protocol) {
			return fmt.Errorf("invalid %s value: %s", annotations.AnnotationPrefix+annotations.ProtocolsKey, protocol)
		}
	}
	return nil
}
