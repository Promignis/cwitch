package logger

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/shibukawa/configdir"
)

var (
	Log      zerolog.Logger
	fileName = "cwitch.log"
)

const (
	vendor  = "Promignis"
	appName = "cwitch"
)

// initialize Log
// variable
func init() {
	configDirs := configdir.New(vendor, appName)
	folders := configDirs.QueryFolders(configdir.Global)
	configPath := folders[0].Path
	logPath := path.Join(configPath, fileName)
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Errorf("Error while setting up Logger %s", err)
	}
	Log = zerolog.New(file).With().Timestamp().Logger()
}
