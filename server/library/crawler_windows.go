// +build windows

package library

import (
	"errors"
	"io/ioutil"
	"net/http"
)

const jsonUrl = "https://mediathekdirekt.de/good.json"

func DownloadJson(crawlerDirName string, jsonPath string) error {
	resp, err := http.Get(jsonUrl)
	if err != nil {
		return errors.New("failed sending get request: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("failed reading request body: " + err.Error())
	}
	err = ioutil.WriteFile(jsonPath, body, 0644)
	if err != nil {
		return errors.New("failed writing to file: " + err.Error())
	}
	return nil
}
