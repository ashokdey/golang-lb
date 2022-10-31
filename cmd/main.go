package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

var (
	lastServedIndex = 0
	serverList      = make(map[string][]string)
)

func main() {
	fileName, _ := filepath.Abs("conf.yml")
	conf, err := readConfFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// populate serverList using conf data
	for _, service := range conf.Microservices {
		for _, property := range service.Service {
			serverList[property.ApiPrefix] = property.Servers
		}
	}

	fmt.Println(serverList)

	http.HandleFunc("/", forwardRequests)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// forward requests
func forwardRequests(res http.ResponseWriter, req *http.Request) {
	// get the api prefix
	serviceName := strings.Split(req.URL.Path, "/")

	url := getServer(serviceName[1])
	rProxy := httputil.NewSingleHostReverseProxy(url)
	rProxy.ServeHTTP(res, req)
}

func getServer(name string) *url.URL {
	nextIndex := (lastServedIndex + 1) % 2

	fmt.Println(serverList[name])
	fmt.Println(serverList[name][lastServedIndex])

	url, _ := url.Parse(serverList[name][lastServedIndex])
	lastServedIndex = nextIndex
	return url
}
