"""
    Script used to parse entries from the ip_extractor and engine_parse.
    Please make sure to run ip_extractor and engine_parse first.
    
    This script outputs :
    - an extended version of "engine_ids" output of the engine_parse with the matching vendor in case the engine_id matches a IANA enterprise number.
    - the list of unindentified routers.
    - the list of identified routers.


"""

import json
from enum import Enum

engine_id_file = open("../output/engine_ids.json")
engine_id_data = json.load(engine_id_file)

output_headers = ["Domain","IP","AuthoritativeEngineBoots","AuthoritativeEngineTimes","EngineID","ScanTime", "Vendor"]

unique_ips = []
unique_unknown_ips = []
class DataColumn(Enum):
    DOMAIN = 0
    IP = 1
    AUTHORATIVE_ENGINE_BOOTS = 2
    AUTHORATIVE_ENGINE_TIMES = 3
    ENGINE_ID = 4
    SCAN_TIME = 5
    VENDOR = 6

def get_vendor(engine_id):
    if engine_id_data.get(engine_id):
        return engine_id_data[engine_id]["vendor"]["name"]
    return None

def parse_data_line(columns):
    return ",".join(columns) + '\n'

output_file = open('output.csv', 'w')
output_file.write(parse_data_line(output_headers))                                          # write header

print("Parsing and filtering SNMPV3 entries...")

with open("../output/testing_snmp_results.csv") as file:
    lines = [line.rstrip() for (line_index, line) in enumerate(file) if line_index != 0]
    for line in lines:
        if "Error" in line:
            columns = line.split(",")
            if columns[DataColumn.IP._value_] in unique_unknown_ips:
                continue
            unique_unknown_ips.append(columns[DataColumn.IP._value_])
        else:
            columns = line.split(",")
            if columns[DataColumn.IP._value_] in unique_ips:
                continue
            unique_ips.append(columns[DataColumn.IP._value_])
            vendor = get_vendor(engine_id=columns[DataColumn.ENGINE_ID._value_])                     
            if vendor is None:
                continue
            columns.append(vendor)                                                          # append vendor data

            output_file.write(parse_data_line(columns=columns))

with open("unique_ips.txt", "w") as file:
    for unique_ip in unique_ips:
        file.writelines(unique_ip + "\n")

with open("unknown_unique_ips.txt", "w") as file:
    for unique_ip in unique_unknown_ips:
        file.writelines(unique_ip + "\n")

print("Done")



