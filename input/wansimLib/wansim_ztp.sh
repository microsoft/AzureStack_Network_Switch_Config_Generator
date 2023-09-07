# Init
sudo apt-get update
sudo apt-get install -y net-tools frr iperf traceroute lldpd
HOSTNAME="WAN-SIM"
sudo hostnamectl set-hostname $HOSTNAME
# sudo nano /etc/hosts

# Config Interface
sudo cp ./30_netplan_wansim.yaml /etc/netplan/30_netplan_wansim.yaml
sudo netplan apply

# Config FRR
sudo cp ./wansim_daemons /etc/frr/daemons
sudo service frr restart
sudo cp ./wansim_frr.conf /etc/frr/frr.conf
sudo service frr restart

# Post-Validation
sudo chmod +x postValidation.sh
sudo ./postValidation.sh