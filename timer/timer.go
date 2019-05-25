package timer

import (
	"time"

	"github.com/getlantern/systray"
)

type Timer struct {
	Id            uint32 `json:"-"`
	Mode          string
	IsEnabled     bool
	OnStart       time.Time
	Elapsed       time.Duration
	PrettyElapsed string
	MenuItem      *systray.MenuItem `json:"-"`
}

// enable it, start timer
// and add check mark
func (t *Timer) Begin() {
	t.IsEnabled = true
	t.OnStart = time.Now()
	t.MenuItem.Check()
	// Anytime a timer starts
	// show it's current time
	// systray.SetTitle(t.PrettyElapsed)
}

// mark it IsEnabled false,
// calculate time elapsed
// store easy to read timer value
// uncheck
func (t *Timer) End() {
	t.IsEnabled = false
	t.Elapsed += time.Since(t.OnStart)
	t.PrettyElapsed = t.Elapsed.String()
	t.MenuItem.Uncheck()
}

// This type is used to store
// timers
type TimerMap map[uint32]*Timer
