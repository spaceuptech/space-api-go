package main

import (
	"encoding/json"
	"fmt"

	api "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/utils"
)

func main() {
	api := api.New("testapigo", "localhost:4122", false)
	db := api.DB("mongo")

	subscription := db.LiveQuery("books").Where(utils.Cond("age", "==", 20)).Subscribe()

	fmt.Println("subscribed")
	for value := range subscription.C {
		if err := value.Err(); err != nil {
			fmt.Println("error:", err)
			break
		}

		if value.Type() == "initial" {
			continue
		}
		var v interface{}
		_ = value.Unmarshal(&v)
		data, _ := json.MarshalIndent(v, "", " ")
		fmt.Print("value", string(data))
	}

}
