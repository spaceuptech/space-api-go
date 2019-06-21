package main

import (
	"github.com/spaceuptech/space-api-go"
	"fmt"
	"os"
)

func main() {
	api, err := api.New("books-app", "localhost:8081", false)
	if(err != nil) {
		fmt.Println(err)
	}
	filestore := api.Filestore()

	file, err := os.Create("test1.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	err = filestore.DownloadFile("\\Folder\\text.txt", file)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
