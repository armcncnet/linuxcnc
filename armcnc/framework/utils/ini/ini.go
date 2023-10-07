/**
 ******************************************************************************
 * @file    ini.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package IniUtils

import "gopkg.in/ini.v1"

func Empty() *ini.File {
	iniFile := ini.Empty()
	return iniFile
}

func ReflectFrom(cfg *ini.File, v interface{}) error {
	err := ini.ReflectFrom(cfg, v)
	return err
}

func Load(source interface{}, others ...interface{}) (*ini.File, error) {
	iniFile, err := ini.Load(source, others)
	return iniFile, err
}

func SaveTo(cfg *ini.File, filename string) error {
	err := cfg.SaveTo(filename)
	return err
}

func MapTo(cfg *ini.File, v interface{}) error {
	err := cfg.MapTo(v)
	return err
}
