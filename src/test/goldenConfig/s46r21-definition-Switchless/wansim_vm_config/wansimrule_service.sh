#!/bin/bash

Start() {
    # Remove All Rules
    sudo tc qdisc del dev eth0 root > /dev/null 2>&1
    sudo tc qdisc del dev gre1 root > /dev/null 2>&1
    sudo tc qdisc del dev gre2 root > /dev/null 2>&1
    # Upload BW 1Gbit for all
    sudo tc qdisc add dev eth0 root handle 1a1a: htb default 1
    sudo tc class add dev eth0 parent 1a1a: classid 1a1a:1 htb rate 1Gbit
    # Download BW 1Gbit for all
    # TC Rule for gre1
    sudo tc qdisc add dev gre1 root handle 1a1a: htb default 1
    sudo tc class add dev gre1 parent 1a1a: classid 1a1a:1 htb rate 1Gbit
    # TC Rule for gre2
    sudo tc qdisc add dev gre2 root handle 1a1a: htb default 1
    sudo tc class add dev gre2 parent 1a1a: classid 1a1a:1 htb rate 1Gbit
    sudo tc qdisc show dev eth0

    # Load all ENV profiles if any
    PROFILES_FOLDER="/etc/wansimrule/profiles/"
    if [ -d "$PROFILES_FOLDER" ]; then
        # Iterate through each profiles in the folder
        for script in "$PROFILES_FOLDER"/*.sh; do
            # Ensure the script has execute permissions
            chmod +x "$script"
            # Run the script
            . "$script"
        done
    else
        echo "$PROFILES_FOLDER is not created."
    fi

    # Validate
    sudo tc class show dev eth0
    sudo tc class show dev gre1
    sudo tc class show dev gre2
}

Stop() {
    # Remove All Rules
    sudo tc qdisc del dev eth0 root > /dev/null 2>&1
    sudo tc qdisc del dev gre1 root > /dev/null 2>&1
    sudo tc qdisc del dev gre2 root > /dev/null 2>&1
}

case $1 in
  Start|Stop) "$1" ;;
esac