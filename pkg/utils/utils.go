package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	keyRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/jakehl/goid"
	devDTO "github.com/lamassuiot/lamassuiot/pkg/device-manager/common/dto"
)

func GenrateRandKey(logger log.Logger, device devDTO.Device) ([]byte, *x509.CertificateRequest, error) {
	var ecdsaKey *ecdsa.PrivateKey
	var rsaKey *rsa.PrivateKey
	var privKey []byte
	var err error
	var csr *x509.CertificateRequest
	var KeyType string
	var keybit int
	KeyType, keybit = RandomKeyTypeBits()
	if KeyType == "rsa" {
		rsaKey, _ = rsa.GenerateKey(keyRand.Reader, keybit)
		privKey, err = x509.MarshalPKCS8PrivateKey(rsaKey)
		if err != nil {
			return nil, nil, err
		}
		csr, err = GenerateCSR(rsaKey, KeyType, device)
		if err != nil {
			level.Error(logger).Log("err", err)
			return nil, nil, err
		}
		level.Info(logger).Log("msg", "Generated "+KeyType+" "+strconv.Itoa(keybit)+" bits private key")
		return privKey, csr, nil
	} else {
		switch keybit {
		case 224:
			ecdsaKey, _ = ecdsa.GenerateKey(elliptic.P224(), keyRand.Reader)
		case 256:
			ecdsaKey, _ = ecdsa.GenerateKey(elliptic.P256(), keyRand.Reader)
		case 384:
			ecdsaKey, _ = ecdsa.GenerateKey(elliptic.P384(), keyRand.Reader)
		}
		privKey, err = x509.MarshalPKCS8PrivateKey(ecdsaKey)
		if err != nil {
			return nil, nil, err
		}
		csr, err = GenerateCSR(ecdsaKey, KeyType, device)
		if err != nil {
			level.Error(logger).Log("err", err)
			return nil, nil, err
		}
		level.Info(logger).Log("msg", "Generated "+KeyType+" "+strconv.Itoa(keybit)+" bits private key")
		return privKey, csr, nil
	}
}

func GenerateCSR(key interface{}, Keytype string, device devDTO.Device) (*x509.CertificateRequest, error) {

	var subj pkix.Name
	if device.Id == "" {
		subj = pkix.Name{
			CommonName: goid.NewV4UUID().String(),
		}
	} else {
		subj = pkix.Name{
			CommonName: device.Id,
		}
	}

	rawSubject := subj.ToRDNSequence()
	asn1Subj, _ := asn1.Marshal(rawSubject)
	var template x509.CertificateRequest
	if Keytype == "rsa" {
		template = x509.CertificateRequest{
			RawSubject:         asn1Subj,
			SignatureAlgorithm: x509.SHA512WithRSA,
		}
	} else {
		template = x509.CertificateRequest{
			RawSubject:         asn1Subj,
			SignatureAlgorithm: x509.ECDSAWithSHA512,
		}
	}

	csrBytes, err := x509.CreateCertificateRequest(keyRand.Reader, &template, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate request: %v", err)
	}

	csrNew, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate request: %v", err)
	}
	return csrNew, nil
}
func ToHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
}

func InsertNth(s string, n int) string {
	if len(s)%2 != 0 {
		s = "0" + s
	}
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}

func AliasName() string {
	Alias := []string{"Smart TV", "Smart thermostats", "E-reader", "Smart lock", "Mobile robot", "Smart Light Switch", "Camera", "Security sensor"}
	randomIndex := rand.Intn(len(Alias))
	return Alias[randomIndex]

}
func IconName() string {
	Names := []string{"Cg/CgSmartHomeBoiler", "Cg/CgSmartHomeCooker", "Cg/CgSmartHomeHeat", "Cg/CgSmartHomeLight", "Cg/CgSmartHomeRefrigerator", "Cg/CgSmartHomeWashMachine", "Cg/CgSmartphoneChip", "Cg/CgSmartphoneRam", "Cg/CgSmartphoneShake", "Cg/CgSmartphone"}
	randomIndex := rand.Intn(len(Names))
	return Names[randomIndex]
}
func IconColor() string {
	Colors := []string{"#FF8A65", "#B968C7", "#DCE775", "#697689", "#2F657B", "#66BB6A", "#02B6DC"}
	randomIndex := rand.Intn(len(Colors))
	return Colors[randomIndex]
}
func Tags() []string {
	tags := []string{"ES North Fleet", "Device v2", "TPM", "5G", "Linux OS", "Battery powered", "Sensor", "IPv6"}
	randomIndex := rand.Intn(len(tags))
	var tag []string
	tag = append(tag, tags[randomIndex])
	return tag
}
func RandomKeyTypeBits() (string, int) {
	var keybit int
	KeyTypes := []string{"rsa", "ec"}
	randomIndex := rand.Intn(len(KeyTypes))
	KeyType := KeyTypes[randomIndex]
	if KeyType == "rsa" {
		KeyBits := []int{2048, 3072, 4096}
		Index := rand.Intn(len(KeyBits))
		keybit = KeyBits[Index]
		return KeyType, keybit
	} else {
		KeyBits := []int{224, 256, 384}
		Index := rand.Intn(len(KeyBits))
		keybit = KeyBits[Index]
		return KeyType, keybit
	}

}
func DecodeB64(message string) (string, error) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	_, err := base64.StdEncoding.Decode(base64Text, []byte(message))
	return string(base64Text), err
}
func ReadCertPool(path string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}
func ReadCert(path string) (*x509.Certificate, error) {
	certContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cpb, _ := pem.Decode(certContent)

	crt, err := x509.ParseCertificate(cpb.Bytes)
	if err != nil {
		return nil, err
	}
	return crt, nil
}

func ReadKey(path string) ([]byte, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return key, nil
}
func CheckIfNull(field []string) string {
	var result = ""
	if field != nil {
		result = field[0]
	}
	return result
}
