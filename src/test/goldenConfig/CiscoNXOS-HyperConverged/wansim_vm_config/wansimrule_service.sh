#!/bin/bash

Start() {
    # Remove All Rules
    sudo tc qdisc del dev eth0 root > /dev/null 2>&1
    sudo tc qdisc del dev gre1 root > /dev/null 2>&1
    sudo tc qdisc del dev gre2 root > /dev/null 2>&1
    # Upload Profile from Definition JSON
    sudo tc qdisc add dev eth0 root netem delay 0ms loss 0% rate 1Gbit
    # Download Profile from Definition JSON
    # TC Rule for gre1
    sudo tc qdisc add dev gre1 root netem delay 0ms loss 0% rate 1Gbit
    # TC Rule for gre2
    sudo tc qdisc add dev gre2 root netem delay 0ms loss 0% rate 1Gbit

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
    sudo tc qdisc show dev eth0
    sudo tc qdisc show dev gre1
    sudo tc qdisc show dev gre2
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