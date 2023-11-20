/**
 ******************************************************************************
 * @file    program.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package ProgramPackage

import (
	"armcnc/framework/config"
	"armcnc/framework/utils/file"
	"bufio"
	"github.com/djherbis/times"
	"github.com/goccy/go-json"
	"os"
	"sort"
	"strings"
	"time"
)

type Program struct {
	Path string `json:"path"`
}

type Data struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Describe string    `json:"describe"`
	Version  string    `json:"version"`
	Line     []string  `json:"line"`
	Content  string    `json:"content"`
	Time     time.Time `json:"-"`
}

func Init() *Program {
	return &Program{
		Path: Config.Get.Basic.Workspace + "/files/",
	}
}

func (program *Program) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(program.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		if !file.IsDir() {
			item.Path = file.Name()
			timeData, _ := times.Stat(program.Path + file.Name())
			item.Time = timeData.BirthTime()
			if file.Name() == "armcnc.ngc" || file.Name() == "demo.ngc" || file.Name() == "linuxcnc.ngc" {
				item.Time = item.Time.Add(-525600 * time.Minute)
			}
			firstLine := program.ReadFirstLine(file.Name())
			if firstLine.Version != "" {
				item.Name = firstLine.Name
				item.Describe = firstLine.Describe
				item.Version = firstLine.Version
				item.Line = make([]string, 0)
				item.Content = ""
				data = append(data, item)
			}
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Time.After(data[j].Time)
	})
	return data
}

func (program *Program) ReadContent(fileName string) string {
	content := ""
	exists, _ := FileUtils.PathExists(program.Path + fileName)
	if exists {
		contentByte, err := FileUtils.ReadFile(program.Path + fileName)
		if err == nil {
			content = string(contentByte)
		}
	}
	return content
}

func (program *Program) ReadLine(fileName string) Data {
	data := Data{}
	data.Line = make([]string, 0)

	file, err := os.Open(program.Path + fileName)
	if err != nil {
		return data
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data.Line = append(data.Line, line)
	}
	return data
}

func (program *Program) ReadFirstLine(fileName string) Data {
	data := Data{}

	file, err := os.Open(program.Path + fileName)
	if err != nil {
		return data
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return data
	}

	line = strings.TrimSpace(line)
	if len(line) > 0 && line[0] == '(' && line[len(line)-1] == ')' {
		jsonStr := line[1 : len(line)-1]
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			return data
		}
	}
	return data
}

func (program *Program) UpdateContent(fileName string, content string) bool {
	status := false
	exists, _ := FileUtils.PathExists(program.Path + fileName)
	if exists {
		write := FileUtils.WriteFile(content, program.Path+fileName)
		if write == nil {
			status = true
		}
	}
	return status
}
