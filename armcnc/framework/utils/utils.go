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
	"regexp"
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

func GetIP(name string) string {
	ip := "0.0.0.0"
	iface, err := net.InterfaceByName(name)
	if err == nil {
		addr, err := iface.Addrs()
		if err == nil {
			for _, address := range addr {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ip = ipNet.IP.String()
				}
			}
		}
	}
	return ip
}

func EmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
