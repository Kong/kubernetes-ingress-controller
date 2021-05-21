//+build integration_tests

package integration

const (
	admissionWebhookCert = `-----BEGIN CERTIFICATE-----
MIIDIDCCAgigAwIBAgIUCRzHAvZF4+Roul3QHLsqCLs9ZLMwDQYJKoZIhvcNAQEL
BQAwMzELMAkGA1UEBhMCUEwxJDAiBgNVBAMMG3ZhbGlkYXRpb25zLmtvbmctc3lz
dGVtLnN2YzAgFw0yMTA1MjExOTI5MjlaGA8yMTIxMDQyNzE5MjkyOVowMzELMAkG
A1UEBhMCUEwxJDAiBgNVBAMMG3ZhbGlkYXRpb25zLmtvbmctc3lzdGVtLnN2YzCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOCAjc9P/NgSaCi7+xjLJDSA
7+czwc8fxr90MvvRbsj8Lrl2wkyKJMBIP7dRlLcrBx0/1Zu8d9d8pFRRZBttXcfA
SAGNM+J1t2dQP2RaUilZJN7Tw3UYF8iy0JD81QW5ykDXlGdwbieKRlpeb/L2uVDu
PbHgLTk4FCkZS2TIh+7xlqisINuzPuQKe3jT5sPsljbiLJRxXSa+LnwexMcSlxTQ
TUtHdeQbcGjfQXuzse7sehDynazsUiGcgs5hKbCjMdgrgHkMJRSqV1m/I4ZrIfon
gIaoUuV9Wx1y95gpPKZ/3fs3F/nW7vxtf//C3wWWnA0dkIJoWYnx10gvMk/wo7EC
AwEAAaMqMCgwJgYDVR0RBB8wHYIbdmFsaWRhdGlvbnMua29uZy1zeXN0ZW0uc3Zj
MA0GCSqGSIb3DQEBCwUAA4IBAQBzTrDxES0gsC1GuLlyb0eeri32An4ujW62wtCE
rYmzPpQyhX3LPGr6KZM54PRGh5PGgdf68Gls0qeTgXq8uI49i2ttwUsBA9BWi46B
RqhUw3m0Z6p0N0y+JgIMZXzqct6gBeJ4Xg134FHeNRLLSQb3udO/o9gZv9bgmGNr
mTjWhBeQdYUn8vcUx1XpmipClkIwvI+8kxcmC8LizuNXFmRy57nwfYq3XRL0T09j
L9rSmC3ehyVLndD8x6hXkOkMGofLqhIzRfMrCFem9g4Xi1yo1HkwH2hZvw+3dwJJ
UuCPFug0koLCC1XyYw0YT9vjZemy/z0JM+FCNo5I8wfTtgvA
-----END CERTIFICATE-----`
	admissionWebhookCertBase64 = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJRENDQWdpZ0F3SUJBZ0lVQ1J6SEF2WkY0K1JvdWwzUUhMc3FDTHM5WkxNd0RRWUpLb1pJaHZjTkFRRUwKQlFBd016RUxNQWtHQTFVRUJoTUNVRXd4SkRBaUJnTlZCQU1NRzNaaGJHbGtZWFJwYjI1ekxtdHZibWN0YzNsegpkR1Z0TG5OMll6QWdGdzB5TVRBMU1qRXhPVEk1TWpsYUdBOHlNVEl4TURReU56RTVNamt5T1Zvd016RUxNQWtHCkExVUVCaE1DVUV3eEpEQWlCZ05WQkFNTUczWmhiR2xrWVhScGIyNXpMbXR2Ym1jdGMzbHpkR1Z0TG5OMll6Q0MKQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFPQ0FqYzlQL05nU2FDaTcreGpMSkRTQQo3K2N6d2M4ZnhyOTBNdnZSYnNqOExybDJ3a3lLSk1CSVA3ZFJsTGNyQngwLzFadThkOWQ4cEZSUlpCdHRYY2ZBClNBR05NK0oxdDJkUVAyUmFVaWxaSk43VHczVVlGOGl5MEpEODFRVzV5a0RYbEdkd2JpZUtSbHBlYi9MMnVWRHUKUGJIZ0xUazRGQ2taUzJUSWgrN3hscWlzSU51elB1UUtlM2pUNXNQc2xqYmlMSlJ4WFNhK0xud2V4TWNTbHhUUQpUVXRIZGVRYmNHamZRWHV6c2U3c2VoRHluYXpzVWlHY2dzNWhLYkNqTWRncmdIa01KUlNxVjFtL0k0WnJJZm9uCmdJYW9VdVY5V3gxeTk1Z3BQS1ovM2ZzM0Yvblc3dnh0Zi8vQzN3V1duQTBka0lKb1dZbngxMGd2TWsvd283RUMKQXdFQUFhTXFNQ2d3SmdZRFZSMFJCQjh3SFlJYmRtRnNhV1JoZEdsdmJuTXVhMjl1WnkxemVYTjBaVzB1YzNaagpNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUJ6VHJEeEVTMGdzQzFHdUxseWIwZWVyaTMyQW40dWpXNjJ3dENFCnJZbXpQcFF5aFgzTFBHcjZLWk01NFBSR2g1UEdnZGY2OEdsczBxZVRnWHE4dUk0OWkydHR3VXNCQTlCV2k0NkIKUnFoVXczbTBaNnAwTjB5K0pnSU1aWHpxY3Q2Z0JlSjRYZzEzNEZIZU5STExTUWIzdWRPL285Z1p2OWJnbUdOcgptVGpXaEJlUWRZVW44dmNVeDFYcG1pcENsa0l3dkkrOGt4Y21DOExpenVOWEZtUnk1N253ZllxM1hSTDBUMDlqCkw5clNtQzNlaHlWTG5kRDh4NmhYa09rTUdvZkxxaEl6UmZNckNGZW05ZzRYaTF5bzFIa3dIMmhadncrM2R3SkoKVXVDUEZ1ZzBrb0xDQzFYeVl3MFlUOXZqWmVteS96MEpNK0ZDTm81STh3ZlR0Z3ZBCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
	admissionWebhookKey        = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4ICNz0/82BJoKLv7GMskNIDv5zPBzx/Gv3Qy+9FuyPwuuXbC
TIokwEg/t1GUtysHHT/Vm7x313ykVFFkG21dx8BIAY0z4nW3Z1A/ZFpSKVkk3tPD
dRgXyLLQkPzVBbnKQNeUZ3BuJ4pGWl5v8va5UO49seAtOTgUKRlLZMiH7vGWqKwg
27M+5Ap7eNPmw+yWNuIslHFdJr4ufB7ExxKXFNBNS0d15BtwaN9Be7Ox7ux6EPKd
rOxSIZyCzmEpsKMx2CuAeQwlFKpXWb8jhmsh+ieAhqhS5X1bHXL3mCk8pn/d+zcX
+dbu/G1//8LfBZacDR2QgmhZifHXSC8yT/CjsQIDAQABAoIBAB7wmJqho28D2mcC
wTBBjtPNkUKD14n8DyADm6Mo0ePRHX9h5pU11KrLSjyxeZVk0K4vRfkYmEuSWfNk
5C8De5Ez5riQBT6IiqYqYRIrgHdCWdp7xMw2bdCzFBtnPNR1LnKRQ1qeHBBG9jsu
GK+bYR7ONqJ1CsZ//AdN/R3+pP1A/CMPwvq43j84KSoqPJrgMHJXullAbUlKlY29
jtrG4GUvt3BiGXC/kKFYI41xAuVI5SdEXu+ycZcU8Ifx2Gpg8Fjp3f4l7GRwvY5K
OE0zP9GSahB6rJF7GZCaKWxOBBix7GjzFPogfcuHCa9Fjwm4IQmC6w3/rYf4+9+v
8M25uKECgYEA+2ki69Hq4jEHUW/B24j+jYKBMlRa3Cd1ci4CP36R1K0080UCz2/4
FmSUyrecIL8C3UZtHQg3EeQrPSJ8h+FrYB48XfLcwDXDWlTIzxPmiWUlG8YK3KLA
LN7kvMuRmOeASBFJIVQuGo9kH+xnFmR7vbYuZ8HUJiMLoWUl8vNuqusCgYEA5Jmr
9JItPPOuWVww6qj+fhbIeBGGjzvUciYdhi5DrIWI2pap9Jn6C08YqsKaKkqPwCjO
VwhwSFwYZbs+pC8YOUxZ9vzkCo623GUcPKe0EjW2GkvufsJohE3ecQkmeve96Qyi
99wmrYmyYfABz+ff1edC7VFk79cuqUiC/3v0TNMCgYAhkAaObsapjZwJfh7mHOLG
p25x9prumwHtzUCVk2MKflj8RPE8GhmHe8P1UA+yu205dwZoAsm/RLOVBL6VMT2x
Zjfu3tYjfsnmjD0GkASNwQf0LjsS+1Mmalck8RQt0nHorQ4TOfaxqwTV0ixs69st
F14YkeKteK47zJIFXgQfIwKBgB8kUCihQUhsafQCeyd8ni7PK8AvowUgQXDLgHon
E1ENX/dnTv/jegzQWavpltbsEWk8Jd/1ZlZ1NV2mhIIZaFNl81uSV/6YMpETtSUO
M5nHd2ddsL/T/CkJ8qOze2qFFXoKHqlldF9vwr1U1OpdzEB3oMZzsCx8Q/8LwczM
NhvBAoGBAPtiu+/HDixPpuFeUsoGqN46rjjyM6sG6AFjjARie+d+jZi5Q7rEfrtY
7Oh25R61i2QlO0fTi01jK6VH6O1G+D2AciGK7cWgLCAJv4ppAA3ysbRbfL1mXS3B
O5PMpLb3dETohTyDk7+r1UgPmKSFDs4OnO5mS1dvQGM1f3OpcgJQ
-----END RSA PRIVATE KEY-----`
)
