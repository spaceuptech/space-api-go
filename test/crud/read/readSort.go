package main

import (
	"github.com/spaceuptech/space-api-go/api"
	"github.com/spaceuptech/space-api-go/api/utils"
	"fmt"
)

func main() {
	api, err := api.Init("books-app", "localhost", "8081", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	condition := utils.Cond("id", "==", 1)
	resp, err := db.Get("books").Where(condition).Sort("name", "-id").Apply()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if resp.Status == 200 {
			var v []map[string]interface{}
			err:= resp.Unmarshal(&v)
			if err != nil {
				fmt.Println("Error Unmarshalling:", err)
			} else {
				fmt.Println("Result:", v)
			}
		} else {
			fmt.Println("Error Processing Request:", resp.Error)
		}
	}
}
