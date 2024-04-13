#!/bin/bash
# /etc/telegraf/custom_metrics.sh

echo "# HELP bgp_nbr_uptime_ms uptime of BGP neighbors in frr"
echo "# TYPE bgp_nbr_uptime_ms counter"
sudo vtysh -c 'show ip bgp neigh json' | jq -r 'to_entries[] | "\(.key) \(.value.bgpTimerUpMsec)"' | while read nbrip bgpTimerUpMsec; do
    # This line runs the command, pipes the output into jq for processing
    # The -r option tells jq to output raw strings, not JSON-encoded strings
    # The 'to_entries[]' part converts the JSON object into an array of key-value pairs
    # The '| "\(.nbrip) \(.value.bgpTimerUpMsec)"' part formats the output as a string with the nbrip and bgpTimerUpMsec value
    # The output is then piped into a while loop
    echo "bgp_nbr_uptime_ms{nbrip=\"$nbrip\"} $bgpTimerUpMsec"
done

echo "# HELP tc_bw_rate_bits TC Rule for BW by ENV"
echo "# TYPE tc_bw_rate_bits gauge"
intfs=("eth0" "gre1" "gre2")
for intf in "${intfs[@]}"; do
    while read -r line; do
        if [[ $line =~ "class htb" ]]; then
            read class rate <<< $(echo "$line" | awk '{split($3, a, ":"); print a[2], $8}')
            # Extract the numeric value and unit from the rate
            value=$(echo $rate | grep -oP '\d+')
            unit=$(echo $rate | grep -oP '[GMK]bit')
            # Convert the value to bits if it's in Gbits
            if [[ $unit == "Gbit" ]]; then
                value=$((value * 1000 * 1000 * 1000))
            # Convert the value to bits if it's in Mbits
            elif [[ $unit == "Mbit" ]]; then
                value=$((value * 1000 * 1000))
            # Convert the value to bits if it's in Kbits
            elif [[ $unit == "Kbit" ]]; then
                value=$((value * 1000))
            fi
            # Update the rate with the new value
            rate="${value}"
        fi
        # Export class and rate as a single Prometheus metric
        echo "tc_bw_rate_bits{intf=\"$intf\", class=\"$class\", rate=\"$rate\"} $rate"
    done <<< "$(sudo tc class show dev "$intf")"
done