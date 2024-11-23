package main

import (
	"chzzk/command"
	"chzzk/config"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Request struct {
	Bdy   RequestBody     `json:"bdy"`
	Cid   string          `json:"cid"`
	Cmd   command.ChatCmd `json:"cmd"`
	Svcid string          `json:"svcid"`
	Tid   int             `json:"tid"`
	Ver   string          `json:"ver"`
}

type RequestBody struct {
	Uid      interface{} `json:"uid"`
	DevType  int         `json:"devType"`
	AccTkn   string      `json:"accTkn"`
	Auth     string      `json:"auth"`
	LibVer   string      `json:"libVer"`
	OsVer    string      `json:"osVer"`
	DevName  string      `json:"devName"`
	Locale   string      `json:"locale"`
	Timezone string      `json:"timezone"`
}

type InitialResponse struct {
	Bdy     InitialResponseBody `json:"bdy"`
	Cid     string              `json:"cid"`
	Cmd     command.ChatCmd     `json:"cmd"`
	Svcid   string              `json:"svcid"`
	Tid     string              `json:"tid"`
	RetCode int                 `json:"retCode"`
	ResMsg  string              `json:"retMsg"`
}

type InitialResponseBody struct {
	AccTkn string      `json:"accTkn"`
	Auth   string      `json:"auth"`
	UUID   interface{} `json:"uuid"`
	SID    string      `json:"sid"`
}

type Response struct {
	Bdy   []ResponseBody  `json:"bdy"`
	Cid   string          `json:"cid"`
	Cmd   command.ChatCmd `json:"cmd"`
	Svcid string          `json:"svcid"`
	Tid   string          `json:"tid"`
	Ver   string          `json:"ver"`
}

type ResponseBody struct {
	Svcid         string  `json:"svcid"`
	Cid           string  `json:"cid"`
	MbrCnt        int     `json:"mbrCnt"`
	Uid           string  `json:"uid"`
	Profile       string  `json:"profile"`
	Msg           string  `json:"msg"`
	MsgTypeCode   int     `json:"msgTypeCode"`
	MsgStatusType string  `json:"msgStatusType"`
	Extras        string  `json:"extras"`
	Ctime         int64   `json:"ctime"`
	Utime         int64   `json:"utime"`
	MsgTid        *string `json:"msgTid"`
	MsgTime       int64   `json:"msgTime"`
}

func main() {
	url := "wss://kr-ss1.chat.naver.com/chat"

	// WebSocket 연결
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("WebSocket 연결 실패:", err)
	}
	defer conn.Close()

	// 전송할 메시지 구조체 정의
	request := Request{
		Ver:   "3",
		Cmd:   100,
		Svcid: "game",
		Cid:   "N1V_SF",
		Bdy: RequestBody{
			Uid:      nil,
			DevType:  2001,
			AccTkn:   config.Config.AccessToken,
			Auth:     "READ",
			LibVer:   "4.9.3",
			OsVer:    "Linux/",
			DevName:  "Google Chrome/131.0.0.0",
			Locale:   "en-US",
			Timezone: "Asia/Seoul",
		},
		Tid: 1,
	}
	jsonMessage, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		panic(err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}

	var initRes InitialResponse
	err = json.Unmarshal([]byte(message), &initRes)
	if err != nil {
		panic(err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("메시지 읽기 오류:", err)
			continue
		}
		var res Response
		json.Unmarshal([]byte(message), &res)

		if res.Cmd == command.Ping {
			json, err := json.Marshal(command.PongInstance)
			if err != nil {
				panic(err)
			}
			err = conn.WriteMessage(websocket.TextMessage, json)
			if err != nil {
				panic(err)
			}
		}

		for i, body := range res.Bdy {
			var profile map[string]interface{}
			json.Unmarshal([]byte(body.Profile), &profile)
			log.Println(i, profile["nickname"], body.Msg)
		}

	}
}
