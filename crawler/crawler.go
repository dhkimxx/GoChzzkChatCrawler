package crawler

import (
	"chzzk/api"
	"chzzk/command"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const WS_URL string = "wss://kr-ss1.chat.naver.com/chat"

type ChzzkChatCrawler struct {
	StreamerID string
	ChatChan   chan ChzzkChatMessage
}

type ChzzkChatMessage struct {
	StreamerID string
	UserHashID string
	Nickname   string
	Content    string
	Timestamp  int64
}

type wsConnRequest struct {
	Bdy   wsConnRequestBody `json:"bdy"`
	Cid   string            `json:"cid"`
	Cmd   command.ChatCmd   `json:"cmd"`
	Svcid string            `json:"svcid"`
	Tid   int               `json:"tid"`
	Ver   string            `json:"ver"`
}

type wsConnRequestBody struct {
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

type wsConnInitialResponse struct {
	Bdy     wsConnInitialResponseBody `json:"bdy"`
	Cid     string                    `json:"cid"`
	Cmd     command.ChatCmd           `json:"cmd"`
	Svcid   string                    `json:"svcid"`
	Tid     string                    `json:"tid"`
	RetCode int                       `json:"retCode"`
	ResMsg  string                    `json:"retMsg"`
}

type wsConnInitialResponseBody struct {
	AccTkn string      `json:"accTkn"`
	Auth   string      `json:"auth"`
	UUID   interface{} `json:"uuid"`
	SID    string      `json:"sid"`
}

type wsResponse struct {
	Bdy   []wsResponseBody `json:"bdy"`
	Cid   string           `json:"cid"`
	Cmd   command.ChatCmd  `json:"cmd"`
	Svcid string           `json:"svcid"`
	Tid   string           `json:"tid"`
	Ver   string           `json:"ver"`
}

type wsResponseBody struct {
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

func NewCrawler(streamerID string, bufferSize int) ChzzkChatCrawler {
	newCrawler := ChzzkChatCrawler{
		StreamerID: streamerID,
		ChatChan:   make(chan ChzzkChatMessage, bufferSize),
	}
	return newCrawler
}

func (crawler *ChzzkChatCrawler) Run() error {
	defer close(crawler.ChatChan)
	cid, err := api.FetchChatChannelID(crawler.StreamerID)
	if err != nil {
		return fmt.Errorf("failed to fetch chat channel id: %s", err)
	}

	accTkn, err := api.FetchAccessToken(cid)
	if err != nil {
		return fmt.Errorf("failed to fetch access token: %s", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(WS_URL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := wsConnRequest{
		Ver:   "3",
		Cmd:   100,
		Svcid: "game",
		Cid:   cid,
		Bdy: wsConnRequestBody{
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
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		return err
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	var initRes wsConnInitialResponse
	err = json.Unmarshal([]byte(message), &initRes)
	if err != nil {
		return err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)

		}
		var res wsResponse
		err = json.Unmarshal([]byte(message), &res)
		if err != nil {
			return err
		}

		if res.Cmd == command.Ping {
			json, err := json.Marshal(command.PongInstance)
			if err != nil {
				return err
			}
			err = conn.WriteMessage(websocket.TextMessage, json)
			if err != nil {
				return err
			}
		}
		for _, body := range res.Bdy {
			var profile map[string]interface{}
			err = json.Unmarshal([]byte(body.Profile), &profile)
			if err != nil {
				return err
			}

			msg := ChzzkChatMessage{
				StreamerID: crawler.StreamerID,
				Content:    body.Msg,
				Nickname:   profile["nickname"].(string),
				Timestamp:  body.MsgTime,
				UserHashID: profile["userIdHash"].(string),
			}
			crawler.ChatChan <- msg
		}
	}
}
