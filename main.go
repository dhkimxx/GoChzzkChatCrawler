package main

import (
	"chzzk/crawler"
	"fmt"
)

func main() {

	defaultStreamerID := "e9210c823439e7add5006cb4b93826a0"

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
