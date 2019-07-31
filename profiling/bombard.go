package main

import (
	"github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/api/utils"
	"fmt"
	"time"
	"math/rand"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.Mongo()
	db.Delete("books").Apply()
	for {
		sel := rand.Intn(4)
		id := rand.Intn(10)
		start := time.Now()
		switch sel {
		case 0:
			db.Insert("books").Doc(
				map[string]interface{}{"_id":id, "name":"SomeBook"+string(id), "author":"Author"+string(id)}).Apply()
		case 1:
			condition := utils.Cond("_id", "==", id)
			set := map[string]interface{}{"author":"auth"+string(id)}
			db.Update("books").Where(condition).Set(set).Apply()
		case 2:
			db.DeleteOne("books").Where(utils.Cond("_id", "==", id)).Apply()
		default:
			db.Get("books").Apply()
		}
		fmt.Println(time.Now().Sub(start))
	}
}
