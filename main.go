package main

import (
	"fmt"

	"github.com/dhkimxx/GoChzzkChatCrawler/crawler"
)

func main() {

	defaultStreamerID := "c7ded8ea6b0605d3c78e18650d2df83b"

	crawlerClient := crawler.NewCrawlerClient(defaultStreamerID, 1, nil)

	go func() {
		crawlerClient.SetMessageHandler(func(msg crawler.ChzzkChatMessage) {
			fmt.Printf("[%d] %s: %s\n", msg.Timestamp, msg.Nickname, msg.Content)
		})
		err := crawlerClient.Run()
		if err != nil {
			panic(err)
		}
	}()

	for msg := range crawlerClient.ChatChan {
		fmt.Println("from chan:", msg)
	}
}
