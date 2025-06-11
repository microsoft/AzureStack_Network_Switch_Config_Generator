import json
import os

class StandardJSONBuilder:
    def __init__(self, input_data):
        self.input_data = input_data
        self.sections = {}

    def build_switch(self, target_types=None):
        if target_types is None:
            target_types = {"TOR1", "TOR2", "BMC"}

        switches = self.input_data.get("InputData", {}).get("Switches", [])
        self.sections["switch"] = {}

        for sw in switches:
            sw_type = sw.get("Type")
            if sw_type in target_types:
                self.sections["switch"][sw_type] = {
                    "make": sw.get("Make", "").lower(),
                    "model": sw.get("Model", "").lower(),
                    "type": sw_type,
                    "hostname": sw.get("Hostname", "").lower(),
                    "firmware": sw.get("Firmware", "").lower()
                }

        if not self.sections["switch"]:
            print("[!] No valid switches found in input data.")

    def build_vlans(self, switch_type):
        vlans = []

        for net in self.input_data.get("InputData", {}).get("Supernets", []):
            ipv4 = net.get("IPv4", {})
            vlan_id = ipv4.get("VlanId") or ipv4.get("VLANID") or 0

            if switch_type == "BMC" and vlan_id != 125:
                continue
            if vlan_id == 0:
                continue

            vlan_entry = {
                "name": ipv4.get("Name"),
                "vlan_id": vlan_id,
            }

            ip_present = False
            for assignment in ipv4.get("Assignment", []):
                name = assignment.get("Name", "").upper()
                ip = assignment.get("IP")

                if name == "GATEWAY":
                    vlan_entry["gateway"] = ip
                elif name == switch_type.upper() or name.startswith(switch_type.upper()):
                    vlan_entry["ip"] = ip
                    ip_present = True

            if ip_present and ipv4.get("Cidr"):
                vlan_entry["cidr"] = ipv4["Cidr"]

            vlans.append(vlan_entry)

        self.sections["vlans"] = sorted(vlans, key=lambda v: v["vlan_id"])

        if not self.sections["vlans"]:
            print("[!] No VLANs found in input data.")

    def build_interfaces(self):
        self.sections["interfaces"] = []

    def build_bgp(self):
        self.sections["bgp"] = []

    def generate(self, switch_type):
        self.build_switch()
        self.build_vlans(switch_type)
        self.build_interfaces()
        self.build_bgp()
        return self.sections


def convert_switch_input_json(input_data, output_dir="output"):
    output_path = os.path.abspath(output_dir)
    os.makedirs(output_path, exist_ok=True)

    builder = StandardJSONBuilder(input_data)

    for sw_type in ["TOR1", "TOR2", "BMC"]:
        result = builder.generate(sw_type)
        switch = builder.sections["switch"].get(sw_type)
        if not switch:
            continue

        hostname = switch.get("hostname", sw_type)
        content = {
            "switch": switch,
            "vlans": result.get("vlans", []),
            "interfaces": result.get("interfaces", []),
            "bgp": result.get("bgp", [])
        }

        out_path = os.path.join(output_path, f"{hostname.lower()}.json")
        with open(out_path, "w", encoding="utf-8") as f:
            json.dump(content, f, indent=2)

        print(f"[✓] Output {hostname} ({sw_type}) config to: {out_path}")
