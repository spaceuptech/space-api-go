package main

import (
	"fmt"
	"time"

	api "github.com/spaceuptech/space-api-go"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if err != nil {
		fmt.Println(err)
	}
	pubsub := api.Pubsub()
	subscription := pubsub.Subscribe("/subject/", func(subject string, msg interface{}) {
		fmt.Println("received", subject, msg)
	})
	for i := 0; i < 30; i++ {
		pubsub.Publish("/subject/a", i*7)
		time.Sleep(2000)
	}
	time.Sleep(2000)
	subscription.Unsubscribe()
}
