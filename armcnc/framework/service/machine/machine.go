/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachineService

import (
	"armcnc/framework/config"
	"armcnc/framework/package/launch"
	"armcnc/framework/package/machine"
	"armcnc/framework/utils"
	"armcnc/framework/utils/file"
	"armcnc/framework/utils/socket"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type responseSelect struct {
	Machine []MachinePackage.Data `json:"machine"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Machine = make([]MachinePackage.Data, 0)

	machine := MachinePackage.Init()
	returnData.Machine = machine.Select()

	Utils.Success(c, 0, "", returnData)
	return
}

type responseGet struct {
	IsDefault bool                `json:"is_default"`
	Path      string              `json:"path"`
	User      MachinePackage.USER `json:"user"`
	Ini       MachinePackage.INI  `json:"ini"`
	Table     string              `json:"table"`
	Launch    string              `json:"launch"`
	Hal       string              `json:"hal"`
	Xml       string              `json:"xml"`
}

func Get(c *gin.Context) {
	returnData := responseGet{}

	path := c.DefaultQuery("path", "")

	machine := MachinePackage.Init()
	returnData.IsDefault = false
	if path != "" {
		if strings.Contains(path, "default_") {
			returnData.IsDefault = true
		}
		returnData.Path = path
		returnData.User = machine.GetUser(path)
		if returnData.User.Base.Name == "" {
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		returnData.Ini = machine.GetIni(path)
		returnData.Table = machine.GetTable(path)
		returnData.Launch = machine.GetLaunch(path)
		returnData.Hal = machine.GetHal(path)
		returnData.Xml = machine.GetXml(path)
	}

	if path == "" {
		returnData.Path = ""
		returnData.User = machine.DefaultUser(returnData.User)
		returnData.Ini = machine.DefaultIni(returnData.Ini)
		returnData.Table = ""
		returnData.Launch = ""
		returnData.Hal = ""
		returnData.Xml = ""
	}

	Utils.Success(c, 0, "", returnData)
	return
}

type requestUpdate struct {
	Path   string                  `json:"path"`
	User   MachinePackage.UserJson `json:"user"`
	Ini    MachinePackage.IniJson  `json:"ini"`
	Table  string                  `json:"table"`
	Launch string                  `json:"launch"`
	Hal    string                  `json:"hal"`
	Xml    string                  `json:"xml"`
}

func Update(c *gin.Context) {

	requestJson := requestUpdate{}
	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()

	if requestJson.Path == "" {
		requestJson.Path = time.Now().Format("20060102150405")
		mkdir, _ := FileUtils.PathMkdirAll(machine.Path + requestJson.Path + "/launch")
		if !mkdir {
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		writeHal := FileUtils.WriteFile("", machine.Path+requestJson.Path+"/machine.hal")
		if writeHal != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		writeIni := FileUtils.WriteFile("", machine.Path+requestJson.Path+"/machine.ini")
		if writeIni != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		writePosition := FileUtils.WriteFile("", machine.Path+requestJson.Path+"/machine.position")
		if writePosition != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		writeTbl := FileUtils.WriteFile("", machine.Path+requestJson.Path+"/machine.tbl")
		if writeTbl != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		writeUser := FileUtils.WriteFile("", machine.Path+requestJson.Path+"/machine.user")
		if writeUser != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		writeVar := FileUtils.WriteFile("", machine.Path+requestJson.Path+"/machine.var")
		if writeVar != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		writeXml := FileUtils.WriteFile("", machine.Path+requestJson.Path+"/machine.xml")
		if writeXml != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		launch := "#!/usr/bin/env python\n# -*- coding: utf-8 -*-\n\nimport armcnc as sdk\n\ndef armcnc_start(cnc):\n    while True:\n        pass\n\ndef armcnc_message(cnc, message):\n    pass\n\ndef armcnc_exit(cnc):\n    pass\n\nif __name__ == '__main__':\n    sdk.Init()"
		writeLaunch := FileUtils.WriteFile(launch, machine.Path+requestJson.Path+"/launch/launch.py")
		if writeLaunch != nil {
			os.RemoveAll(machine.Path + requestJson.Path)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
	}

	updateIni := machine.UpdateIni(requestJson.Path, requestJson.Ini)
	if !updateIni {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	updateUser := machine.UpdateUser(requestJson.Path, requestJson.User)
	if !updateUser {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	updateTable := machine.UpdateTable(requestJson.Path, requestJson.Table)
	if !updateTable {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	if requestJson.Path == Config.Get.Machine.Path {
		message := SocketUtils.SocketMessageFormat{}
		message.Command = "service:package:status"
		message.Data = struct {
			Package string `json:"package"`
			Status  string `json:"status"`
		}{Package: "handwheel", Status: requestJson.User.HandWheel.Status}
		SocketUtils.SendMessage(message.Command, message.Message, message.Data)
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

type requestUpdateLaunch struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func UpdateLaunch(c *gin.Context) {

	requestJson := requestUpdateLaunch{}
	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	if requestJson.Path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	update := machine.UpdateLaunch(requestJson.Path, requestJson.Content)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

type requestUpdateHal struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func UpdateHal(c *gin.Context) {

	requestJson := requestUpdateHal{}
	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	if requestJson.Path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	update := machine.UpdateHal(requestJson.Path, requestJson.Content)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

type requestUpdateXml struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func UpdateXml(c *gin.Context) {

	requestJson := requestUpdateXml{}
	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	if requestJson.Path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	update := machine.UpdateXml(requestJson.Path, requestJson.Content)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

type responseDownload struct {
	File string `json:"file"`
}

func Download(c *gin.Context) {

	returnData := responseDownload{}

	path := c.DefaultQuery("path", "")
	if path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	exists, _ := FileUtils.PathExists(machine.Path + path)
	if !exists {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	fileName := "machine_" + path + ".zip"
	filePath := Config.Get.Basic.Runtime + "/" + fileName
	zip := FileUtils.ZipFile(machine.Path+path+"/", filePath)
	if !zip {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	returnData.File = fileName

	Utils.Success(c, 0, "", returnData)
	return
}

func Delete(c *gin.Context) {

	path := c.DefaultQuery("path", "")
	if path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	status := machine.Delete(path)
	if !status {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

func Default(c *gin.Context) {

	path := c.DefaultQuery("path", "")
	if path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	if path != Config.Get.Machine.Path {
		machine := MachinePackage.Init()
		check := machine.GetIni(path)
		if check.Emc.Version == "" {
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		Config.Get.Machine.Path = path
		save := Config.Update()
		if !save {
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		launch := LaunchPackage.Init()
		launch.Start(Config.Get.Machine.Path)
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}
