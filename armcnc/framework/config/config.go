/**
 ******************************************************************************
 * @file    config.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Config

import (
	"armcnc/framework/utils/file"
	"armcnc/framework/utils/ini"
	"github.com/gookit/color"
	"log"
)

var Get = &Data{}

type Data struct {
	Basic         DataBasic         `ini:"basic"`
	Authorization DataAuthorization `ini:"authorization"`
	Machine       DataMachine       `ini:"machine"`
}

type DataBasic struct {
	Name      string `ini:"name"`
	Version   string `ini:"version"`
	Workspace string `ini:"workspace"`
	Port      int    `ini:"port"`
}

type DataAuthorization struct {
	Getaway string `ini:"getaway"`
	AppId   string `ini:"app_id"`
	AppKey  string `ini:"app_key"`
}

type DataMachine struct {
	Path string `ini:"path"`
}

func Init() {

	Get.Basic.Name = "armcnc"
	Get.Basic.Version = "1.0.0"
	Get.Basic.Workspace = "/opt/armcnc"
	Get.Basic.Port = 10081

	Get.Authorization.Getaway = "https://getaway.geekros.com"

	exists, _ := FileUtils.PathExists(Get.Basic.Workspace + "/armcnc.ini")
	if !exists {
		iniFile := IniUtils.Empty()
		err := IniUtils.ReflectFrom(iniFile, Get)
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}

		err = IniUtils.SaveTo(iniFile, Get.Basic.Workspace+"/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}
	}

	if exists {
		iniFile, err := IniUtils.Load(Get.Basic.Workspace + "/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration information load failed."))
			return
		}

		var iniConfig Data
		err = IniUtils.MapTo(iniFile, iniConfig)
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration information mapTo failed."))
			return
		}

		Get.Authorization = iniConfig.Authorization
		Get.Machine = iniConfig.Machine

		iniFile.Section("basic").Key("version").SetValue(Get.Basic.Version)
		err = IniUtils.SaveTo(iniFile, Get.Basic.Workspace+"/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}
	}
}
