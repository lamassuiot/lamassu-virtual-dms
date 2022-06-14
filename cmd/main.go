// Demo code for the Flex primitive.
package main

import (
	"flag"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-default-dms/pkg/config"
	filestore "github.com/lamassuiot/lamassu-default-dms/pkg/device/store/file"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	deviceview "github.com/lamassuiot/lamassu-default-dms/pkg/ui/device_view"
	enrollerdevicesview "github.com/lamassuiot/lamassu-default-dms/pkg/ui/enroller_devices_view"
	"github.com/lamassuiot/lamassuiot/pkg/dms-enroller/common/dto"
	"github.com/rivo/tview"
)

var DMS_B64_PRIVATE_KEY = flag.String("DMS_B64_PRIVATE_KEY", "", "privkey")
var DMS_ID = flag.String("DMS_ID", "", "dmsid")

func main() {
	var logger log.Logger
	f, _ := os.OpenFile("./vdms.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(f))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	flag.Parse()
	cfg, _ := config.NewConfig()
	deviceFile := filestore.NewFile(cfg.Dms.DeviceStore, logger)
	level.Info(logger).Log("msg", "Devices CSRs, CERTs and KEY filesystem home path created")
	dmsFile := filestore.NewFile(cfg.Dms.DmsStore, logger)
	level.Info(logger).Log("msg", "DMS CSRs, CERTs and KEY filesystem home path created")

	app := tview.NewApplication()

	obs := observer.DeviceState{
		Devices:    make([]observer.EnrolledDeviceData, 0),
		Config:     cfg,
		Aps:        "",
		DmsPrivKey: *DMS_B64_PRIVATE_KEY,
		DmsId:      *DMS_ID,
		DeviceFile: deviceFile,
		DmsFile:    dmsFile,
		Dms:        dto.DMS{},
	}

	flex := tview.NewFlex().
		AddItem(deviceview.GetItem(logger, &obs, app), 70, 1, false).
		AddItem(enrollerdevicesview.GetEnrolledDevicesItem(logger, &obs, app), 0, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
