package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-api-go/api/model"
)

func Send(endpoint string, obj interface{}, meta *model.Meta) {
	var url string
	if endpoint != "batch" {
		url = "/v1/api/" + meta.Project + "/crud/" + meta.DbType + "/" + meta.Col + "/" + endpoint
	} else {
		url = "/v1/api/" + meta.Project + "/crud/" + meta.DbType + "/" + endpoint
	}
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("Error in convertin to json ", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", meta.Token)

	req = mux.SetURLVars(req, map[string]string{"project": meta.Project, "dbType": meta.DbType, "col": meta.Col})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

}
