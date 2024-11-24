package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var HEADERS = map[string]string{
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

func FetchChatChannelID(streamerID string) (string, error) {
	url := fmt.Sprintf("https://api.chzzk.naver.com/polling/v2/channels/%s/live-status", streamerID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Content struct {
			ChatChannelID string `json:"chatChannelId"`
		} `json:"content"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if data.Content.ChatChannelID == "" {
		return "", fmt.Errorf("chatChannelId is empty")
	}

	return data.Content.ChatChannelID, nil
}

func FetchChannelName(streamerID string) (string, error) {
	url := fmt.Sprintf("https://api.chzzk.naver.com/service/v1/channels/%s", streamerID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Content struct {
			ChannelName string `json:"channelName"`
		} `json:"content"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Content.ChannelName, nil
}

func FetchAccessToken(chatChannelID string) (string, error) {
	url := fmt.Sprintf("https://comm-api.game.naver.com/nng_main/v1/chats/access-token?channelId=%s&chatType=STREAMING", chatChannelID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Content struct {
			AccessToken string `json:"accessToken"`
		} `json:"content"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Content.AccessToken, nil
}
