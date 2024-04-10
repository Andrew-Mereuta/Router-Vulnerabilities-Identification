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
            ip = row[1]
            engineId = row[4]

            if engineId == "Error":
                continue
            
            if engineId not in engToIps:
                engToIps[engineId] = set()
            # writer.writerow([engineId, reset_date])
            engToIps[engineId].add(ip)

    for key in engToIps:
        engToIps[key] = list(engToIps[key])

    with open(join(output_folder, "engine_to_ips.json"), "w") as f:
        json.dump(engToIps, f)
            



if __name__ == "__main__":
    parser = ArgumentParser()
    parser.add_argument("--snmp_results", type=str, default="testing_snmp_results.csv")

    args = parser.parse_args()

    snmp_results = args.snmp_results

    main(snmp_results)