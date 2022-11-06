package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	TOR    = "TOR"
	BMC    = "BMC"
	BORDER = "BORDER"
	MUX    = "MUX"
)

func init() {
	// Set Log Output Options - example: 2022/08/24 21:51:10 main.go:58:
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/lab_input.json", "File path of switch deploy input.json")
	flag.Parse()
	// Covert input.json to Go Object, structs are defined in model.go
	inputObj := parseInputJson(*inputJsonFile)
	inputData := inputObj.InputData
	// Create random credential for switch config
	// randomUsername := "aszadmin-" + generateRandomString(5, 0, 0, 0)
	// randomPassword := generateRandomString(16, 3, 3, 3)

	for _, device := range inputData.Switches {
		fmt.Println(device)
		outputObj := &OutputType{}
	}
}
