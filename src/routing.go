package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func (o *OutputType) ParseRouting(frameworkFolder string, inputData InputData) {
	routingJsonPath := fmt.Sprintf("%s/%s", frameworkFolder, BGPROUTINGJSON)
	routingJsonObj := parseRoutingJson(routingJsonPath)
	if inputData.SwitchUplink == "BGP" {
		o.ParseBGP(routingJsonObj.BGP)
		o.ParsePrefixList(routingJsonObj.PrefixList)
	}
}

func parseRoutingJson(routingJsonPath string) *RoutingType {
	routingJsonObj := &RoutingType{}
	bytes, err := ioutil.ReadFile(routingJsonPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, routingJsonObj)
	if err != nil {
		log.Fatalln(err)
	}
	return routingJsonObj
}

func (o *OutputType) ParseBGP(BGPObj BGPType) {
	o.Routing.BGP = BGPObj
}

func (o *OutputType) ParsePrefixList(RoutingPolicyObj []PrefixListType) {
	o.Routing.PrefixList = RoutingPolicyObj
}
