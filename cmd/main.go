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

type Microservices struct {
	Name    string
	Servers []ServiceRProxy
}

type ServiceRProxy struct {
	RProxy      *httputil.ReverseProxy
	Healthy     bool
	HealthRoute *url.URL
	Url         *url.URL
}

var (
	lastServedIndex = 0
	serverList      = make(map[string]Microservices)
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
			ms := Microservices{
				Name:    property.Name,
				Servers: []ServiceRProxy{},
			}
			for _, urlStr := range property.Servers {
				// create and store the reverse proxy for each url
				u, _ := url.Parse(urlStr)
				healthUrl, _ := url.Parse(urlStr + property.HealthRoute)
				ms.Servers = append(ms.Servers, ServiceRProxy{
					RProxy:      createReverseProxy(u),
					Healthy:     true,
					Url:         u,
					HealthRoute: healthUrl,
				})
			}
			serverList[property.ApiPrefix] = ms
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
	nextIndex := (lastServedIndex + 1) % len(serverList[name].Servers)
	server := serverList[name].Servers[lastServedIndex]
	lastServedIndex = nextIndex
	return server.RProxy
}

func createReverseProxy(url *url.URL) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(url)
}
