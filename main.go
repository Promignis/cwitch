package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/getlantern/systray"
	"github.com/promignis/cwitch/icon"
	"github.com/promignis/cwitch/utils"
)

var menuMap map[string]*Timer

const (
	menuDataPath = "./data.json"
	dataSavePath = "./saved.json"
)

type Timer struct {
	IsEnabled     bool
	OnStart       time.Time
	Elapsed       time.Duration
	PrettyElapsed string
	MenuItem      *systray.MenuItem `json:"-"`
}

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

type Menu struct {
	Mode    string `json:"mode"`
	ToolTip string `json:"tooltip"`
}

type Menus struct {
	Modes []*Menu `json:"modes"`
}

var menus *Menus
var prevModeSelected string

func InitMenu() {
	menuMap = make(map[string]*Timer)

	isExist, err := utils.FileExists(dataSavePath)

	if isExist && err != nil {
		// some error while checking file stat
		log.Fatalf("Error while checking file existence %s", err.Error())
	} else if isExist && err == nil {
		// actually file exists
		menuMapData, err := ioutil.ReadFile(dataSavePath)
		if err != nil {
			log.Printf("Error reading file %s", dataSavePath)
		}
		json.Unmarshal(menuMapData, &menuMap)
	}

	menuData, err := ioutil.ReadFile(menuDataPath)
	if err != nil {
		log.Fatal("Error in reading file %s", menuDataPath)
	}
	json.Unmarshal(menuData, &menus)
}

func main() {
	InitMenu()
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)

}

func HandleMenuItem(menuStr string, ch chan struct{}) {
	for msg := range ch {
		_ = msg
		timer := menuMap[menuStr]

		var prevTimer *Timer
		if prevModeSelected == "" {
			// first time
			prevModeSelected = menuStr
		} else {
			// prevTimer exists
			prevTimer = menuMap[prevModeSelected]
		}

		if timer.IsEnabled {
			// clicked on what is already enabled
			// end it
			timer.End()
		} else {
			timer.Begin()
			// end previous
			if prevTimer != nil {
				prevTimer.End()
			}
		}
		// store the current selected value
		prevModeSelected = menuStr
	}

}

func createMenuItems(menus *Menus) {
	for _, menuItem := range menus.Modes {
		item := systray.AddMenuItem(menuItem.Mode, menuItem.ToolTip)
		_, ok := menuMap[menuItem.Mode]
		if !ok {
			menuMap[menuItem.Mode] = &Timer{false, time.Time{}, 0, "", item}
		} else {
			// already present update item
			prevItem := menuMap[menuItem.Mode]
			if prevItem != nil {
				prevItem.MenuItem = item
			}
		}
		go HandleMenuItem(menuItem.Mode, item.ClickedCh)
	}
}

func HandleInterrupts() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Sig %s recieved\n", sig.String())
			onExit()
			os.Exit(1)
		}
	}()
}

func onReady() {
	createMenuItems(menus)
	// systray.SetTitle("cwitch")
	systray.SetIcon(icon.Data)

	systray.SetTooltip("Switch between activities consciously")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		log.Println("Quitting")
	}()

	go HandleInterrupts()
}

func onExit() {
	if prevModeSelected != "" {
		lastTimer := menuMap[prevModeSelected]
		lastTimer.End()
	}
	b, err := json.Marshal(menuMap)
	if err != nil {
		log.Fatal("Error while Marshaling menuMap")
	}

	ioutil.WriteFile(dataSavePath, b, 0644)
}
