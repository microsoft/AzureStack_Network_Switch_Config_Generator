#!/bin/bash
#####
# Pre-Config
# # Config Interface
# sudo cp ./30_netplan.yaml /etc/netplan/00-installer-config.yaml
# sudo netplan apply

# # Config FRR
# sudo cp ./daemons /etc/frr/daemons
# sudo cp ./frr.conf /etc/frr/frr.conf
# sudo service frr restart

# # Config Log
# sudo tr -d '\r' < ./wansim_log.sh | sudo tee ./wansim_log.sh > /dev/null
# # Execute Every Half Hour
# # Minutes (0,30): The script will run when the minute is either 0 or 30.
# Define using crobtab
# sudo crontab -e
# MAILTO=""
# 0,30 * * * *  sudo /home/wansimadmin/wansim_config/wansim_log.sh
#####

# Define the function
function convert_hex_to_ipv4() {
    # Read the input from standard input into a variable
    file_contents=$(cat)

    # Use sed to extract the hex values
    hex_values=$(echo "$file_contents" | sed -n "s|.*\([0-9a-f]\{8\}\/[0-9a-f]\{8\}\).*|\1|p")

    # Loop through the hex values and convert them to IPv4 format
    for hex_value in $hex_values; do
        # Extract the two hex values from the string
        hex1=${hex_value:0:8}
        hex2=${hex_value:9:8}

        # Convert the hex values to decimal and format as an IPv4 address
        ipv4=$(printf "%d.%d.%d.%d" 0x${hex1:0:2} 0x${hex1:2:2} 0x${hex1:4:2} 0x${hex1:6:2})
        ipv4+="/"
        ipv4+=$(printf "%d.%d.%d.%d" 0x${hex2:0:2} 0x${hex2:2:2} 0x${hex2:4:2} 0x${hex2:6:2})

        # Replace the hex value with the IPv4 address in the file contents
        file_contents=$(echo "$file_contents" | sed "s|$hex_value|$ipv4|g")
    done

    # Output the modified file contents
    echo "$file_contents"
}

# START
logger -t WANSIM "WANSIM LOG - START"
# FRR
echo "show ip bgp summary" | sudo vtysh | logger -t WANSIM-FRR-BGP
echo "show ip route" | sudo vtysh | logger -t WANSIM-FRR-Route
# NETEM
# sudo tc class show dev eth0 | sed 's/^/eth0 - /' | logger -t WANSIM-NETEM
# for interface in eth0 gre1 gre2; do
for interface in eth0 gre1 gre2; do
    sudo tc class show dev $interface | sed "s/^/$interface - /" | logger -t WANSIM-NETEM-Class
    sudo tc filter show dev $interface | convert_hex_to_ipv4 | logger -t WANSIM-NETEM2-Filter
done
# END
logger -t WANSIM "WANSIM LOG - END"
# tail -n 50 /var/log/syslog