# RESULTS FROM RIPE ATLAS IPS

This is a branch containing the results from running our code with 192 737 non unique IPs extracted from a daily traceroute dump provided by RIPE Atlas (https://data-store.ripe.net/datasets/atlas-daily-dumps/) on March 18, 2024.

* The tested IPs are in `testing_snmp_results.csv`
* The mapping from IP to CVEs is in `ips_to_cves.csv`
* The mapping from EngineID to CVEs is in `engine_to_cves.csv`
* Full results regarding severities and CVEs for each engine ID are in the folder `cves_per_engId` (each subfolder is an EngineID)