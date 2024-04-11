#!/bin/bash
# /etc/telegraf/custom_metrics.sh
bgp_total_nbrs=$(sudo vtysh -c "show ip bgp summary" | grep "Total number of neighbors" | awk '{print $NF}')
echo "# HELP bgp_total_nbrs The total number of BGP neighbors in frr"
echo "# TYPE bgp_total_nbrs gauge"
echo "bgp_total_nbrs $bgp_total_nbrs"


echo "# HELP bw_rule_by_class TC Rule for BW by ENV"
echo "# TYPE bw_rule_by_class gauge"
intfs=("eth0" "gre1" "gre2")

for intf in "${intfs[@]}"; do
    while read -r line; do
        if [[ $line =~ "class htb" ]]; then
            read class rate <<< $(echo "$line" | awk '{split($3, a, ":"); print a[2], $8}')
        fi
    done <<< "$(sudo tc class show dev "$intf")"

    # Export class and rate as a single Prometheus metric
    echo "bw_rule_by_class{intf=\"$intf\", class=\"$class\", rate=\"$rate\"} $class"
done