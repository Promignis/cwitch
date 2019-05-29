package timer

import (
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/hako/durafmt"
)

type Timer struct {
	Mode          string
	IsEnabled     bool
	OnStart       time.Time
	LastUpdate    time.Time
	Elapsed       time.Duration
	PrettyElapsed string
	MenuItem      *systray.MenuItem `json:"-"`
}

func NewTimer(mode string, item *systray.MenuItem) *Timer {
	return &Timer{mode, false, time.Time{}, time.Time{}, 0, "", item}
}

// enable it, start timer
// and add check mark
func (t *Timer) Begin() {
	t.IsEnabled = true
	t.OnStart = time.Now()
	t.MenuItem.Check()
	t.LastUpdate = time.Now()
	// Anytime a timer starts
	// show it's current time
}

// Update the timer values
func (t *Timer) Update() {
	t.Elapsed += time.Since(t.LastUpdate)
	durTime := durafmt.Parse(t.Elapsed).String()
	splitTime := strings.Split(durTime, " ")
	untilLast := len(splitTime) - 2
	t.PrettyElapsed = strings.Join(splitTime[:untilLast], " ")
	t.LastUpdate = time.Now()
}

// mark it IsEnabled false,
// calculate time elapsed
// store easy to read timer value
// uncheck
func (t *Timer) End() {
	t.IsEnabled = false
	t.Update()
	t.MenuItem.Uncheck()
}

// This type is used to store
// timers
type TimerMap map[uint32]*Timer
