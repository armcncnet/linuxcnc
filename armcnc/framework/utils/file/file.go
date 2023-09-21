/**
 ******************************************************************************
 * @file    file.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package FileUtils

import "os"

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathMkdir(path string) (bool, error) {
	err := os.Mkdir(path, 0666)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func PathMkdirAll(path string) (bool, error) {
	err := os.MkdirAll(path, 0666)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func WriteFile(data string, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return err
}
