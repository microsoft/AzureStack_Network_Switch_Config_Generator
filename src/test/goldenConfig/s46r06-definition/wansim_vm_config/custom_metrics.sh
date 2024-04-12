#!/bin/bash
# /etc/telegraf/custom_metrics.sh

echo "# HELP bgp_total_nbrs The total number of BGP neighbors in frr"
echo "# TYPE bgp_total_nbrs gauge"
bgp_total_nbrs=$(sudo vtysh -c "show ip bgp summary" | grep "Total number of neighbors" | awk '{print $NF}')
echo "bgp_total_nbrs $bgp_total_nbrs"

echo "# HELP tc_bw_kilobyte TC Rule for BW by ENV"
echo "# TYPE tc_bw_kilobyte gauge"
intfs=("eth0" "gre1" "gre2")
for intf in "${intfs[@]}"; do
    while read -r line; do
        if [[ $line =~ "class htb" ]]; then
            read class rate <<< $(echo "$line" | awk '{split($3, a, ":"); print a[2], $8}')
            # Extract the numeric value and unit from the rate
            value=$(echo $rate | grep -oP '\d+')
            unit=$(echo $rate | grep -oP '[GMK]bit')
            # Convert the value to KB if it's in Gbits
            if [[ $unit == "Gbit" ]]; then
                value=$((value * 1000 * 1000 / 8))
            # Convert the value to KB if it's in Mbits
            elif [[ $unit == "Mbit" ]]; then
                value=$((value * 1000 / 8))
            # Convert the value to KB if it's in Kbits
            elif [[ $unit == "Kbit" ]]; then
                value=$((value / 8))
            fi
            # Update the rate with the new value
            rate="${value}"
        fi
        # Export class and rate as a single Prometheus metric
        echo "tc_bw_kilobyte{intf=\"$intf\", class=\"$class\", rate=\"$rate\"} $rate"
    done <<< "$(sudo tc class show dev "$intf")"
done