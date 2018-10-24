/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	oidExtensionSubjectAltName = asn1.ObjectIdentifier{2, 5, 29, 17}
)

// ParseX509Certificate returns a *X509.Certifcate from a cert and key pair.
func ParseX509Certificate(cert, key []byte) (*x509.Certificate, error) {
	pemCerts := append(cert, '\n')

	pemCerts = append(pemCerts, key...)

	pemBlock, _ := pem.Decode(pemCerts)
	if pemBlock == nil {
		return nil, fmt.Errorf("no valid PEM formatted block found")
	}

	// If the file does not start with 'BEGIN CERTIFICATE' it's invalid and must not be used.
	if pemBlock.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("certificate contains invalid data, and must be created with 'kubectl create secret tls'")
	}

	pemCert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	//Ensure that certificate and private key have a matching public key
	if _, err := tls.X509KeyPair(cert, key); err != nil {
		return nil, err
	}
	return pemCert, nil
}

// ParseCommonNamesFromCert returns common names in cert and SAN extension
func ParseCommonNamesFromCert(cert *x509.Certificate) []string {
	cn := sets.NewString(cert.Subject.CommonName)
	for _, dns := range cert.DNSNames {
		if !cn.Has(dns) {
			cn.Insert(dns)
		}
	}

	if len(cert.Extensions) > 0 {
		glog.V(3).Info("parsing ssl certificate extensions")
		for _, ext := range getExtension(cert, oidExtensionSubjectAltName) {
			dns, _, _, err := parseSANExtension(ext.Value)
			if err != nil {
				glog.Warningf("unexpected error parsing certificate extensions: %v", err)
				continue
			}

			for _, dns := range dns {
				if !cn.Has(dns) {
					cn.Insert(dns)
				}
			}
		}
	}

	return cn.List()
}

func getExtension(c *x509.Certificate, id asn1.ObjectIdentifier) []pkix.Extension {
	var exts []pkix.Extension
	for _, ext := range c.Extensions {
		if ext.Id.Equal(id) {
			exts = append(exts, ext)
		}
	}
	return exts
}

func parseSANExtension(value []byte) (dnsNames, emailAddresses []string, ipAddresses []net.IP, err error) {
	// RFC 5280, 4.2.1.6

	// SubjectAltName ::= GeneralNames
	//
	// GeneralNames ::= SEQUENCE SIZE (1..MAX) OF GeneralName
	//
	// GeneralName ::= CHOICE {
	//      otherName                       [0]     OtherName,
	//      rfc822Name                      [1]     IA5String,
	//      dNSName                         [2]     IA5String,
	//      x400Address                     [3]     ORAddress,
	//      directoryName                   [4]     Name,
	//      ediPartyName                    [5]     EDIPartyName,
	//      uniformResourceIdentifier       [6]     IA5String,
	//      iPAddress                       [7]     OCTET STRING,
	//      registeredID                    [8]     OBJECT IDENTIFIER }
	var seq asn1.RawValue
	var rest []byte
	if rest, err = asn1.Unmarshal(value, &seq); err != nil {
		return
	} else if len(rest) != 0 {
		err = errors.New("x509: trailing data after X.509 extension")
		return
	}
	if !seq.IsCompound || seq.Tag != 16 || seq.Class != 0 {
		err = asn1.StructuralError{Msg: "bad SAN sequence"}
		return
	}

	rest = seq.Bytes
	for len(rest) > 0 {
		var v asn1.RawValue
		rest, err = asn1.Unmarshal(rest, &v)
		if err != nil {
			return
		}
		switch v.Tag {
		case 1:
			emailAddresses = append(emailAddresses, string(v.Bytes))
		case 2:
			dnsNames = append(dnsNames, string(v.Bytes))
		case 7:
			switch len(v.Bytes) {
			case net.IPv4len, net.IPv6len:
				ipAddresses = append(ipAddresses, v.Bytes)
			default:
				err = errors.New("x509: certificate contained IP address of length " + strconv.Itoa(len(v.Bytes)))
				return
			}
		}
	}

	return
}
