package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type ServerStatus struct {
	UpdateTime string `json:"uptime"`
	Status     string `json:"status"`
}

func main() {
	ID := "localhost"
	url := strings.Join([]string{"http://", ID, ":8000/status"}, "")
	response, err := http.Get(url)
	// response, err := http.Get("http://192.168.2.0:8000/status")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))

	var responseServer ServerStatus
	err1 := json.Unmarshal(responseData, &responseServer)
	if err1 != nil {
		log.Fatal(err1)
	}

	fmt.Println(responseServer.UpdateTime)
	fmt.Println(responseServer.Status)

}
