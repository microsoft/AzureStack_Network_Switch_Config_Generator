package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Policy struct {
	PrefixList []PrefixListType `json:"PrefixList"`
}

type PrefixListType struct {
	Name   string `json:"Name"`
	Config []struct {
		Name      string `json:"Name"`
		Action    string `json:"Action"`
		Supernet  string `json:"Supernet"`
		IPAddress string `json:"IPAddress"`
		Operation string `json:"Operation"`
		Prefix    int    `json:"Prefix"`
	} `json:"Config"`
}

// type PrefixList map[string][]string

func main() {
	p := &Policy{}
	bytes, err := ioutil.ReadFile("./test.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, p)
	if err != nil {
		log.Fatalln(err)
	}
	// (*p)["test"] = append((*p)["test"], PrefixListType{
	// 	Name:   "name",
	// 	Action: "allow",
	// })
	fmt.Printf("%+v", p)
	fmt.Println()
}
