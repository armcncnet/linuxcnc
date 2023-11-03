/**
 ******************************************************************************
 * @file    request.go
 * @author  GEEKROS,  site:www.geekros.com
 ******************************************************************************
 */

package RequestUtils

import (
	"armcnc/framework/config"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type ResponseService struct {
	Code int                 `json:"code"`
	Data responseServiceData `json:"data"`
}

type responseServiceData struct {
	Token string `json:"token"`
}

func Service(path string, method string, parameters map[string]string, data map[string]string) (*http.Response, ResponseService, error) {

	responseData := ResponseService{}
	responseData.Code = 10000

	var bodyData []byte
	if data != nil {
		bodyData, _ = json.Marshal(data)
	}

	request, err := http.NewRequest(method, Config.Get.Authorization.Getaway+path, bytes.NewReader(bodyData))
	if err != nil {
		return nil, responseData, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Account-Token", Config.Get.Authorization.Token)

	query := request.URL.Query()
	if parameters != nil {
		for key, val := range parameters {
			query.Add(key, val)
		}
		request.URL.RawQuery = query.Encode()
	}

	client := &http.Client{}

	response, _ := client.Do(request)

	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &responseData)
	return response, responseData, err
}

func Upload(path string, filePath string, parameters map[string]string) (*http.Response, ResponseService, error) {

	responseData := ResponseService{}
	responseData.Code = 10000

	file, err := os.Open(filePath)
	if err != nil {
		return nil, responseData, err
	}
	defer file.Close()

	bodyData := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyData)

	for key, value := range parameters {
		_ = writer.WriteField(key, value)
	}

	fileField, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return nil, responseData, err
	}

	_, err = io.Copy(fileField, file)
	if err != nil {
		return nil, responseData, err
	}
	writer.Close()

	request, err := http.NewRequest("POST", Config.Get.Authorization.Getaway+path, bodyData)
	if err != nil {
		return nil, responseData, err
	}

	request.Header.Set("Content-type", writer.FormDataContentType())
	request.Header.Set("Account-Token", Config.Get.Authorization.Token)

	client := &http.Client{}

	response, _ := client.Do(request)

	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &responseData)
	return response, responseData, err
}
