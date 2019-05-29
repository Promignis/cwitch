package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

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

// stores last selected cwitchItem
var prevCwitchItem *menu.CwitchItem

func HandleMenuItem(hashMode uint32, ch chan struct{}, cwitchItem *menu.CwitchItem) {
	for m := range ch {
		// just an empty struct
		_ = m
		currentTimer := menu.MenuMap[hashMode]
		menu.CTray.UpdateItem(cwitchItem)

		if prevCwitchItem == nil {
			// first time
			prevCwitchItem = cwitchItem
		}

		if currentTimer.IsEnabled {
			// clicked on what is already enabled
			// end it
			cwitchItem.EndItem()
		} else {
			cwitchItem.StartItem()
			if prevCwitchItem != cwitchItem {
				prevCwitchItem.EndItem()
			}
		}
		prevCwitchItem = cwitchItem
	}
}

// create Menu items, timers
// and stores both in timerMap
func createMenuItems(menus *menu.Menus) {
	for _, menuItem := range menus.Modes {
		// menuItem := menu.MenuItem{menu.Mode, menu, timer}
		item := systray.AddMenuItem(menuItem.Mode, menuItem.ToolTip)

		// generate hash from mode
		hashMode := utils.HashMode(menuItem.Mode)
		_, ok := menu.MenuMap[hashMode]

		cwitchItem := &menu.CwitchItem{menuItem.Mode, menuItem, nil}
		if !ok {
			newTimer := timer.NewTimer(menuItem.Mode, item)
			menu.MenuMap[hashMode] = newTimer
			cwitchItem.Timer = newTimer
		} else {
			// already present update item
			// and cwitch item timer
			prevTimer := menu.MenuMap[hashMode]
			if prevTimer != nil {
				prevTimer.MenuItem = item
				cwitchItem.Timer = prevTimer
			}
		}
		cwitchItem.Update()
		go HandleMenuItem(hashMode, item.ClickedCh, cwitchItem)
	}
}

// for things like ctrl+c
// and similar aborts
// before exiting save changes
func handleInterrupts() {
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

	go handleInterrupts()

	go menu.CTray.PerSecondUpdates()
}

func onExit() {
	menu.CTray.Exit()
	// last selected timer value
	// to be saved
	if prevCwitchItem != nil {
		prevCwitchItem.EndItem()
	}
	b, err := json.Marshal(menu.MenuMap)
	utils.FailOnError("Error while Marshaling menu.MenuMap", err)

	config.SaveTimerMap(b)
}
