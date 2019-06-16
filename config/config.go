package config

import (
	"encoding/json"

	"github.com/promignis/cwitch/timer"
	"github.com/promignis/cwitch/utils"
	"github.com/shibukawa/configdir"
)

const (
	configFileName = "timer_map.json"
	vendor         = "Promignis"
	appName        = "cwitch"
)

var timerConfig timer.TimerMap

// unmarshal timermap
func GetTimerMap() timer.TimerMap {
	var config timer.TimerMap
	configDirs := configdir.New(vendor, appName)
	folder := configDirs.QueryFolderContainsFile(configFileName)
	if folder != nil {
		data, _ := folder.ReadFile(configFileName)
		json.Unmarshal(data, &config)
		return config
	}
	// if not found, default
	return make(timer.TimerMap)
}

// marshal and save timermap
func SaveTimerMap(timerData []byte) {
	configDirs := configdir.New(vendor, appName)
	folders := configDirs.QueryFolders(configdir.Global)
	err := folders[0].WriteFile(configFileName, timerData)

	// crash if you can't save
	// as you lose data
	utils.FailOnError(utils.AddErrorIfExists("Error while saving TimerMap : %s", err), err)
}

// not used in logger package
// to prevent cyclic dependency
func GetConfigPath() string {
	configDirs := configdir.New(vendor, appName)
	folders := configDirs.QueryFolders(configdir.Global)
	return folders[0].Path
	//return fmt.Sprintf("%s/%s/%s", configdir.Global, vendor, appName)
}
