# sudo rm -rf /etc/wansimrule/
# sudo tc qdisc del dev eth0 root
# sudo tc qdisc del dev gre1 root
# sudo tc qdisc del dev gre2 root
# sudo rm /etc/wansimrule/wansimrule_service.sh

sudo mkdir -p /etc/wansimrule/
sudo chmod 777 /etc/wansimrule/
sudo cp /home/wansimadmin/wansim_vm_config/wansimrule_service.sh /etc/wansimrule/wansimrule_service.sh
sudo chmod +x /etc/wansimrule/wansimrule_service.sh
sudo nano /etc/systemd/system/wansimrule.service