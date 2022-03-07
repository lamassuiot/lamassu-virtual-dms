package service

import (
	"context"
	"errors"

	"github.com/go-kit/log"

	lamassuestclient "github.com/lamassuiot/lamassu-est/pkg/client"

	filestore "github.com/lamassuiot/lamassu-default-dms/pkg/server/models/device/store"
)

type Service interface {
	Health(ctx context.Context) bool
}

type dmsService struct {
	fileStore        filestore.File
	homePath         string
	logger           log.Logger
	lamassuEstClient lamassuestclient.LamassuEstClient
}

var (
	// Client errors
	ErrInvalidID = errors.New("invalid CSR ID, does not exist") //404

	//Server errors
	ErrInvalidOperation = errors.New("invalid operation")
	ErrGetCSR           = errors.New("unable to get CSR")
	ErrGetCert          = errors.New("unable to get certificate")
	ErrInsertCert       = errors.New("unable to insert certificate")
	ErrInsertCSR        = errors.New("unable to insert CSR")
	ErrInsertKey        = errors.New("unable to insert Key")
	ErrResponseEncode   = errors.New("error encoding response")
)

func NewDMSService(fileStore filestore.File, homePath string, lamassuEstClient *lamassuestclient.LamassuEstClient, logger log.Logger) Service {
	return &dmsService{
		fileStore:        fileStore,
		homePath:         homePath,
		lamassuEstClient: *lamassuEstClient,
		logger:           logger,
	}
}

func (s *dmsService) Health(ctx context.Context) bool {
	return true
}
