package registerdmsview

import (

	//"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/lamassuiot/lamassuiot/pkg/dms-enroller/common/dto"
	"github.com/rivo/tview"
)

type DeviceForm struct {
	DeviceID    string
	DeviceAlias string
}

var (
	grid   *tview.Grid
	logger *log.Logger
)

type concreteObserver struct {
}

func (s *concreteObserver) Update(t interface{}) {
	// do something
	level.Info(*logger).Log("Observer updated", t)
	grid.Clear()
	grid = DrawGrid(t.(*observer.DeviceState).Dms)
}

func DrawGrid(dms dto.DMS) *tview.Grid {

	grid.AddItem(tview.NewTextView().SetText("DMS ID").SetTextColor(tcell.ColorYellow), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(dms.Id), 0, 1, 1, 1, 0, 0, false)

	grid.AddItem(tview.NewTextView().SetText("DMS Name").SetTextColor(tcell.ColorYellow), 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(dms.Name), 1, 1, 1, 1, 0, 0, false)

	grid.AddItem(tview.NewTextView().SetText("Creation Timestamp").SetTextColor(tcell.ColorYellow), 2, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(dms.CreationTimestamp), 2, 1, 1, 1, 0, 0, false)

	grid.AddItem(tview.NewTextView().SetText("Common Name").SetTextColor(tcell.ColorYellow), 3, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(dms.Subject.CN), 3, 1, 1, 1, 0, 0, false)

	return grid
}

func GetDecodedCertItem(inlog log.Logger, deviceData *observer.DeviceState) *tview.Grid {
	logger = &inlog
	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	grid = tview.NewGrid().
		SetRows(2, 2, 2, 2).
		SetColumns(30, 0).
		SetBorders(true)
	return DrawGrid(deviceData.Dms)
}
