package service

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"net/url"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jakehl/goid"
	"github.com/lamassuiot/lamassu-default-dms/pkg/device/store"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/lamassuiot/lamassu-default-dms/pkg/utils"
	lamassuDevClient "github.com/lamassuiot/lamassuiot/pkg/device-manager/client"
	devDTO "github.com/lamassuiot/lamassuiot/pkg/device-manager/common/dto"
	lamassuDMSClient "github.com/lamassuiot/lamassuiot/pkg/dms-enroller/client"
	"github.com/lamassuiot/lamassuiot/pkg/dms-enroller/common/dto"
	estclient "github.com/lamassuiot/lamassuiot/pkg/est/client"
	"github.com/lamassuiot/lamassuiot/pkg/utils/client"
)

func Enroll(lamassuEstClient estclient.LamassuEstClient, data *observer.DeviceState, file store.File, aps string, dmsname string, logger log.Logger) (deviceAlias string, deviceId string, certSN string, CAname string, err error) {
	level.Info(logger).Log("msg", "Enroll New Device... ")
	var ctx context.Context
	var device devDTO.Device

	devclient, _ := lamassuDevClient.NewLamassuDeviceManagerClient(client.ClientConfiguration{
		URL: &url.URL{
			Scheme: "https",
			Host:   data.Config.Domain,
			Path:   "/api/devmanager/",
		},
		AuthMethod: client.JWT,
		AuthMethodConfig: &client.JWTConfig{
			Username: "enroller",
			Password: "enroller",
			URL: &url.URL{
				Scheme: "https",
				Host:   data.Config.Auth.Endpoint,
			},
			CACertificate: data.Config.DevManager.DevCrt,
		},
		CACertificate: data.Config.DevManager.DevCrt,
	})
	deviceregister := []string{"yes", "no"}
	randomIndex := rand.Intn(len(deviceregister))
	mode := deviceregister[randomIndex]

	if mode == "yes" {
		level.Info(logger).Log("msg", "Register Device... ")
		deviceID := goid.NewV4UUID().String()
		device, err = devclient.CreateDevice(ctx, utils.AliasName(), deviceID, data.Dms.Id, "", utils.Tags(), utils.IconName(), utils.IconColor())
		if err != nil {
			level.Error(logger).Log("err", err)
		}
	}

	privateKeyBytes, csr, err := utils.GenrateRandKey(logger, device)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", "", "", "", err
	}
	if aps == "" {
		dmsClient, _ := lamassuDMSClient.NewLamassuEnrollerClient(client.ClientConfiguration{
			URL: &url.URL{
				Scheme: "https",
				Host:   data.Config.Domain,
				Path:   "/api/dmsenroller/",
			},
			AuthMethod: client.JWT,
			AuthMethodConfig: &client.JWTConfig{
				Username: "enroller",
				Password: "enroller",
				URL: &url.URL{
					Scheme: "https",
					Host:   data.Config.Auth.Endpoint,
				},
				CACertificate: data.Config.DevManager.DevCrt,
			},
			CACertificate: data.Config.DevManager.DevCrt,
		})
		dms, err := dmsClient.GetDMSbyID(ctx, data.Dms.Id)
		if err != nil {
			level.Error(logger).Log("err", err)
			return "", "", "", "", err
		}
		index := rand.Intn(len(dms.AuthorizedCAs))
		aps = dms.AuthorizedCAs[index]
	}
	level.Info(logger).Log("msg", aps)
	devCert, err := lamassuEstClient.Enroll(ctx, aps, csr)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", "", "", "", err
	} else if devCert != nil {
		file.InsertCSR(csr.Subject.CommonName, csr.Raw, "device", dmsname)
		file.InsertKEY(csr.Subject.CommonName, privateKeyBytes, "device", dmsname)
		file.InsertCERT(devCert.Subject.CommonName, devCert.Raw, "device", dmsname, utils.InsertNth(utils.ToHexInt(devCert.SerialNumber), 2))

		level.Info(logger).Log("msg", "Certificate content: "+devCert.Subject.String()+" Issuer: "+devCert.Issuer.String())
	}
	return device.Alias, csr.Subject.CommonName, utils.InsertNth(utils.ToHexInt(devCert.SerialNumber), 2), aps, nil
}

func RegisterDMS(data *observer.DeviceState, file store.File, name string, subject dto.Subject, PrivateKeyMetadata dto.PrivateKeyMetadata, dmsClient lamassuDMSClient.LamassuEnrollerClient, logger log.Logger) (string, error) {
	key, dms, err := dmsClient.CreateDMSForm(context.Background(), subject, PrivateKeyMetadata, name)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	privkey, _ := base64.StdEncoding.DecodeString(key)

	err = ioutil.WriteFile(data.Config.Dms.DmsStore+"/dms-"+dms.Id+".key", privkey, 0644)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	level.Info(logger).Log("msg", "KEY with ID "+dms.Id+" inserted in file system")
	return dms.Id, nil
}
