# HackingLab


The goal of this project is to extract possible vunerabilities in routers on the routes to destinations given by domain names. To do this we created a pipeline including a series of scripts that in order:

1. Extracts the list of IPs of routers from the traceroutes on the domains
2. Send snmp requests and wait for answers that contain EngineID (an identifier for the router) and last reboot time (a good identifier for the last software update)
3. Parse EngineIDs into a vendor name: as most engineIDs have a specific format that ties them to a vendor by including the Private Enterprise Number (PEN) as assigned by the Internet Assigned Numbers Authority (IANA) in their code
4. For each EngineID from which we managed to extract the vendor (step 3) and reboot time (step 2) extract CVEs listed after the reboot time for routers of the found vendor.

## How to use 

### Add input
Add your list of domains that you want to test in `./input/domains.txt`.
Each line represents a domain. Example:
```
google.com
youtube.com
facebook.com
twitter.com
wikipedia.org
```
### Run
To run the code with `make`

### Check outputs
The code outputs 2 files (`./output/engine_to_cves.csv` and `./output/ips_to_cves.csv`) and a folder (`./output/cves_per_engId`)

1. `./output/engine_to_cves.csv` contains for each EngineId a max severity (highest severity among all the correlated CVEs), CVEs with max severity (IDs of these cves), All CVEs (all correlated cves), IPs correlated (IPs that returned that specific EngineId)
2. `./output/ips_to_cves.csv` contains the same as above but thus time each tested IP with its CVEs
3. `./output/cves_per_engId` contains for each EngineID a folder with two files:
    * `cves_important_info.json` important information from each CVE extracted (CVE ID, severity using different metrics) 
    * `cves_full_format.json` full CVEs extracted from EngineID

