package api_test

import (
	"chzzk/api"
	"testing"
)

func TestFetchChatChannelID(t *testing.T) {
	cid, err := api.FetchChatChannelID("7ce8032370ac5121dcabce7bad375ced")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cid)
}

func TestFetchAccessToken(t *testing.T) {
	accTkn, err := api.FetchAccessToken("N1V_SF")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accTkn)
}

func TestFetchChannelName(t *testing.T) {
	channel, err := api.FetchChannelName("7ce8032370ac5121dcabce7bad375ced")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(channel)
}
