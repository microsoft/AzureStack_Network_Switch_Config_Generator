package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

var ()

func TestMain(t *testing.T) {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cwd)

	switchLibFolder := "../input/switchLib/"
	testInputFolder := cwd + "/test/testInput/"
	testOutputFolder := cwd + "/test/testOutput/"
	testGoldenFolder := cwd + "/test/goldenConfig/"

	type test struct {
		inputTestFileName string
	}
	testCases := map[string]test{
		"cisco_nobmc_bgp": {
			inputTestFileName: "cisco_nobmc_bgp_input.json",
		},
		"cisco_bmc_bgp": {
			inputTestFileName: "cisco_bmc_bgp_input.json",
		},
		"cisco_nobmc_static": {
			inputTestFileName: "cisco_nobmc_static_input.json",
		},
		"cisco_bmc_static": {
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
						t.Errorf("name: %s VLAN failed \n want: %#v \n got: %#v", name, goldenConfigObj.Vlans, testOutputObj.Vlans)
					}
					if !reflect.DeepEqual(goldenConfigObj.Ports, testOutputObj.Ports) {
						t.Errorf("name: %s Interface failed \n want: %#v \n got: %#v", name, goldenConfigObj.Ports, testOutputObj.Ports)
					}
					if len(goldenConfigObj.Routing.BGP.IPv4Network) != len(testOutputObj.Routing.BGP.IPv4Network) {
						t.Errorf("name: %s BGP routing failed \n want: %#v \n got: %#v", name, len(goldenConfigObj.Routing.BGP.IPv4Network), len(testOutputObj.Routing.BGP.IPv4Network))
					}
					if len(goldenConfigObj.Routing.PrefixList) != len(testOutputObj.Routing.PrefixList) {
						t.Errorf("name: %s Routing PrefixList failed \n want: %#v \n got: %#v", name, len(goldenConfigObj.Routing.PrefixList), len(testOutputObj.Routing.PrefixList))
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
