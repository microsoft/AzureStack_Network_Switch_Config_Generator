package main

import "fmt"

func (o *OutputType) UpdateGlobalSetting(inputData InputData) {
	o.GlobalSetting.Username = Username
	o.GlobalSetting.Password = Password
	o.GlobalSetting.TimeServer = inputData.Setting.TimeServer
	o.GlobalSetting.SyslogServer = inputData.Setting.SyslogServer
	o.GlobalSetting.DNSForwarder = inputData.Setting.DNSForwarder
	for _, v := range o.Vlans {
		if v.GroupID == BMC {
			o.GlobalSetting.OOB = fmt.Sprintf("vlan%d", v.VlanID)
			break
		}
	}
}
