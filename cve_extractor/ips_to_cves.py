from argparse import ArgumentParser
import csv
from os.path import join
from datetime import datetime
import json

output_folder = "./output"
input_folder = "./input"
    

def main():
    with open(join(output_folder, "ips_to_cves.csv"), "w", newline='') as w:
        writer = csv.writer(w)
        writer.writerow(["IP", "EngineId",  "Max severity", "CVEs with max severity", "All CVEs"])

        with open(join(output_folder, "engine_to_cves.csv"), "r") as f:
            reader = csv.reader(f)
            next(reader)
            for i, row in enumerate(reader):
                ips = row[4].strip('][').replace("'", "").split(', ')

                for ip in ips:
                    writer.writerow([ip, row[0], row[1], row[2], row[3]])

if __name__ == "__main__":
    main()