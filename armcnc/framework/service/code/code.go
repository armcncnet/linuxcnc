/**
 ******************************************************************************
 * @file    code.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package CodeService

import (
	"armcnc/framework/package/code"
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
)

type responseSelect struct {
	Code []CodePackage.Data `json:"code"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Code = make([]CodePackage.Data, 0)

	code := CodePackage.Init()
	returnData.Code = code.Select()

	Utils.Success(c, 0, "", returnData)
	return
}
