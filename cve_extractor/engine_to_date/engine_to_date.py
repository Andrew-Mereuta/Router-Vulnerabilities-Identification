from argparse import ArgumentParser
import csv
from os.path import join
from datetime import datetime
import json

output_folder = "./output"
input_folder = "./input"
    

def main(snmp_results):
    engToResetDate = {}
    with open(join(output_folder, "engine_to_reset_date.csv"), "w") as w:
        writer = csv.writer(w)
        writer.writerow(["EngineId", "date"])
        with open(join(output_folder, snmp_results), "r") as f:
            reader = csv.reader(f)
            next(reader)
            for i, row in enumerate(reader):
                reset_date = row[3]
                # final_date = datetime(row[5])
                engineId = row[4]

                if engineId == "Error":
                    continue

                writer.writerow([engineId, reset_date])
                engToResetDate[engineId] = reset_date

    with open(join(output_folder, "engine_to_reset_date.json"), "w") as f:
        json.dump(engToResetDate, f)
            



if __name__ == "__main__":
    parser = ArgumentParser()
    parser.add_argument("--snmp_results", type=str, required=True)

    args = parser.parse_args()

    snmp_results = args.snmp_results

    main(snmp_results)