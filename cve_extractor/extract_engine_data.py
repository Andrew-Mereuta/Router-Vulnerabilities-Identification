from argparse import ArgumentParser
import csv
from os.path import join
from datetime import datetime
import json

output_folder = "./output"
input_folder = "./input"
    

def main(snmp_results):
    engToIps = {}
    with open(join(output_folder, snmp_results), "r") as f:
        reader = csv.reader(f)
        next(reader)
        for i, row in enumerate(reader):
            engineId = row[4]
            ip = row[1]
            reset_date = row[3]

            if engineId == "Error":
                continue
            
            if engineId not in engToIps:
                engToIps[engineId] = {"ips": set(), "reset_date": None}
            # writer.writerow([engineId, reset_date])
            engToIps[engineId]["ips"].add(ip)
            engToIps[engineId]["reset_date"] = reset_date

    for key in engToIps:
        engToIps[key]["ips"] = list(engToIps[key]["ips"])

    with open(join(output_folder, "engine_data.json"), "w") as f:
        json.dump(engToIps, f)
            



if __name__ == "__main__":
    parser = ArgumentParser()
    parser.add_argument("--snmp_results", type=str)

    args = parser.parse_args()

    snmp_results = args.snmp_results

    main(snmp_results)