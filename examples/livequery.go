package main

import (
	"fmt"
	"time"

	"github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/model"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	subscription := db.LiveQuery("books").Options(&model.LiveQueryOptions{ChangesOnly: false}).
	Subscribe(func(liveData *model.LiveData, changeType string, changedData *model.ChangedData) () {
		fmt.Println("type", changeType)
		var v []interface{}
		liveData.Unmarshal(&v)
		fmt.Println("data", v)
		var v2 interface{}
		changedData.Unmarshal(&v2)
		fmt.Println("chagned", v2)
		fmt.Println()
	}, func(err error) () {
		fmt.Println(err)
	})
	time.Sleep(1000)
	var v []interface{}
	subscription.GetSnapshot().Unmarshal(&v)
	fmt.Println(v)
	for {}
	subscription.Unsubscribe()
}

