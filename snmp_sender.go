package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"regexp"
	"encoding/csv"

	"github.com/gosnmp/gosnmp"
)

func main() {

	// Open the file containing the IP addresses
	file, err := os.Open("input/router_ips3.txt")
	if err != nil {
		// Print the error and return if there was an issue opening the file
		fmt.Println(err)
		return
	}
	defer file.Close()

	var ips []string

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)
	// Loop over all lines in the file
	for scanner.Scan() {
		// Append each line (IP address) to the slice
		ips = append(ips, scanner.Text())
	}

	// Check if there were errors during the scanning
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	// Create a new CSV file to store the results
	csvFile, err := os.Create("output/snmp_results.csv")
	if err != nil {
		// Print the error and return if there was an issue creating the file
		fmt.Println("Error creating the file: ", err)
		return
	}
	defer csvFile.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write the CSV header
	writer.Write([]string{"IP", "AuthoritativeEngineBoots", "AuthoritativeEngineTimes", "EngineID", "ScanTime"})
	writer.Flush()

	// Loop over all IP addresses
	for _, ip := range ips {

		fmt.Println(ip)

		// Create a new SNMP object with missing authentication parameters
		snmp := &gosnmp.GoSNMP{
			Target: ip,
			Port: 161,
			Version: gosnmp.Version3,
			SecurityModel: gosnmp.UserSecurityModel,
			Timeout:	time.Duration(2) * time.Second,
			Retries: 0,
			SecurityParameters: &gosnmp.UsmSecurityParameters{
				UserName: " ",
			},
		}

		// Connect to the SNMP server
		err = snmp.Connect()
		if err != nil {
			// Print the error and return if there was an issue connecting to the server
			fmt.Println("Error connecting to the SNMP server: ", err)
			return
		}
		defer snmp.Conn.Close()

		// OIDs to get
		oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
		// Get the OIDs from the SNMP server
		result, err2 := snmp.Get(oids)
		if err2 != nil {
			// Check if there was a result and if it has security parameters
			if(result != nil && result.SecurityParameters != nil){
				// Extract the authoritative engine boots, times, and ID from the security parameters
				reBoots := regexp.MustCompile(`AuthoritativeEngineBoots:(\d+)`)
				matchBoots := reBoots.FindStringSubmatch(result.SecurityParameters.SafeString())
				reTime := regexp.MustCompile(`AuthoritativeEngineTimes:(\d+)`)
				matchTime := reTime.FindStringSubmatch(result.SecurityParameters.SafeString())
				reEngineID := regexp.MustCompile(`engine=\((\w+)\)`)
				matchEngineID := reEngineID.FindStringSubmatch(result.SecurityParameters.Description())
				// Check if all values were found
				if len(matchBoots) > 0 && len(matchTime) > 0 && len(matchEngineID) > 0 {
					// Get the current time
					scanTime := time.Now().Format(time.RFC3339)
					// Write the values to the CSV file
					writer.Write([]string{ip, matchBoots[1], matchTime[1], matchEngineID[1], scanTime})
					writer.Flush()
				}
				// Print a success message
				fmt.Println("Success getting the SNMP values")
			} else {
				// Print the error if there was an issue getting the SNMP values
				fmt.Println("Error getting the SNMP values: ", err2)
			}
		}
	}
}