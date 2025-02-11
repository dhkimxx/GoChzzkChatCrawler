package main

import (
	"fmt"

	"github.com/dhkimxx/GoChzzkChatCrawler/crawler"
)

func main() {
	streamerId := "75cbf189b3bb8f9f687d2aca0d0a382b"

	crawlerClient := crawler.NewCrawlerClient(streamerId, 1, func(msg crawler.ChzzkChatMessage) {
		fmt.Printf("[%d] %s: %s\n", msg.Timestamp, msg.Nickname, msg.Content)
	})

	err := crawlerClient.Run()
	if err != nil {
		panic(err)
	}
}
