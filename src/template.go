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


func (o *OutputType) parseSelectedTemplate(templateFolder, outputFolder string) {
	// Create the path for the new config file in the output folder
	newConfigOnlyFilePath := outputFolder + "/" + o.GlobalSetting.OOBIP + CONFIGExtension

	// Open the new config file with write-only access, create it if it doesn't exist, and truncate it if it does
	f, err := os.OpenFile(newConfigOnlyFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// If there's an error opening the file, log the error and exit the program
	if err != nil {
		log.Fatalln(err)
	}
	// Ensure the file gets closed once the function finishes
	defer f.Close()

	// Define the template files to be parsed
	selectedTemplateFiles := []string{
		"WANSIMConfigOnly.go.tmpl",
		"wansim.go.tmpl",
		"prefixlist.go.tmpl",
	}

	// Initialize a slice to hold the full file paths of the template files
	var filePaths []string
	// For each template file, join the template folder path and the template file name, and append it to the file paths slice
	for _, file := range selectedTemplateFiles {
		filePaths = append(filePaths, filepath.Join(templateFolder, file))
	}

	// Parse the template files
	t, err := template.ParseFiles(filePaths...)
	// If there's an error parsing the files, panic and exit the program
	if err != nil {
		panic(err)
	}

	// Execute the template with the OutputType data
	err = t.Execute(f, o)
	// If there's an error executing the template, panic and exit the program
	if err != nil {
		panic(err)
	}
}