import requests
from datetime import datetime, timedelta
import json
from argparse import ArgumentParser
import csv
from os.path import join
import pytz
import os

utc=pytz.UTC


output_folder = "./output"
input_folder = "./input"
api_key = "e9df04f8-5b8e-4c88-8089-0649f8fcef63"

def craftReqURL(part, company, timeFrom, timeTo):
    fromstr = timeFrom.strftime('%Y-%m-%dT%H:%M:%S.%f')[:-3]
    toStr = timeTo.strftime('%Y-%m-%dT%H:%M:%S.%f')[:-3]
    requestURL = f"https://services.nvd.nist.gov/rest/json/cves/2.0?virtualMatchString=cpe:2.3:{part}:{company}:*:*:*:*:*:*:*:*:*&pubStartDate={fromstr}&pubEndDate={toStr}"
    return requestURL


def transform(x):
    id = x["cve"]["id"]

    baseSeverityV31 = None
    baseScoreV31 = None
    if "cvssMetricV31" in x["cve"]["metrics"]:
        cvssMetricV31_list = list(filter(lambda y: y["type"] == "Primary", x["cve"]["metrics"]["cvssMetricV31"]))
        if len(cvssMetricV31_list) == 0:
            cvssMetricV31_list = list(filter(lambda y: y["type"] == "Secondary", x["cve"]["metrics"]["cvssMetricV31"]))
        
        if len(cvssMetricV31_list) > 0:
            cvssMetricV31 = cvssMetricV31_list[0]
            baseSeverityV31, baseScoreV31 = [cvssMetricV31["cvssData"]["baseSeverity"], cvssMetricV31["cvssData"]["baseScore"]]

    baseSeverityV30 = None
    baseScoreV30 = None
    if "cvssMetricV30" in x["cve"]["metrics"]:
        cvssMetricV30_list = list(filter(lambda y: y["type"] == "Primary", x["cve"]["metrics"]["cvssMetricV30"]))
        if len(cvssMetricV30_list) == 0:
            cvssMetricV30_list = list(filter(lambda y: y["type"] == "Secondary", x["cve"]["metrics"]["cvssMetricV30"]))

        if len(cvssMetricV30_list) > 0:
            cvssMetricV30 = cvssMetricV30_list[0]
            baseSeverityV30, baseScoreV30 = [cvssMetricV30["cvssData"]["baseSeverity"], cvssMetricV30["cvssData"]["baseScore"]]

    baseSeverityV2 = None
    baseScoreV2 = None
    if "cvssMetricV2" in x["cve"]["metrics"]:
        cvssMetricV2_list = list(filter(lambda y: y["type"] == "Primary", x["cve"]["metrics"]["cvssMetricV2"]))
        if len(cvssMetricV2_list) == 0:
            cvssMetricV2_list = list(filter(lambda y: y["type"] == "Secondary", x["cve"]["metrics"]["cvssMetricV2"]))

        if len(cvssMetricV2_list) > 0:
            cvssMetricV2 = cvssMetricV2_list[0]
        baseSeverityV2, baseScoreV2 = [cvssMetricV2["baseSeverity"], cvssMetricV2["cvssData"]["baseScore"]]

    publishedDate = x["cve"]["published"]
    modifiedDate = x["cve"]["lastModified"]

    return {"id": id,
            "baseSeverityV31": baseSeverityV31, 
            "baseScoreV31": baseScoreV31, 
            "baseSeverityV30": baseSeverityV30, 
            "baseScoreV30": baseScoreV30, 
            "baseSeverityV2": baseSeverityV2, 
            "baseScoreV2": baseScoreV2, 
            "publishedDate": publishedDate, 
            "modifiedDate": modifiedDate}


def filter_routers(x):
    # print(x)
    return "router" in json.dumps(x)

def main(part, vendor, timeFrom, engId):
    windowS = timeFrom
    windowE = windowS + timedelta(days = 100)

    vul = []
    cves = []
    while windowE <= datetime.now():
        url = craftReqURL(part, vendor, windowS, windowE)
        print(url)
        
        headers = {"apiKey": api_key}

        response = requests.get(url, headers=headers)

        if response:
            routers = list(filter(filter_routers, response.json()["vulnerabilities"]))
            vul.extend(map(transform, routers))
            cves.extend(routers)
                
        
        windowS = windowE
        windowE = windowE + timedelta(days = 100)

    # print(vul)
    try:
        os.mkdir(f"./output/cves_per_engId/{engId}")
    except:
        pass

    with open(f"./output/cves_per_engId/{engId}/cves_full_format.json", "w") as f:
        f.write(json.dumps(cves))

    with open(f"./output/cves_per_engId/{engId}/cves_important_info.json", "w") as f:
        f.write(json.dumps(vul))

if __name__ == "__main__":
    parser = ArgumentParser()
    parser.add_argument("--engine_ids_file", type=str, required=True)
    parser.add_argument("--engine_to_reset", type=str, required=True)

    args = parser.parse_args()
    engine_ids_file = args.engine_ids_file
    engine_to_reset = args.engine_to_reset

    with open(join(output_folder, engine_ids_file), "r") as f:
        data = json.load(f)
    
    with open(join(output_folder, engine_to_reset), "r") as f:
        reset_data = json.load(f)

    with open(join(input_folder, "mapping.json"), "r") as f:
        mapping = json.load(f)

    try:
        os.rmdir("./output/cves_per_engId")
    except:
        print("error")
        pass

    os.mkdir("./output/cves_per_engId")


    for key in data:
        if "vendor" not in data[key]:
            print("No vendor for ", key)
            continue
        if "name" not in data[key]["vendor"]:
            print("No name for vendor for ", key)
            continue

        company = data[key]["vendor"]["name"]

        if company not in mapping:
            print("No mapping for ", company)
            continue

        vendor =  mapping[company]

        if key not in reset_data:
            print("No reset date for ", key) 
            continue

        reset_date = datetime.fromisoformat(reset_data[key]).replace(tzinfo=None)

        main("h", vendor, reset_date, key)