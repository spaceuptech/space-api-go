package main

import (
	"fmt"

	"github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/service"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	api.SetToken("my_secret")
	service, _ := api.Service("service")
	service.RegisterFunc("echo_func", Echo)
	service.Start()
	
}

func Echo(params, auth *model.Message, fn service.CallBackFunction) {
	var i interface{}
	params.Unmarshal(&i)
	fn("response", i)
}
