package deviceview

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/lamassuiot/lamassu-default-dms/pkg/service"
	"github.com/lamassuiot/lamassu-default-dms/pkg/utils"
	"github.com/lamassuiot/lamassuiot/pkg/est/client"
	"github.com/rivo/tview"
)

func GetEnrollItem(logger log.Logger, data *observer.DeviceState, app *tview.Application) tview.Primitive {
	devManagerEndpoint := data.Config.DevManager.DevAddr
	var flex *tview.Flex
	statusTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText(" ")

	form := tview.NewForm().
		AddInputField("Device Manager endpoint", devManagerEndpoint, 50, nil, func(text string) {
			devManagerEndpoint = text
		}).
		AddInputField("DMS ID", data.Dms.Id, 50, nil, func(text string) {
			data.Dms.Id = text
		}).
		AddButton("Auto-Enroll", func() {
			statusTextView.SetText("AUTO-ENROLLING...")
			app.ForceDraw()
			dmsCrtFile := data.Config.Dms.DmsStore + "/dms-" + data.Dms.Id + ".crt"
			dmsKeyFile := data.Config.Dms.DmsStore + "/dms-" + data.Dms.Id + ".key"

			serverCert, err := utils.ReadCertPool(data.Config.DevManager.DevCrt)
			if err != nil {
				level.Error(logger).Log("err", err)
			}
			dmsCert, err := utils.ReadCert(dmsCrtFile)
			if err != nil {
				level.Error(logger).Log("err", err)
			}
			dmsKey, err := utils.ReadKey(dmsKeyFile)
			if err != nil {
				level.Error(logger).Log("err", err)
			}
			lamassuEstClient, err := client.NewLamassuEstClient(devManagerEndpoint, serverCert, dmsCert, dmsKey, logger)
			if err != nil {
				level.Error(logger).Log("err", err)
			}
			data.AddStop(false, logger)
			go func() {

				for {
					if data.Stop {
						level.Info(logger).Log("msg", "STOP")
						app.Unlock()
						statusTextView.SetText("AUTO-ENROLL STOPPED")
						app.ForceDraw()
						break
					}
					token, err := service.RequestToken(data, logger)
					if err != nil {
						level.Error(logger).Log("err", err)

					}
					certContent, err := ioutil.ReadFile(dmsCrtFile)
					cpb, _ := pem.Decode(certContent)

					crt, err := x509.ParseCertificate(cpb.Bytes)

					alias, id, sn, ca, err := service.Enroll(lamassuEstClient, data, data.DeviceFile, data.Aps, token, crt.Subject.CommonName, logger)
					if err != nil {
						level.Error(logger).Log("err", err)
					} else {
						enrolledInfo := observer.EnrolledDeviceData{
							DeviceAlias:              alias,
							DeviceID:                 id,
							EnrolledCertSerialNumber: sn,
							EnrolledCA:               ca,
						}
						data.AddDevice(enrolledInfo, logger)
						app.ForceDraw()
					}
					time.Sleep(10 * time.Second)

				}

			}()

		}).AddButton("STOP", func() {
		data.AddStop(true, logger)
		statusTextView.SetText("STOPPING...")
		app.ForceDraw()
		app.Lock()

	}).AddButton("QUIT", func() {
		level.Info(logger).Log("msg", "QUIT... ")
		app.Stop()
	})

	flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 7, 0, false).
		AddItem(statusTextView, 7, 1, false)
	return flex
}
