package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/jakehl/goid"
)

const (
	PublicKeyHeader = "-----BEGIN PUBLIC KEY-----"
	PublicKeyFooter = "-----END PUBLIC KEY-----"
)

func CreateCAPool(CAPath string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(CAPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func ParseKeycloakPublicKey(data []byte) (*rsa.PublicKey, error) {
	pubPem, _ := pem.Decode(data)
	parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		return nil, errors.New("Unable to parse public key")
	}
	pubKey := parsedKey.(*rsa.PublicKey)
	return pubKey, nil
}
func VerifyPeerCertificate(certClient *x509.Certificate, certCA *x509.Certificate) ([][]*x509.Certificate, error) {
	clientCAs := x509.NewCertPool()
	clientCAs.AddCert(certCA)

	opts := x509.VerifyOptions{
		Roots:     clientCAs,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	candidateCa, err := certClient.Verify(opts)
	if err != nil {
		return nil, err
	}
	return candidateCa, err
}
func GenerateCSR(key interface{}) (*x509.CertificateRequest, error) {
	CommonName := goid.NewV4UUID()
	subj := pkix.Name{
		Country:            []string{"ES"},
		Province:           []string{"Gipuzkoa"},
		Organization:       []string{"IKERLAN"},
		OrganizationalUnit: []string{"ZPD"},
		Locality:           []string{"Arrasate"},
		CommonName:         CommonName.String(),
	}
	rawSubject := subj.ToRDNSequence()
	asn1Subj, _ := asn1.Marshal(rawSubject)
	template := x509.CertificateRequest{
		RawSubject: asn1Subj,
		//EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.ECDSAWithSHA512,
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate request: %v", err)
	}

	csrNew, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate request: %v", err)
	}
	return csrNew, nil
}
