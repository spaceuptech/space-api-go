package main

import (
	"fmt"

	"github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/utils"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	condition := utils.Cond("id", "==", 1)
	set := map[string]interface{}{"name":"ABook"}
	resp, err := db.UpdateOne("books").Where(condition).Set(set).Apply()
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
