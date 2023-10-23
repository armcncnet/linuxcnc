/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachinePackage

import (
	"armcnc/framework/config"
	"armcnc/framework/utils/file"
	"armcnc/framework/utils/ini"
	"github.com/djherbis/times"
	"os"
	"sort"
	"strings"
	"time"
)

type Machine struct {
	Path string `json:"path"`
}

type Data struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Describe    string    `json:"describe"`
	Version     string    `json:"version"`
	ControlType int       `json:"control_type"`
	Time        time.Time `json:"-"`
}

type INI struct {
	EMC struct {
		MACHINE      string `ini:"MACHINE"`
		DESCRIBE     string `ini:"DESCRIBE"`
		CONTROL_TYPE int    `ini:"CONTROL_TYPE"`
		DEBUG        string `ini:"DEBUG"`
		VERSION      string `ini:"VERSION"`
	} `ini:"EMC"`
}

func Init() *Machine {
	return &Machine{
		Path: Config.Get.Basic.Workspace + "/configs/",
	}
}

func (machine *Machine) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(machine.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		if file.IsDir() {
			item.Path = file.Name()
			timeData, _ := times.Stat(machine.Path + file.Name())
			item.Time = timeData.BirthTime()
			if strings.Contains(file.Name(), "default_") {
				item.Time = item.Time.Add(-10 * time.Minute)
			}
			info := machine.Get(file.Name())
			if info.EMC.VERSION != "" {
				item.Name = info.EMC.MACHINE
				item.Describe = info.EMC.DESCRIBE
				item.Version = info.EMC.VERSION
				item.ControlType = info.EMC.CONTROL_TYPE
				data = append(data, item)
			}
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Time.After(data[j].Time)
	})

	return data
}

func (machine *Machine) Get(path string) INI {
	data := INI{}
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + path + "/machine.ini")
		if err == nil {
			err = IniUtils.MapTo(iniFile, &data)
		}
	}
	return data
}

func (machine *Machine) GetContent(path string) string {
	content := ""
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		contentByte, err := FileUtils.ReadFile(machine.Path + path + "/machine.ini")
		if err == nil {
			content = string(contentByte)
		}
	}

	return content
}

func (machine *Machine) UpdateContent(path string, content string) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		write := FileUtils.WriteFile(content, machine.Path+path+"/machine.ini")
		if write == nil {
			status = true
		}
	}
	return status
}
