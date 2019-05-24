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
	"github.com/promignis/cwitch/timer"
	"github.com/promignis/cwitch/utils"
)

const menuDataPath = "./data.json"

type Menu struct {
	Mode    string `json:"mode"`
	ToolTip string `json:"tooltip"`
}

type Menus struct {
	Modes []*Menu `json:"modes"`
}

var menus *Menus
var prevModeSelected string

var menuMap timer.TimerMap

func InitMenu() {
	menuMap = config.GetTimerMap()

	// isExist, err := utils.FileExists(dataSavePath)

	// if isExist && err != nil {
	// 	// some error while checking file stat
	// 	log.Fatalf("Error while checking file existence %s", err.Error())
	// } else if isExist && err == nil {
	// 	// actually file exists
	// 	menuMapData, err := ioutil.ReadFile(dataSavePath)
	// 	if err != nil {
	// 		log.Printf("Error reading file %s", dataSavePath)
	// 	}
	// 	json.Unmarshal(menuMapData, &timer.MenuMap)
	// }

	// fetch value from data.json
	menuData, err := ioutil.ReadFile(menuDataPath)

	utils.FailOnError(fmt.Sprintf("Error in reading file %s", menuDataPath), err)

	json.Unmarshal(menuData, &menus)
}

func main() {
	InitMenu()
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)

}

func HandleMenuItem(menuStr string, ch chan struct{}) {
	for m := range ch {
		// just an empty struct
		_ = m
		currentTimer := menuMap[menuStr]

		var prevTimer *timer.Timer
		if prevModeSelected == "" {
			// first time
			prevModeSelected = menuStr
		} else {
			// prevTimer exists
			prevTimer = menuMap[prevModeSelected]
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
		prevModeSelected = menuStr
	}

}

func createMenuItems(menus *Menus) {
	for _, menuItem := range menus.Modes {
		item := systray.AddMenuItem(menuItem.Mode, menuItem.ToolTip)
		_, ok := menuMap[menuItem.Mode]
		if !ok {
			menuMap[menuItem.Mode] = &timer.Timer{false, time.Time{}, 0, "", item}
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
	// last selected timer value
	// to be saved
	if prevModeSelected != "" {
		lastTimer := menuMap[prevModeSelected]
		lastTimer.End()
	}
	b, err := json.Marshal(menuMap)
	utils.FailOnError("Error while Marshaling menuMap", err)

	config.SaveTimerMap(b)
}
