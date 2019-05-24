package timer

import (
	"time"

	"github.com/getlantern/systray"
)

type Timer struct {
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
}

func (t *Timer) End() {
	t.IsEnabled = false
	t.Elapsed += time.Since(t.OnStart)
	t.PrettyElapsed = t.Elapsed.String()
	t.MenuItem.Uncheck()
}

// This type is used to store
// timers
type TimerMap map[string]*Timer
