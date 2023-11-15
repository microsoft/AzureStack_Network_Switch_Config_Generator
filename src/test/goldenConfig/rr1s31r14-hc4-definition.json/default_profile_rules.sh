# Upload BW 1Gbit for all
sudo tc qdisc add dev eth0 root handle 1a1a: htb default 10
sudo tc class add dev eth0 parent 1a1a: classid 1a1a:1 htb rate 1gbit
sudo tc class add dev eth0 parent 1a1a:1 classid 1a1a:10 htb rate 1gbit
sudo tc filter add dev eth0 protocol ip parent 1a1a: prio 1 u32 match ip dst 0.0.0.0/0 flowid 1a1a:10

# Download BW 1Gbit for all
# TC Rule for gre1
sudo tc qdisc add dev gre1 root handle 1a1a: htb default 10
sudo tc class add dev gre1 parent 1a1a: classid 1a1a:1 htb rate 1gbit
sudo tc class add dev gre1 parent 1a1a:1 classid 1a1a:10 htb rate 1gbit
sudo tc filter add dev gre1 protocol ip parent 1a1a: prio 1 u32 match ip dst 0.0.0.0/0 flowid 1a1a:10
# TC Rule for gre2
sudo tc qdisc add dev gre2 root handle 1a1a: htb default 10
sudo tc class add dev gre2 parent 1a1a: classid 1a1a:1 htb rate 1gbit
sudo tc class add dev gre2 parent 1a1a:1 classid 1a1a:10 htb rate 1gbit
sudo tc filter add dev gre2 protocol ip parent 1a1a: prio 1 u32 match ip dst 0.0.0.0/0 flowid 1a1a:10

# Check All Rules
# sudo tc -s class show dev eth0
# sudo tc -s class show dev gre1
# sudo tc -s class show dev gre2

# Remove All Rules
# sudo tc qdisc del dev eth0 root
# sudo tc qdisc del dev gre1 root
# sudo tc qdisc del dev gre2 root