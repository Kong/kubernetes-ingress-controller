package parser

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
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
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIDEzCCAfugAwIBAgIUOwYJvXJ+s0qX9uAKFjW0zExV51IwDQYJKoZIhvcNAQEF
BQAwJTEKMAgGA1UEAwwBIzEKMAgGA1UECgwBIzELMAkGA1UEBhMCQ04wHhcNMjEw
NDMwMDAwMDAwWhcNMjEwNDMwMDAwMDAwWjBDMQswCQYDVQQGEwJDTjEKMAgGA1UE
CAwBIzEKMAgGA1UEBwwBIzEKMAgGA1UECgwBIzEQMA4GA1UEAwwHZm9vLmNvbTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPRFV285ScP1ntF3zlj60GC0
DJyCEX6Ji38gBf+6An6Zk7D+3Aif9C/3e7V0811x0VoO4o9ZdQUSdxmE9fj/ADOU
OM/AYf62L51d/zdqXVaF89vpqPk8em4179wo6jg2IiCewGVLTtuMa/5Mud8XZOly
tcOXS7ZnCbfm/XklwGL1rAmWhOTSDHlIH5bbC46tmi3E9Cjp+VTiwzVCgVtrLkzY
0cjs72m2wb5uZ9TlT7n1TKYjdX74FvYp4X70YEcFYEUmFMxMV7otkJ7wTWWVNah/
ZsojaiJ48ueJFQR1S9utYA/h6LcA4T6UQJxw7+6SjJElLCHGht5UHFvQkjQvxZkC
AwEAAaMdMBswCwYDVR0RBAQwAoIAMAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQEF
BQADggEBAHE4U9SlCIVNjpfOyfH0NPhxLTAqH83GQKJc7TgQFmhby1dfQE7MOTaN
ayA1RJ0qKcNGlHP70M/Xc8TIF+E7pOASqa+zNztiv14zHIgJC9oGJcwt1sh8GADz
4EJSQ1mIRxbgs39BA9FDY91HBa3RfLxkmyTbQK1rhKdh8aBYr0/6R1oAdKEQF/vQ
HxD4NCpJruxp7g+RSet1PB12GOao1Ntfb7kOLAHzYW3yvTsCaQ7EdeueOs8dv+G/
Ncy+4n/l3audbi+WQFfEvb1bwyADPpp90C9OczHzpR4+dtuR4oUXB6ZXimB3MljC
BhoOkUOMjrKl/QDkB5pxa/IxURffFDs=
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQD0RVdvOUnD9Z7R
d85Y+tBgtAycghF+iYt/IAX/ugJ+mZOw/twIn/Qv93u1dPNdcdFaDuKPWXUFEncZ
hPX4/wAzlDjPwGH+ti+dXf83al1WhfPb6aj5PHpuNe/cKOo4NiIgnsBlS07bjGv+
TLnfF2TpcrXDl0u2Zwm35v15JcBi9awJloTk0gx5SB+W2wuOrZotxPQo6flU4sM1
QoFbay5M2NHI7O9ptsG+bmfU5U+59UymI3V++Bb2KeF+9GBHBWBFJhTMTFe6LZCe
8E1llTWof2bKI2oiePLniRUEdUvbrWAP4ei3AOE+lECccO/ukoyRJSwhxobeVBxb
0JI0L8WZAgMBAAECggEAS3gBA4QXnkuMvfrd7e/P4ZC/NLua3BVy29vw/olhq+uX
DeAv6xpAP3Ge7vkrF3vKyqA+rztcRCzoTyIdPMjMLyNkSguOOzveiig4ac6a99h6
9z7Bzf85dEOFz+d0NdnDwYBKwQD7ZCmGVBMwehSoQNgIAF4GLly3S/I57ewT/H6A
GjknY/jCmk+L9388hjcL0jrEJR/br2O/o6f1zdRYWqqb9A1wDW4zkiu4Wrq5in9s
cxQ/7667eGkciD3HkJvwcbi7xg9ZJHxCWScVYGRBX9ek6fMKxML5hUsITZjt5zxF
p+HmOKJcii6hlR1RWaUsbrpQOHVui3US7CAJNR/gAQKBgQD7zE0z2XGv81XjIvMk
sS+IvtsSGpvoUI2QRbdnelC8ahCdKj5PmVQyfhPxrgNCTu6k7VzQUDL+wZDGKoRL
NEaRkoHz7tVzBE7DY7Y2SD0yfjT477w98iaF/nwortmhpXms0KyzPhZOF5d5166q
PDR31NIFvmy2H+Hh9bVM5BYaGQKBgQD4WOIgocc+pXn+3fehNT3qedrvGuXGYX7I
PAO/4zM/oP/0TtxKTz5wGAFR9heBKfogW2jYUBBOofraLMJq3X+T9jEOXuQ9+UQq
HaybHdQycxpTIWhtiAs9khvSbuEBs2SXyKussPGW8Do5uVfi4/KWWu/wcTzMlfEv
w207iaN3gQKBgQDAh7u0XJx4PCi871lZAf5logGiOyRhI07LNPOCtN0M5FDly4ov
lP7zSMH5NuQZDH+fLjucsOX9M4Z+b74OPt+CqbKiEUm2k2GiNxj5Mo1QkX3xpmWa
PBDGvgqzlNalqgB6amjS+TNW7OUO7iMI2dYIlnsslylKrOArxZOmQnS/6QKBgQCs
ZVcj++nKDSjwybk6yTDf8hMO5IcY/Vj7Ot4HeHp88xB60buOQhA/1AomkUSjvzYI
/Ct97aZET6FJjsSvVm9XkRFgvnKGquCss8i8LSq+krR1fL13O3dCGIkDvUCo45Uy
4HR7/qDWfJCOvaDKuh4OTbY+HP1tr7CrzWeoatV1AQKBgFFmlMWrIThfjtjVWDTg
+QPLQTofB1A3lrCmB52iBdUMi0qGnExLn8aiy54wPz/I7rEplsLzg2hmDNuBPM7q
QLAtVaZd9SSi4Z/RX6B4L3Rj0Mwfn+tbrtYO5Pyhi40hiXf4aMgbVDFYMR0MMmH0
4uiYeQPmK6USKjntOFQ0eNOe
-----END PRIVATE KEY-----`,
		},
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIDEzCCAfugAwIBAgIUBXt+uyE/nkPDO+LmBHyau3e7QpIwDQYJKoZIhvcNAQEF
BQAwJTEKMAgGA1UEAwwBIzEKMAgGA1UECgwBIzELMAkGA1UEBhMCQ04wHhcNMjEw
NDMwMDAwMDAwWhcNMjEwNDMwMDAwMDAwWjBDMQswCQYDVQQGEwJDTjEKMAgGA1UE
CAwBIzEKMAgGA1UEBwwBIzEKMAgGA1UECgwBIzEQMA4GA1UEAwwHZm9vLmNvbTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANWQF6xf6EAYSdSUnsOQlojg
6zjI+U26JpfJqkqsZJ4aOQdUnlDUY+nYEKb7VSH10rszmHltJfKXUBWlXAK/qzZs
8fnnSa7mI1BejVDQpZ+QYPmxgSfmKlxwgYbHxCB5dklzZv8OGUIeKWau+pnc8Bsr
gofM+MDdHV4Njcfdf9+g23gNGNQN1MxRdbKRUmmTV3Yiq2YAYO62HkYemd7KrdO3
QLsKMdodOxN8of6E6654ecctfiae3xkdSYDcBr8ig0T1KVhremFNAAw9mEOzmpR+
tRKB0Z7hwsR5wYviKmkMa3MyakYAb1EE8R/dt7WedqXDC+JYOjLWh1vVDcrib3cC
AwEAAaMdMBswCwYDVR0RBAQwAoIAMAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQEF
BQADggEBAGkM8IFbeonb8hguAccoIM3sZgvUnFq3jbdfi4xD4AgjqqLnGRSZxobh
jDQIX+4wEbGillWspZvg4W8Eb20LOwWpg9J19ORTFLzz09+rkRieawINbs6diirh
smkEPTI18kKZb1virUoSipviiXFRLra9YOG16YTCoZyeToSe1UKfvBJJPro451tU
RjIaqFeHIC9t1SoMIHJS1H4jjpdLNunEgNDUdYOMXfFqX2HgYU9626cjxt0otrlw
6AN08UAvnh72gjaMa7YG0/SQJIJvKFPrC/C4A6vcLd+RMf3y4uNBsMPLUX/ksNs+
IqaxjG7sWj87o/uWzTCujD7PjdQ2P/o=
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDVkBesX+hAGEnU
lJ7DkJaI4Os4yPlNuiaXyapKrGSeGjkHVJ5Q1GPp2BCm+1Uh9dK7M5h5bSXyl1AV
pVwCv6s2bPH550mu5iNQXo1Q0KWfkGD5sYEn5ipccIGGx8QgeXZJc2b/DhlCHilm
rvqZ3PAbK4KHzPjA3R1eDY3H3X/foNt4DRjUDdTMUXWykVJpk1d2IqtmAGDuth5G
Hpneyq3Tt0C7CjHaHTsTfKH+hOuueHnHLX4mnt8ZHUmA3Aa/IoNE9SlYa3phTQAM
PZhDs5qUfrUSgdGe4cLEecGL4ippDGtzMmpGAG9RBPEf3be1nnalwwviWDoy1odb
1Q3K4m93AgMBAAECggEAWqZBBi/ne11T6WH0SfqTiyI9desIt09uljtJh0zJoLps
vonScWjVyCqoVtiT+EhJ3KA39kK4OtKzcZyxA5Gr8PwpcfQUuKKDFtgrj8QgqSw1
nYcU1YTDMl2m/TjKLSahaOgNLfIkEoaO+EEJqkm3uTCsxUvhwquYjZWfOtjwMFFm
Oa5m1kb+QvIRTle3wumzN6yIUXdQRFsK0D8UQUQw5UWGZ8xhwa1TJRGdGwxFT+Kk
hH+/+PQO62+bWJ27EQFZ9mTMuleG2YGa20u05tXlj3qUmws53yuzwALCBGQYmolh
2I+pfQ4vAPuXWSL33QIAAkOwPHExbKNyD/w6r79wsQKBgQD4cSejoNFcEzIqlONj
RNMh2UhpwhMlboU7Xs7x5BAggImk6/wPzOdCL1ytfBJoaNs8MD6tYWU8sa38Vsdh
DwcbRRTSFzCdZGZw7E1g8/hOqY6dcirJYmOOTAVM0UBy5M2nL/W385B8+9IpzSTN
ylsJTMZNBh1U9da36D2lqdSbAwKBgQDcD01CW7rDnXS6jnn03KAN2DGXAR3JxvKO
GSdhcmojORgKYD9tN/AnNriWzUVTB7L0hzsLzuyqFW8g6nY8PKV+IAYBmq12XlpO
llppwWoDz7JXcO3pKhdwfjYLDag92fWChO4pFctwZhIRurD0FDdvlzM9Ou5LKyzE
px+sfJ6VfQKBgQDmxNyoDfpKR35itTfd/pELMPQlYGM+csGI21DouQoN1reEHLte
xdrNzIaOkt/aYgO6jam6jOxnizdsXAMI8deNCgcD+wxqNlc9bxGVDClKkA7ryp9Z
vm1xQMXvi7MMxeEM+eyJONGudo7Jy0bzbJJROiI8a4CVaWFQJIYWuoDElwKBgQC9
BTPCrRImVoheemVNK3kbizlFUMMqf4X3Aqot7N44NSFuQDAa+3J/7GPvvJAweqt/
mOzh/rKQgeq7pkk7AojQZmdiV19qDi+Z01IEBwuuDGhO7YSdw/bwPKjlI60Au8hD
fTUo+zyM5k/dBLRcY0UeyAxOKuFmlcZVgIwXV8/L7QKBgF/uzlzkW7kujr9QLuas
b8mVWKEUxmwD4ppmNDrszztQRel0ujl/15bMmXav2UaAzIvnuo0JqVnTc+e2/xu+
5/cjaN4hxgSSdxsYgDjk727A5L0jb85rVF559qeZck5PatMDt/Lbaz7BeMLw0c3y
t/0TgcCP3Nl7JDtqRP6PrnZp
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
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"foo1": "bar1"}`),
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
		IngressesV1beta1: []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
						annotations.AnnotationPrefix + annotations.PluginsKey: "bar-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.net",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
					ConfigFrom: &configurationv1.ConfigSource{
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
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					ConfigFrom: &configurationv1.ConfigSource{
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
					ConfigFrom: &configurationv1.ConfigSource{
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
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &configurationv1.ConfigSource{
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
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &configurationv1.ConfigSource{
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
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"fake": true}`),
					},
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					ConfigFrom: &configurationv1.ConfigSource{
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
					ConfigFrom: &configurationv1.ConfigSource{
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
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
					ConfigFrom: &configurationv1.NamespacedConfigSource{
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
					TLS: []networkingv1beta1.IngressTLS{
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
						"konghq.com/client-cert": "secret1",
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
					TLS: []networkingv1beta1.IngressTLS{
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
						"konghq.com/client-cert": "secret1",
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						"konghq.com/strip-path":     "trUe",
						annotations.IngressClassKey: "kong",
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: "kong",
						"konghq.com/strip-path":     "false",
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
			ingresses := []*networkingv1beta1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							annotations.IngressClassKey:             "kong",
							"konghq.com/https-redirect-status-code": "301",
						},
					},
					Spec: networkingv1beta1.IngressSpec{
						Rules: []networkingv1beta1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networkingv1beta1.IngressRuleValue{
									HTTP: &networkingv1beta1.HTTPIngressRuleValue{
										Paths: []networkingv1beta1.HTTPIngressPath{
											{
												Path: "/",
												Backend: networkingv1beta1.IngressBackend{
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
			ingresses := []*networkingv1beta1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							annotations.IngressClassKey:             "kong",
							"konghq.com/https-redirect-status-code": "whoops",
						},
					},
					Spec: networkingv1beta1.IngressSpec{
						Rules: []networkingv1beta1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networkingv1beta1.IngressRuleValue{
									HTTP: &networkingv1beta1.HTTPIngressRuleValue{
										Paths: []networkingv1beta1.HTTPIngressPath{
											{
												Path: "/",
												Backend: networkingv1beta1.IngressBackend{
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
			ingresses := []*networkingv1beta1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/preserve-host":  "faLsE",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networkingv1beta1.IngressSpec{
						Rules: []networkingv1beta1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networkingv1beta1.IngressRuleValue{
									HTTP: &networkingv1beta1.HTTPIngressRuleValue{
										Paths: []networkingv1beta1.HTTPIngressPath{
											{
												Path: "/",
												Backend: networkingv1beta1.IngressBackend{
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
			ingresses := []*networkingv1beta1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
							"konghq.com/preserve-host":  "wiggle wiggle wiggle",
						},
					},
					Spec: networkingv1beta1.IngressSpec{
						Rules: []networkingv1beta1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networkingv1beta1.IngressRuleValue{
									HTTP: &networkingv1beta1.HTTPIngressRuleValue{
										Paths: []networkingv1beta1.HTTPIngressPath{
											{
												Path: "/",
												Backend: networkingv1beta1.IngressBackend{
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
			ingresses := []*networkingv1beta1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/regex-priority": "10",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networkingv1beta1.IngressSpec{
						Rules: []networkingv1beta1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networkingv1beta1.IngressRuleValue{
									HTTP: &networkingv1beta1.HTTPIngressRuleValue{
										Paths: []networkingv1beta1.HTTPIngressPath{
											{
												Path: "/",
												Backend: networkingv1beta1.IngressBackend{
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
			ingresses := []*networkingv1beta1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/regex-priority": "IAmAString",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networkingv1beta1.IngressSpec{
						Rules: []networkingv1beta1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networkingv1beta1.IngressRuleValue{
									HTTP: &networkingv1beta1.HTTPIngressRuleValue{
										Paths: []networkingv1beta1.HTTPIngressPath{
											{
												Path: "/",
												Backend: networkingv1beta1.IngressBackend{
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
	t.Run("route buffering options are processed (true)", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route-buffering-test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:     annotations.DefaultIngressClass,
						"konghq.com/request-buffering":  "True",
						"konghq.com/response-buffering": "True",
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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

		assert.Equal(1, len(state.Services), "expected one service to be rendered")
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

		assert.Equal(1, len(state.Services[0].Routes), "expected one route to be rendered")
		assert.Equal(kong.Route{
			Name:              kong.String("default.route-buffering-test.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(0),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RequestBuffering:  kong.Bool(true),
			ResponseBuffering: kong.Bool(true),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route buffering options are processed (false)", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route-buffering-test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:     annotations.DefaultIngressClass,
						"konghq.com/request-buffering":  "False",
						"konghq.com/response-buffering": "False",
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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

		assert.Equal(1, len(state.Services), "expected one service to be rendered")
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

		assert.Equal(1, len(state.Services[0].Routes), "expected one route to be rendered")
		assert.Equal(kong.Route{
			Name:              kong.String("default.route-buffering-test.00"),
			StripPath:         kong.Bool(false),
			RegexPriority:     kong.Int(0),
			Hosts:             kong.StringSlice("example.com"),
			PreserveHost:      kong.Bool(true),
			Paths:             kong.StringSlice("/"),
			Protocols:         kong.StringSlice("http", "https"),
			RequestBuffering:  kong.Bool(false),
			ResponseBuffering: kong.Bool(false),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route buffering options are not processed with bad annotation values", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route-buffering-test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey:     annotations.DefaultIngressClass,
						"konghq.com/request-buffering":  "invalid-value",
						"konghq.com/response-buffering": "invalid-value",
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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

		assert.Equal(1, len(state.Services), "expected one service to be rendered")
		assert.Equal(1, len(state.Services[0].Routes), "expected one route to be rendered")
		assert.Empty(state.Services[0].Routes[0].Route.RequestBuffering)
		assert.Empty(state.Services[0].Routes[0].Route.ResponseBuffering)
	})
}

func TestKongProcessClasslessIngress(t *testing.T) {
	assert := assert.New(t)
	t.Run("Kong classless ingress evaluated (true)", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
	t.Run("knative ingress annotated with konghq.com/override", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress-with-override",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						"networking.knative.dev/ingress.class":                      annotations.DefaultIngressClass,
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "https-only",
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
		kongIngresses := []*configurationv1.KongIngress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "https-only",
					Namespace: "foo-ns",
				},
				Route: &kong.Route{
					Protocols:               kong.StringSlice("https"),
					HTTPSRedirectStatusCode: kong.Int(308),
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "foo-ns",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			KongIngresses:    kongIngresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services), "expected one knative service")
		svc := state.Services[0]

		assert.Equal(1, len(svc.Routes), "expected one route in knative service")
		route := svc.Routes[0]

		assert.Equal(kong.StringSlice("https"), route.Protocols, "expected https protocol after override")
		assert.Equal(kong.Int(308), route.HTTPSRedirectStatusCode, "expected 308 after override")
	})
	t.Run("knative ingress without konghq.com/override", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress-without-override",
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
		kongIngresses := []*configurationv1.KongIngress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "https-only",
					Namespace: "foo-ns",
				},
				Route: &kong.Route{
					Protocols:               kong.StringSlice("https"),
					HTTPSRedirectStatusCode: kong.Int(308),
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "foo-ns",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			KongIngresses:    kongIngresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services), "expected one knative service")
		svc := state.Services[0]

		assert.Equal(1, len(svc.Routes), "expected one route in knative service")
		route := svc.Routes[0]

		assert.NotEqual(kong.StringSlice("https"), route.Protocols)
		assert.Nil(route.HTTPSRedirectStatusCode)
	})
	t.Run("knative ingress with multiple konghq.com annotations", func(t *testing.T) {
		ingresses := []*knative.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-ingress-with-annotations",
					Namespace: "foo-ns",
					Annotations: map[string]string{
						"networking.knative.dev/ingress.class":                          annotations.DefaultIngressClass,
						annotations.AnnotationPrefix + annotations.ProtocolsKey:         "https",
						annotations.AnnotationPrefix + annotations.HTTPSRedirectCodeKey: "308",
						annotations.AnnotationPrefix + annotations.StripPathKey:         "true",
						annotations.AnnotationPrefix + annotations.MethodsKey:           "POST,PUT",
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
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: ingresses,
			Services:         services,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)

		assert.Equal(1, len(state.Services), "expected one knative service")
		svc := state.Services[0]

		assert.Equal(1, len(svc.Routes), "expected one route in knative service")
		route := svc.Routes[0]

		assert.Equal(kong.StringSlice("https"), route.Protocols, "expected https after konghq.com/protocols")
		assert.Equal(kong.Int(308), route.HTTPSRedirectStatusCode, "expected 308 after konghq.com/https-redirect-status-code")
		assert.Equal(kong.Bool(true), route.StripPath, "expected true after konghq.com/strip-path")
		assert.Equal(kong.StringSlice("POST", "PUT"), route.Methods, "expected POST and PUT after konghq.com/methods")
	})
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
						annotations.AnnotationPrefix + annotations.PluginsKey: "knative-key-auth",
						"networking.knative.dev/ingress.class":                annotations.DefaultIngressClass,
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
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar", "knative": "yo"}`),
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
						"konghq.com/path": "/baz",
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
						"konghq.com/host-header": "example.com",
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
			ingresses := []*networkingv1beta1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/methods":        "POST,GET",
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: networkingv1beta1.IngressSpec{
						Rules: []networkingv1beta1.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: networkingv1beta1.IngressRuleValue{
									HTTP: &networkingv1beta1.HTTPIngressRuleValue{
										Paths: []networkingv1beta1.HTTPIngressPath{
											{
												Path: "/",
												Backend: networkingv1beta1.IngressBackend{
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ing-with-default-backend",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Backend: &networkingv1beta1.IngressBackend{
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
					TLS: []networkingv1beta1.IngressTLS{
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
						"konghq.com/client-cert": "secret1",
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
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
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
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
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
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
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
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

func TestParserSNI(t *testing.T) {
	assert := assert.New(t)
	t.Run("route includes SNI when TLS info present, but not for wildcard hostnames", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"example.com", "*.example.com"},
						},
					},
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
						{
							Host: "*.example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: ingresses,
			Secrets:          secrets,
		})
		assert.Nil(err)
		state, err := Build(logrus.New(), store)
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(kong.Route{
			Name:          kong.String("default.foo.00"),
			StripPath:     kong.Bool(false),
			RegexPriority: kong.Int(0),
			Hosts:         kong.StringSlice("example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
		}, state.Services[0].Routes[0].Route)
		assert.Equal(kong.Route{
			Name:          kong.String("default.foo.10"),
			StripPath:     kong.Bool(false),
			RegexPriority: kong.Int(0),
			Hosts:         kong.StringSlice("*.example.com"),
			SNIs:          nil,
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
		}, state.Services[0].Routes[1].Route)
	})
	t.Run("route does not include SNI when TLS info absent", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		assert.Equal(kong.Route{
			Name:          kong.String("default.foo.00"),
			StripPath:     kong.Bool(false),
			RegexPriority: kong.Int(0),
			Hosts:         kong.StringSlice("example.com"),
			SNIs:          nil,
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
		}, state.Services[0].Routes[0].Route)
	})
}

func TestParserHostAliases(t *testing.T) {
	assert := assert.New(t)
	annHostAliasesKey := annotations.AnnotationPrefix + annotations.HostAliasesKey
	t.Run("route Hosts includes Host-Aliases when Host-Aliases are present", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
						annHostAliasesKey:           "*.example.com,*.sample.com,*.illustration.com",
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		assert.Equal(kong.Route{
			Name:          kong.String("default.foo.00"),
			StripPath:     kong.Bool(false),
			RegexPriority: kong.Int(0),
			Hosts:         kong.StringSlice("example.com", "*.example.com", "*.sample.com", "*.illustration.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route Hosts remain unmodified when Host-Aliases are not present", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		assert.Equal(kong.Route{
			Name:          kong.String("default.foo.00"),
			StripPath:     kong.Bool(false),
			RegexPriority: kong.Int(0),
			Hosts:         kong.StringSlice("example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
		}, state.Services[0].Routes[0].Route)
	})
	t.Run("route Hosts will not contain duplicates when Host-Aliases duplicates the host", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
						annHostAliasesKey:           "example.com,*.example.com",
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		assert.Equal(kong.Route{
			Name:          kong.String("default.foo.00"),
			StripPath:     kong.Bool(false),
			RegexPriority: kong.Int(0),
			Hosts:         kong.StringSlice("example.com", "*.example.com"),
			PreserveHost:  kong.Bool(true),
			Paths:         kong.StringSlice("/"),
			Protocols:     kong.StringSlice("http", "https"),
		}, state.Services[0].Routes[0].Route)
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{
					"foo": "bar",
					"add": {
						"headers": [
							"header1:value1",
							"header2:value2"
							]
						}
					}`),
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar"}`),
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
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar"}`),
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "foo-plugin",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"foo": "bar"}`),
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
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "does-not-exist",
						annotations.IngressClassKey:                           annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1beta1.IngressRuleValue{
								HTTP: &networkingv1beta1.HTTPIngressRuleValue{
									Paths: []networkingv1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: networkingv1beta1.IngressBackend{
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
		result []util.Endpoint
	}{
		{
			"no service should return 0 endpoints",
			nil,
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]util.Endpoint{},
		},
		{
			"no service port should return 0 endpoints",
			&corev1.Service{},
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]util.Endpoint{},
		},
		{
			"a service without endpoints should return 0 endpoints",
			&corev1.Service{},
			&corev1.ServicePort{Name: "default"},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]util.Endpoint{},
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
			[]util.Endpoint{},
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
			[]util.Endpoint{
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
			[]util.Endpoint{
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
			[]util.Endpoint{},
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
			[]util.Endpoint{},
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
			[]util.Endpoint{},
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
			[]util.Endpoint{},
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
			[]util.Endpoint{
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
			[]util.Endpoint{
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

func TestPickPort(t *testing.T) {
	assert := assert.New(t)
	svc0 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-0",
			Namespace: "foo-namespace",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "port1", Port: 111, TargetPort: intstr.FromInt(1111)},
				{Name: "port2", Port: 222, TargetPort: intstr.FromString("port1")},
				{Name: "port3", Port: 333, TargetPort: intstr.FromString("potato")},
				{Port: 444},
			},
		},
	}

	svc1 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-1",
			Namespace: "foo-namespace",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "port1", Port: 9999},
			},
		},
	}

	svc2 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-2",
			Namespace: "foo-namespace",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: "external.example.com",
		},
	}

	endpointList := []*corev1.Endpoints{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "service-0", Namespace: "foo-namespace"},
			Subsets: []corev1.EndpointSubset{{
				Addresses: []corev1.EndpointAddress{{IP: "1.1.1.1"}},
				Ports: []corev1.EndpointPort{
					{Name: "port1", Port: 111, Protocol: "TCP"},
					{Name: "port2", Port: 222, Protocol: "TCP"},
					{Name: "port3", Port: 333, Protocol: "TCP"},
				},
			}},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "service-1", Namespace: "foo-namespace"},
			Subsets: []corev1.EndpointSubset{{
				Addresses: []corev1.EndpointAddress{{IP: "2.2.2.2"}},
				Ports: []corev1.EndpointPort{
					{Name: "port1", Port: 9999, Protocol: "TCP"},
				},
			}},
		},
	}

	for _, tt := range []struct {
		name string
		objs store.FakeObjects
		port networkingv1.ServiceBackendPort

		wantTarget string
	}{
		{
			name: "port by number",
			objs: store.FakeObjects{
				Services:  []*corev1.Service{&svc0},
				Endpoints: endpointList,

				IngressesV1: []*networkingv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: networkingv1.IngressSpec{
							Rules: []networkingv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: networkingv1.IngressRuleValue{
										HTTP: &networkingv1.HTTPIngressRuleValue{
											Paths: []networkingv1.HTTPIngressPath{
												{
													Path: "/",
													Backend: networkingv1.IngressBackend{
														Service: &networkingv1.IngressServiceBackend{
															Name: "service-0",
															Port: networkingv1.ServiceBackendPort{Number: 111},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "1.1.1.1:111",
		},
		{
			name: "port by number external name",
			objs: store.FakeObjects{
				Services:  []*corev1.Service{&svc2},
				Endpoints: endpointList,

				IngressesV1: []*networkingv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: networkingv1.IngressSpec{
							Rules: []networkingv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: networkingv1.IngressRuleValue{
										HTTP: &networkingv1.HTTPIngressRuleValue{
											Paths: []networkingv1.HTTPIngressPath{
												{
													Path: "/externalname",
													Backend: networkingv1.IngressBackend{
														Service: &networkingv1.IngressServiceBackend{
															Name: "service-2",
															Port: networkingv1.ServiceBackendPort{Number: 222},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "external.example.com:222",
		},
		{
			name: "port by name",
			objs: store.FakeObjects{
				Services:  []*corev1.Service{&svc0},
				Endpoints: endpointList,

				IngressesV1: []*networkingv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: networkingv1.IngressSpec{
							Rules: []networkingv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: networkingv1.IngressRuleValue{
										HTTP: &networkingv1.HTTPIngressRuleValue{
											Paths: []networkingv1.HTTPIngressPath{
												{
													Path: "/",
													Backend: networkingv1.IngressBackend{
														Service: &networkingv1.IngressServiceBackend{
															Name: "service-0",
															Port: networkingv1.ServiceBackendPort{Name: "port3"},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "1.1.1.1:333",
		},
		{
			name: "port implicit",
			objs: store.FakeObjects{
				Services:  []*corev1.Service{&svc1},
				Endpoints: endpointList,

				IngressesV1: []*networkingv1.Ingress{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "foo",
							Namespace:   "foo-namespace",
							Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
						},
						Spec: networkingv1.IngressSpec{
							Rules: []networkingv1.IngressRule{
								{
									Host: "example.com",
									IngressRuleValue: networkingv1.IngressRuleValue{
										HTTP: &networkingv1.HTTPIngressRuleValue{
											Paths: []networkingv1.HTTPIngressPath{
												{
													Path: "/",
													Backend: networkingv1.IngressBackend{
														Service: &networkingv1.IngressServiceBackend{
															Name: "service-1",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantTarget: "2.2.2.2:9999",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			store, err := store.NewFakeStore(tt.objs)
			assert.NoError(err)

			state, err := Build(logrus.New(), store)
			assert.NoError(err)

			assert.Equal(tt.wantTarget, *state.Upstreams[0].Targets[0].Target.Target)
		})
	}
}

func TestCertificate(t *testing.T) {
	assert := assert.New(t)
	t.Run("same host with multiple namespace return the first namespace/secret by asc ", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "ns3",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "ns2",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
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
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID("6392jz73-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "ns2",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[1].Cert),
					"tls.key": []byte(tlsPairs[1].Key),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID("72x2j56k-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "ns3",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[2].Cert),
					"tls.key": []byte(tlsPairs[2].Key),
				},
			},
		}
		fooCertificate := kongstate.Certificate{
			Certificate: kong.Certificate{
				ID:   kong.String("7428fb98-180b-4702-a91f-61351a33c6e4"),
				Cert: kong.String(tlsPairs[0].Cert),
				Key:  kong.String(tlsPairs[0].Key),
				SNIs: []*string{kong.String("foo.com")},
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
		assert.Equal(3, len(state.Certificates))
		//foo.com with cert should be fixed
		assert.Contains(state.Certificates, fooCertificate)
	})
	t.Run("SNIs slice with same certificate should be ordered by asc", func(t *testing.T) {
		ingresses := []*networkingv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo3",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
						{
							SecretName: "secret",
							Hosts:      []string{"foo3.xxx.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo2",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
						{
							SecretName: "secret",
							Hosts:      []string{"foo2.xxx.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo1",
					Namespace: "ns1",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: networkingv1beta1.IngressSpec{
					TLS: []networkingv1beta1.IngressTLS{
						{
							SecretName: "secret",
							Hosts:      []string{"foo1.xxx.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret",
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
		}
		fooCertificate := kongstate.Certificate{
			Certificate: kong.Certificate{
				ID:   kong.String("7428fb98-180b-4702-a91f-61351a33c6e4"),
				Cert: kong.String(tlsPairs[0].Cert),
				Key:  kong.String(tlsPairs[0].Key),
				SNIs: []*string{
					kong.String("foo1.xxx.com"),
					kong.String("foo2.xxx.com"),
					kong.String("foo3.xxx.com"),
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
		assert.Equal(1, len(state.Certificates))
		assert.Equal(state.Certificates[0], fooCertificate)
	})
}
