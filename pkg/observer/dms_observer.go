package observer

import (
	"errors"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-default-dms/pkg/config"
	"github.com/lamassuiot/lamassu-default-dms/pkg/device/store"
	"github.com/lamassuiot/lamassuiot/pkg/dms-enroller/common/dto"
)

type EnrolledDeviceData struct {
	DeviceAlias              string
	DeviceID                 string
	EnrolledCertSerialNumber string
	EnrolledCA               string
}

type DeviceState struct {
	// internal state
	Devices []EnrolledDeviceData
	Config  config.Config
	Aps     string
	//Auto       string
	DeviceFile store.File
	DmsPrivKey string
	DmsId      string
	DmsFile    store.File
	Dms        dto.DMS
	Bits       []string
	Stop       bool

	observers []Observer
}

func (s *DeviceState) Attach(o Observer) (bool, error) {

	for _, observer := range s.observers {
		if observer == o {
			return false, errors.New("Observer already exists")
		}
	}
	s.observers = append(s.observers, o)
	return true, nil
}

func (s *DeviceState) Detach(o Observer) (bool, error) {

	for i, observer := range s.observers {
		if observer == o {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return true, nil
		}
	}
	return false, errors.New("Observer not found")
}

func (s *DeviceState) Notify(logger log.Logger) (bool, error) {
	level.Info(logger).Log("msg", "Obserer notify... "+strconv.Itoa(len(s.observers)))
	for _, observer := range s.observers {
		observer.Update(s)
	}
	return true, nil
}

func (s *DeviceState) AddDevice(device EnrolledDeviceData, logger log.Logger) {
	s.Devices = append(s.Devices, device)
	s.Notify(logger)
}
func (s *DeviceState) EditDMS(dms dto.DMS, logger log.Logger) {
	s.Dms = dms
	s.Notify(logger)
}
func (s *DeviceState) AddKeyBits(bit []string, logger log.Logger) {
	s.Bits = bit
	s.Notify(logger)
}
func (s *DeviceState) AddStop(stop bool, logger log.Logger) {
	s.Stop = stop
	s.Notify(logger)
}
