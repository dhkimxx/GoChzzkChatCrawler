package api

import (
	"chzzk/url"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
)

var HEADERS = map[string]string{
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

type httpClient struct {
	header                 map[string]string
	acceptedStatusCodeList []int
}

func NewHttpClient() *httpClient {
	hc := &httpClient{
		header:                 HEADERS,
		acceptedStatusCodeList: []int{200},
	}
	return hc
}

func (hc *httpClient) sendRequest(url, method string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range hc.header {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !slices.Contains(hc.acceptedStatusCodeList, resp.StatusCode) {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FetchLiveChannelIdOfStreamer(streamerId string) (string, error) {

	url, method := url.GetLiveChannelIdOfStreamer(streamerId)
	body, err := NewHttpClient().sendRequest(url, method)
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

func FetchChennelNameOfStreamer(streamerId string) (string, error) {
	url, method := url.GetChennelNameOfStreamer(streamerId)
	body, err := NewHttpClient().sendRequest(url, method)
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

func FetchChatAccessToken(chatChannelId string) (string, error) {
	url, method := url.GetChatAccessToken(chatChannelId)
	body, err := NewHttpClient().sendRequest(url, method)
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
