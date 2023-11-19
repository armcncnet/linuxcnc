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
	"fmt"
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

	now := time.Now()
	millis := now.UnixNano() / 1000000
	timestamp := now.Format("20060102150405") + fmt.Sprintf("%03d", millis%1000)
	ext := filepath.Ext(file.Filename)
	newFileName := timestamp + ext
	filepath := Config.Get.Basic.Workspace + "/uploads/" + newFileName

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}
