DATA_DIR=engine_parse/data

all: build engine_parse

build:
	cd engine_parse && cargo b -r

engine_parse: build
	engine_parse/target/release/engine_parse $(DATA_DIR)/snmp_results.csv $(DATA_DIR)/engine_ids.json $(DATA_DIR)/enterprise-numbers $(DATA_DIR)/mac-vendors-export.json
	cat $(DATA_DIR)/engine_ids.json

clean:
	rm -rf engine_parse/target
