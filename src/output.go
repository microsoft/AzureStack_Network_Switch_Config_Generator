package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

func createFolder(folderPath string) {
	_, err := os.Stat(folderPath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderPath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
}

func writeToYaml(yamlFile string, outputResult interface{}) {
	b, err := yaml.Marshal(outputResult)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile(yamlFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}

func writeToJson(jsonFile string, outputResult interface{}) {
	b, err := json.MarshalIndent(outputResult, "", " ")
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

func (o *OutputType) parseTemplate(templatePath, outputConfigName string) {
	t, err := template.ParseFiles(
		templatePath+"/allConfig.go.tmpl",
		templatePath+"/header.go.tmpl",
		templatePath+"/stig.go.tmpl",
		templatePath+"/port.go.tmpl",
		templatePath+"/vlan.go.tmpl",
		templatePath+"/default.go.tmpl",
		templatePath+"/bgp.go.tmpl",
		templatePath+"/stp.go.tmpl",
		templatePath+"/external.go.tmpl",
	)
	if err != nil {
		log.Fatalln(err)
	}

	// err = t.Execute(os.Stdout, o)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	f, err := os.OpenFile(outputConfigName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln("create file: ", err)
		return
	}

	err = t.Execute(f, o)
	if err != nil {
		log.Fatalln("execute: ", err)
		return
	}
	f.Close()
}
