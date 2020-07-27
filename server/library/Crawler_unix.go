// +build !windows

package library

import (
	"errors"
	"os/exec"
)

func DownloadJson(crawlerDirName string, jsonPath string) error {
	err := exec.Command("/bin/bash", "-c", "cd", crawlerDirName, "&&", "python3", "mediathek.py").Run()
	if err != nil {
		return errors.New("error executing python script: " + err.Error())
	}
	return nil
}
