package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/getlantern/systray"
	"github.com/promignis/cwitch/config"
	"github.com/promignis/cwitch/icon"
	"github.com/promignis/cwitch/menu"
	"github.com/promignis/cwitch/timer"
	"github.com/promignis/cwitch/utils"
)

// currently the file used to store
const menuDataPath = "./data.json"

func InitMenu() {
	menu.MenuMap = config.GetTimerMap()

	// fetch value from data.json
	menuData, err := ioutil.ReadFile(menuDataPath)

	utils.FailOnError(fmt.Sprintf("Error in reading file %s", menuDataPath), err)

	json.Unmarshal(menuData, &menu.AllMenus)
}

func main() {
	InitMenu()
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)

}

func HandleMenuItem(hashMode uint32, ch chan struct{}) {
	for m := range ch {
		// just an empty struct
		_ = m
		currentTimer := menu.MenuMap[hashMode]

		var prevTimer *timer.Timer
		if menu.PrevSelected == 0 {
			// first time
			menu.PrevSelected = hashMode
		} else {
			// prevTimer exists
			prevTimer = menu.MenuMap[menu.PrevSelected]
		}

		if currentTimer.IsEnabled {
			// clicked on what is already enabled
			// end it
			currentTimer.End()
		} else {
			currentTimer.Begin()
			// end previous
			if prevTimer != nil {
				prevTimer.End()
			}
		}
		// store the current selected value
		menu.PrevSelected = hashMode
	}

}

// create Menu items, timers
// and stores both in timerMap
func createMenuItems(menus *menu.Menus) {
	for _, menuItem := range menus.Modes {
		item := systray.AddMenuItem(menuItem.Mode, menuItem.ToolTip)

		// generate hash from mode
		hashMode := utils.HashMode(menuItem.Mode)
		_, ok := menu.MenuMap[hashMode]

		if !ok {
			menu.MenuMap[hashMode] = &timer.Timer{hashMode, menuItem.Mode, false, time.Time{}, 0, "", item}
		} else {
			// already present update item
			prevItem := menu.MenuMap[hashMode]
			if prevItem != nil {
				prevItem.MenuItem = item
			}
		}
		go HandleMenuItem(hashMode, item.ClickedCh)
	}
}

// for things like ctrl+c
// and similar aborts
// before exiting save changes
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
	createMenuItems(menu.AllMenus)
	systray.SetIcon(icon.Data)

	systray.SetTooltip("Switch between activities consciously")
	mQuit := systray.AddMenuItem("Quit", "Quit cwitch")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		log.Println("Quitting cwitch")
	}()

	go HandleInterrupts()
}

func onExit() {
	// last selected timer value
	// to be saved
	if menu.PrevSelected != 0 {
		lastTimer := menu.MenuMap[menu.PrevSelected]
		lastTimer.End()
	}
	b, err := json.Marshal(menu.MenuMap)
	utils.FailOnError("Error while Marshaling menu.MenuMap", err)

	config.SaveTimerMap(b)
}
