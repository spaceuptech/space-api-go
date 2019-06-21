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
	db.LiveQuery("books").Subscribe(func(liveData *model.LiveData, changeType string) () {
		fmt.Println(changeType)
		var v []interface{}
		liveData.Unmarshal(&v)
		fmt.Println(v)
	}, func(err error) () {
		fmt.Println(err)
	})
	for {}
}

