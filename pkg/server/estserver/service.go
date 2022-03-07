package estserver

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	keyRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/rand"
	"net/http"

	"github.com/go-kit/log"

	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/utils"
	lamassuestclient "github.com/lamassuiot/lamassu-est/pkg/client"
	lamassuest "github.com/lamassuiot/lamassu-est/pkg/server/api/service"
)

type EstService struct {
	logger           log.Logger
	lamassuEstClient lamassuestclient.LamassuEstClient
}

func NewEstService(lamassuEstClient *lamassuestclient.LamassuEstClient, logger log.Logger) lamassuest.Service {
	return &EstService{
		lamassuEstClient: *lamassuEstClient,
		logger:           logger,
	}
}

func (s *EstService) Health(ctx context.Context) bool {
	return true
}

func (s *EstService) CACerts(ctx context.Context, aps string, r *http.Request) ([]*x509.Certificate, error) {
	certs, err := s.lamassuEstClient.CACerts(ctx)
	if err != nil {
		level.Error(s.logger).Log("err", err, "msg", "Error in client request")
		return nil, err
	}
	level.Info(s.logger).Log("msg", "Certificates sent CACerts method")
	return certs, nil
}

func (s *EstService) Enroll(ctx context.Context, csr *x509.CertificateRequest, aps string, cert *x509.Certificate, r *http.Request) (*x509.Certificate, error) {
	KeyTypes := []string{"rsa", "ec"}
	randomIndex := rand.Intn(len(KeyTypes))
	KeyType := KeyTypes[randomIndex]
	var ecdsaKey *ecdsa.PrivateKey
	var rsaKey *rsa.PrivateKey
	var privKey []byte
	var err error
	if KeyType == "rsa" {
		KeyBits := []int{2048, 3072, 7680}
		Index := rand.Intn(len(KeyBits))
		keybit := KeyBits[Index]
		rsaKey, _ = rsa.GenerateKey(keyRand.Reader, keybit)
		privKey, err = x509.MarshalPKCS8PrivateKey(rsaKey)
		if err != nil {
			return nil, err
		}
	} else {
		KeyBits := []int{224, 256, 384}
		Index := rand.Intn(len(KeyBits))
		keybit := KeyBits[Index]
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
			return nil, err
		}
	}
	csr, err = utils.GenerateCSR(privKey)
	if err != nil {
		return nil, err
	}

	dataCert, err := s.lamassuEstClient.Enroll(ctx, aps, csr)
	if err != nil {
		level.Error(s.logger).Log("err", err, "msg", "Error in client request")
		return &x509.Certificate{}, err
	}
	level.Info(s.logger).Log("msg", "Certificate sent ENROLL method")
	return dataCert, nil
}

func (s *EstService) Reenroll(ctx context.Context, cert *x509.Certificate, csr *x509.CertificateRequest, aps string, r *http.Request) (*x509.Certificate, error) {
	dataCert, err := s.lamassuEstClient.Reenroll(ctx, csr)
	if err != nil {
		level.Error(s.logger).Log("err", err, "msg", "Error in client request")
		return &x509.Certificate{}, err
	}
	level.Info(s.logger).Log("msg", "Certificate sent REENROLL method")
	return dataCert, nil
}
func (s *EstService) ServerKeyGen(ctx context.Context, csr *x509.CertificateRequest, aps string, r *http.Request) (*x509.Certificate, []byte, error) {
	dataCert, key, err := s.lamassuEstClient.ServerKeyGen(ctx, aps, csr)
	if err != nil {
		level.Error(s.logger).Log("err", err, "msg", "Error in client request")
		return &x509.Certificate{}, nil, err
	}
	level.Info(s.logger).Log("msg", "Certificate and key sent ServerKeyGen method")
	return dataCert, key, nil
}
