#!/bin/bash
# /etc/telegraf/monitor_frr_bgp.sh

# Function to disable static redistribute
disable_static_redistribute() {
    sudo vtysh -c "configure terminal" -c "router bgp" -c "address-family ipv4 unicast" -c "no redistribute static" -c "end" > /dev/null 2>&1
}

# Function to enable static redistribute
enable_static_redistribute() {
    sudo vtysh -c "configure terminal" -c "router bgp" -c "address-family ipv4 unicast" -c "redistribute static" -c "end" > /dev/null 2>&1
}

# Read the JSON file
bgp_json=$(sudo vtysh -c 'show ip bgp summary json')

echo "# HELP bgp_failed_peers total failed peer nbr number"
echo "# TYPE bgp_failed_peers gauge"
failedPeers=$(echo "$bgp_json" | jq -r '.ipv4Unicast.failedPeers')
echo "bgp_failed_peers $failedPeers"
# Disable Static Redistribute if $failedPeers > 0
if [ "$failedPeers" -gt 0 ]; then
    disable_static_redistribute
fi

echo "# HELP bgp_nbr_uptime_ms uptime of BGP neighbors in frr"
echo "# TYPE bgp_nbr_uptime_ms counter"

echo "# HELP bgp_nbr_pfxSnt sent prefix routes of BGP neighbors in frr"
echo "# TYPE bgp_nbr_pfxSnt gauge"

peers=$(echo "$bgp_json" | jq -r '.ipv4Unicast.peers | to_entries[] | "\(.key) \(.value.peerUptimeMsec) \(.value.state) \(.value.pfxSnt) \(.value.desc)"')
# Iterate over each peer in the variable
echo "$peers" | while read nbrip peerUptimeMsec state pfxSnt desc; do
    echo "bgp_nbr_uptime_ms{nbrip=\"$nbrip\", state=\"$state\"} $peerUptimeMsec"
    echo "bgp_nbr_pfxSnt{nbrip=\"$nbrip\", desc=\"$desc\"} $pfxSnt"
    # "$pfxSnt" -eq 1 for uplink nbr means static redistribute is disabled
    if [ "$failedPeers" -eq 0 ] && [ "$pfxSnt" -eq 1 ] && [[ "$desc" == *"To_Uplink"* ]]; then
        enable_static_redistribute
    fi
done