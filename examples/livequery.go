package main

import (
	"github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/api/model"
	"fmt"
)

func main() {
	api, err := api.New("books-app", "localhost:8081", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	db.LiveQuery("books").Options(&model.LiveQueryOptions{ChangesOnly:false}).Subscribe(func(liveData *model.LiveData, changeType string, changedData *model.ChangedData) () {
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
	for {}
}

