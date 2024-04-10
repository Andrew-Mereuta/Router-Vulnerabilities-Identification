IN_DIR=input
OUT_DIR=output

all: build main engine_parse engine_data get_cves engine_to_cves ips_to_cves

build:
	cd engine_parse && cargo b -r

main:
	go build -o snmp_results
	sudo ./snmp_results 

engine_parse: build main
	engine_parse/target/release/engine_parse $(OUT_DIR)/snmp_results.csv $(OUT_DIR)/engine_ids.json $(IN_DIR)/enterprise-numbers $(IN_DIR)/mac-vendors-export.json
	cat $(OUT_DIR)/engine_ids.json

engine_data: main
	python3 cve_extractor/extract_engine_data.py --snmp_results snmp_results.csv

get_cves: engine_parse
	python3 cve_extractor/get_cves.py --engine_ids_file engine_ids.json --engine_data engine_data.json

engine_to_cves: get_cves
	python3 cve_extractor/engine_to_cves.py

ips_to_cves: engine_to_cves
	python3 cve_extractor/ips_to_cves.py

clean:
	rm -rf engine_parse/target
