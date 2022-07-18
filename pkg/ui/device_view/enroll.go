package deviceview

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"strings"
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
	var dmsCrtFile, dmsKeyFile string
	var flex *tview.Flex
	statusTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText(" ")
	form2 := tview.NewForm()
	form := tview.NewForm().
		AddInputField("Device Manager endpoint", devManagerEndpoint, 50, nil, func(text string) {
			devManagerEndpoint = text
		}).
		AddButton("Auto-Enroll", func() {
			statusTextView.SetText("AUTO-ENROLLING...")
			app.ForceDraw()

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

					certContent, err := ioutil.ReadFile(dmsCrtFile)
					cpb, _ := pem.Decode(certContent)

					crt, err := x509.ParseCertificate(cpb.Bytes)

					alias, id, sn, ca, err := service.Enroll(lamassuEstClient, data, data.DeviceFile, data.Aps, crt.Subject.CommonName, logger)
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
	form2.AddDropDown("DMS ID", getDmsIds(data.Config.Dms.DmsStore, logger), 0, func(option string, optionIndex int) {
		if len(option) > 0 && optionIndex > 0 {
			level.Info(logger).Log("msg", option)
			data.Dms.Id = strings.TrimPrefix(option, "dms-")
			dmsKeyFile = data.Config.Dms.DmsStore + "/" + option + ".key"
			dmsCrtFile = data.Config.Dms.DmsStore + "/" + option + ".crt"
		}
	}).AddButton("Refresh", func() {
		form2.Clear(false)
		form2.AddDropDown("DMS ID", getDmsIds(data.Config.Dms.DmsStore, logger), 0, func(option string, optionIndex int) {
			if len(option) > 0 && optionIndex > 0 {
				dmsKeyFile = data.Config.Dms.DmsStore + "/" + option + ".key"
				dmsCrtFile = data.Config.Dms.DmsStore + "/" + option + ".crt"
			}
		})

	})
	flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form2, 7, 0, false).
		AddItem(form, 7, 0, false).
		AddItem(statusTextView, 7, 1, false)
	return flex
}
func getDmsIds(path string, logger log.Logger) []string {
	files, _ := ioutil.ReadDir(path)
	var s []string
	s = append(s, "")

	for _, file := range files {
		fileName := file.Name()
		if len(fileName) > 0 {
			last := len(fileName) - 4
			fileName = fileName[:last]
			if !contains(s, fileName) {

				s = append(s, fileName)
			}
		}

	}
	return s
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
