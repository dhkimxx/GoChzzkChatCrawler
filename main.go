package main

import (
	"chzzk/crawler"
	"fmt"
)

func main() {

	defaultStreamerID := "4d39d99252f247f06de349ccc0d444a7"

	crawler := crawler.NewCrawler(defaultStreamerID, 1)

	go func() {
		err := crawler.Run()
		if err != nil {
			panic(err)
		}
	}()

	for msg := range crawler.ChatChan {
		fmt.Println(msg)
	}
}
