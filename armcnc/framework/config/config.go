/**
 ******************************************************************************
 * @file    config.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Config

import (
	"armcnc/framework/utils/file"
	"github.com/gookit/color"
	"gopkg.in/ini.v1"
	"log"
)

var Get = &Data{}

type Data struct {
	Basic         DataBasic         `ini:"basic"`
	Authorization DataAuthorization `ini:"authorization"`
}

type DataBasic struct {
	Name      string `ini:"name"`
	Version   string `ini:"version"`
	Workspace string `ini:"workspace"`
}

type DataAuthorization struct {
	Getaway string `ini:"getaway"`
	AppId   string `ini:"app_id"`
	AppKey  string `ini:"app_key"`
}

func Init() {

	Get.Basic.Name = "armcnc"
	Get.Basic.Version = "1.0.0"
	Get.Basic.Workspace = "/opt/armcnc"

	Get.Authorization.Getaway = "https://getaway.geekros.com"

	exists, _ := FileUtils.PathExists(Get.Basic.Workspace + "/armcnc.ini")
	if !exists {
		iniFile := ini.Empty()
		err := ini.ReflectFrom(iniFile, Get)
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}

		err = iniFile.SaveTo(Get.Basic.Workspace + "/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}
	}

	if exists {
		iniFile, err := ini.Load(Get.Basic.Workspace + "/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration information read failed."))
			return
		}

		var intConfig Data
		err = iniFile.MapTo(&intConfig)
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration information read failed."))
			return
		}

		Get.Authorization = intConfig.Authorization
	}
}
