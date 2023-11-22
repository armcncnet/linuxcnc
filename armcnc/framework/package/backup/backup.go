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
	"sort"
	"strings"
	"time"
)

type Backup struct {
	Path string `json:"path"`
}

type Data struct {
	Name string    `json:"name"`
	Path string    `json:"path"`
	Date string    `json:"date"`
	Time time.Time `json:"-"`
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
		if strings.Contains(file.Name(), ".zip") {
			item := Data{}
			item.Name = file.Name()
			item.Path = file.Name()
			timeData, _ := times.Stat(backup.Path + file.Name())
			item.Time = timeData.BirthTime()
			item.Date = item.Time.Format("2006-01-02 15:04:05")
			data = append(data, item)
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Time.After(data[j].Time)
	})

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

func (backup *Backup) Delete(fileName string) bool {
	status := false
	exists, _ := FileUtils.PathExists(backup.Path + fileName)
	if exists {
		err := os.RemoveAll(backup.Path + fileName)
		if err == nil {
			status = true
		}
	}
	return status
}
