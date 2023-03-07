package v1

// KongProtocol is a valid Kong protocol.
// This alias is necessary to deal with https://github.com/kubernetes-sigs/controller-tools/issues/342
// +kubebuilder:validation:Enum=http;https;grpc;grpcs;tcp;tls;udp
// +kubebuilder:object:generate=true
type KongProtocol string

// KongProtocolsToStrings converts a slice of KongProtocol to plain strings
func KongProtocolsToStrings(protocols []KongProtocol) (res []string) {
	for _, protocol := range protocols {
		res = append(res, string(protocol))
	}
	return
}

// StringsToKongProtocols converts a slice of strings to KongProtocols
func StringsToKongProtocols(strings []string) (res []KongProtocol) {
	for _, protocol := range strings {
		res = append(res, KongProtocol(protocol))
	}
	return
}

// ProtocolSlice converts a slice of string to a slice of *KongProtocol
func ProtocolSlice(elements ...string) []*KongProtocol {
	var res []*KongProtocol
	for _, element := range elements {
		e := KongProtocol(element)
		res = append(res, &e)
	}
	return res
}
