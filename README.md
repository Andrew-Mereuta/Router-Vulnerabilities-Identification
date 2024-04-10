# Router-Vulnerabilities-Identification

The goal of this project is to extract possible vunerabilities in routers on the routes to destinations given by domain names. To do this we created a pipeline including a series of scripts that in order:

1. Extracts the list of IPs of routers from the traceroutes on the domains
2. Send snmp requests and wait for answers that contain EngineID (an identifier for the router) and last reboot time (a good identifier for the last software update)
3. Parse EngineIDs into a vendor name: as most engineIDs have a specific format that ties them to a vendor by including the Private Enterprise Number (PEN) as assigned by the Internet Assigned Numbers Authority (IANA) in their code
4. For each EngineID from which we managed to extract the vendor (step 3) and reboot time (step 2) extract CVEs listed after the reboot time for routers of the found vendor.

### Active probing - addition to the main pipeline

This projects integrates [LFP](https://github.com/routerfingerprinting/lfp) to increase the number of routers identified via vendor.

We leverage the results from the SNMPv3 technique in order to discover new vendor-specific signatures that we use as an addition to the already existing set of LFP signatures.
This set we further make use of in order to identify routers that could not have been identified neither by SNMPv3 or by LFP with the default signatures dataset.

In "/active_probing" one can find a number of helpfull scripts to ease the integration of [LFP](https://github.com/routerfingerprinting/lfp) with the currents approach:

- setup.py: Parse results from "ip_extractor" and "engine_parse", and extend them with the correct vendor. Also exports the ip addresses that were identified and the ones that were not identified.
- create_signatures.py: Create the list of signatures based on the results of setup.py output.
- compare_signature.py: Analyse signature matches between two signature datasets.
- probe.sh: Send LFP-style probes to a list of IPs.
- signatures/signatures_new: A folder containing the resulted set of signatures from various routers identified using SNMPv3.

Altough active probing managed to identify 54% of the routers that could not be identified by SNMPv3 in our dataset, it will not be integrated into the main pipeline for the moment as it is incompatible with the current CVE extractor technique. The CVE extractor filters CVEs with respect with the vendor and the last reboot time, the latest being unavailable using this method.

## How to use - Pipeline returning CVEs from domains tracerouting

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

Run the code with `make` and enter the sudo password when asked (traceroutes need sudo access)
As dependencies, the pipeline uses a combination of Golang, Rust and Python.

### Check outputs

The code outputs 2 files (`./output/engine_to_cves.csv` and `./output/ips_to_cves.csv`) and a folder (`./output/cves_per_engId`)

1. `./output/engine_to_cves.csv` contains for each EngineId a max severity (highest severity among all the correlated CVEs), CVEs with max severity (IDs of these cves), All CVEs (all correlated cves), IPs correlated (IPs that returned that specific EngineId)
2. `./output/ips_to_cves.csv` contains the same as above but thus time each tested IP with its CVEs
3. `./output/cves_per_engId` contains for each EngineID a folder with two files:
   - `cves_important_info.json` important information from each CVE extracted (CVE ID, severity using different metrics)
   - `cves_full_format.json` full CVEs extracted from EngineID




## Results from our testing
We ran our pipeline with 192 737 non unique IPs extracted from a daily traceroute dump provided by RIPE Atlas (https://data-store.ripe.net/datasets/atlas-daily-dumps/) on March 18, 2024. The reults are on `results_ripe_atlas_ips` branch 
