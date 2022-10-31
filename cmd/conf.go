package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type YamlConf struct {
	Microservices []YamlConfServices `yaml:"microservices"`
}

type YamlConfServices struct {
	Service map[string]YamlConfServicesServiceProperties `yaml:"service"`
}

type YamlConfServicesServiceProperties struct {
	Name        string   `yaml:"name"`
	Servers     []string `yaml:"servers"`
	Version     string   `yaml:"version"`
	ApiPrefix   string   `yaml:"apiPrefix"`
	HealthRoute string   `yaml:"healthRoute"`
}

func readConfFile(filename string) (*YamlConf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err != nil {
		fmt.Println(err)
	}

	c := &YamlConf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
