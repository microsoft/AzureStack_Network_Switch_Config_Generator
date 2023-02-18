package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"
)

var (
	switchLibFolder  = "/workspaces/AzureStack_Network_Switch_Framework/input/switchLib/"
	testInputFolder  = "/workspaces/AzureStack_Network_Switch_Framework/src/test/testInput/"
	testOutputFolder = "/workspaces/AzureStack_Network_Switch_Framework/src/test/testOutput/"
	testGoldenFolder = "/workspaces/AzureStack_Network_Switch_Framework/src/test/goldenConfig/"
)

func TestMain(t *testing.T) {

	type test struct {
		inputTestFileName string
	}
	testCases := map[string]test{
		"cisco_bgp_nobmc": {
			inputTestFileName: "cisco_nobmc_bgp_input.json",
		},
		"cisco_bgp_bmc": {
			inputTestFileName: "cisco_bmc_bgp_input.json",
		},
		"cisco_static_nobmc": {
			inputTestFileName: "cisco_nobmc_static_input.json",
		},
		"cisco_static_bmc": {
			inputTestFileName: "cisco_bmc_static_input.json",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testInputData := parseInputJson(testInputFolder + tc.inputTestFileName)
			testDeviceTypeMap := testInputData.createDeviceTypeMap()
			generateSwitchConfig(testInputData, switchLibFolder, testOutputFolder+tc.inputTestFileName, testDeviceTypeMap)
			outputFiles := getFilesInFolder(testOutputFolder + tc.inputTestFileName)
			for _, file := range outputFiles {
				if strings.Contains(file, ".json") {
					relativePath := fmt.Sprintf("%s/%s", tc.inputTestFileName, file)
					goldenConfigObj := parseOutputJson(testGoldenFolder + relativePath)
					testOutputObj := parseOutputJson(testOutputFolder + relativePath)
					if !reflect.DeepEqual(goldenConfigObj.Vlans, testOutputObj.Vlans) {
						t.Errorf("name: %s vlan failed \n want: %#v \n got: %#v", name, goldenConfigObj.Vlans, testOutputObj.Vlans)
					}
					if !reflect.DeepEqual(goldenConfigObj.Ports, testOutputObj.Ports) {
						t.Errorf("name: %s interface failed \n want: %#v \n got: %#v", name, goldenConfigObj.Ports, testOutputObj.Ports)
					}
				}
			}
		})
	}
}

func getFilesInFolder(foldername string) []string {
	fileList := []string{}
	files, err := ioutil.ReadDir(foldername)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}
	return fileList
}

func parseOutputJson(outputJsonFile string) *OutputType {
	outputObj := &OutputType{}
	bytes, err := ioutil.ReadFile(outputJsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, outputObj)
	if err != nil {
		log.Fatalln(err)
	}
	return outputObj
}
