import sys
import nmap
import requests
import argparse
import json


vulnDbApiKey = 'a11673aa01708a6b49b86e8d9f1ec8f7'


"""
Example of result:

[{'name': 'Cisco 870 router or 2960 switch (IOS 12.2 - 12.4)', 'accuracy': '100', 'line': '16424', 'osclass': 
[{'type': 'router', 'vendor': 'Cisco', 'osfamily': 'IOS', 'osgen': '12.X', 'accuracy': '100', 'cpe': 
['cpe:/h:cisco:870_router', 'cpe:/o:cisco:ios:12']}, {'type': 'switch', 'vendor': 'Cisco', 'osfamily': 'IOS', 
'osgen': '12.X', 'accuracy': '100', 'cpe': ['cpe:/h:cisco:2960_switch', 'cpe:/o:cisco:ios:12']}]},

 {'name': 'Cisco Aironet 1250 WAP (IOS 12.4) or IOS XE', 'accuracy': '100', 'line': '18050', 
 'osclass': [{'type': 'WAP', 'vendor': 'Cisco', 'osfamily': 'IOS', 'osgen': '12.X', 'accuracy': '
 100', 'cpe': ['cpe:/h:cisco:aironet_ap1250', 'cpe:/o:cisco:ios:12.4']}, 
 {'type': 'router', 'vendor': 'Cisco', 'osfamily': 'IOS XE', 'osgen': None, 'accuracy': '100', 'cpe': ['cpe:/o:cisco:ios_xe']}]}]

 Example of search string for vulndb: 'Cisco IOS 12.2 - 12.4'

@todo do a vuln search for all possible vendors/os versions.
"""
def perform_nmap_os_scan(target, asVulnDBSearch=True):
    nm = nmap.PortScanner()
    nm.scan(hosts=target, arguments='-O')

    if asVulnDBSearch:
        return nm[target]['osmatch'][0]["osclass"][0]["vendor"] + " " + nm[target]['osmatch'][0]["osclass"][0]["osfamily"] + " " + nm[target]['osmatch'][0]["osclass"][0]["osgen"] 
    return nm[target]['osmatch'][0]

def searchVuln(search):
    userAgent = 'VulDB API Advanced Python Demo Agent'
    headers = {'User-Agent': userAgent, 'X-VulDB-ApiKey': vulnDbApiKey}

    url = 'https://vuldb.com/?api'

    postData = {'search': search}

    response = requests.post(url,headers=headers,data=postData)

    if response.status_code == 200:
        responseJson = json.loads(response.content)	
        for i in responseJson['result']:		
            print(i['entry'])


if len(sys.argv) != 2:
    print("Usage: python os_fetch.py <ip_address>")
    sys.exit(1)

target_platform = perform_nmap_os_scan(sys.argv[1], True)

print(f"PLATFORM DETECTED: {target_platform}")
searchVuln(target_platform);


