package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestGetMaxIPSize(t *testing.T) {
	type test struct {
		input string
		want  int
	}
	ipNets := map[string]test{
		"/24": {"10.0.0.0/24", 256},
		"/30": {"10.0.0.0/30", 4},
		"/32": {"10.0.0.0/32", 1},
	}

	for name, tc := range ipNets {
		t.Run(name, func(t *testing.T) {
			got := getMaxIPSize(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("name: %s failed, want: %v, got: %v", name, tc.want, got)
			}
		})
	}
}

func TestGenerateIPList(t *testing.T) {
	type test struct {
		input string
		want  []string
	}
	ipNets := map[string]test{
		"/30": {"10.0.0.0/30", []string{"10.0.0.0", "10.0.0.1", "10.0.0.2", "10.0.0.3"}},
		"/32": {"10.0.0.0/32", []string{"10.0.0.0"}},
	}

	for name, tc := range ipNets {
		t.Run(name, func(t *testing.T) {
			got := generateIPList(tc.input)
			if !reflect.DeepEqual(*got, tc.want) {
				t.Errorf("name: %s failed, want: %v, got: %v", name, tc.want, *got)
			}
		})
	}
}

func TestIPCaculator2(t *testing.T) {
	type test struct {
		inputJsonFile  string
		resultJsonFile string
	}
	testJsonFolder := "./testcases/ipcaculator/"
	testCases := map[string]test{
		"standard_input": {
			testJsonFolder + "input1.json",
			testJsonFolder + "result1.json",
		},
		"empty": {
			testJsonFolder + "input2.json",
			testJsonFolder + "result2.json",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bytes, err := ioutil.ReadFile(tc.inputJsonFile)
			if err != nil {
				log.Fatalln(err)
			}
			got := parseSupernetSection(bytes)
			expect := parseJsonFile(tc.resultJsonFile)
			if !reflect.DeepEqual(*got, *expect) {
				t.Errorf("name: %s failed \n want: %#v \n got: %#v", name, expect, got)
			}
		})
	}
}

func parseJsonFile(jsonFile string) *[]SupernetOutputType {
	b, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	outputResult := []SupernetOutputType{}
	err = json.Unmarshal(b, &outputResult)
	if err != nil {
		log.Fatalln(err)
	}
	return &outputResult
}
