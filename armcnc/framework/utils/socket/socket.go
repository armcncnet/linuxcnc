/**
 ******************************************************************************
 * @file    socket.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package SocketUtils

import "net/http"

import (
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"sync"
)

var SocketGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var SocketStruct = &Struct{}

type Struct struct {
	mutex  sync.Mutex
	User   map[*websocket.Conn]bool
	Status bool `json:"status"`
}

type SocketMessageFormat struct {
	Command string      `json:"command"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendMessage(command string, message string, data interface{}) {

	sendData := SocketMessageFormat{}
	sendData.Command = command
	sendData.Message = message
	sendData.Data = data
	messageJson, _ := json.Marshal(sendData)

	SocketStruct.mutex.Lock()
	defer SocketStruct.mutex.Unlock()

	for user := range SocketStruct.User {
		err := user.WriteMessage(websocket.TextMessage, messageJson)
		if err != nil {
			user.Close()
			delete(SocketStruct.User, user)
		}
	}
}
