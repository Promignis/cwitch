package utils

import (
	"fmt"
	"log"
	"os"
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

// fail if error exists
// func FailOnError(err error) {
// 	if err != nil {
// 		log.Fatalf("Error : %s", err.Error())
// 	}
// }

// If error log and panic
func FailOnError(errorStr string, err error) {
	if err != nil {
		log.Fatalf(errorStr)
	}
}

// log the error and don't fail
func LogOnError(err error) {
	if err != nil {
		log.Printf("Error : %s\n", err.Error())
	}
}

// expects one %s (Error)
func LogfOnError(formatStr string, err error) {
	if err != nil {
		log.Printf(formatStr, err.Error())
	}
}

func AddErrorIfExists(formatStr string, err error) string {
	if err != nil {
		return fmt.Sprintf(formatStr, err.Error())
	}
	return ""
}
