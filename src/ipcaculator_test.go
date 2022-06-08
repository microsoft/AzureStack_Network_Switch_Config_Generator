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
	testJsonFolder := "./testcases/"
	testCases := map[string]test{
		"standardInput": {
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
			got := parseInterfaceSection(bytes)
			expect := parseJsonFile(tc.resultJsonFile)
			if !reflect.DeepEqual(*got, *expect) {
				t.Errorf("name: %s failed \n want: %#v \n got: %#v", name, expect, got)
			}
		})
	}
}

func parseJsonFile(jsonFile string) *[]NetworkOutputType {
	b, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	outputResult := []NetworkOutputType{}
	err = json.Unmarshal(b, &outputResult)
	if err != nil {
		log.Fatalln(err)
	}
	return &outputResult
}

// func TestIPCaculator1(t *testing.T) {
// 	type test struct {
// 		inputAssignCIDR    string
// 		inputJsonInputFile string
// 		want               *OutputSubnet
// 	}
// 	testJsonFolder := "./testJsonFiles/"
// 	testCases := map[string]test{
// 		"10.0.0.0/24": {
// 			inputAssignCIDR:    "10.0.0.0/24",
// 			inputJsonInputFile: testJsonFolder + "input1.json",
// 			want:               &OutputSubnet{IPList: []IPItem{IPItem{Name: "P2P_TOR1_To_Border1", IPNet: "10.0.0.0/31"}, IPItem{Name: "P2P_TOR1", IPNet: "10.0.0.0/31"}, IPItem{Name: "P2P_Border2", IPNet: "10.0.0.1/31"}, IPItem{Name: "P2P_TOR1_To_Border2", IPNet: "10.0.0.2/31"}, IPItem{Name: "P2P_TOR2", IPNet: "10.0.0.2/31"}, IPItem{Name: "P2P_Border2", IPNet: "10.0.0.3/31"}, IPItem{Name: "P2P_TOR2_To_Border1", IPNet: "10.0.0.4/31"}, IPItem{Name: "P2P_TOR1", IPNet: "10.0.0.4/31"}, IPItem{Name: "P2P_Border2", IPNet: "10.0.0.5/31"}, IPItem{Name: "P2P_TOR2_To_Border2", IPNet: "10.0.0.6/31"}, IPItem{Name: "P2P_TOR2", IPNet: "10.0.0.6/31"}, IPItem{Name: "P2P_Border2", IPNet: "10.0.0.7/31"}, IPItem{Name: "Loopback_TOR1", IPNet: "10.0.0.8/32"}, IPItem{Name: "Loopback_TOR2", IPNet: "10.0.0.9/32"}}},
// 		},
// 		"10.0.0.0/28": {
// 			inputAssignCIDR:    "10.0.0.0/28",
// 			inputJsonInputFile: testJsonFolder + "input2.json",
// 			want:               &OutputSubnet{IPList: []IPItem{IPItem{Name: "S1", IPNet: "10.0.0.0/31"}, IPItem{Name: "S1_IP1", IPNet: "10.0.0.0/31"}, IPItem{Name: "S1_IP2", IPNet: "10.0.0.1/31"}, IPItem{Name: "S2", IPNet: "10.0.0.2/31"}, IPItem{Name: "S2_IP1", IPNet: "10.0.0.2/31"}, IPItem{Name: "S2_IP1", IPNet: "10.0.0.3/31"}, IPItem{Name: "S3", IPNet: "10.0.0.4/31"}, IPItem{Name: "S3_IP1", IPNet: "10.0.0.4/31"}, IPItem{Name: "S3_IP2", IPNet: "10.0.0.5/31"}, IPItem{Name: "S4", IPNet: "10.0.0.6/31"}, IPItem{Name: "S4_IP1", IPNet: "10.0.0.6/31"}, IPItem{Name: "S4_IP2", IPNet: "10.0.0.7/31"}, IPItem{Name: "S5", IPNet: "10.0.0.8/32"}, IPItem{Name: "S6", IPNet: "10.0.0.9/32"}}},
// 		},
// 	}
// 	for name, tc := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			got := IPCaculator(tc.inputAssignCIDR, tc.inputJsonInputFile)
// 			if !reflect.DeepEqual(*got, *tc.want) {
// 				t.Errorf("name: %s failed, want: %v, got: %v", name, tc.want, got)
// 			}
// 		})
// 	}
// }
