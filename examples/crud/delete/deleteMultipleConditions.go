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
	condition1 := utils.Cond("id", "==", 1)
	condition2 := utils.Cond("id", "==", 2)
	condition := utils.Or(condition1, condition2)
	resp, err := db.Delete("books").Where(condition).Apply()
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
