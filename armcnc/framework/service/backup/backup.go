/**
 ******************************************************************************
 * @file    backup.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package BackupService

import (
	"armcnc/framework/package/backup"
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
)

type responseSelect struct {
	Backup []BackupPackage.Data `json:"backup"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Backup = make([]BackupPackage.Data, 0)

	backup := BackupPackage.Init()
	returnData.Backup = backup.Select()

	Utils.Success(c, 0, "", returnData)
	return
}
