package deviceview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/rivo/tview"
)

func GetItem(logger log.Logger, deviceData *observer.DeviceState, app *tview.Application) tview.Primitive {

	pages := tview.NewPages()

	activePage := 1
	pages.AddPage("Register DMS page", GetRegisterDMSItem(logger, deviceData, app, pages), true, activePage == 1)
	pages.AddPage("automatic enroll-page", GetEnrollItem(logger, deviceData, app), true, activePage == 2)

	actionsList := tview.NewList().
		ShowSecondaryText(false)

	actionsList.AddItem("Register DMS", "", '1', func() {
		activePage = 1
		pages.SwitchToPage("Register DMS page")
	})
	actionsList.AddItem("Auto Enroll", "", '2', func() {
		activePage = 2
		pages.SwitchToPage("automatic enroll-page")
	})
	test := tview.NewBox().
		SetBorder(false).
		SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			// Draw a horizontal line across the middle of the box.
			centerY := y + height/2
			for cx := x + 1; cx < x+width-1; cx++ {
				screen.SetContent(cx, centerY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.NewHexColor(0xffffff)))
			}
			// Write some text along the horizontal line.
			tview.Print(screen, " Operation ", x+1, centerY, width-2, tview.AlignCenter, tcell.NewHexColor(0xF9F1A5))
			// Space for other content.
			return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(test, 2, 0, false).
		AddItem(pages, 0, 1, false).
		AddItem(actionsList, 4, 0, true)

	flex.SetBorder(true)
	return flex
}
