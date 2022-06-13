package deviceview

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/lamassuiot/lamassu-default-dms/pkg/service"
	enrollerdevicesview "github.com/lamassuiot/lamassu-default-dms/pkg/ui/enroller_devices_view"
	"github.com/lamassuiot/lamassu-default-dms/pkg/utils"
	lamassuDMSClient "github.com/lamassuiot/lamassuiot/pkg/dms-enroller/client"
	"github.com/lamassuiot/lamassuiot/pkg/dms-enroller/common/dto"
	"github.com/lamassuiot/lamassuiot/pkg/utils/client"
	"github.com/rivo/tview"
)

var view *tview.Grid

type concreteObserver struct {
}

func GetRegisterDMSItem(logger log.Logger, data *observer.DeviceState, app *tview.Application, pages *tview.Pages) tview.Primitive {
	dmsEndpoint := data.Config.Dms.Endpoint
	var key_type, key_bits string
	dmsName := data.Config.Dms.Name
	common_name := data.Config.Dms.CN
	country := data.Config.Dms.C
	locality := data.Config.Dms.L
	organization := data.Config.Dms.O
	organization_unit := data.Config.Dms.OU
	state := data.Config.Dms.ST
	var flex *tview.Flex
	statusTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText(" ")
	var keybits []string
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
	if data.DmsPrivKey != "" {
		view = tview.NewGrid().
			SetRows(1, 1).
			SetColumns(10, 0)
		view = DrawDMS(data)
		form := tview.NewForm().
			AddButton("Check if the DMS is approved", func() {
				dmsKey, err := utils.ReadKey(data.DmsPrivKey)
				if err != nil {
					level.Error(logger).Log("err", err)
				}
				for {
					statusTextView.SetText("DMS with ID " + data.DmsId + " is registered, Pending approval...")
					app.ForceDraw()
					dms, err := dmsClient.GetDMSbyID(context.Background(), data.DmsId)
					if err != nil {
						level.Error(logger).Log("err", err)
					}
					if dms.Status == "APPROVED" {
						data.EditDMS(dms, logger)
						level.Info(logger).Log("msg", data.Dms.CerificateBase64)
						cert, _ := base64.StdEncoding.DecodeString(data.Dms.CerificateBase64)
						block, _ := pem.Decode([]byte(cert))
						level.Error(logger).Log("err", err)
						data.DmsFile.InsertCERT(data.DmsId, block.Bytes, "dms", "", "")
						err = ioutil.WriteFile(data.Config.Dms.DmsStore+"/dms-"+dms.Id+".key", dmsKey, 0644)
						if err != nil {
							level.Error(logger).Log("err", err)
						}
						break
					}
					time.Sleep(10 * time.Second)
				}
				statusTextView.SetText("DMS with ID " + data.Dms.Id + " and name " + data.Dms.Name + " approved...")
				app.ForceDraw()
				modal := tview.NewModal().
					SetText("DMS successfully approved. Do you want to auto-enroll?").
					AddButtons([]string{"Quit", "OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "OK" {
							flex := tview.NewFlex().
								AddItem(GetEnrollItem(logger, data, app), 70, 1, false).
								AddItem(enrollerdevicesview.GetEnrolledDevicesItem(logger, data, app), 0, 1, false)
							flex.SetBorder(true)
							if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
								panic(err)
							}

						} else {
							flex := tview.NewFlex().
								AddItem(GetItem(logger, data, app), 70, 1, false).
								AddItem(enrollerdevicesview.GetEnrolledDevicesItem(logger, data, app), 0, 1, false)

							if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
								panic(err)
							}
						}
					})
				app.SetRoot(modal, false)
			})
		flex = tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(view, 3, 1, false).
			AddItem(form, 4, 1, false).
			AddItem(statusTextView, 10, 2, false)
		return flex

	} else {
		form := tview.NewForm().
			AddInputField("DMS Register endpoint", dmsEndpoint, 50, nil, func(text string) {
				dmsEndpoint = text
			}).
			AddInputField("DMS Name", dmsName, 50, nil, func(text string) {
				dmsName = text
			}).
			AddInputField("Common Name", common_name, 50, nil, func(text string) {
				common_name = text
			}).
			AddInputField("Country", country, 50, nil, func(text string) {
				country = text
			}).
			AddInputField("Locality", locality, 50, nil, func(text string) {
				locality = text
			}).
			AddInputField("Organization", organization, 50, nil, func(text string) {
				organization = text
			}).
			AddInputField("Organization unit", organization_unit, 50, nil, func(text string) {
				organization_unit = text
			}).
			AddInputField("State", state, 50, nil, func(text string) {
				state = text
			})
		form3 := tview.NewForm().
			AddDropDown("Key Bits", data.Bits, 10, func(option string, optionIndex int) {
				key_bits = option
			})
		form2 := tview.NewForm().
			AddDropDown("Key Type", []string{"RSA", "EC"}, 10, func(option string, optionIndex int) {
				key_type = option
				form3.Clear(false)
				level.Info(logger).Log("msg", key_type)
				if key_type == "RSA" {
					keybits = []string{"2048", "3072", "4096"}
				} else if key_type == "EC" {
					keybits = []string{"256", "384"}
				} else {
					keybits = []string{}
				}
				data.AddKeyBits(keybits, logger)
				form3.AddDropDown("Key Bits", data.Bits, 10, func(option string, optionIndex int) {
					key_bits = option
				})
			})
		form4 := tview.NewForm().
			AddButton("Register DMS", func() {
				statusTextView.SetText("Registering DMS...")
				key_bit, _ := strconv.Atoi(key_bits)
				Subject := dto.Subject{
					CN: common_name,
					C:  country,
					L:  locality,
					O:  organization,
					OU: organization_unit,
					ST: state,
				}
				PrivateKeyMetadata := dto.PrivateKeyMetadata{
					KeyType: key_type,
					KeyBits: key_bit,
				}
				id, err := service.RegisterDMS(data, data.DmsFile, dmsName, Subject, PrivateKeyMetadata, dmsClient, logger)
				if err != nil {
					level.Error(logger).Log("err", err)
				}

				for {
					statusTextView.SetText("DMS with ID " + id + " is registered, Pending approval...")
					app.ForceDraw()
					dms, err := dmsClient.GetDMSbyID(context.Background(), id)
					if err != nil {
						level.Error(logger).Log("err", err)
					}
					if dms.Status == "APPROVED" {
						data.EditDMS(dms, logger)
						level.Info(logger).Log("msg", data.Dms.CerificateBase64)
						cert, _ := base64.StdEncoding.DecodeString(data.Dms.CerificateBase64)
						block, _ := pem.Decode([]byte(cert))
						level.Error(logger).Log("err", err)
						data.DmsFile.InsertCERT(id, block.Bytes, "dms", "", "")

						break
					}
					time.Sleep(10 * time.Second)
				}
				statusTextView.SetText("DMS with ID " + data.Dms.Id + "and name " + data.Dms.Name + " approved...")
				app.ForceDraw()
				modal := tview.NewModal().
					SetText("DMS successfully approved. Do you want to auto-enroll?").
					AddButtons([]string{"Quit", "OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "OK" {
							flex := tview.NewFlex().
								AddItem(GetEnrollItem(logger, data, app), 70, 1, false).
								AddItem(enrollerdevicesview.GetEnrolledDevicesItem(logger, data, app), 0, 1, false)
							flex.SetBorder(true)
							if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
								panic(err)
							}

						} else {
							flex := tview.NewFlex().
								AddItem(GetItem(logger, data, app), 70, 1, false).
								AddItem(enrollerdevicesview.GetEnrolledDevicesItem(logger, data, app), 0, 1, false)

							if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
								panic(err)
							}
						}
					})
				app.SetRoot(modal, false)

			}).AddButton("STOP", func() {
			level.Info(logger).Log("msg", "STOP... ")
			flex := tview.NewFlex().
				AddItem(GetItem(logger, data, app), 70, 1, false).
				AddItem(enrollerdevicesview.GetEnrolledDevicesItem(logger, data, app), 0, 1, false)

			if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
				panic(err)
			}

		}).AddButton("QUIT", func() {
			level.Info(logger).Log("msg", "QUIT... ")
			app.Stop()
		})
		flex = tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(form, 18, 1, false).
			AddItem(form2, 4, 1, false).
			AddItem(form3, 4, 1, false).
			AddItem(form4, 4, 1, false).
			AddItem(statusTextView, 7, 1, false)
		return flex
	}

}
func (s *concreteObserver) Update(t interface{}) {
	view.Clear()
	view = DrawDMS(t.(*observer.DeviceState))
}

func DrawDMS(obs *observer.DeviceState) *tview.Grid {
	view.AddItem(tview.NewTextView().SetText("DMS ID").SetTextColor(tcell.ColorYellow), 0, 0, 1, 1, 0, 0, false)
	view.AddItem(tview.NewTextView().SetText(obs.DmsId), 0, 1, 1, 1, 0, 0, false)

	return view
}
