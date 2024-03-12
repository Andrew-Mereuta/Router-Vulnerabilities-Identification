package snmp_sender

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"regexp"
	"encoding/csv"
	"strings"

	"github.com/gosnmp/gosnmp"
)

// readIPsFromFile reads the IPs from the file
func readIPsFromFile() ([][]string, error) {
	// Open the file
	file, err := os.Open("input/ips.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records [][]string
	scanner := bufio.NewScanner(file)
	// Read each line and split it by comma, first is the domain, and then it is the IPs
	for scanner.Scan() {
		records = append(records, strings.Split(scanner.Text(), ","))
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// createCSVFile creates a new CSV file and writes the header
func createCSVFile() (*csv.Writer, error) {
	// Create the file
	csvFile, err := os.Create("output/snmp_results.csv")
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(csvFile)
	// Write the header
	writer.Write([]string{"Domain","IP", "AuthoritativeEngineBoots", "AuthoritativeEngineTimes", "EngineID", "ScanTime"})
	writer.Flush()

	return writer, nil
}

// getSNMPValues gets the SNMP values for the given IP
func getSNMPValues(ip string) ([]string, error) {
	// Setup the SNMP connection with the given IP and missing authentication
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
	err := snmp.Connect()
	if err != nil {
		return nil, err
	}
	defer snmp.Conn.Close()

	// Define the OIDs to get
	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
	// Get the SNMP values
	result, err := snmp.Get(oids)
	if err != nil && err != gosnmp.ErrUnknownUsername {
		return nil, err
	}

	// Parse the SNMP values
	reBoots := regexp.MustCompile(`AuthoritativeEngineBoots:(\d+)`)
	matchBoots := reBoots.FindStringSubmatch(result.SecurityParameters.SafeString())
	reTime := regexp.MustCompile(`AuthoritativeEngineTimes:(\d+)`)
	matchTime := reTime.FindStringSubmatch(result.SecurityParameters.SafeString())
	reEngineID := regexp.MustCompile(`engine=\((\w+)\)`)
	matchEngineID := reEngineID.FindStringSubmatch(result.SecurityParameters.Description())

	// Check if all values were found
	if len(matchBoots) > 0 && len(matchTime) > 0 && len(matchEngineID) > 0 {
		scanTime := time.Now().Format(time.RFC3339)
		return []string{ip, matchBoots[1], matchTime[1], matchEngineID[1], scanTime}, nil
	}

	return nil, fmt.Errorf("Error getting the SNMP values")
}

// SnmpSend is the main function that reads the IPs, gets the SNMP values and writes them to the CSV file
func SnmpSend() {
	// Read the IPs from the file
	records, err := readIPsFromFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the CSV file
	writer, err := createCSVFile()
	if err != nil {
		fmt.Println("Error creating the file: ", err)
		return
	}
	defer writer.Flush()

	// Loop over each record
	for _, record := range records {
		domain := record[0]
		ips := record[1:]

		// Loop over each IP
		for _, ip := range ips {
			ip := strings.TrimSpace(ip)
			// Get the SNMP values
			values, err := getSNMPValues(ip)
			if err != nil && err != gosnmp.ErrUnknownUsername {
				// Write an error row if there was an error
				writer.Write([]string{domain, ip, "Error", "Error", "Error", "Error"})
				writer.Flush()
				fmt.Println("Error getting the SNMP values: ", err)
			} else {
				// Write the SNMP values
				writer.Write(append([]string{domain}, values...))
				writer.Flush()
				fmt.Println("Success getting the SNMP values")
			}
		}
	}
}