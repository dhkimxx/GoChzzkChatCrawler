package crawler

import (
	"encoding/json"
	"fmt"

	"github.com/dhkimxx/GoChzzkChatCrawler/api"
	"github.com/dhkimxx/GoChzzkChatCrawler/command"
	"github.com/dhkimxx/GoChzzkChatCrawler/url"

	"github.com/gorilla/websocket"
)

type ChzzkChatCrawler struct {
	StreamerId string
	ChatChan   chan ChzzkChatMessage
	onMessage  func(ChzzkChatMessage)
}

type ChzzkChatMessage struct {
	StreamerId string
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

func NewCrawlerClient(streamerID string, bufferSize int, onMessage func(ChzzkChatMessage)) ChzzkChatCrawler {
	newCrawler := ChzzkChatCrawler{
		StreamerId: streamerID,
		ChatChan:   make(chan ChzzkChatMessage, bufferSize),
		onMessage:  onMessage,
	}
	return newCrawler
}

func (crawler *ChzzkChatCrawler) SetMessageHandler(onMessage func(ChzzkChatMessage)) {
	crawler.onMessage = onMessage
}

func (crawler *ChzzkChatCrawler) Run() error {
	defer close(crawler.ChatChan)
	cid, err := api.FetchLiveChannelIdOfStreamer(crawler.StreamerId)
	if err != nil {
		return fmt.Errorf("failed to fetch chat channel id: %s", err)
	}

	accTkn, err := api.FetchChatAccessToken(cid)
	if err != nil {
		return fmt.Errorf("failed to fetch access token: %s", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(url.ChatChannelWebSocket(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.CloseHandler()

	wsConnRequestJson, err := json.Marshal(wsConnRequest{
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
	})

	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, wsConnRequestJson)
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
		if err == nil {
			if res.Cmd == command.Ping {
				pongResjson, err := json.Marshal(command.PongInstance)
				if err != nil {
					return err
				}
				err = conn.WriteMessage(websocket.TextMessage, pongResjson)
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
					StreamerId: crawler.StreamerId,
					Content:    body.Msg,
					Nickname:   profile["nickname"].(string),
					Timestamp:  body.MsgTime,
					UserHashID: profile["userIdHash"].(string),
				}

				/* send to channel */
				go func() {
					crawler.ChatChan <- msg
				}()

				/* execute callback method of client */
				go func() {
					if crawler.onMessage != nil {
						crawler.onMessage(msg)
					}
				}()

			}
		}
	}
}
