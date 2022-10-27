package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// type ServerStatus struct {
// 	UpdateTime string `json:"uptime"`
// 	Status     string `json:"status"`
// }

// func main() {
// 	ID := "localhost"
// 	url := strings.Join([]string{"http://", ID, ":8000/status"}, "")
// 	response, err := http.Get(url)
// 	// response, err := http.Get("http://192.168.2.0:8000/status")

// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}

// 	responseData, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(string(responseData))

// 	var responseServer ServerStatus
// 	err1 := json.Unmarshal(responseData, &responseServer)
// 	if err1 != nil {
// 		log.Fatal(err1)
// 	}

// 	fmt.Println(responseServer.UpdateTime)
// 	fmt.Println(responseServer.Status)

// }

// func (sc ServerController) CheckStatusServer() (int, int, int, float64) {

// 	var servers []models.Server
// 	sc.DB.Find(&servers)
// 	totalServer := len(servers)
// 	if totalServer == 0 {
// 		log.Fatal("No server exists")
// 	}

// 	countServerOn := 0
// 	countServerOff := 0
// 	totalUptime := 0.0
// 	for _, server := range servers {
// 		url := strings.Join([]string{"http://", server.ID, ":8000/status"}, "")
// 		response, err := http.Get(url)
// 		// response, err := http.Get("http://192.168.2.0:8000/status")
// 		if err != nil {
// 			// fmt.Print(err.Error())
// 			// os.Exit(1)
// 			countServerOff++
// 			continue
// 		}
// 		countServerOn++
// 		responseData, err := ioutil.ReadAll(response.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		var responseServer models.ServerStatus
// 		err1 := json.Unmarshal(responseData, &responseServer)
// 		if err1 != nil {
// 			log.Fatal(err1)
// 		}
// 		IntUptime, _ := strconv.ParseFloat(responseServer.UpdateTime, 8)
// 		totalUptime += IntUptime
// 	}

// 	var avgUptime float64
// 	avgUptime = totalUptime / float64(totalServer)

// 	return totalServer, countServerOn, countServerOff, avgUptime
// }
