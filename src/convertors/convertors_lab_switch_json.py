import json
from copy import deepcopy
from pathlib import Path

# ── Static config ─────────────────────────────────────────────────────────
TARGET_SWITCH_TYPES          = ["TOR1", "TOR2", "BMC"]
TOR1, TOR2                   = "TOR1", "TOR2"
BMC                          = "BMC"

CISCO, NXOS                  = "cisco", "nxos"
DELL,  OS10                  = "dell",  "os10"

OUTPUT_FILE_EXTENSION        = ".json"
DEFAULT_OUTPUT_DIR           = "output"

JUMBO_MTU                    = 9216
REDUNDANCY_PRIORITY_ACTIVE   = 150
REDUNDANCY_PRIORITY_STANDBY  = 140
HSRP                         = "hsrp"
VRRP                         = "vrrp"

UNUSED_VLAN       = "UNUSED_VLAN"  
NATIVE_VLAN      = "NATIVE_VLAN"


# ── In-code templates (never mutate these!) ───────────────────────────────
SWITCH_TEMPLATE = {
    "make"    : "",
    "model"   : "",
    "type"    : "",
    "hostname": "",
    "version" : "",
    "firmware": ""
}

SVI_TEMPLATE = {
    "ip"        : "",
    "cidr"      : 0,
    "mtu"       : JUMBO_MTU,
    "redundancy": {          # this whole block is removed for BMC
        "type"      : "",
        "group"     : 0,
        "priority"  : 0,
        "virtual_ip": ""
    }
}

VLAN_TEMPLATE = {
    "vlan_id" : 0,
    "name"    : "",
    # "interface" key is added later only if needed
}

