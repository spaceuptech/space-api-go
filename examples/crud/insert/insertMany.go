package main

import (
	"github.com/spaceuptech/space-api-go"
	"fmt"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	docs := make([]map[string]interface{}, 2)
	docs[0] = map[string]interface{}{"name": "SomeBook"}
	docs[1] = map[string]interface{}{"name": "SomeOtherBook"}
	resp, err := db.Insert("books").Docs(docs).Apply()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if resp.Status == 200 {
			fmt.Println("Success")
		} else {
			fmt.Println("Error Processing Request:", resp.Error)
		}
	}
}
