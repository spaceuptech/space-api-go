package main

import (
	"github.com/spaceuptech/space-api-go"
	"fmt"
)

func main() {
	api, err := api.New("books-app", "localhost:8081", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	doc := map[string]interface{}{"name":"SomeBook"}
	resp, err := db.Insert("books").Doc(doc).Apply()
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
