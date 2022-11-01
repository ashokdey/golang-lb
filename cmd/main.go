package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	lastServedIndex = 0
	serverList      = make(map[string]Microservices)
)

func main() {
	err := populateMicroservices("conf.yml")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Populated config")

	go checkServers()
	fmt.Println("Started CRON")

	http.HandleFunc("/", forwardRequests)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// forward requests
func forwardRequests(res http.ResponseWriter, req *http.Request) {
	// get the api prefix
	serviceName := strings.Split(req.URL.Path, "/")
	server, err := getHealthyServer(serviceName[1])

	// fmt.Println("Server value ", server)
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}

	if server.Healthy {
		server.RProxy.ServeHTTP(res, req)
	} else {
		log.Println("Server not healthy returning error")
		res.Write([]byte(err.Error()))
		return
	}
}

func getHealthyServer(name string) (Server, error) {
	for i := 0; i < len(serverList[name].Servers); i++ {
		server := getServer(name)
		if server.Healthy {
			// log.Printf("Serving from : %s", server.Url)
			return server, nil
		}
	}
	return Server{}, errors.New("all servers are down")
}

func getServer(name string) Server {
	nextIndex := (lastServedIndex + 1) % len(serverList[name].Servers)
	server := serverList[name].Servers[lastServedIndex]
	lastServedIndex = nextIndex
	return server
}
