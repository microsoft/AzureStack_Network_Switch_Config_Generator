# /etc/telegraf/monitor_tc_rules.py

import re
import subprocess

# sudo tc class show dev gre1
# class htb 1a1a:50 root leaf 50: prio 0 rate 10Mbit ceil 10Mbit burst 1600b cburst 1600b 
# class htb 1a1a:50 root prio 0 rate 5Mbit ceil 5Mbit burst 1600b cburst 1600b 
# class htb 1a1a:1 root prio 0 rate 1Gbit ceil 1Gbit burst 1375b cburst 1375b 
intfs = ["eth0", "gre1"]
tc_rule_dict = {}

for intf in intfs:
    # BW
    class_output = subprocess.check_output(f"sudo tc class show dev {intf}", shell=True).decode('utf-8')
    lines = class_output.split('\n')

    for line in lines:
        if "class htb" in line:
            class_id, rate = re.search(r'(\w+:\d+).*rate (\d+[GMK]bit)', line).groups()
            value = int(re.search(r'\d+', rate).group())
            unit = re.search(r'[GMK]bit', rate).group()
            if unit == "Gbit":
                value *= 1000
            elif unit == "Mbit":
                value = value
            elif unit == "Kbit":
                value /= 1000

            rate = value
            env_id = 0
            if int(class_id.split(':')[1]) == 1:
                env_id = 0
            elif int(class_id.split(':')[1]) >= 10:
                env_id = int(class_id.split(':')[1]) // 10

            if env_id not in tc_rule_dict:
                tc_rule_dict[env_id]= {
                    "env_id": env_id,
                }

            if intf == "eth0":
                tc_rule_dict[env_id]["upload_bw"] = rate
            elif intf == "gre1":
                tc_rule_dict[env_id]["download_bw"] = rate

    # Delay and Loss
    qdisc_output = subprocess.check_output(f"sudo tc qdisc show dev {intf}", shell=True).decode('utf-8')
    lines = qdisc_output.split('\n')

    for line in lines:
        if "qdisc netem" in line:
            pattern = r'netem (\d+):(?:.*delay (\d+)ms)?(?:.*loss (\d+)%)?'
            matches = re.findall(pattern, line)
            class_id = matches[0][0] if matches else None
            delay = next((matches[0][1] for match in matches if matches[0][1]), 0)
            loss = next((matches[0][2] for match in matches if matches[0][2]), 0)
            env_id = int(class_id) // 10


            if intf == "eth0":
                tc_rule_dict[env_id]["upload_delay"] = delay
                tc_rule_dict[env_id]["upload_loss"] = loss
            elif intf == "gre1":
                tc_rule_dict[env_id]["download_delay"] = delay
                tc_rule_dict[env_id]["download_loss"] = loss

print("# HELP tc_rule_by_env TC Rule by ENV, upload_bw bit, download_bw bit, loss %, delay ms")
print("# TYPE tc_rule_by_env counter")
# Iterate tc_rule_dict and print value
for env_id, tc_rule in tc_rule_dict.items():
    upload_bw = tc_rule["upload_bw"]
    download_bw = tc_rule["download_bw"]
    upload_delay = tc_rule["upload_delay"] if "upload_delay" in tc_rule else 0
    download_delay = tc_rule["download_delay"] if "download_delay" in tc_rule else 0
    upload_loss = tc_rule["upload_loss"] if "upload_loss" in tc_rule else 0
    download_loss = tc_rule["download_loss"] if "download_loss" in tc_rule else 0
    print(f"tc_rule_by_env{{env=\"{env_id}\", upload_bw=\"{upload_bw}\", download_bw=\"{download_bw}\", upload_loss=\"{upload_loss}\", upload_delay=\"{upload_delay}\", download_loss=\"{download_loss}\", download_delay=\"{download_delay}\"}} 1")