package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-default-dms/pkg/config"
	"github.com/lamassuiot/lamassu-default-dms/pkg/device/store"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/lamassuiot/lamassu-default-dms/pkg/utils"
	"github.com/lamassuiot/lamassuiot/pkg/est/client"
)

func Enroll(lamassuEstClient client.LamassuEstClient, data *observer.DeviceState, file store.File, aps string, token config.Token, dmsname string, logger log.Logger) (deviceAlias string, deviceId string, certSN string, CAname string, err error) {
	level.Info(logger).Log("msg", "Enroll New Device... ")
	var ctx context.Context
	var device config.Device
	var dmss []config.DMS

	deviceregister := []string{"yes", "no"}
	randomIndex := rand.Intn(len(deviceregister))
	mode := deviceregister[randomIndex]

	if mode == "yes" {
		level.Info(logger).Log("msg", "Register Device... ")

		body := utils.CreateRequestBody(data.Dms.Id)

		req, err := utils.NewRequest(http.MethodPost, "/v1/devices", data.Config.DevManager.DevAddr, "application/json", "application/json", "", "", "", "", body, token.AccesToken)
		if err != nil {
			level.Error(logger).Log("err", err)
		}
		_, resp, err := utils.Do(req, data)
		if err != nil {
			level.Error(logger).Log("err", err)
			return "", "", "", "", err
		}
		jsonString, err := json.Marshal(resp)
		json.Unmarshal(jsonString, &device)
	}

	privateKeyBytes, csr, err := utils.GenrateRandKey(logger, device)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", "", "", "", err
	}
	if aps == "" {
		req, err := utils.NewRequest(http.MethodGet, "/v1/", data.Config.Dms.Endpoint, "", "application/json", "", "", "", "", nil, token.AccesToken)
		if err != nil {
			level.Error(logger).Log("err", err)
			return "", "", "", "", err
		}
		_, resp, err := utils.Do(req, data)
		if err != nil {
			level.Error(logger).Log("err", err)
			return "", "", "", "", err
		}
		jsonString, _ := json.Marshal(resp)
		json.Unmarshal(jsonString, &dmss)
		var CAs []string
		for _, dms := range dmss {
			if data.Dms.Id == dms.Id {
				CAs = dms.AuthorizedCAs
			}
		}
		index := rand.Intn(len(CAs))
		aps = CAs[index]
	}

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
func RequestToken(data *observer.DeviceState, logger log.Logger) (config.Token, error) {
	var token config.Token
	level.Info(logger).Log("msg", "Request Access Token... ")
	req, err := utils.NewRequest(http.MethodPost, "/auth/realms/lamassu/protocol/openid-connect/token", data.Config.Auth.Endpoint, "application/x-www-form-urlencoded", "", "password", "frontend", data.Config.Auth.Username, data.Config.Auth.Password, nil, "")
	if err != nil {
		level.Error(logger).Log("err", err)
		return config.Token{}, err
	}

	_, resp, err := utils.Do(req, data)
	if err != nil {
		level.Error(logger).Log("err", err)
		return config.Token{}, err
	}
	jsonString, _ := json.Marshal(resp)
	json.Unmarshal(jsonString, &token)
	return token, nil
}

func RegisterDMS(data *observer.DeviceState, file store.File, name string, subject config.Subject, PrivateKeyMetadata config.PrivateKeyMetadata, logger log.Logger) (string, error) {
	var dms config.DmsCreationResponse

	token, err := RequestToken(data, logger)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	body := utils.CreateDmsRequestBody(name, subject, PrivateKeyMetadata)

	req, err := utils.NewRequest(http.MethodPost, "/v1/"+name+"/form", data.Config.Dms.Endpoint, "application/json", "application/json", "", "", "", "", body, token.AccesToken)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	_, resp, err := utils.Do(req, data)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	jsonString, err := json.Marshal(resp)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	json.Unmarshal(jsonString, &dms)
	key, _ := base64.StdEncoding.DecodeString(dms.PrivKey)

	err = ioutil.WriteFile(data.Config.Dms.DmsStore+"/dms-"+dms.Dms.Id+".key", key, 0644)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	level.Info(logger).Log("msg", "KEY with ID "+dms.Dms.Id+" inserted in file system")
	return dms.Dms.Id, nil
}
