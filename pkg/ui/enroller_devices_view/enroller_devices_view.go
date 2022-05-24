package enrollerdevicesview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/rivo/tview"
)

type DeviceForm struct {
	DeviceID    string
	DeviceAlias string
}

type concreteObserver struct {
}

var (
	grid   *tview.Grid
	logger *log.Logger
	app    *tview.Application
)

func (s *concreteObserver) Update(t interface{}) {
	// do something

	grid = DrawView(t.(*observer.DeviceState).Devices)
}

func DrawView(devices []observer.EnrolledDeviceData) *tview.Grid {

	grid.AddItem(tview.NewTextView().SetText("Device Alias").SetTextColor(tcell.ColorYellow), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText("Device ID").SetTextColor(tcell.ColorYellow), 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText("Enrolled CA").SetTextColor(tcell.ColorYellow), 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText("Enrolled Ceritifcate Serial Number").SetTextColor(tcell.ColorYellow), 0, 3, 1, 1, 0, 0, false)

	for idx, device := range devices {
		grid.AddItem(tview.NewTextView().SetText(device.DeviceAlias), idx+1, 0, 1, 1, 0, 0, false)
		grid.AddItem(tview.NewTextView().SetText(device.DeviceID), idx+1, 1, 1, 1, 0, 0, false)
		grid.AddItem(tview.NewTextView().SetText(device.EnrolledCA), idx+1, 2, 1, 1, 0, 0, false)
		grid.AddItem(tview.NewTextView().SetText(device.EnrolledCertSerialNumber), idx+1, 3, 1, 1, 0, 0, false)
	}
	return grid

}

func GetEnrolledDevicesItem(inlogger log.Logger, deviceData *observer.DeviceState, app *tview.Application) tview.Primitive {
	logger = &inlogger

	grid = tview.NewGrid().
		SetRows(2).
		SetColumns(30, 0, 0).
		SetBorders(true)

	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	return DrawView(deviceData.Devices)
}
