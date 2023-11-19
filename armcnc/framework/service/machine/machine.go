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
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io/ioutil"
	"strings"
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
	if path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	returnData.IsDefault = false
	if strings.Contains(path, "default_") {
		returnData.IsDefault = true
	}
	returnData.Path = path
	returnData.User = machine.GetUser(path)
	returnData.Ini = machine.GetIni(path)
	returnData.Table = machine.GetTable(path)
	returnData.Launch = machine.GetLaunch(path)
	returnData.Hal = machine.GetHal(path)
	returnData.Xml = machine.GetXml(path)

	Utils.Success(c, 0, "", returnData)
	return
}

type responseNew struct {
	IsDefault bool                `json:"is_default"`
	Path      string              `json:"path"`
	User      MachinePackage.USER `json:"user"`
	Ini       MachinePackage.INI  `json:"ini"`
	Table     string              `json:"table"`
	Launch    string              `json:"launch"`
	Hal       string              `json:"hal"`
	Xml       string              `json:"xml"`
}

func New(c *gin.Context) {

	returnData := responseNew{}

	machine := MachinePackage.Init()
	returnData.IsDefault = false
	returnData.Path = ""
	returnData.User = machine.DefaultUser(returnData.User)
	returnData.Ini = machine.DefaultIni(returnData.Ini)
	returnData.Table = "T1 P1 D2.000 X0.000 Y0.000 Z0.000;"
	returnData.Launch = "#!/usr/bin/env python\n# -*- coding: utf-8 -*-\n\nimport armcnc as sdk\n\ndef armcnc_start(cnc):\n    while True:\n        pass\n\ndef armcnc_message(cnc, message):\n    pass\n\ndef armcnc_exit(cnc):\n    pass\n\nif __name__ == '__main__':\n    sdk.Init()"
	returnData.Hal = ""
	returnData.Xml = ""

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

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

type requestUpdateUser struct {
	Path string                  `json:"path"`
	User MachinePackage.UserJson `json:"user"`
}

func UpdateUser(c *gin.Context) {

	requestJson := requestUpdateUser{}
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
	update := machine.UpdateUser(requestJson.Path, requestJson.User)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
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
	filePath := Config.Get.Basic.Workspace + "/uploads/" + fileName
	zip := FileUtils.Zip(machine.Path+path+"/", filePath)
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

func SetCurrentMachine(c *gin.Context) {

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
		save := Config.Save()
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
