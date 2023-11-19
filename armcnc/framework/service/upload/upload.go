/**
 ******************************************************************************
 * @file    upload.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package UploadService

import (
	"armcnc/framework/config"
	"armcnc/framework/utils"
	"armcnc/framework/utils/file"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"time"
)

func UploadMachine(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	timestamp := time.Now().Format("20060102150405")
	ext := filepath.Ext(file.Filename)
	newFileName := timestamp + ext
	filePath := Config.Get.Basic.Workspace + "/uploads/" + newFileName
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	zip := FileUtils.Unzip(filePath, Config.Get.Basic.Workspace+"/configs/"+timestamp+"/")
	if !zip {
		os.Remove(filePath)
		os.Remove(Config.Get.Basic.Workspace + "/configs/" + timestamp)
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	time.Sleep(1 * time.Second)
	os.Remove(filePath)

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}