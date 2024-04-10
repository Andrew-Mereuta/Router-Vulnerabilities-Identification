IN_DIR=input
OUT_DIR=output

all: build main engine_parse engine_to_date engine_to_ips get_cves engine_to_cves ips_to_cves

build:
	cd engine_parse && cargo b -r

main:
	go build -o snmp_results
	sudo ./snmp_results 

engine_parse: build main
	engine_parse/target/release/engine_parse $(OUT_DIR)/snmp_results.csv $(OUT_DIR)/engine_ids.json $(IN_DIR)/enterprise-numbers $(IN_DIR)/mac-vendors-export.json
	cat $(OUT_DIR)/engine_ids.json

engine_to_date: main
	python3 cve_extractor/engine_to_date.py --snmp_results snmp_results.csv

engine_to_ips: main
	python3 cve_extractor/engine_to_ips.py --snmp_results snmp_results.csv

get_cves: engine_parse
	python3 cve_extractor/get_cves.py --engine_ids_file engine_ids.json --engine_to_reset engine_to_reset_date.json

engine_to_cves: get_cves
	python3 cve_extractor/engine_to_cves.py

ips_to_cves: engine_to_cves
	python3 cve_extractor/ips_to_cves.py

clean:
	rm -rf engine_parse/target
