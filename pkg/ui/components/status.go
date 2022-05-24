package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//status == 0	ongoing
//status == 1	finished OK
//status == -1	finished ERROR

func NewStatus(action string, status int) tview.Primitive {
	status_string := "Ongoing action"
	if status == 1 {
		status_string = "Action sucesfully executed"
	} else if status == -1 {
		status_string = "Error while executing action"

	}

	f := tcell.ColorBlack
	b := tcell.ColorWhite
	if status == 1 {
		f = tcell.ColorBlack
		b = tcell.ColorGreen

	} else if status == -1 {
		f = tcell.ColorBlack
		b = tcell.ColorRed
	}
	a := tview.NewTextView().SetTextColor(f).SetText("Action: " + action)
	a.SetBorderPadding(1, 0, 2, 0)
	a.SetBackgroundColor(b)
	s := tview.NewTextView().SetTextColor(f).SetText("Status: " + status_string)
	s.SetBackgroundColor(b)
	s.SetBorderPadding(0, 1, 2, 0)
	status_view := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a, 2, 0, false).
		AddItem(s, 2, 0, false)
	status_view.SetBackgroundColor(b)

	return status_view

}
