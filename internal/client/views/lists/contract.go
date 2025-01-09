package lists

import "time"

type refresher interface {
	Refresh()
	RefreshIn(t time.Duration)
}
