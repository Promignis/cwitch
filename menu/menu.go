package menu

import (
	"fmt"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/promignis/cwitch/timer"
	"github.com/rs/zerolog/log"
)

type Menu struct {
	Mode    string `json:"mode"`
	ToolTip string `json:"tooltip"`
	Emoji   string `json:"emoji"`
}

type Menus struct {
	Modes []*Menu `json:"modes"`
}

var (
	MenuMap  timer.TimerMap
	AllMenus *Menus
)

var CTray = &CwitchTray{"", nil, time.NewTicker(time.Second)}

// abstraction of current
// systray title
// +emoji and pretty time
type CwitchTray struct {
	Title           string
	CurrentMenuItem *CwitchItem
	perSecTicker    *time.Ticker
}

func (c *CwitchTray) UpdateItem(item *CwitchItem) {
	c.CurrentMenuItem = item
}

func (c *CwitchTray) UpdateTitle() {
	elapsed := c.CurrentMenuItem.Timer.PrettyElapsed
	emoji := c.CurrentMenuItem.Menu.Emoji
	var title string
	// if emoji is present
	// use it as well
	if emoji != "" {
		title = fmt.Sprintf("%s %s", emoji, elapsed)
	} else {
		title = elapsed
	}
	c.Title = title
	systray.SetTitle(title)
}

func (c *CwitchTray) PerSecondUpdates() {
	log.Info().Msg("Starting PerSecondUpdates")
	for t := range c.perSecTicker.C {
		log.Debug().Msg("PerSecondTick Tick")
		_ = t
		if c.CurrentMenuItem != nil && c.CurrentMenuItem.Timer.IsEnabled {
			log.Debug().Msgf("Running %s", c.CurrentMenuItem.Timer.Mode)
			c.CurrentMenuItem.Timer.Update()
			c.CurrentMenuItem.Update()
		}
	}
	log.Info().Msg("Exiting PerSecondUpdates")
}

// clean all resources
func (c *CwitchTray) Exit() {
	log.Info().Msg("Exiting CwitchTray")
	c.StopTicker()
}

func (c *CwitchTray) StopTicker() {
	log.Info().Msg("CwitchTray StopTicker called")
	c.perSecTicker.Stop()
}

func (c *CwitchTray) StartTicker() {
	log.Info().Msg("CwitchTray StartTicker called")
	c.perSecTicker = time.NewTicker(time.Second)
	c.PerSecondUpdates()
}

// Type is convoluted for now
// later will be flattened and simplified
type CwitchItem struct {
	Title string
	Menu  *Menu
	Timer *timer.Timer
	// menuItem
}

func (m *CwitchItem) StartItem() {
	log.Info().Msg("Starting CwitchItem")
	m.Timer.Begin()
	m.Update()
}

func (m *CwitchItem) EndItem() {
	log.Info().Msg("Ending CwitchItem")
	m.Timer.End()
	m.Update()
}

func (m *CwitchItem) GetTitle() string {
	var title string
	var elapsed string

	// Only show time for enabled
	// items
	if m.Timer.IsEnabled {
		elapsed = strings.Repeat(" ", 3) + m.Timer.PrettyElapsed
	} else {
		elapsed = ""
	}

	if m.Menu.Emoji == "" {
		// if emoji wasn't there
		title = fmt.Sprintf("%s%s", m.Menu.Mode, elapsed)
	} else {
		title = fmt.Sprintf("%s %s%s", m.Menu.Emoji, m.Menu.Mode, elapsed)
	}
	return title
}

func (m *CwitchItem) UpdateTitle() {
	title := m.GetTitle()
	log.Debug().Msgf("Updating CwitchItem Title %s", title)
	m.Title = title
	m.Timer.MenuItem.SetTitle(title)
}

// update cwitch title
// and tooltip
func (m *CwitchItem) Update() {
	log.Debug().Msg("CwitchItem Update")
	m.UpdateTitle()
	log.Debug().Msg("Updating tooltip")
	m.Timer.MenuItem.SetTooltip(m.Timer.PrettyElapsed)
}
