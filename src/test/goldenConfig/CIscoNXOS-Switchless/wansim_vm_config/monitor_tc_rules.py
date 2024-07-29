# /etc/telegraf/monitor_tc_rules.py

import re
import subprocess

# sudo tc qdisc show dev eth0
# qdisc netem 8003: root refcnt 65 limit 1000 delay 50ms loss 1% rate 50Mbit

intfs = ["eth0", "gre1"]
tc_rule_dict = {}

for intf in intfs:

    qdisc_output = subprocess.check_output(f"sudo tc qdisc show dev {intf}", shell=True).decode('utf-8')
    lines = qdisc_output.split('\n')

    for line in lines:
        if "qdisc netem" in line:
            pattern = r'qdisc netem (\d+): root (?:refcnt \d+ )?limit \d+(?: delay (\d+ms))?(?: loss (\d+%))? rate (\d+[GMK]?bit)'
            matches = re.findall(pattern, line)
            qdisc_id = matches[0][0] if matches else None
            delay = next((matches[0][1] for match in matches if matches[0][1]), 0)
            loss = next((matches[0][2] for match in matches if matches[0][2]), 0)
            rate_and_unit = matches[0][3]
            rate_match = re.match(r'(\d+)([GMK]?bit)', rate_and_unit)
            if rate_match:
                rate_value = int(rate_match.group(1))  # Extract numeric part
                rate_unit = rate_match.group(2)  # Extract rate_unit (e.g., 'Gbit', 'Mbit', 'Kbit', 'bit')
                if rate_unit == "Gbit":
                    rate_value *= 1000
                elif rate_unit == "Mbit":
                    rate_value = rate_value
                elif rate_unit == "Kbit":
                    rate_value /= 1000
                rate = rate_value

                if intf == "eth0":
                    tc_rule_dict["upload_bw"] = rate
                    tc_rule_dict["upload_delay"] = delay
                    tc_rule_dict["upload_loss"] = loss
                elif intf == "gre1":
                    tc_rule_dict["download_delay"] = delay
                    tc_rule_dict["download_loss"] = loss
                    tc_rule_dict["download_bw"] = rate

print("# HELP tc_qdisc_show NETEM Rule upload_bw Mbit, download_bw Mbit, loss %, delay ms")
print("# TYPE tc_qdisc_show counter")
upload_bw = tc_rule_dict["upload_bw"]
download_bw = tc_rule_dict["download_bw"]
upload_delay = tc_rule_dict["upload_delay"] if "upload_delay" in tc_rule_dict else 0
download_delay = tc_rule_dict["download_delay"] if "download_delay" in tc_rule_dict else 0
upload_loss = tc_rule_dict["upload_loss"] if "upload_loss" in tc_rule_dict else 0
download_loss = tc_rule_dict["download_loss"] if "download_loss" in tc_rule_dict else 0
print(f"tc_qdisc_show{{upload_bw=\"{upload_bw}\", download_bw=\"{download_bw}\", upload_loss=\"{upload_loss}\", upload_delay=\"{upload_delay}\", download_loss=\"{download_loss}\", download_delay=\"{download_delay}\"}} 1")