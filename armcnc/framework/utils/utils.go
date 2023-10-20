/**
 ******************************************************************************
 * @file    utils.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Utils

import (
	"github.com/gin-gonic/gin"
	"net"
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

func GetWlanIP() string {
	ip := "0.0.0.0"
	iface, err := net.InterfaceByName("wlan0")
	if err == nil {
		addr, err := iface.Addrs()
		if err == nil {
			for _, item := range addr {
				ip = item.String()
			}
		}
	}
	return ip
}
