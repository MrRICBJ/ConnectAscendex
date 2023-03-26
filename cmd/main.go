package main

import (
	"fmt"
	project "test"
	"test/internal"
)

func main() {
	client := internal.NewAPIClient()

	if err := client.Connection(); err != nil {
		panic(err)
	}
	defer client.Disconnect()

	ch := make(chan project.BestOrderBook, 10)
	if err := client.SubscribeToChannel("USDT_BTC"); err != nil {
		panic(err)
	}

	go client.ReadMessagesFromChannel(ch)

	for book := range ch {
		fmt.Println(book)
	}
}
