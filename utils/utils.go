package utils

import "os"

func FileExists(path string) (bool, error) {
	// there can be cases when file may or may not exist
	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go/22467409#22467409
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}
