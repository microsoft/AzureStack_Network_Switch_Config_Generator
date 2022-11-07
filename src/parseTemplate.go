package main

import (
	"encoding/json"
	"log"
	"os"
)

func (o *OutputType) UpdateSwitch(switchType SwitchType) {

}

func (o *OutputType) writeToJson(outputFolder string) {
	// Create folder if not existing
	createFolder(outputFolder)
	jsonFile := outputFolder + "/" + o.Switch.Hostname + ".json"
	b, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile(jsonFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}
