package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"
)

func (o *OutputType) writeToJson(outputFolder string) {
	// Create folder if not existing
	createFolder(outputFolder)
	jsonFile := outputFolder + "/" + o.Switch.Hostname + JSONExtension
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

func (o *OutputType) parseTemplate(templateFolder, outputFolder string) {
	configFilePath := outputFolder + "/" + o.Switch.Hostname + CONFIGExtension
	// Parse the whole template folder based on .go.tmpl files
	t := template.Must(template.ParseGlob(templateFolder + "/*.go.tmpl"))

	f, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Execute(f, o)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}
