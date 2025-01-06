package v1

// KongProtocol is a valid Kong protocol.
// This alias is necessary to deal with https://github.com/kubernetes-sigs/controller-tools/issues/342
// +kubebuilder:validation:Enum=http;https;grpc;grpcs;tcp;tls;udp
// +kubebuilder:object:generate=true
type KongProtocol string

// KongProtocolsToStrings converts a slice of KongProtocol to plain strings.
func KongProtocolsToStrings(protocols []KongProtocol) []string {
	res := make([]string, 0, len(protocols))
	for _, protocol := range protocols {
		res = append(res, string(protocol))
	}
	return res
}

// StringsToKongProtocols converts a slice of strings to KongProtocols.
func StringsToKongProtocols(strings []string) []KongProtocol {
	res := make([]KongProtocol, 0, len(strings))
	for _, protocol := range strings {
		res = append(res, KongProtocol(protocol))
	}
	return res
}

// ProtocolSlice converts a slice of string to a slice of *KongProtocol.
func ProtocolSlice(elements ...string) []*KongProtocol {
	res := make([]*KongProtocol, 0, len(elements))
	for _, element := range elements {
		e := KongProtocol(element)
		res = append(res, &e)
	}
	return res
}
