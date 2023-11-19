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
	"github.com/gin-gonic/gin"
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

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}
