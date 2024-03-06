#!/bin/bash

ips=$(traceroute www.discord.com | grep -oP '\d+\.\d+\.\d+\.\d+')

for ip in $ips; do
    echo "Querying SNMP for IP: $ip"
    ./snmpQuery $ip
done
