package mtls

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	stdhttp "net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/configs"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/utils"
)

type contextKey string

const (
	PeerCertificatesContextKey contextKey = "PeerCertificatesContextKey"
)

var (
	ErrPeerCertificatesContextMissing = errors.New("token up for parsing was not passed through the context")
)

func HTTPToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			return context.WithValue(ctx, PeerCertificatesContextKey, r.TLS.PeerCertificates[0])
		} else {
			return ctx
		}
	}
}
func NewParser(cfg configs.Config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			peerCert, ok := ctx.Value(PeerCertificatesContextKey).(*x509.Certificate)
			if !ok {
				return nil, ErrPeerCertificatesContextMissing
			}
			_ = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: peerCert.Raw})
			var certContent []byte
			if cfg.MutualTLSEnabled {
				certContent, err = ioutil.ReadFile(cfg.MutualTLSClientCA)
				if err != nil {
					return nil, err
				}

			} else {
				certContent, err = ioutil.ReadFile(cfg.BootstrapCert)
				if err != nil {
					return nil, err
				}
			}
			cpb, _ := pem.Decode(certContent)
			crt, err := x509.ParseCertificate(cpb.Bytes)
			if err != nil {
				return nil, err
			}
			_, err = utils.VerifyPeerCertificate(peerCert, crt)
			if err != nil {
				return nil, err
			}
			return next(ctx, request)
		}
	}
}
