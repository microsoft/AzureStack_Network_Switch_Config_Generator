package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func (o *OutputType) writeToYaml(outputFolder string) {
	// Create folder if not existing
	createFolder(outputFolder)
	yamlFile := outputFolder + "/" + o.Switch.Hostname + YAMLExtension
	b, err := yaml.Marshal(o)
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

func containsString(s, substring string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substring))
}

func (o *OutputType) parseCombineTemplate(templateFolder, outputFolder, configFilename string) {
	configFilePath := outputFolder + "/" + configFilename + CONFIGExtension
	templateFiles := templateFolder + "/*.go.tmpl"
	// Parse the whole template folder based on .go.tmpl files

	t := template.Must(template.ParseGlob(templateFiles))
	// .Funcs(template.FuncMap{
	// 	"contains":  containsString,
	// 	"hasPrefix": strings.HasPrefix,
	// 	"hasSuffix": strings.HasSuffix,
	// })

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

func (o *OutputType) parseEachTemplate(templateFolder, outputFolder string) {
	// Parse all .go.tmpl files in the specified folder
	templateFiles, err := filepath.Glob(filepath.Join(templateFolder, "*.go.tmpl"))
	if err != nil {
		log.Fatalln(err)
	}

	// Define the templates
	t := template.New("")

	for _, templateFile := range templateFiles {
		// Parse the template file
		t = template.Must(t.ParseFiles(templateFile))
	}

	// Generate the output files
	for _, templateFile := range templateFiles {
		configFilename := strings.TrimSuffix(filepath.Base(templateFile), ".go.tmpl")
		configFilePath := filepath.Join(outputFolder, configFilename)

		f, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalln(err)
		}

		err = t.ExecuteTemplate(f, configFilename, o.WANSIM)
		if err != nil {
			log.Fatalln(err)
		}
		f.Close()
	}
}
