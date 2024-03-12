package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"
)

type Part string

const (
	UNDEFINED        Part = ""
	APPLICATION      Part = "a"
	OPERATING_SYSTEM Part = "o"
	HARDWARE         Part = "h"
)

type Severity int

const (
	LOW    Severity = iota
	MEDIUM Severity = iota
	HIGH   Severity = iota
)

type CVE struct {
	cveId         string   `json:"id"`
	severity      Severity `json:""`
	company       string
	publishedDate time.Time
	lastModified  time.Time
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func craftReqURL(part Part, company string, from time.Time, to time.Time) string {
	requestURL := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?virtualMatchString=cpe:2.3:%s:%s:*:*:*:*:*:*:*:*:*&pubStartDate=%s&pubEndDate=%s", part, company, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
	return requestURL
}

func get_cves(part Part, company string, from time.Time, sev int) {

	windowS := from
	windowE := from.Add(time.Hour * 24 * 100)

	// var jsonMap map[string]interface{}

	for windowS.Before(time.Now()) {
		requestURL := craftReqURL(part, company, windowS, windowE)
		fmt.Print("\n\n\n")
		fmt.Println(requestURL)

		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("client: got response!\n")
		fmt.Printf("client: status code: %d\n", res.StatusCode)

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("client: response body: %s\n", resBody)
		fmt.Printf("t1: %s\n", reflect.TypeOf(resBody))

		windowS = windowE
		windowE = windowE.Add(time.Hour * 24 * 100)
	}

	// requestURL := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=Cisco&resultsPerPage=300&startIndex=20&cvssV3Severity=CRITICAL")

	// f, err := os.Create("./dat2.txt")
	// check(err)

	// defer f.Close()

	// // var jsonMap []string
	// json.Unmarshal([]byte(resBody), &jsonMap)
	// asdf, _ := json.Marshal(jsonMap["vulnerabilities"])
	// // aux := jsonMap["vulnerabilities"]

	// n3, err := f.WriteString(string(asdf))
	// check(err)
	// fmt.Printf("wrote %d bytes\n", n3)
	// fmt.Println(requestURL)
	// // fmt.Printf("%d", gg)
	// // fmt.Println(len(jsonMap))
}
