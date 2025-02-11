package url

import "fmt"

const (
	CHZZK_BASE_URL = "https://api.chzzk.naver.com"
	NAVER_GAME_URL = "https://comm-api.game.naver.com"
)

func GetLiveChannelIdOfStreamer(streamerId string) (url, method string) {
	return fmt.Sprintf("%s/polling/v2/channels/%s/live-status", CHZZK_BASE_URL, streamerId), "GET"
}

func GetChennelNameOfStreamer(streamerId string) (url, method string) {
	return fmt.Sprintf("%s/service/v1/channels/%s", CHZZK_BASE_URL, streamerId), "GET"
}

func GetChatAccessToken(chatChannelId string) (url, method string) {
	return fmt.Sprintf("%s/nng_main/v1/chats/access-token?channelId=%s&chatType=STREAMING", NAVER_GAME_URL, chatChannelId), "GET"
}
