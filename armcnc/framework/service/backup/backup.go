/**
 ******************************************************************************
 * @file    backup.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package BackupService

import (
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
)

type responseSelect struct {
}

func Select(c *gin.Context) {

	returnData := responseSelect{}

	Utils.Success(c, 0, "", returnData)
	return
}
