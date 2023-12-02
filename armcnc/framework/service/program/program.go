/**
 ******************************************************************************
 * @file    program.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package ProgramService

import (
	"armcnc/framework/package/program"
	"armcnc/framework/utils"
	"armcnc/framework/utils/file"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type responseSelect struct {
	Code []ProgramPackage.Data `json:"code"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Code = make([]ProgramPackage.Data, 0)

	program := ProgramPackage.Init()
	returnData.Code = program.Select()

	Utils.Success(c, 0, "", returnData)
	return
}

type responseReadLine struct {
	Line []string `json:"lines"`
}

func ReadLine(c *gin.Context) {

	returnData := responseReadLine{}
	returnData.Line = make([]string, 0)

	fileName := c.DefaultQuery("file_name", "")
	if fileName == "" {
		Utils.Error(c, 10000, "", returnData)
		return
	}

	program := ProgramPackage.Init()
	read := program.ReadLine(fileName)
	returnData.Line = read.Line

	Utils.Success(c, 0, "", returnData)
	return
}

type responseReadContent struct {
	IsDefault bool   `json:"is_default"`
	Content   string `json:"content"`
}

func ReadContent(c *gin.Context) {

	returnData := responseReadContent{}

	fileName := c.DefaultQuery("file_name", "")
	if fileName == "" {
		Utils.Error(c, 10000, "", returnData)
		return
	}

	program := ProgramPackage.Init()
	returnData.IsDefault = false
	if fileName == "armcnc.ngc" || fileName == "linuxcnc.ngc" {
		returnData.IsDefault = true
	}
	returnData.Content = program.ReadContent(fileName)

	Utils.Success(c, 0, "", returnData)
	return
}

type requestUpdateContent struct {
	FileName string `json:"file_name"`
	Content  string `json:"content"`
}

func UpdateContent(c *gin.Context) {

	requestJson := requestUpdateContent{}
	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	program := ProgramPackage.Init()

	if requestJson.FileName == "" {
		requestJson.FileName = time.Now().Format("20060102150405") + ".ngc"
		writeFile := FileUtils.WriteFile(requestJson.Content+"\n", program.Path+requestJson.FileName)
		if writeFile != nil {
			os.RemoveAll(program.Path + requestJson.FileName)
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
	} else {
		if !strings.Contains(requestJson.Content, "\n") {
			requestJson.Content = requestJson.Content + "\n"
		}
		update := program.UpdateContent(requestJson.FileName, requestJson.Content)
		if !update {
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

func Delete(c *gin.Context) {

	fileName := c.DefaultQuery("file_name", "")
	if fileName == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	program := ProgramPackage.Init()
	status := program.Delete(fileName)
	if !status {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}
