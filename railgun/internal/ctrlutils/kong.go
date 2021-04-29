package ctrlutils

import "github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"

// ProxyUpdateParams defines all the attrs needed to perform a full configuration update on the Kong Admin API.
type ProxyUpdateParams struct {
	IngressClassName               string
	KongConfig                     sendconfig.Kong
	ProcessClasslessIngressV1Beta1 bool
	ProcessClasslessIngressV1      bool
	ProcessClasslessKongConsumer   bool
}
