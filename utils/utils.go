package utils

import (
	"fmt"
	"hash/fnv"
	"os"

	"github.com/rs/zerolog/log"
)

func FileExists(path string) (bool, error) {
	// there can be cases when file may or may not exist
	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go/22467409#22467409
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

// If error log and panic
func FailOnError(errorStr string, err error) {
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf(errorStr)
	}
}

// log the error and don't fail
func LogOnError(err error) {
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Error occured : %s", err.Error())
	}
}

// Add error to format string
func AddErrorIfExists(formatStr string, err error) string {
	if err != nil {
		return fmt.Sprintf(formatStr, err.Error())
	}
	return ""
}

// generate a hash of the mode
// used as an Id and used
// to find previously same mode
func HashMode(mode string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(mode))
	return h.Sum32()
}
