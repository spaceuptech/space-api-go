package main

import (
	api "github.com/spaceuptech/space-api-go"
	"fmt"
	"os"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	filestore := api.Filestore()

	file, err := os.Open("a.txt")
	if err != nil {
		panic(err)
	}
	resp, err := filestore.UploadFile("\\Folder", "hello1.txt", file)
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
