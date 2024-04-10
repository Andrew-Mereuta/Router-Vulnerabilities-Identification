
import csv
import os
from os.path import join
import json

output_folder = "./output"

checkOrder = ["baseSeverityV31", "baseSeverityV30", "baseSeverityV2"]

metricToNumber = {
    "CRITICAL": 5,
    "HIGH": 4,
    "MEDIUM": 3,
    "LOW": 2,
    "INFORMATIONAL": 1
}



def main():
    analysedEngIds = [name for name in os.listdir(join(output_folder, "cves_per_engId"))]

    with open(join(output_folder, "engine_to_ips.json"), "r") as f:
        engToIps = json.load(f)

    with open(join(output_folder, "engine_to_cves.csv"), 'w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(["engineId", "Max severity", "CVEs with max severity", "All CVEs", "IPs correlated"])

        for engineId in analysedEngIds:
            with open(join(output_folder, "cves_per_engId", engineId, "cves_important_info.json"), "r") as f:
                data = json.load(f)

            max = 0
            maxSev = "NO VULNERABILITY FOUND"
            maxSevCVEs = []
            allCVEs = []
            for cve in data:
                noMetriv = True
                allCVEs.append(cve["id"])
                for metric in checkOrder:
                    if(cve[metric] != None):
                        noMetriv = False
                        if metricToNumber[cve[metric]] > max:
                            max = metricToNumber[cve[metric]]
                            maxSev = cve[metric]
                            maxSevCVEs.clear()

                        if metricToNumber[cve[metric]] == max:
                            maxSevCVEs.append(cve["id"])

                        break
                
                if noMetriv:
                    print("No metric for CVE: " + cve["id"] + " in engineId: " + engineId)
            
            
            writer.writerow([engineId, maxSev, maxSevCVEs, allCVEs, engToIps[engineId]])

if __name__ == "__main__":
    main()