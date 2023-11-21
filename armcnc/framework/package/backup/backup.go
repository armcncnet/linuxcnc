/**
 ******************************************************************************
 * @file    backup.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package BackupPackage

import (
	"armcnc/framework/config"
	"armcnc/framework/utils/file"
	"github.com/djherbis/times"
	"os"
	"time"
)

type Backup struct {
	Path string `json:"path"`
}

type Data struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Data string `json:"data"`
}

func Init() *Backup {
	return &Backup{
		Path: Config.Get.Basic.Workspace + "/backups/",
	}
}

func (backup *Backup) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(backup.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		item.Name = file.Name()
		item.Path = file.Name()
		timeData, _ := times.Stat(backup.Path + file.Name())
		createTime := timeData.BirthTime()
		item.Data = createTime.Format("2006-01-02 15:04:05")
		data = append(data, item)
	}

	return data
}

func (backup *Backup) Pack(Type string) bool {
	status := true

	files := make([]string, 0)

	if Type == "all" {
		files = append(files, Config.Get.Basic.Workspace+"/configs")
		files = append(files, Config.Get.Basic.Workspace+"/plugins")
		files = append(files, Config.Get.Basic.Workspace+"/programs")
		files = append(files, Config.Get.Basic.Workspace+"/scripts")
	}

	if Type == "machine" {
		files = append(files, Config.Get.Basic.Workspace+"/configs")
	}

	if Type == "program" {
		files = append(files, Config.Get.Basic.Workspace+"/programs")
	}

	if Type == "plugin" {
		files = append(files, Config.Get.Basic.Workspace+"/plugins")
	}

	if Type == "script" {
		files = append(files, Config.Get.Basic.Workspace+"/scripts")
	}

	if len(files) > 0 {
		fileName := Type + "_" + time.Now().Format("20060102150405") + ".zip"
		pack := FileUtils.ZipFiles(files, backup.Path+fileName)
		if !pack {
			status = false
		}
	}

	return status
}
