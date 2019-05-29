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
}

// mark it IsEnabled false,
// calculate time elapsed
// store easy to read timer value
// uncheck
func (t *Timer) End() {
	t.IsEnabled = false
	t.Elapsed += time.Since(t.OnStart)
	time := durafmt.Parse(t.Elapsed).String()
	splitTime := strings.Split(time, " ")
	untilLast := len(splitTime) - 2
	t.PrettyElapsed = strings.Join(splitTime[:untilLast], " ")
	t.MenuItem.Uncheck()
}

// This type is used to store
// timers
type TimerMap map[uint32]*Timer
