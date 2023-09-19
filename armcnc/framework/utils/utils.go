/**
 ******************************************************************************
 * @file    utils.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmptyData struct {
}

func Success(c *gin.Context, code int, msg string, data interface{}) {

	if msg == "" {
		msg = "success"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"data":    data,
	})

	c.Set("code", code)
}

func Error(c *gin.Context, code int, msg string, data interface{}) {

	if msg == "" {
		msg = "error"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"data":    data,
	})
}
