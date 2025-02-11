package api_test

import (
	"chzzk/api"
	"testing"
)

func TestFetchChatChannelID(t *testing.T) {
	cid, err := api.FetchLiveChannelIdOfStreamer("7ce8032370ac5121dcabce7bad375ced")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cid)
}

func TestFetchAccessToken(t *testing.T) {
	accTkn, err := api.FetchChatAccessToken("N1V_SF")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accTkn)
}

func TestFetchChannelName(t *testing.T) {
	channel, err := api.FetchChennelNameOfStreamer("7ce8032370ac5121dcabce7bad375ced")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(channel)
}
