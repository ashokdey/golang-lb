package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

type Microservices struct {
	Name    string
	Servers []Server
}

type Server struct {
	RProxy      *httputil.ReverseProxy
	Healthy     bool
	HealthRoute string
	Url         string
}

func (s *Server) Ping() bool {
	res, err := http.Head(s.HealthRoute)
	if err != nil {
		s.Healthy = false
		return s.Healthy
	}
	if res.StatusCode != http.StatusOK {
		s.Healthy = false
		return s.Healthy
	}
	s.Healthy = true
	return s.Healthy
}

func createReverseProxy(url *url.URL) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(url)
}

func populateMicroservices(fileName string) error {
	file, _ := filepath.Abs(fileName)
	conf, err := readConfFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// populate serverList using conf data
	for _, service := range conf.Microservices {
		for _, property := range service.Service {
			ms := Microservices{
				Name:    property.Name,
				Servers: []Server{},
			}
			for _, urlStr := range property.Servers {
				// create and store the reverse proxy for each url
				u, _ := url.Parse(urlStr)
				ms.Servers = append(ms.Servers, Server{
					RProxy:      createReverseProxy(u),
					Healthy:     true,
					Url:         urlStr,
					HealthRoute: urlStr + property.HealthRoute,
				})
			}
			serverList[property.ApiPrefix] = ms
		}
	}
	return nil
}
