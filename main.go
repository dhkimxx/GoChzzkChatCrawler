package main

import (
	"chzzk/api"
	"chzzk/command"
	"encoding/json"
	"fmt"
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
	Svcid         string `json:"svcid"`
	Cid           string `json:"cid"`
	MbrCnt        int    `json:"mbrCnt"`
	Uid           string `json:"uid"`
	Profile       string `json:"profile"`
	Msg           string `json:"msg"`
	MsgTypeCode   int    `json:"msgTypeCode"`
	MsgStatusType string `json:"msgStatusType"`
	Extras        string `json:"extras"`
	Ctime         int64  `json:"ctime"`
	Utime         int64  `json:"utime"`
	MsgTid        string `json:"msgTid"`
	MsgTime       int64  `json:"msgTime"`
}

func main() {

	url := "wss://kr-ss1.chat.naver.com/chat"

	defaultStreamerID := "f00f6d46ecc6d735b96ecf376b9e5212"

	cid, err := api.FetchChatChannelID(defaultStreamerID)
	if err != nil {
		panic(fmt.Errorf("failed to fetch chat channel id: %s", err))
	}

	accTkn, err := api.FetchAccessToken(cid)
	if err != nil {
		panic(fmt.Errorf("failed to fetch access token: %s", err))
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	request := Request{
		Ver:   "3",
		Cmd:   100,
		Svcid: "game",
		Cid:   cid,
		Bdy: RequestBody{
			Uid:      nil,
			DevType:  2001,
			AccTkn:   accTkn,
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
	log.Println(initRes)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message:", err)
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
		for _, body := range res.Bdy {
			var profile map[string]interface{}
			json.Unmarshal([]byte(body.Profile), &profile)
			log.Println(profile["nickname"], ":", body.Msg)
		}
	}
}
