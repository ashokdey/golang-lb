package main

import (
	"log"
	"time"
)

func checkServers() {
	for {
		for key := range serverList {
			for i, s := range serverList[key].Servers {
				isHealthy := s.Ping()
				if isHealthy {
					log.Printf("%s is healthy", s.Url)
					// update server list and mark the server as true
					serverList[key].Servers[i].Healthy = true
				} else {
					log.Printf("%s is down", s.Url)
					// update server list and mark the server as false
					serverList[key].Servers[i].Healthy = false
				}
				time.Sleep(time.Second * 2)
			}
		}
	}
}
