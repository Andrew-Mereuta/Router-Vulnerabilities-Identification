from scapy.all import *
import bz2
import json
import re

# This script converts the RIPE Atlas traceroute JSON format to our desired format accepted by the snmp sender

# Regular expression pattern for IPv4 address
ipv4_pattern = re.compile(r'^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$')

# Function to convert the provided JSON format to the desired format
def convert_to_desired_format(json_data):
    converted_data = []
    for item in json_data:
        domain = item.get("dst_addr", "Unknown")  # Use "Unknown" if "dst_addr" is missing
        # Skip domains that are not IPv4 addresses
        if not ipv4_pattern.match(domain):
            continue
        pathIPs = []
        hop_ip_set = set()  # Set to store unique IPs for each hop, because the json data contains 3 entries for each hop
        for hop_info in item.get("result", []):
            for result in hop_info.get("result", []):
                if "from" in result and result["from"] not in hop_ip_set:
                    ip_address = result["from"]
                    if ipv4_pattern.match(ip_address):
                        pathIPs.append({"hop": hop_info.get("hop", -1), "ip": ip_address})
                        hop_ip_set.add(ip_address)  # Add IP to set to ensure uniqueness
        converted_data.append({"domain": domain, "pathIPs": pathIPs})
    return converted_data

with open("./input/traceroute.json", "r") as file:
    input_json = file.readlines()

# Parse JSON data from each line as one traceroute entry is on one line
json_data = []
for line in input_json:
    try:
        json_data.append(json.loads(line.strip()))
    except json.JSONDecodeError:
        print("Error decoding JSON:", line)

converted_data = convert_to_desired_format(json_data)

with open("./input/testing_ips.json", "w") as file:
    json.dump(converted_data, file, indent=2)