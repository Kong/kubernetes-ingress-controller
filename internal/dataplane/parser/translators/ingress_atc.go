package translators

import (
	"strings"

	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
)

func MatcherFromIngressPath(ingressPath netv1.HTTPIngressPath, regexPathPrefix string) atc.Matcher {
	switch *ingressPath.PathType {
	case netv1.PathTypeExact:
		relative := strings.TrimLeft(ingressPath.Path, "/")
		return atc.NewPredicateHTTPPath(atc.OpEqual, "/"+relative)
	case netv1.PathTypePrefix:
		base := "/" + strings.Trim(ingressPath.Path, "/")
		return atc.Or(
			atc.NewPredicateHTTPPath(atc.OpEqual, base),
			atc.NewPredicateHTTPPath(atc.OpPrefixMatch, base+"/"),
		)
	case netv1.PathTypeImplementationSpecific:
		path := ingressPath.Path
		if strings.HasPrefix(path, regexPathPrefix) {
			regexPath := strings.TrimPrefix(path, regexPathPrefix)
			if !strings.HasPrefix(regexPath, "^") {
				regexPath = "^" + regexPath
			}
			return atc.NewPredicateHTTPPath(atc.OpRegexMatch, regexPath)
		}
		return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, path)
	}

	return nil
}

func MatcherFromIngressHost()
