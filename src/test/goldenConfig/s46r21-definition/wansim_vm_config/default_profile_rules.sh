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