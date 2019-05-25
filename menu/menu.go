package menu

import (
	"github.com/promignis/cwitch/timer"
)

type Menu struct {
	Mode    string `json:"mode"`
	ToolTip string `json:"tooltip"`
}

type Menus struct {
	Modes []*Menu `json:"modes"`
}

var (
	PrevSelected uint32
	MenuMap      timer.TimerMap
	AllMenus     *Menus
)
