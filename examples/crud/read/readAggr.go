package main

import (
	"fmt"
	"github.com/spaceuptech/space-api-go"
)

func main() {
	api, err := api.New("books-app", "localhost:8081", false)
	if err != nil {
		fmt.Println(err)
	}
	db := api.Mongo()
	pipe := []interface{}{
		map[string]interface{}{"$match": map[string]interface{}{"status": "A"}},
		map[string]interface{}{"$group": map[string]interface{}{"_id": "$cust_id", "total": map[string]interface{}{"$sum": "$amount"}}},
	}
	resp, err := db.Aggr("books").Pipe(pipe).Apply()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if resp.Status == 200 {
			var v []map[string]interface{}
			err := resp.Unmarshal(&v)
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
