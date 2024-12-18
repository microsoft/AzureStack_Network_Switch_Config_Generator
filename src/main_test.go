package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

var ()

func TestMain(t *testing.T) {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	switchLibFolder := "../input/switchLib"
	testInputFolder := cwd + "/test/testInput/"
	testOutputFolder := cwd + "/test/testOutput/"
	testGoldenFolder := cwd + "/test/goldenConfig/"
	debugYamlFolderName := "debug_yaml"

	type test struct {
		inputTestFileName string
	}
	testCases := map[string]test{
		"DellEMC-Switched": {
			inputTestFileName: "DellEMC-Switched.json",
		},
		"CiscoNXOS-HyperConverged": {
			inputTestFileName: "CiscoNXOS-HyperConverged.json",
		},
		"CiscoNXOS-Switched": {
			inputTestFileName: "CiscoNXOS-Switched.json",
		},
		"CIscoNXOS-Switchless": {
			inputTestFileName: "CIscoNXOS-Switchless.json",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testInputData := parseInputJson(testInputFolder + tc.inputTestFileName)
			testDeviceTypeMap := testInputData.createDeviceTypeMap()
			generateSwitchConfig(testInputData, switchLibFolder, testOutputFolder+name, testDeviceTypeMap)
			outputFiles := getFilesInFolder(testOutputFolder + name + "/" + debugYamlFolderName)
			// fmt.Println(outputFiles)
			for _, file := range outputFiles {
				if strings.Contains(file, YAMLExtension) {
					relativePath := fmt.Sprintf("%s/%s/%s", name, debugYamlFolderName, file)
					// fmt.Println(relativePath)
					goldenConfigObj := parseOutputYaml(testGoldenFolder + relativePath)
					testOutputObj := parseOutputYaml(testOutputFolder + relativePath)
					if !reflect.DeepEqual(goldenConfigObj.Vlans, testOutputObj.Vlans) {
						t.Errorf("name: %s VLAN failed \n want: %#v \n got: %#v", name, goldenConfigObj.Vlans, testOutputObj.Vlans)
					}
					if !reflect.DeepEqual(goldenConfigObj.Ports, testOutputObj.Ports) {
						t.Errorf("name: %s Interface failed \n want: %#v \n got: %#v", name, goldenConfigObj.Ports, testOutputObj.Ports)
					}
					if len(goldenConfigObj.Routing.BGP.IPv4Network) != len(testOutputObj.Routing.BGP.IPv4Network) {
						t.Errorf("name: %s BGP routing failed \n want: %#v \n got: %#v", name, len(goldenConfigObj.Routing.BGP.IPv4Network), len(testOutputObj.Routing.BGP.IPv4Network))
					}
					// if len(goldenConfigObj.Routing.PrefixList) != len(testOutputObj.Routing.PrefixList) {
					// 	t.Errorf("name: %s Routing PrefixList failed \n want: %#v \n got: %#v", name, len(goldenConfigObj.Routing.PrefixList), len(testOutputObj.Routing.PrefixList))
					// }
				}
			}
		})
	}
}

func getFilesInFolder(foldername string) []string {
	fileList := []string{}
	files, err := os.ReadDir(foldername)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}
	return fileList
}

func parseOutputYaml(outputFile string) *OutputType {
	outputObj := &OutputType{}
	bytes, err := os.ReadFile(outputFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(bytes, outputObj)
	if err != nil {
		log.Fatalln(err)
	}
	return outputObj
}
