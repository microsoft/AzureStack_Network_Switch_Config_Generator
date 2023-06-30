package main

import (
	"fmt"
	"strings"
)

func (o *OutputType) UpdateGlobalSetting(inputData InputData) {
	// Create random credential for switch config if no input values
	if Username == "" || Password == "" {
		Username = CRED_SCAN_PLACEHOLDER
		Password = CRED_SCAN_PLACEHOLDER
		// // Generate Random Username Password based on requirement
		// Username = "azureadmin-" + generateRandomString(5, 0, 0, 0)
		// Password = generateRandomString(16, 3, 3, 3)
	}
	o.GlobalSetting.Username = Username
	o.GlobalSetting.Password = Password
	o.GlobalSetting.TimeServer = inputData.Setting.TimeServer
	o.GlobalSetting.SyslogServer = inputData.Setting.SyslogServer
	o.GlobalSetting.DNSForwarder = inputData.Setting.DNSForwarder
	// Update Deployment Pattern
	o.DeploymentPattern = inputData.DeploymentPattern

	for _, v := range o.Vlans {
		if v.GroupName == BMC {
			o.GlobalSetting.OOB = fmt.Sprintf("vlan%d", v.VlanID)
			break
		}
	}
}

func (o *OutputType) UpdateDHCPIps(inputData InputData) {
	for _, supernet := range inputData.Supernets {
		if supernet.GroupName == BMC {
			DHCPInfra := []string{}
			DHCPTenant := []string{}
			for _, assignIP := range supernet.IPv4.Assignment {
				if strings.Contains(assignIP.Name, "DVM") {
					DHCPInfra = append(DHCPInfra, assignIP.IP)
				} else if strings.EqualFold(assignIP.Name, "HLH-OS") {
					DHCPTenant = append(DHCPTenant, assignIP.IP)
				}
			}
			o.GlobalSetting.DHCPInfra = DHCPInfra
			o.GlobalSetting.DHCPTenant = DHCPTenant
		}
	}
}
