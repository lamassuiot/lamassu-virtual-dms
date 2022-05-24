package observer

import "log"

type Subject interface {
	Attach(o Observer) (bool, error)
	Detach(o Observer) (bool, error)
	Notify(logger log.Logger) (bool, error)
}
