IN_DIR=input
OUT_DIR=output

all: build engine_parse

build:
	cd engine_parse && cargo b -r

engine_parse: build
	engine_parse/target/release/engine_parse $(OUT_DIR)/snmp_results.csv $(OUT_DIR)/engine_ids.json $(IN_DIR)/enterprise-numbers $(IN_DIR)/mac-vendors-export.json
	cat $(OUT_DIR)/engine_ids.json

clean:
	rm -rf engine_parse/target
