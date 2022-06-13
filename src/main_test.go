package main

import (
	"reflect"
	"testing"
)

func TestOutputNetwork(t *testing.T) {
	type test struct {
		inputJsonFile  string
		outputJsonFile string
	}
	testFolder := "./testcases/"
	testCases := map[string]test{
		"cisco93180yc-fx": {
			testFolder + "cisco93180yc-fx/input.json",
			testFolder + "cisco93180yc-fx/S31R28-TOR1.json",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			inputObj := parseInputJSON(tc.inputJsonFile)
			outputObj := parseOutputJSON(tc.outputJsonFile)

			got := inputObj.outputNetwork()
			expect := outputObj.Network
			if !reflect.DeepEqual(*got, *expect) {
				t.Errorf("name: %s failed \n want: %#v \n got: %#v", name, expect, got)
			}
		})
	}
}

func TestUpdateOutputObj(t *testing.T) {
	type test struct {
		inputJsonFile string
		frameworkPath string
		templatePath  string
	}
	testFolder := "./testcases/"
	testCases := map[string]test{
		"cisco93180yc-fx": {
			inputJsonFile: testFolder + "cisco93180yc-fx/input.json",
			frameworkPath: "../input/switchfolder/cisco/93180yc-fx/9.3/framework",
			templatePath:  "../input/switchfolder/cisco/93180yc-fx/9.3/template",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			inputObj := parseInputJSON(tc.inputJsonFile)
			outputObj := parseOutputJSON(tc.outputJsonFile)

			got := inputObj.outputNetwork()
			expect := outputObj.Network
			if !reflect.DeepEqual(*got, *expect) {
				t.Errorf("name: %s failed \n want: %#v \n got: %#v", name, expect, got)
			}
		})
	}
}