# ── Builder class ─────────────────────────────────────────────────────────
class StandardJSONBuilder:
    def __init__(self, input_data: dict):
        self.input_data = input_data
        self.sections   = {}
        self.vlan_map = {}

    # ------------------------------------------------------------------ #
    # Build switch section
    # ------------------------------------------------------------------ #
    def build_switch(self, target_types=None):
        if target_types is None:
            target_types = set(TARGET_SWITCH_TYPES)

        switches_json = self.input_data.get("InputData", {}).get("Switches", [])
        self.sections["switch"] = {}

        for sw in switches_json:
            sw_type = sw.get("Type")
            if sw_type not in target_types:
                continue

            sw_make = sw.get("Make", "").lower()
            firmware = (
                NXOS if sw_make == CISCO else
                OS10 if sw_make == DELL  else
                sw.get("Firmware", "").lower()
            )

            sw_entry = deepcopy(SWITCH_TEMPLATE)
            sw_entry.update(
                make     = sw_make,
                model    = sw.get("Model", "").lower(),
                type     = sw_type,
                hostname = sw.get("Hostname", "").lower(),
                version  = sw.get("Firmware", "").lower(),
                firmware = firmware
            )

            self.sections["switch"][sw_type] = sw_entry

        if not self.sections["switch"]:
            print("[!] No valid switches found in input.")

    # ------------------------------------------------------------------ #
    # Build VLAN section
    # ------------------------------------------------------------------ #
    def build_vlans(self, switch_type: str):
        vlans_out = []
        supernets = self.input_data.get("InputData", {}).get("Supernets", [])

        for net in supernets:
            group_name  = net.get("GroupName", "").upper()
            ipv4        = net.get("IPv4", {})
            vlan_id     = ipv4.get("VlanId") or ipv4.get("VLANID") or 0
            if vlan_id == 0:
                continue                                # skip invalid IDs

            # BMC: only keep interface for GroupName "BMC" or "UNUSED_VLAN" or "NATIVEVLAN"
            if switch_type == BMC and group_name.upper() not in (BMC, UNUSED_VLAN, NATIVE_VLAN):
                continue

            # Construct M,C,S Mapping, GroupName hardcode defined
            if group_name.startswith("HNVPA"):
                self.vlan_map["HNVPA"].append(vlan_id)
            elif group_name.startswith("INFRA"):
                self.vlan_map["M"].append(vlan_id)
            elif group_name.startswith("TENANT"):
                self.vlan_map["C"].append(vlan_id)
            elif group_name.startswith("L3FORWARD"):
                self.vlan_map["C"].append(vlan_id)

            # collect IP / GW
            ip          = ""
            virtual_ip  = ""
            for assign in ipv4.get("Assignment", []):
                a_name = assign.get("Name", "").upper()
                a_ip   = assign.get("IP", "")
                if a_name == "GATEWAY":
                    virtual_ip = a_ip
                elif a_name == switch_type.upper() or a_name.startswith(switch_type.upper()):
                    ip = a_ip

            vlan_entry = deepcopy(VLAN_TEMPLATE)
            vlan_entry.update(vlan_id=vlan_id, name=ipv4.get("Name"))

            # add interface only if IP present & non-blank
            if ip.strip() and ipv4.get("Cidr"):
                iface = deepcopy(SVI_TEMPLATE)
                iface["ip"]   = ip
                iface["cidr"] = ipv4["Cidr"]

                if switch_type in (TOR1, TOR2):
                    sw_make = self.sections["switch"][switch_type]["make"]
                    redundancy_type = HSRP if sw_make == CISCO else VRRP
                    priority        = (
                        REDUNDANCY_PRIORITY_ACTIVE if switch_type == TOR1
                        else REDUNDANCY_PRIORITY_STANDBY
                    )

                    iface["redundancy"].update(
                        type       = redundancy_type,
                        group      = vlan_id,
                        priority   = priority,
                        virtual_ip = virtual_ip
                    )
                else:
                    # BMC: remove redundancy block entirely
                    iface.pop("redundancy", None)

                vlan_entry["interface"] = iface    # attach interface

            vlans_out.append(vlan_entry)

        self.sections["vlans"] = sorted(vlans_out, key=lambda v: v["vlan_id"])
        if not self.sections["vlans"]:
            print(f"[!] No VLANs produced for {switch_type}")

    # ------------------------------------------------------------------ #
    # Build interface and port-channels section
    # ------------------------------------------------------------------ #
    def build_interfaces(self, switch_type: str):
        """
        Load and build both interfaces and port-channels from model-based templates.
        """
        switch = self.sections["switch"].get(switch_type)
        if not switch:
            print(f"[!] No switch of type {switch_type} found for interface build.")
            return

        make = switch.get("make", "").lower()
        model = switch.get("model", "").lower()
        template_path = Path("input/switch_interface_templates") / make / f"{model}.json"

        if not template_path.exists():
            raise FileNotFoundError(f"[!] Interface template not found: {template_path}")

        with open(template_path) as f:
            template_data = json.load(f)

        self._build_interfaces_from_template(template_data)
        self._build_port_channels_from_template(template_data)

    def _build_interfaces_from_template(self, template_data: dict):
        vlan_map = {
            "M": "7",
            "C": "6,201,301,401",
            "S1": "501-516",
            "S2": "711,712"
        }

        def expand_vlans(raw):
            if not raw:
                return ""
            return ",".join([
                vlan_map.get(part.strip(), part.strip())
                for part in raw.replace("|", ",").split(",")
            ])

        templates = template_data.get("interface_templates", {})
        profile = templates.get("FullyConverged", [])
        common = templates.get("Common", {})

        interfaces = []
        for name in ["Unused", "Trunk_TO_BMC_SWITCH", "Loopback0", "P2P_Border2", "P2P_Border1"]:
            iface = deepcopy(common[name])
            if name == "Loopback0":
                iface["ipv4"] = "100.71.85.21/32"
            elif name == "P2P_Border1":
                iface["ipv4"] = "100.71.85.2/30"
            elif name == "P2P_Border2":
                iface["ipv4"] = "100.71.85.10/30"
            interfaces.append(iface)

        for tmpl in profile:
            new_iface = deepcopy(tmpl)
            new_iface["name"] = "Hyperconverged"
            new_iface["native_vlan"] = vlan_map.get(new_iface["native_vlan"], new_iface["native_vlan"])
            new_iface["tagged_vlans"] = expand_vlans(new_iface["tagged_vlans"])
            interfaces.insert(1, new_iface)

        self.sections["interfaces"] = interfaces

    def _build_port_channels_from_template(self, template_data: dict):
        port_channels = template_data.get("port_channels", [])
        enriched_pcs = []

        for pc in port_channels:
            pc_copy = deepcopy(pc)

            # Enrich with values based on ID
            if pc_copy["id"] == 50:
                pc_copy["ipv4"] = "100.71.85.17/30"
            if pc_copy["id"] == 101:
                pc_copy["description"] = "L2_LACP_Peer"
                pc_copy["tagged_vlans"] = "7,125"

            enriched_pcs.append(pc_copy)

        self.sections["port_channels"] = enriched_pcs

    # ------------------------------------------------------------------ #
    # Build BGP section
    # ------------------------------------------------------------------ #
    def build_bgp(self, switch_type: str):
        """
        Build BGP config with neighbor and network structure.
        IPs are abstracted with placeholders for future parsing.
        """
        switch = self.sections["switch"].get(switch_type)
        if not switch:
            print(f"[!] No switch info for BGP")
            return

        bgp = {
            "asn": 65242,
            "router_id": "{{ loopback_ip }}",  # Replace with parsed loopback IP
            "networks": [
                "{{ p2p_border1_subnet }}",
                "{{ p2p_border2_subnet }}",
                "{{ ibgp_subnet }}",
                "{{ loopback_ip }}/32",
                "{{ pod_pool_prefix }}"
            ],
            "neighbors": [
                {
                    "ip": "{{ p2p_border1_ip }}",
                    "description": "TO_Border1",
                    "remote_as": 64846,
                    "af_ipv4_unicast": {
                        "prefix_list_in": "DefaultRoute"
                    }
                },
                {
                    "ip": "{{ p2p_border2_ip }}",
                    "description": "TO_Border2",
                    "remote_as": 64846,
                    "af_ipv4_unicast": {
                        "prefix_list_in": "DefaultRoute"
                    }
                },
                {
                    "ip": "{{ ibgp_peer_ip }}",
                    "description": "iBGP_PEER",
                    "remote_as": 65242,
                    "af_ipv4_unicast": {}
                },
                {
                    "ip": "{{ hnvpa_subnet }}",
                    "description": "TO_HNVPA",
                    "remote_as": 65112,
                    "update_source": "Loopback0",
                    "ebgp_multihop": 3,
                    "af_ipv4_unicast": {
                        "prefix_list_out": "DefaultRoute"
                    }
                }
            ]
        }

        self.sections["bgp"] = bgp


    # ------------------------------------------------------------------ #
    # Build Prefix-List section (Hardcoded for now)
    # ------------------------------------------------------------------ #
    def build_prefix_lists(self):
        """
        Build basic prefix-lists. Values can be reused or extended easily.
        """
        self.sections["prefix_lists"] = {
            "DefaultRoute": [
                {
                    "seq": 10,
                    "action": "permit",
                    "prefix": "0.0.0.0/0"
                },
                {
                    "seq": 50,
                    "action": "deny",
                    "prefix": "0.0.0.0/0",
                    "prefix_filter": "le 32"
                }
            ]
        }

    # ------------------------------------------------------------------ #
    # Build QoS section (Hardcoded for now)
    # ------------------------------------------------------------------ #
    def build_qos(self):
        """
        Build basic QoS policies. Values can be reused or extended easily.
        """
        self.sections["qos"] = True


    # ------------------------------------------------------------------ #
    # Generate all sections for a specific switch type
    def generate_for_switch(self, switch_type):
        self.build_switch()
        self.build_vlans(switch_type)
        self.build_interfaces(switch_type)
        return self.sections



