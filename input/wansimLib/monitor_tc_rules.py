# /etc/telegraf/monitor_tc_rules.py
import re
import subprocess

print("# HELP tc_bw_rate_bits TC Rule for BW by ENV")
print("# TYPE tc_bw_rate_bits gauge")
# sudo tc class show dev gre1
# class htb 1a1a:50 root leaf 50: prio 0 rate 10Mbit ceil 10Mbit burst 1600b cburst 1600b 
# class htb 1a1a:50 root prio 0 rate 5Mbit ceil 5Mbit burst 1600b cburst 1600b 
# class htb 1a1a:1 root prio 0 rate 1Gbit ceil 1Gbit burst 1375b cburst 1375b 
intfs = ["eth0", "gre1", "gre2"]
for intf in intfs:
    output = subprocess.check_output(f"sudo tc class show dev {intf}", shell=True).decode('utf-8')
    lines = output.split('\n')

    for line in lines:
        if "class htb" in line:
            class_id, rate = re.search(r'(\w+:\d+).*rate (\d+[GMK]bit)', line).groups()
            value = int(re.search(r'\d+', rate).group())
            unit = re.search(r'[GMK]bit', rate).group()
            if unit == "Gbit":
                value *= 1000 * 1000 * 1000
            elif unit == "Mbit":
                value *= 1000 * 1000
            elif unit == "Kbit":
                value *= 1000

            rate = value
            env_id = 0
            if int(class_id.split(':')[1]) == 1:
                env_id = 0
            elif int(class_id.split(':')[1]) >= 10:
                env_id = int(class_id.split(':')[1]) // 10
            print(f"tc_bw_rate_bits{{intf=\"{intf}\", env_id=\"{env_id}\", rate=\"{rate}\"}} {rate}")