/**
 ******************************************************************************
 * @file    code.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package CodePackage

import (
	"armcnc/framework/config"
	"github.com/djherbis/times"
	"os"
	"strings"
	"time"
)

type Code struct {
	Path string `json:"path"`
}

type Data struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Describe string    `json:"describe"`
	Version  string    `json:"version"`
	Time     time.Time `json:"-"`
}

func Init() *Code {
	return &Code{
		Path: Config.Get.Basic.Workspace + "/files/",
	}
}

func (code *Code) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(code.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		if !file.IsDir() {
			item.Path = file.Name()
			timeData, _ := times.Stat(code.Path + file.Name())
			item.Time = timeData.BirthTime()
			if strings.Contains(file.Name(), "demo") {
				item.Time = item.Time.Add(-10 * time.Minute)
			}
		}
	}

	return data
}
