package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

func writeToYaml(yamlFile string, outputResult interface{}) {
	b, err := yaml.Marshal(outputResult)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile(yamlFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

	f, err := os.OpenFile(jsonFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}

func (o *OutputJsonType) parseTemplate(templatePath, outputConfigName string) {
	t, err := template.ParseFiles(
		templatePath+"/allConfig.go.tmpl",
		templatePath+"/header.go.tmpl",
		templatePath+"/inBandPort.go.tmpl",
		templatePath+"/vlan.go.tmpl",
		templatePath+"/bgp.go.tmpl",
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
