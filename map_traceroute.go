package main

import (
	"fmt"

	"hackinglab/ip_extractor"
	"hackinglab/snmp_sender"
)

func main() {
	fmt.Println("Extracting IPs...")
	// Call the ExtractIPs function from the ip_extractor package to extract IPs
	ip_extractor.ExtractIPs()

	fmt.Println("Sending SNMP requests...")
	// Call the SnmpSend function from the snmp_sender package to send SNMP requests
	snmp_sender.SnmpSend()
}
