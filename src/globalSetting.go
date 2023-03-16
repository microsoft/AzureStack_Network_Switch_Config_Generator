package main

import "fmt"

func (o *OutputType) UpdateGlobalSetting(inputData InputData) {
	// Create random credential for switch config if no input values
	if Username == "" || Password == "" {
		Username = "azureadmin-" + generateRandomString(5, 0, 0, 0)
		Password = generateRandomString(16, 3, 3, 3)
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
