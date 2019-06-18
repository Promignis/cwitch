package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/promignis/cwitch/config"
	"github.com/promignis/cwitch/icon"
	"github.com/promignis/cwitch/logger"
	"github.com/promignis/cwitch/menu"
	"github.com/promignis/cwitch/timer"
	"github.com/promignis/cwitch/utils"
	"github.com/rs/zerolog"
)

// currently the default file as data
var menuDataPath = "data.json"

// initialize flags
// data file
// logger
func InitMenu() {
	currExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	utils.FailOnError("Error getting executable path", err)
	fullPath := path.Join(currExecPath, menuDataPath)
	dataFilePath := flag.String("data", fullPath, "Cwitch json file path")
	// to show debugging logs
	debug := flag.Bool("debug", false, "Run Cwitch in debug mode")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	menu.MenuMap = config.GetTimerMap()

	// fetch value from data.json
	menuData, err := ioutil.ReadFile(*dataFilePath)

	utils.FailOnError(fmt.Sprintf("Error in reading file %s", menuDataPath), err)

	json.Unmarshal(menuData, &menu.AllMenus)
	logger.Log.Info().Msg("Data file loaded")
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

		logger.Log.Info().Msgf("%s mode selected", currentTimer.Mode)

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

		cwitchItem := menu.NewCwitchItem(menuItem.Mode, menuItem)

		if !ok {
			newTimer := timer.NewTimer(menuItem.Mode, item)
			menu.MenuMap[hashMode] = newTimer
			cwitchItem.Timer = newTimer
		} else {
			// already present, update item
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
			logger.Log.Info().Msgf("Sig %s recieved", sig.String())
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
		logger.Log.Info().Msg("Quitting cwitch")
	}()

	go handleInterrupts()

	menu.CTray.StartTicker()
}

func onExit() {
	logger.Log.Info().Msg("Systray onExit called")
	menu.CTray.Exit()
	// last selected timer value
	// to be saved
	if prevCwitchItem != nil {
		prevCwitchItem.EndItem()
	}
	b, err := json.Marshal(menu.MenuMap)
	utils.FailOnError("Error while Marshaling menu.MenuMap", err)

	config.SaveTimerMap(b)
	logger.Log.Debug().Msgf("Saving timer %s", string(b))
	logger.Log.Info().Msg("Saving timermap")
}
