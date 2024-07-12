# /etc/telegraf/monitor_frr_bgp.py
import json
import subprocess

def run_command(command):
    process = subprocess.Popen(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    return output.decode('utf-8'), error.decode('utf-8')

def disable_static_redistribute():
    command = 'sudo vtysh -c "configure terminal" -c "router bgp" -c "address-family ipv4 unicast" -c "no redistribute static" -c "end"'
    run_command(command)

def enable_static_redistribute():
    command = 'sudo vtysh -c "configure terminal" -c "router bgp" -c "address-family ipv4 unicast" -c "redistribute static" -c "end"'
    run_command(command)

bgp_json_output, error = run_command("sudo vtysh -c 'show ip bgp summary json'")
bgp_json = json.loads(bgp_json_output)

print("# HELP bgp_failed_peers total failed peer nbr number")
print("# TYPE bgp_failed_peers gauge")
failedPeers = bgp_json['ipv4Unicast']['failedPeers']
print(f"bgp_failed_peers {failedPeers}")

if failedPeers > 0:
    disable_static_redistribute()

print("# HELP bgp_nbr_uptime_ms uptime of BGP neighbors in frr")
print("# TYPE bgp_nbr_uptime_ms counter")

print("# HELP bgp_nbr_pfxSnt sent prefix routes of BGP neighbors in frr")
print("# TYPE bgp_nbr_pfxSnt gauge")

peers = bgp_json['ipv4Unicast']['peers']

for nbrip, data in peers.items():
    peerUptimeMsec = data['peerUptimeMsec']
    state = data['state']
    pfxSnt = data['pfxSnt']
    desc = data['desc']

    print(f"bgp_nbr_uptime_ms{{nbrip=\"{nbrip}\", state=\"{state}\"}} {peerUptimeMsec}")
    print(f"bgp_nbr_pfxSnt{{nbrip=\"{nbrip}\", desc=\"{desc}\"}} {pfxSnt}")

    if failedPeers == 0 and pfxSnt == 1 and "To_Uplink" in desc:
        enable_static_redistribute()