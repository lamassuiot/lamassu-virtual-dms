package csrview

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-default-dms/pkg/observer"
	"github.com/lamassuiot/lamassu-default-dms/pkg/utils"
	"github.com/lamassuiot/lamassuiot/pkg/dms-enroller/common/dto"
	"github.com/rivo/tview"
)

type DeviceForm struct {
	DeviceID    string
	DeviceAlias string
}

type concreteObserver struct {
}

var (
	view   *tview.TextView
	logger *log.Logger
)

func (s *concreteObserver) Update(t interface{}) {
	// do something
	level.Info(*logger).Log("Observer updated", t)
	view = DrawView(t.(*observer.DeviceState).Dms)
}

func DrawView(dms dto.DMS) *tview.TextView {
	txt := "This DMS has not been yet registered"
	if dms.CerificateBase64 != "" {
		txt, _ = utils.DecodeB64(dms.CsrBase64)
	}

	view.SetText(txt)
	view.SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	return view
}

func GetRawCertItem(inlogger log.Logger, deviceData *observer.DeviceState) tview.Primitive {
	logger = &inlogger

	view = tview.NewTextView()

	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	return DrawView(deviceData.Dms)
}