# ── Helper to create per-switch JSON files ────────────────────────────────
def convert_switch_input_json(input_data: dict, output_dir: str = DEFAULT_OUTPUT_DIR):
    out_path = Path(output_dir)
    out_path.mkdir(exist_ok=True, parents=True)

    builder = StandardJSONBuilder(input_data)

    for sw_type in TARGET_SWITCH_TYPES:
        builder.sections = {}               # reset between runs
        builder.build_switch()
        builder.build_vlans(sw_type)
        builder.build_interfaces(sw_type)
        builder.build_bgp(sw_type)
        builder.build_prefix_lists()
        builder.build_qos()
        
        sw_info = builder.sections.get("switch", {}).get(sw_type)
        if not sw_info:                     # no matching switch of this type
            continue

        hostname = sw_info["hostname"] or sw_type.lower()
        final_json = {
            "switch": sw_info,
            "vlans" : builder.sections["vlans"],
            "interfaces": builder.sections["interfaces"],
            "port_channels": builder.sections["port_channels"],
            "bgp": builder.sections["bgp"],
            "prefix_lists": builder.sections["prefix_lists"],
            "qos": builder.sections["qos"]
        }

        out_file = out_path / f"{hostname}{OUTPUT_FILE_EXTENSION}"
        with out_file.open("w", encoding="utf-8") as f:
            json.dump(final_json, f, indent=2)

        print(f"[✓] Wrote {out_file}")
