package parser

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser/kongstate"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

type TLSPair struct {
	Key, Cert string
}

var (
	tlsPairs = []TLSPair{
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIC2DCCAcACCQC32eFOsWpKojANBgkqhkiG9w0BAQsFADAuMRcwFQYDVQQDDA5z
ZWN1cmUtZm9vLWJhcjETMBEGA1UECgwKa29uZ2hxLm9yZzAeFw0xODEyMTgyMTI4
MDBaFw0xOTEyMTgyMTI4MDBaMC4xFzAVBgNVBAMMDnNlY3VyZS1mb28tYmFyMRMw
EQYDVQQKDAprb25naHEub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAqhl/HSwV6PbMv+cMFU9X+HuM7QbNNPh39GKa4pkxzFgiAnuuJ4jw9V/bzsEy
S+ZIyjzo+QKB1LzmgdcX4vkdI22BjxUd9HPHdZxtv3XilbNmSk9UOl2Hh1fORJoS
7YH+VbvVwiz5lo7qKRepbg/jcKkbs6AUE0YWFygtDLTvhP2qkphQkxZ0m8qroW91
CWgI73Ar6U2W/YQBRI3+LwtsKo0p2ASDijvqxElQBgBIiyGIr0RZc5pkCJ1eQdDB
2F6XaMfpeEyBj0MxypNL4S9HHfchOt55J1KOzYnUPkQnSoxp6oEjef4Q/ZCj5BRL
EGZnTb3tbwzHZCxGtgl9KqO9pQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAKQ5BX
kkBL+alERL31hsOgWgRiUMw+sPDtRS96ozUlPtVvAg9XFdpY4ldtWkxFcmBnhKzp
UewjrHkf9rR16NISwUTjlGIwaJu/ACQrY15v+r301Crq2DV+GjiUJFVuT495dp/l
0LZbt2Sh/uD+r3UNTcJpJ7jb1V0UP7FWXFj8oafsoFSgmxAPjpKQySTC54JK4AYb
QSnWu1nQLyohnrB9qLZhe2+jOQZnkKuCcWJQ5njvU6SxT3SOKE5XaOZCezEQ6IVL
U47YCCXsq+7wKWXBhKl4H2Ztk6x3HOC56l0noXWezsMfrou/kjwGuuViGnrjqelS
WQ7uVeNCUBY+l+qY
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCqGX8dLBXo9sy/
5wwVT1f4e4ztBs00+Hf0YprimTHMWCICe64niPD1X9vOwTJL5kjKPOj5AoHUvOaB
1xfi+R0jbYGPFR30c8d1nG2/deKVs2ZKT1Q6XYeHV85EmhLtgf5Vu9XCLPmWjuop
F6luD+NwqRuzoBQTRhYXKC0MtO+E/aqSmFCTFnSbyquhb3UJaAjvcCvpTZb9hAFE
jf4vC2wqjSnYBIOKO+rESVAGAEiLIYivRFlzmmQInV5B0MHYXpdox+l4TIGPQzHK
k0vhL0cd9yE63nknUo7NidQ+RCdKjGnqgSN5/hD9kKPkFEsQZmdNve1vDMdkLEa2
CX0qo72lAgMBAAECggEADxMTYNJ3Xp4Ap0EioQDXGv5YDul7ZiZe+xmCAHLzJtjo
qq+rT3WjZRuJr1kPzAosiT+8pdTDDMdw5jDZvRO2sV0TDksgzHk2RAYI897OpdWw
SwWcwU9oo2X0sb+1zbang5GR8BNsSxt/RQUDzu05itJx0gltvgeIDaVR2L5wO6ja
USa8OVuj/92XtIIve9OtyK9jAzgR6LQOTFrCCEv89/vmy5Bykv4Uz8s8swZmTs3v
XJmAmruHGuSLMfXk8lBRp/gVyNTi3uMsdph5AJbVKnra5TZLguEozZKbLdNUYk0p
+aAc7rxDcH2sPqa/7DwRvei9dvd5oB3VJlxGVgC8AQKBgQDfznRSSKAD15hoSDzt
cKNyhLgWAL+MD0jhHKUy3x+Z9OCvf0DVnmru5HfQKq5UfT0t8VTRPGKmOtAMD4cf
LYjIurvMvpVzQGSJfhtHQuULZTh3dfsM7xivMqSV+9txklMAakM7vGQlOQxhrScM
21Mp5LWDU6+e2pFCrQPop0IPkQKBgQDCkVE+dou2yFuJx3uytCH1yKPSy9tkdhQH
dGF12B5dq8MZZozAz5P9YN/COa9WjsNKDqWbEgLEksEQUq4t8SBjHnSV/D3x7rEF
qgwii0GETYxax6gms8nueIqWZQf+0NbX7Gc5mTqeVb7v3TrhsKr0VNMFRXXQwP2E
M/pxJq8q1QKBgQC3rH7oXLP+Ez0AMHDYSL3LKULOw/RvpMeh/9lQA6+ysTaIsP3r
kuSdhCEUVULXEiVYhBug0FcBp3jAvSmem8cLPb0Mjkim2mzoLfeDJ1JEZODPoaLU
fZEbj4tlj9oLvhOiXpMo/jaOGeCgdPN8aK86zXlt+wtBao0WVFnF4SalEQKBgQC1
uLfi2SGgs/0a8B/ORoO5ZY3s4c2lRMtsMvyb7iBeaIAuByPLKZUVABe89deXxnsL
fiaacPX41wBO2IoqCp2vNdC6DP9mKQNZQPtYgCvPAAbo+rVIgH9HpXn7AZ24FyGy
RfAbUcv3+in9KelGxZTF4zu8HqXtNXMSuOFeMT1FiQKBgF0R+IFDGHhD4nudAQvo
hncXsgyzK6QUzak6HmFji/CMZ6EU9q6A67JkiEWrYoKqIAKZ2Og8+Eucr/rDdGWc
kqlmLPBJAJeUsP/9KidBjTE5mIbn/2n089VPMBvnlt2xIcuB6+zrf2NjvlcZEyKS
Gn+T2uCyOP4a1DTUoPyoNJXo
-----END PRIVATE KEY-----`,
		},
	}

	caCert1 = `-----BEGIN CERTIFICATE-----
MIIEvjCCAqagAwIBAgIJALabx/Nup200MA0GCSqGSIb3DQEBCwUAMBMxETAPBgNV
BAMMCFlvbG80Mi4xMCAXDTE5MDkxNTE2Mjc1M1oYDzIxMTkwODIyMTYyNzUzWjAT
MREwDwYDVQQDDAhZb2xvNDIuMTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC
ggIBANIW67Ay0AtTeBY2mORaGet/VPL5jnBRz0zkZ4Jt7fEq3lbxYaJBnFI8wtz3
bHLtLsxkvOFujEMY7HVd+iTqbJ7hLBtK0AdgXDjf+HMmoWM7x0PkZO+3XSqyRBbI
YNoEaQvYBNIXrKKJbXIU6higQaXYszeN8r3+RIbcTIlZxy28msivEGfGTrNujQFc
r/eyf+TLHbRqh0yg4Dy/U/T6fqamGhFrjupRmOMugwF/BHMH2JHhBYkkzuZLgV2u
7Yh1S5FRlh11am5vWuRSbarnx72hkJ99rUb6szOWnJKKew8RSn3CyhXbS5cb0QRc
ugRc33p/fMucJ4mtCJ2Om1QQe83G1iV2IBn6XJuCvYlyWH8XU0gkRxWD7ZQsl0bB
8AFTkVsdzb94OM8Y6tWI5ybS8rwl8b3r3fjyToIWrwK4WDJQuIUx4nUHObDyw+KK
+MmqwpAXQWbNeuAc27FjuJm90yr/163aGuInNY5Wiz6CM8WhFNAi/nkEY2vcxKKx
irSdSTkbnrmLFAYrThaq0BWTbW2mwkOatzv4R2kZzBUOiSjRLPnbyiPhI8dHLeGs
wMxiTXwyPi8iQvaIGyN4DPaSEiZ1GbexyYFdP7sJJD8tG8iccbtJYquq3cDaPTf+
qv5M6R/JuMqtUDheLSpBNK+8vIe5e3MtGFyrKqFXdynJtfHVAgMBAAGjEzARMA8G
A1UdEwQIMAYBAf8CAQAwDQYJKoZIhvcNAQELBQADggIBAK0BmL5B1fPSMbFy8Hbc
/ESEunt4HGaRWmZZSa/aOtTjhKyDXLLJZz3C4McugfOf9BvvmAOZU4uYjfHTnNH2
Z3neBkdTpQuJDvrBPNoCtJns01X/nuqFaTK/Tt9ZjAcVeQmp51RwhyiD7nqOJ/7E
Hp2rC6gH2ABXeexws4BDoZPoJktS8fzGWdFBCHzf4mCJcb4XkI+7GTYpglR818L3
dMNJwXeuUsmxxKScBVH6rgbgcEC/6YwepLMTHB9VcH3X5VCfkDIyPYLWmvE0gKV7
6OU91E2Rs8PzbJ3EuyQpJLxFUQp8ohv5zaNBlnMb76UJOPR6hXfst5V+e7l5Dgwv
Dh4CeO46exmkEsB+6R3pQR8uOFtubH2snA0S3JA1ji6baP5Y9Wh9bJ5McQUgbAPE
sCRBFoDLXOj3EgzibohC5WrxN3KIMxlQnxPl3VdQvp4gF899mn0Z9V5dAsGPbxRd
quE+DwfXkm0Sa6Ylwqrzu2OvSVgbMliF3UnWbNsDD5KcHGIaFxVC1qkwK4cT3pyS
58i/HAB2+P+O+MltQUDiuw0OSUFDC0IIjkDfxLVffbF+27ef9C5NG81QlwTz7TuN
zeigcsBKooMJTszxCl6dtxSyWTj7hJWXhy9pXsm1C1QulG6uT4RwCa3m0QZoO7G+
6Wu6lP/kodPuoNubstIuPdi2
-----END CERTIFICATE-----`
	caCert2 = `-----BEGIN CERTIFICATE-----
MIIEvjCCAqagAwIBAgIJAPf5iqimiR2BMA0GCSqGSIb3DQEBCwUAMBMxETAPBgNV
BAMMCFlvbG80Mi4yMCAXDTE5MDkxNTE2Mjc1OVoYDzIxMTkwODIyMTYyNzU5WjAT
MREwDwYDVQQDDAhZb2xvNDIuMjCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC
ggIBANIW67Ay0AtTeBY2mORaGet/VPL5jnBRz0zkZ4Jt7fEq3lbxYaJBnFI8wtz3
bHLtLsxkvOFujEMY7HVd+iTqbJ7hLBtK0AdgXDjf+HMmoWM7x0PkZO+3XSqyRBbI
YNoEaQvYBNIXrKKJbXIU6higQaXYszeN8r3+RIbcTIlZxy28msivEGfGTrNujQFc
r/eyf+TLHbRqh0yg4Dy/U/T6fqamGhFrjupRmOMugwF/BHMH2JHhBYkkzuZLgV2u
7Yh1S5FRlh11am5vWuRSbarnx72hkJ99rUb6szOWnJKKew8RSn3CyhXbS5cb0QRc
ugRc33p/fMucJ4mtCJ2Om1QQe83G1iV2IBn6XJuCvYlyWH8XU0gkRxWD7ZQsl0bB
8AFTkVsdzb94OM8Y6tWI5ybS8rwl8b3r3fjyToIWrwK4WDJQuIUx4nUHObDyw+KK
+MmqwpAXQWbNeuAc27FjuJm90yr/163aGuInNY5Wiz6CM8WhFNAi/nkEY2vcxKKx
irSdSTkbnrmLFAYrThaq0BWTbW2mwkOatzv4R2kZzBUOiSjRLPnbyiPhI8dHLeGs
wMxiTXwyPi8iQvaIGyN4DPaSEiZ1GbexyYFdP7sJJD8tG8iccbtJYquq3cDaPTf+
qv5M6R/JuMqtUDheLSpBNK+8vIe5e3MtGFyrKqFXdynJtfHVAgMBAAGjEzARMA8G
A1UdEwQIMAYBAf8CAQAwDQYJKoZIhvcNAQELBQADggIBALNx2xaS5nv1QjEqtiCO
EA/ZTXbs+il6cf6ZyUwFXs7d3OKx6Kk2Nr7wGgM1M5WuTyIGKtZspz9ThzYmsuN/
UBCSKLw3X7U2fLiHJDipXboU1txasTErUTPJs/Vq4v7PWh8sMLCQH/ha4FAOXR0M
Uie+VgSJNKoQSj7G1hzU/LZv0KdvJ45mQBCnBXrUrGgeEcRqubbkDKgdBh7dJQzW
Xgy6rPb6H1aXbsSuRuUVv/xFHJoCdZJmqPH4JTMYRbHNS2km9nHVJzmtL6pQFe32
24wfpue9geFndOE9bDU9/cqoRYA4Pce4V5qDL0wL9W4uPmyPDkulKNQtAvZnDA9V
6ccYYthlTBr62UEnw7zZOnSm0q4fB2o82/6bdPwrT7WhbHZQWN7SeqYNWAbYZ1EE
40f5IpTwZ7E5LaG62qPhKLXame7SPAaqaQ9aCTYxaWR7XSYBsvCBRanjRq0r9Tql
T1I8lwssIgbA3XubokI+IMkLDEpCQ27niWXOZL5y2M3xyutd6PPjmEEmoHMkOrZL
etlxzx2CCoUDXKkYW2gZKEozwBZ+eBgUj8WB5g/8jGDAI0qzYnfAgiahjGwlEUtP
hJiPG/YFADw0m5b/8OMCZ6AXNhxjdweHniDxY2HE734Nwm9mG/7UbkdvhR05tqFh
G4KCViLH0cXt/TgW1sYB2o9Z
-----END CERTIFICATE-----`
)

func TestGlobalPlugin(t *testing.T) {
	assert := assert.New(t)
	t.Run("global plugins are processed correctly", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KongClusterPlugins: []*configurationv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					Config: configurationv1.Configuration{
						"foo1": "bar1",
					},
				},
			},
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Plugins),
			"expected one plugin to be rendered")

		sort.SliceStable(state.Plugins, func(i, j int) bool {
			return strings.Compare(*state.Plugins[i].Name, *state.Plugins[j].Name) > 0
		})

		assert.Equal("basic-auth", *state.Plugins[0].Name)
		assert.Equal(kong.Configuration{"foo1": "bar1"}, state.Plugins[0].Config)
	})
}

func TestSecretConfigurationPlugin(t *testing.T) {
	jwtPluginConfig := `{"run_on_preflight": false}`  // JSON
	basicAuthPluginConfig := "hide_credentials: true" // YAML
	assert := assert.New(t)
	stock := store.FakeObjects{
		Services: []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "foo-svc",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
		},
		IngressesV1beta1: []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.DeprecatedPluginsKey: "foo-plugin",
						annotations.IngressClassKey:      annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.DeprecatedPluginsKey: "bar-plugin",
						annotations.IngressClassKey:      annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.net",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}}
	t.Run("plugins with secret configuration are processed correctly",
		func(t *testing.T) {
			objects := stock
			objects.KongPlugins = []*configurationv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "jwt-config",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*configurationv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-broken-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
						Annotations: map[string]string{
							// explicitly none, this should not get rendered
						},
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)
			assert.Equal(3, len(state.Plugins),
				"expected three plugins to be rendered")

			sort.SliceStable(state.Plugins, func(i, j int) bool {
				return strings.Compare(*state.Plugins[i].Name,
					*state.Plugins[j].Name) > 0
			})
			assert.Equal("jwt", *state.Plugins[0].Name)
			assert.Equal(kong.Configuration{"run_on_preflight": false},
				state.Plugins[0].Config)

			assert.Equal("basic-auth", *state.Plugins[1].Name)
			assert.Equal(kong.Configuration{"hide_credentials": true},
				state.Plugins[2].Config)
			assert.Equal("basic-auth", *state.Plugins[2].Name)
			assert.Equal(kong.Configuration{"hide_credentials": true},
				state.Plugins[2].Config)
		})

	t.Run("plugins with missing secrets or keys are not constructed",
		func(t *testing.T) {
			objects := stock
			objects.KongPlugins = []*configurationv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-foo-plugin",
						Namespace: "default",
						Labels: map[string]string{
							"global": "true",
						},
					},
					PluginName: "jwt",
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*configurationv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)
			assert.Equal(0, len(state.Plugins),
				"expected no plugins to be rendered")
		})

	t.Run("plugins with both config and configFrom are not constructed",
		func(t *testing.T) {
			objects := stock
			objects.KongPlugins = []*configurationv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-foo-plugin",
						Namespace: "default",
						Labels: map[string]string{
							"global": "true",
						},
					},
					PluginName: "jwt",
					Config:     configurationv1.Configuration{"fake": true},
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "jwt-config",
							Secret: "conf-secret",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					Config:     configurationv1.Configuration{"fake": true},
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "jwt-config",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*configurationv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					Config:     configurationv1.Configuration{"fake": true},
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					Config:     configurationv1.Configuration{"fake": true},
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)
			assert.Equal(0, len(state.Plugins),
				"expected no plugins to be rendered")
		})

	t.Run("secretToConfiguration handles valid configuration and "+
		"discards invalid configuration", func(t *testing.T) {
		objects := stock
		jwtPluginConfig := `{"run_on_preflight": false}`  // JSON
		basicAuthPluginConfig := "hide_credentials: true" // YAML
		badJwtPluginConfig := "22222"                     // not JSON
		badBasicAuthPluginConfig := "111111"              // not YAML
		objects.Secrets = []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "conf-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"jwt-config":            []byte(jwtPluginConfig),
					"basic-auth-config":     []byte(basicAuthPluginConfig),
					"bad-jwt-config":        []byte(badJwtPluginConfig),
					"bad-basic-auth-config": []byte(badBasicAuthPluginConfig),
				},
			},
		}
		references := []*configurationv1.SecretValueFromSource{
			{
				Secret: "conf-secret",
				Key:    "jwt-config",
			},
			{
				Secret: "conf-secret",
				Key:    "basic-auth-config",
			},
		}
		badReferences := []*configurationv1.SecretValueFromSource{
			{
				Secret: "conf-secret",
				Key:    "bad-basic-auth-config",
			},
			{
				Secret: "conf-secret",
				Key:    "bad-jwt-config",
			},
		}
		store, err := store.NewFakeStore(objects)
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		for _, testcase := range references {
			config, err := kongstate.SecretToConfiguration(store, *testcase, "default")
			assert.NotEmpty(config)
			assert.Nil(err)
		}
		for _, testcase := range badReferences {
			config, err := kongstate.SecretToConfiguration(store, *testcase, "default")
			assert.Empty(config)
			assert.NotEmpty(err)
		}
	})
	t.Run("plugins with unparsable configuration are not constructed",
		func(t *testing.T) {
			jwtPluginConfig := "22222"        // not JSON
			basicAuthPluginConfig := "111111" // not YAML
			objects := stock
			objects.KongPlugins = []*configurationv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-foo-plugin",
						Namespace: "default",
						Labels: map[string]string{
							"global": "true",
						},
					},
					PluginName: "jwt",
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
					},
					PluginName: "jwt",
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "missing-key",
							Secret: "conf-secret",
						},
					},
				},
			}
			objects.KongClusterPlugins = []*configurationv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "global-bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					ConfigFrom: configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "basic-auth-config",
							Secret:    "missing-secret",
							Namespace: "default",
						},
					},
				},
			}
			objects.Secrets = []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
						Name:      "conf-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"jwt-config":        []byte(jwtPluginConfig),
						"basic-auth-config": []byte(basicAuthPluginConfig),
					},
				},
			}
			store, err := store.NewFakeStore(objects)
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)
			assert.Equal(0, len(state.Plugins),
				"expected no plugins to be rendered")
		})
}

func TestCACertificate(t *testing.T) {
	assert := assert.New(t)
	t.Run("valid CACertificte is processed", func(t *testing.T) {
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eae"),
					"cert": []byte(caCert1),
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Secrets: secrets,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.CACertificates))
		assert.Equal(kong.CACertificate{
			ID:   kong.String("8214a145-a328-4c56-ab72-2973a56d4eae"),
			Cert: kong.String(caCert1),
		}, state.CACertificates[0])
	})
	t.Run("multiple CACertifictes are processed", func(t *testing.T) {
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eae"),
					"cert": []byte(caCert1),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("570c28aa-e784-43c1-8ec7-ae7f4ce40189"),
					"cert": []byte(caCert2),
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Secrets: secrets,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(2, len(state.CACertificates))
	})
	t.Run("invalid CACertifictes are ignored", func(t *testing.T) {
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eae"),
					"cert": []byte(caCert1),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					"id": []byte("570c28aa-e784-43c1-8ec7-ae7f4ce40189"),
					// cert is missing
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "non-default",
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Data: map[string][]byte{
					// id is missing
					"cert": []byte(caCert2),
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Secrets: secrets,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.CACertificates))
		assert.Equal(kong.CACertificate{
			ID:   kong.String("8214a145-a328-4c56-ab72-2973a56d4eae"),
			Cert: kong.String(caCert1),
		}, state.CACertificates[0])
	})
}

func TestServiceClientCertificate(t *testing.T) {
	assert := assert.New(t)
	t.Run("valid client-cert annotation", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Secrets:          secrets,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Certificates),
			"expected one certificates to be rendered")
		assert.Equal("7428fb98-180b-4702-a91f-61351a33c6e4",
			*state.Certificates[0].ID)

		assert.Equal(1, len(state.Services))
		assert.Equal("7428fb98-180b-4702-a91f-61351a33c6e4",
			*state.Services[0].ClientCertificate.ID)
	})
	t.Run("client-cert secret doesn't exist", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered")

		assert.Equal(1, len(state.Services))
		assert.Nil(state.Services[0].ClientCertificate)
	})
}

func TestKongRouteAnnotations(t *testing.T) {
	assert := assert.New(t)
	t.Run("strip-path annotation is correctly processed (true)", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/strip-path": "trUe",
						annotations.IngressClassKey:           "kong",
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
		assert.Equal(kong.Service{
			Name:           kong.String("default.foo-svc.80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
		}, state.Services[0].Service)

		assert.Equal(1, len(state.Services[0].Routes),
			"expected one route to be rendered")
		assert.Equal(kong.Route{
			Name:          kong.String("default.bar.00"),
			StripPath:     kong.Bool(true),
			Hosts:         kong.StringSlice("example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
			RegexPriority: kong.Int(0),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("strip-path annotation is correctly processed (false)", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:           "kong",
						"configuration.konghq.com/strip-path": "false",
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
		assert.Equal(kong.Service{
			Name:           kong.String("default.foo-svc.80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
		}, state.Services[0].Service)

		assert.Equal(1, len(state.Services[0].Routes),
			"expected one route to be rendered")
		assert.Equal(kong.Route{
			Name:          kong.String("default.bar.00"),
			StripPath:     kong.Bool(false),
			Hosts:         kong.StringSlice("example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
			RegexPriority: kong.Int(0),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("https-redirect-status-code annotation is correctly processed",
		func(t *testing.T) {
			ingresses := []*networking.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							annotations.IngressClassKey:             "kong",
							"konghq.com/https-redirect-status-code": "301",
						},
					},
					Spec: networking.IngressSpec{
						Rules: []networking.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networking.IngressRuleValue{
									HTTP: &networking.HTTPIngressRuleValue{
										Paths: []networking.HTTPIngressPath{
											{
												Path: "/",
												Backend: networking.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1beta1: ingresses,
				Services:         services,
			})
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)

			assert.Equal(1, len(state.Services),
				"expected one service to be rendered")
			assert.Equal(kong.Service{
				Name:           kong.String("default.foo-svc.80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
			}, state.Services[0].Service)

			assert.Equal(1, len(state.Services[0].Routes),
				"expected one route to be rendered")
			assert.Equal(kong.Route{
				Name:                    kong.String("default.bar.00"),
				StripPath:               kong.Bool(false),
				HTTPSRedirectStatusCode: kong.Int(301),
				Hosts:                   kong.StringSlice("example.com"),
				PreserveHost:            kong.Bool(true),
				Paths:                   kong.StringSlice("/"),
				Protocols:               kong.StringSlice("http", "https"),
				RegexPriority:           kong.Int(0),
			}, state.Services[0].Routes[0].Route)
		})
	t.Run("bad https-redirect-status-code annotation is ignored",
		func(t *testing.T) {
			ingresses := []*networking.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							annotations.IngressClassKey:             "kong",
							"konghq.com/https-redirect-status-code": "whoops",
						},
					},
					Spec: networking.IngressSpec{
						Rules: []networking.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networking.IngressRuleValue{
									HTTP: &networking.HTTPIngressRuleValue{
										Paths: []networking.HTTPIngressPath{
											{
												Path: "/",
												Backend: networking.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1beta1: ingresses,
				Services:         services,
			})
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)

			assert.Equal(1, len(state.Services),
				"expected one service to be rendered")
			assert.Equal(kong.Service{
				Name:           kong.String("default.foo-svc.80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
			}, state.Services[0].Service)

			assert.Equal(1, len(state.Services[0].Routes),
				"expected one route to be rendered")
			assert.Equal(kong.Route{
				Name:          kong.String("default.bar.00"),
				StripPath:     kong.Bool(false),
				Hosts:         kong.StringSlice("example.com"),
				PreserveHost:  kong.Bool(true),
				Paths:         kong.StringSlice("/"),
				Protocols:     kong.StringSlice("http", "https"),
				RegexPriority: kong.Int(0),
			}, state.Services[0].Routes[0].Route)
		})
	t.Run("preserve-host annotation is correctly processed",
		func(t *testing.T) {
			ingresses := []*networking.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/preserve-host":  "faLsE",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networking.IngressSpec{
						Rules: []networking.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networking.IngressRuleValue{
									HTTP: &networking.HTTPIngressRuleValue{
										Paths: []networking.HTTPIngressPath{
											{
												Path: "/",
												Backend: networking.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1beta1: ingresses,
				Services:         services,
			})
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)

			assert.Equal(1, len(state.Services),
				"expected one service to be rendered")
			assert.Equal(kong.Service{
				Name:           kong.String("default.foo-svc.80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
			}, state.Services[0].Service)

			assert.Equal(1, len(state.Services[0].Routes),
				"expected one route to be rendered")
			assert.Equal(kong.Route{
				Name:          kong.String("default.bar.00"),
				StripPath:     kong.Bool(false),
				Hosts:         kong.StringSlice("example.com"),
				PreserveHost:  kong.Bool(false),
				Paths:         kong.StringSlice("/"),
				Protocols:     kong.StringSlice("http", "https"),
				RegexPriority: kong.Int(0),
			}, state.Services[0].Routes[0].Route)
		})
	t.Run("preserve-host annotation with random string is correctly processed",
		func(t *testing.T) {
			ingresses := []*networking.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
							"konghq.com/preserve-host":  "wiggle wiggle wiggle",
						},
					},
					Spec: networking.IngressSpec{
						Rules: []networking.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networking.IngressRuleValue{
									HTTP: &networking.HTTPIngressRuleValue{
										Paths: []networking.HTTPIngressPath{
											{
												Path: "/",
												Backend: networking.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1beta1: ingresses,
				Services:         services,
			})
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)

			assert.Equal(1, len(state.Services),
				"expected one service to be rendered")
			assert.Equal(kong.Service{
				Name:           kong.String("default.foo-svc.80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
			}, state.Services[0].Service)

			assert.Equal(1, len(state.Services[0].Routes),
				"expected one route to be rendered")
			assert.Equal(kong.Route{
				Name:          kong.String("default.bar.00"),
				StripPath:     kong.Bool(false),
				Hosts:         kong.StringSlice("example.com"),
				PreserveHost:  kong.Bool(true),
				Paths:         kong.StringSlice("/"),
				Protocols:     kong.StringSlice("http", "https"),
				RegexPriority: kong.Int(0),
			}, state.Services[0].Routes[0].Route)
		})
	t.Run("regex-priority annotation is correctly processed",
		func(t *testing.T) {
			ingresses := []*networking.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/regex-priority": "10",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networking.IngressSpec{
						Rules: []networking.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networking.IngressRuleValue{
									HTTP: &networking.HTTPIngressRuleValue{
										Paths: []networking.HTTPIngressPath{
											{
												Path: "/",
												Backend: networking.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1beta1: ingresses,
				Services:         services,
			})
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)

			assert.Equal(1, len(state.Services),
				"expected one service to be rendered")
			assert.Equal(kong.Service{
				Name:           kong.String("default.foo-svc.80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
			}, state.Services[0].Service)

			assert.Equal(1, len(state.Services[0].Routes),
				"expected one route to be rendered")
			assert.Equal(kong.Route{
				Name:          kong.String("default.bar.00"),
				StripPath:     kong.Bool(false),
				RegexPriority: kong.Int(10),
				Hosts:         kong.StringSlice("example.com"),
				PreserveHost:  kong.Bool(true),
				Paths:         kong.StringSlice("/"),
				Protocols:     kong.StringSlice("http", "https"),
			}, state.Services[0].Routes[0].Route)
		})
	t.Run("non-integer regex-priority annotation is ignored",
		func(t *testing.T) {
			ingresses := []*networking.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/regex-priority": "IAmAString",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networking.IngressSpec{
						Rules: []networking.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networking.IngressRuleValue{
									HTTP: &networking.HTTPIngressRuleValue{
										Paths: []networking.HTTPIngressPath{
											{
												Path: "/",
												Backend: networking.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1beta1: ingresses,
				Services:         services,
			})
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)

			assert.Equal(1, len(state.Services),
				"expected one service to be rendered")
			assert.Equal(kong.Service{
				Name:           kong.String("default.foo-svc.80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
			}, state.Services[0].Service)

			assert.Equal(1, len(state.Services[0].Routes),
				"expected one route to be rendered")
			assert.Equal(kong.Route{
				Name:          kong.String("default.bar.00"),
				StripPath:     kong.Bool(false),
				RegexPriority: kong.Int(0),
				Hosts:         kong.StringSlice("example.com"),
				PreserveHost:  kong.Bool(true),
				Paths:         kong.StringSlice("/"),
				Protocols:     kong.StringSlice("http", "https"),
			}, state.Services[0].Routes[0].Route)
		})
}

func TestKongProcessClasslessIngress(t *testing.T) {
	assert := assert.New(t)
	t.Run("Kong classless ingress evaluated (true)", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
	})
	t.Run("Kong classless ingress evaluated (false)", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(0, len(state.Services),
			"expected zero service to be rendered")
	})
}

func TestKnativeIngressAndPlugins(t *testing.T) {
	assert := assert.New(t)
	t.Run("knative ingress rule and service-level plugin", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						"networking.knative.dev/ingress.class": annotations.DefaultIngressClass,
					},
				},
				Spec: knative.IngressSpec{
					Rules: []knative.IngressRule{
						{
							Hosts: []string{"my-func.example.com"},
							HTTP: &knative.HTTPIngressRuleValue{
								Paths: []knative.HTTPIngressPath{
									{
										Path: "/",
										AppendHeaders: map[string]string{
											"foo": "bar",
										},
										Splits: []knative.IngressBackendSplit{
											{
												IngressBackend: knative.IngressBackend{
													ServiceNamespace: "foo-ns",
													ServiceName:      "foo-svc",
													ServicePort:      intstr.FromInt(42),
												},
												Percent: 100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						annotations.DeprecatedPluginsKey:       "knative-key-auth",
						"networking.knative.dev/ingress.class": annotations.DefaultIngressClass,
					},
				},
			},
		}
		plugins := []*configurationv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-key-auth",
					Namespace: "foo-ns",
				},
				PluginName: "key-auth",
				Protocols:  []string{"http"},
				Config: configurationv1.Configuration{
					"foo":     "bar",
					"knative": "yo",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			Services:         services,
			KongPlugins:      plugins,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
		svc := state.Services[0]

		assert.Equal(kong.Service{
			Name:           kong.String("foo-ns.foo-svc.42"),
			Host:           kong.String("foo-svc.foo-ns.42.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
		}, svc.Service)

		assert.Equal(1, len(svc.Plugins), "expected one request-transformer plugin")
		assert.Equal(kong.Plugin{
			Name: kong.String("request-transformer"),
			Config: kong.Configuration{
				"add": map[string]interface{}{
					"headers": []string{"foo:bar"},
				},
			},
		}, svc.Plugins[0])

		assert.Equal(1, len(svc.Routes),
			"expected one route to be rendered")
		assert.Equal(kong.Route{
			Name:          kong.String("foo-ns.knative-ingress.00"),
			StripPath:     kong.Bool(false),
			Hosts:         kong.StringSlice("my-func.example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
			RegexPriority: kong.Int(0),
		}, svc.Routes[0].Route)

		assert.Equal(1, len(state.Plugins), "expected one key-auth plugin")
		assert.Equal(kong.Plugin{
			Name: kong.String("key-auth"),
			Config: kong.Configuration{
				"foo":     "bar",
				"knative": "yo",
			},
			Service: &kong.Service{
				ID: kong.String("foo-ns.foo-svc.42"),
			},
			Protocols: kong.StringSlice("http"),
		}, state.Plugins[0].Plugin)
	})
}

func TestKongServiceAnnotations(t *testing.T) {
	assert := assert.New(t)
	t.Run("path annotation is correctly processed", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/path": "/baz",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
		assert.Equal(kong.Service{
			Name:           kong.String("default.foo-svc.80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/baz"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
		}, state.Services[0].Service)

		assert.Equal(1, len(state.Services[0].Routes),
			"expected one route to be rendered")
		assert.Equal(kong.Route{
			Name:          kong.String("default.bar.00"),
			StripPath:     kong.Bool(false),
			Hosts:         kong.StringSlice("example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
			RegexPriority: kong.Int(0),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("host-header annotation is correctly processed", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/host-header": "example.com",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
		assert.Equal(kong.Service{
			Name:           kong.String("default.foo-svc.80"),
			Host:           kong.String("foo-svc.default.80.svc"),
			Path:           kong.String("/"),
			Port:           kong.Int(80),
			ConnectTimeout: kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			WriteTimeout:   kong.Int(60000),
			Retries:        kong.Int(5),
			Protocol:       kong.String("http"),
		}, state.Services[0].Service)

		assert.Equal(1, len(state.Upstreams),
			"expected one upstream to be rendered")
		assert.Equal(kong.Upstream{
			Name:       kong.String("foo-svc.default.80.svc"),
			HostHeader: kong.String("example.com"),
		}, state.Upstreams[0].Upstream)

		assert.Equal(1, len(state.Services[0].Routes),
			"expected one route to be rendered")
		assert.Equal(kong.Route{
			Name:          kong.String("default.bar.00"),
			StripPath:     kong.Bool(false),
			Hosts:         kong.StringSlice("example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
			RegexPriority: kong.Int(0),
		}, state.Services[0].Routes[0].Route)
	})

	t.Run("methods annotation is correctly processed",
		func(t *testing.T) {
			ingresses := []*networking.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/methods":        "POST,GET",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networking.IngressSpec{
						Rules: []networking.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networking.IngressRuleValue{
									HTTP: &networking.HTTPIngressRuleValue{
										Paths: []networking.HTTPIngressPath{
											{
												Path: "/",
												Backend: networking.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			services := []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-svc",
						Namespace: "default",
					},
				},
			}
			store, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1beta1: ingresses,
				Services:         services,
			})
			assert.Nil(err)
			state, err := Build(logrus.New(), store)
			assert.Nil(err)
			assert.NotNil(state)

			assert.Equal(1, len(state.Services),
				"expected one service to be rendered")
			assert.Equal(kong.Service{
				Name:           kong.String("default.foo-svc.80"),
				Host:           kong.String("foo-svc.default.80.svc"),
				Path:           kong.String("/"),
				Port:           kong.Int(80),
				ConnectTimeout: kong.Int(60000),
				ReadTimeout:    kong.Int(60000),
				WriteTimeout:   kong.Int(60000),
				Retries:        kong.Int(5),
				Protocol:       kong.String("http"),
			}, state.Services[0].Service)

			assert.Equal(1, len(state.Services[0].Routes),
				"expected one route to be rendered")
			assert.Equal(kong.Route{
				Name:          kong.String("default.bar.00"),
				StripPath:     kong.Bool(false),
				RegexPriority: kong.Int(0),
				Hosts:         kong.StringSlice("example.com"),
				PreserveHost:  kong.Bool(true),
				Paths:         kong.StringSlice("/"),
				Protocols:     kong.StringSlice("http", "https"),
				Methods:       kong.StringSlice("POST", "GET"),
			}, state.Services[0].Routes[0].Route)
		})
}

func TestDefaultBackend(t *testing.T) {
	assert := assert.New(t)
	t.Run("default backend is processed correctly", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ing-with-default-backend",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Backend: &networking.IngressBackend{
						ServiceName: "default-svc",
						ServicePort: intstr.FromInt(80),
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "default-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
		assert.Equal("default.default-svc.80", *state.Services[0].Name)
		assert.Equal("default-svc.default.80.svc", *state.Services[0].Host)
		assert.Equal(1, len(state.Services[0].Routes),
			"expected one routes to be rendered")
		assert.Equal("default.ing-with-default-backend", *state.Services[0].Routes[0].Name)
		assert.Equal("/", *state.Services[0].Routes[0].Paths[0])
	})

	t.Run("client-cert secret doesn't exist", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered")

		assert.Equal(1, len(state.Services))
		assert.Nil(state.Services[0].ClientCertificate)
	})
}

func TestParserSecret(t *testing.T) {
	assert := assert.New(t)
	t.Run("invalid TLS secret", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"bar.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(""),
					"tls.key": []byte(""),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Secrets:          secrets,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered with empty secret")
	})
	t.Run("duplicate certificates", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret2",
							Hosts:      []string{"bar.com"},
						},
					},
				},
			},
		}

		t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		t2, _ := time.Parse(time.RFC3339, "2006-01-02T15:05:05Z")
		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       "3e8edeca-7d23-4e02-84c9-437d11b746a6",
					Name:      "secret1",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: t1,
					},
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       "fc28a22c-41e1-4cd6-9099-fd7756ffe58e",
					Name:      "secret2",
					Namespace: "ns1",
					CreationTimestamp: metav1.Time{
						Time: t2,
					},
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Secrets:          secrets,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Certificates),
			"certificates are de-duplicated")

		sort.SliceStable(state.Certificates[0].SNIs, func(i, j int) bool {
			return strings.Compare(*state.Certificates[0].SNIs[i],
				*state.Certificates[0].SNIs[j]) > 0
		})
		assert.Equal(kongstate.Certificate{
			Certificate: kong.Certificate{
				ID:   kong.String("3e8edeca-7d23-4e02-84c9-437d11b746a6"),
				Cert: kong.String(tlsPairs[0].Cert),
				Key:  kong.String(tlsPairs[0].Key),
				SNIs: kong.StringSlice("foo.com", "bar.com"),
			},
		}, state.Certificates[0])
	})
	t.Run("duplicate SNIs", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "ns1",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret2",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret2",
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Secrets:          secrets,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Certificates),
			"SNIs are de-duplicated")
	})
}

func TestPluginAnnotations(t *testing.T) {
	assert := assert.New(t)
	t.Run("simple association", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "foo-svc",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
		}
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.DeprecatedPluginsKey: "foo-plugin",
						annotations.IngressClassKey:      annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		plugins := []*configurationv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "key-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
					"add": map[string]interface{}{
						"headers": []interface{}{
							"header1:value1",
							"header2:value2",
						},
					},
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Services:         services,
			KongPlugins:      plugins,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		p := state.Plugins[0].Plugin
		p.Route = nil
		assert.Equal(p, kong.Plugin{
			Name:      kong.String("key-auth"),
			Protocols: kong.StringSlice("grpc"),
			Config: kong.Configuration{
				"foo": "bar",
				"add": map[string]interface{}{
					"headers": []interface{}{
						"header1:value1",
						"header2:value2",
					},
				},
			},
		})
	})
	t.Run("KongPlugin takes precedence over KongPlugin", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "foo-svc",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
		}
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.DeprecatedPluginsKey: "foo-plugin",
						annotations.IngressClassKey:      annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		clusterPlugins := []*configurationv1.KongClusterPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "basic-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
				},
			},
		}
		plugins := []*configurationv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "key-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1:   ingresses,
			Services:           services,
			KongPlugins:        plugins,
			KongClusterPlugins: clusterPlugins,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		assert.Equal("key-auth", *state.Plugins[0].Name)
		assert.Equal("grpc", *state.Plugins[0].Protocols[0])
	})
	t.Run("KongClusterPlugin association", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "foo-svc",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
		}
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.DeprecatedPluginsKey: "foo-plugin",
						annotations.IngressClassKey:      annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		clusterPlugins := []*configurationv1.KongClusterPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "basic-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1:   ingresses,
			Services:           services,
			KongClusterPlugins: clusterPlugins,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		assert.Equal("basic-auth", *state.Plugins[0].Name)
		assert.Equal("grpc", *state.Plugins[0].Protocols[0])
	})
	t.Run("missing plugin", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.DeprecatedPluginsKey: "does-not-exist",
						annotations.IngressClassKey:      annotations.DefaultIngressClass,
					},
				},
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
	})
}

func TestGetEndpoints(t *testing.T) {
	tests := []struct {
		name   string
		svc    *corev1.Service
		port   *corev1.ServicePort
		proto  corev1.Protocol
		fn     func(string, string) (*corev1.Endpoints, error)
		result []utils.Endpoint
	}{
		{
			"no service should return 0 endpoints",
			nil,
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]utils.Endpoint{},
		},
		{
			"no service port should return 0 endpoints",
			&corev1.Service{},
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service without endpoints should return 0 endpoints",
			&corev1.Service{},
			&corev1.ServicePort{Name: "default"},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service type ServiceTypeExternalName service with an invalid port should return 0 endpoints",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeExternalName,
				},
			},
			&corev1.ServicePort{Name: "default"},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service type ServiceTypeExternalName with a valid port should return one endpoint",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:         corev1.ServiceTypeExternalName,
					ExternalName: "10.0.0.1.xip.io",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{
				{
					Address: "10.0.0.1.xip.io",
					Port:    "80",
				},
			},
		},
		{
			"a service with ingress.kubernetes.io/service-upstream annotation should return one endpoint",
			&corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						"ingress.kubernetes.io/service-upstream": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(2080),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{
				{
					Address: "foo.bar.svc",
					Port:    "2080",
				},
			},
		},
		{
			"should return no endpoints when there is an error searching for endpoints",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when the protocol does not match",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolUDP,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when there is no ready Addresses",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							NotReadyAddresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolUDP,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when the name of the port name do not match any port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolTCP,
									Port:     int32(80),
									Name:     "another-name",
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return one endpoint when the name of the port name match a port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolTCP,
									Port:     int32(80),
									Name:     "default",
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
		{
			"should return one endpoint when the name of the port name match more than one port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromString("port-1"),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "port-1",
				TargetPort: intstr.FromString("port-1"),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Name:     "port-1",
									Protocol: corev1.ProtocolTCP,
									Port:     80,
								},
								{
									Name:     "port-1",
									Protocol: corev1.ProtocolTCP,
									Port:     80,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := getEndpoints(logrus.New(), testCase.svc, testCase.port, testCase.proto, testCase.fn)
			if len(testCase.result) != len(result) {
				t.Errorf("expected %v Endpoints but got %v", testCase.result, len(result))
			}
		})
	}
}

func Test_knativeSelectSplit(t *testing.T) {
	type args struct {
		splits []knative.IngressBackendSplit
	}
	tests := []struct {
		name string
		args args
		want knative.IngressBackendSplit
	}{
		{
			name: "empty ingress",
		},
		{
			name: "no split",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 100,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "foo-ns",
					ServiceName:      "foo-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 100,
			},
		},
		{
			name: "less than 100%% but one split only",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 42,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "foo-ns",
					ServiceName:      "foo-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 42,
			},
		},
		{
			name: "multiple splits with unequal splits",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "bar-ns",
							ServiceName:      "bar-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 42,
					},
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 58,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "foo-ns",
					ServiceName:      "foo-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 58,
			},
		},
		{
			name: "multiple splits with unequal splits",
			args: args{
				splits: []knative.IngressBackendSplit{
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "bar-ns",
							ServiceName:      "bar-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 40,
					},
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "baz-ns",
							ServiceName:      "baz-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 20,
					},
					{
						IngressBackend: knative.IngressBackend{
							ServiceNamespace: "foo-ns",
							ServiceName:      "foo-svc",
							ServicePort:      intstr.FromInt(42),
						},
						Percent: 40,
					},
				},
			},
			want: knative.IngressBackendSplit{
				IngressBackend: knative.IngressBackend{
					ServiceNamespace: "bar-ns",
					ServiceName:      "bar-svc",
					ServicePort:      intstr.FromInt(42),
				},
				Percent: 40,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := knativeSelectSplit(tt.args.splits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("knativeSelectSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}
