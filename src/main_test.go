package main

import (
	"reflect"
	"testing"
)

// Two group unit tests: 1. Json Object Unit Testing (Verify core functions logic); 2. Template Configuration Unit Test (Verify go templates)
// JSON Object Unit Testing
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

			got := inputObj.parseSupernets()
			want := outputObj.Supernets
			if !reflect.DeepEqual(*got, *want) {
				t.Errorf("[Failed] name: %s \n want: %#v \n got: %#v", name, want, got)
			}
		})
	}
}

func TestParseInterfaceObj(t *testing.T) {
	type test struct {
		frameworkPath string
		outputJson    []string
	}
	testFolder := "./testcases/"
	testCases := map[string]test{
		"cisco93180yc-fx": {
			frameworkPath: "../input/switchfolder/cisco/93180yc-fx/9.3/framework",
			outputJson: []string{
				testFolder + "cisco93180yc-fx/S31R28-TOR1.json",
				testFolder + "cisco93180yc-fx/S31R28-TOR2.json"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			for _, v := range tc.outputJson {
				wantoutputObj := parseOutputJSON(v)
				gotoutputObj := parseOutputJSON(v)
				gotoutputObj.Vlan = nil
				gotoutputObj.Port = nil
				gotoutputObj.parseInterfaceObj(tc.frameworkPath)
				// Vlan Unit Test
				wantVlan := wantoutputObj.Vlan
				gotVlan := gotoutputObj.Vlan
				if !reflect.DeepEqual(gotVlan, wantVlan) {
					t.Errorf("[Failed] name: %s - VLAN Test \n got: %#v \n want: %#v", wantoutputObj.Device.Hostname, gotVlan, wantVlan)
				}

				// InBondPort Unit Test
				wantPort := wantoutputObj.Port
				gotPort := gotoutputObj.Port
				if !reflect.DeepEqual(gotPort, wantPort) {
					t.Errorf("[Failed] name: %s - InBondPort Test \n got: %#v \n want: %#v", wantoutputObj.Device.Hostname, gotPort, wantPort)
				}
			}
		})
	}
}

func TestParseBGPFramework(t *testing.T) {
	type test struct {
		inputJson     string
		frameworkPath string
		outputJson    []string
	}
	testFolder := "./testcases/"
	testCases := map[string]test{
		"cisco93180yc-fx": {
			inputJson:     testFolder + "cisco93180yc-fx/input.json",
			frameworkPath: "../input/switchfolder/cisco/93180yc-fx/9.3/framework",
			outputJson: []string{
				testFolder + "cisco93180yc-fx/S31R28-TOR1.json",
				testFolder + "cisco93180yc-fx/S31R28-TOR2.json"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			for _, v := range tc.outputJson {
				inputObj := parseInputJSON(tc.inputJson)
				wantoutputObj := parseOutputJSON(v)
				gotoutputObj := parseOutputJSON(v)
				gotoutputObj.Routing = nil

				// Routing Unit Test
				gotoutputObj.parseRoutingFramework(tc.frameworkPath, inputObj)
				wantBGP := wantoutputObj.Routing
				gotBGP := gotoutputObj.Routing
				if !reflect.DeepEqual(gotBGP, wantBGP) {
					t.Errorf("[Failed] name: %s - BGP Test \n got: %#v \n want: %#v", wantoutputObj.Device.Hostname, gotBGP, wantBGP)
				}
			}
		})
	}
}

// Template Configuration Unit Test
