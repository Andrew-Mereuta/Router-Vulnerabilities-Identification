import requests
from datetime import datetime, timedelta
import json
# The API endpoint



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

def main(part, vendor, timeFrom):
    windowS = timeFrom
    windowE = windowS + timedelta(days = 100)

    vul = []
    cves = []
    while windowE <= datetime.now():
        url = craftReqURL(part, vendor, windowS, windowE)
        print(url)

        response = requests.get(url)

        if response:
            routers = list(filter(filter_routers, response.json()["vulnerabilities"]))
            vul.extend(map(transform, routers))
            cves.extend(routers)
                
        
        windowS = windowE
        windowE = windowE + timedelta(days = 100)

    # print(vul)

    with open("./output/cves_full_format.json", "w") as f:
        f.write(json.dumps(cves))

    with open("./output/cves_important_info.json", "w") as f:
        f.write(json.dumps(vul))

if __name__ == "__main__":
    main("h", "cisco", datetime(2020, 3, 15, 14, 30, 45, 0))