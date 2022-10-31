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
	serverList      = make(map[string][]*httputil.ReverseProxy)
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
			for _, u := range property.Servers {
				// create and store the reverse proxy for each url
				serverList[property.ApiPrefix] =
					append(serverList[property.ApiPrefix], createReverseProxy(u))
			}
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

	rProxy := getServer(serviceName[1])
	rProxy.ServeHTTP(res, req)
}

func getServer(name string) *httputil.ReverseProxy {
	nextIndex := (lastServedIndex + 1) % len(serverList[name])
	server := serverList[name][lastServedIndex]
	lastServedIndex = nextIndex
	return server
}

func createReverseProxy(urlString string) *httputil.ReverseProxy {
	u, _ := url.Parse(urlString)
	return httputil.NewSingleHostReverseProxy(u)
}
