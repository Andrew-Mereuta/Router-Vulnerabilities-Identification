package ip_extractor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	icmpProtocol = 1
)

type PathIP struct {
	Hop int    `json:"hop"`
	IP  string `json:"ip"`
}

type DomainData struct {
	Domain  string   `json:"domain"`
	PathIPs []PathIP `json:"pathIPs"`
}

func ExtractIPs() {
	// getDomains reads the domain.txt file and returns all unique website domains
	hosts := getDomains()
	var domainDataArray []DomainData

	for _, host := range hosts {
		// traceroute(host) performs a traceroute operation in golang and returns set of ips, that were travelled through
		fmt.Println("Tracerouting to", host, "...")
		ips, err := traceroute(host)
		if err != nil {
			fmt.Println("Traceroute error:", err)
			return
		}

		var pathIPs []PathIP
		hop := 1
		for ip := range ips {
			pathIPs = append(pathIPs, PathIP{hop, ip})
			hop = hop + 1
		}

		domainDataArray = append(domainDataArray, DomainData{host, pathIPs})
	}

	jsonData, err := json.MarshalIndent(domainDataArray, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	file, err := os.Create("./input/ips.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("JSON data saved to ips.json")
}

func getDomains() []string {
	file, err := os.Open("./input/domains.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return make([]string, 0)
	}
	defer file.Close()

	domainsSet := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		domains := splitIntoWords(line)
		for _, word := range domains {
			domainsSet[word] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return make([]string, 0)
	}

	domainList := make([]string, 0, len(domainsSet))
	for domain := range domainsSet {
		domainList = append(domainList, domain)
	}

	return domainList
}

func splitIntoWords(line string) []string {
	return splitBySpace(line)
}

func splitBySpace(line string) []string {
	return strings.Fields(line)
}

func traceroute(host string) (map[string]bool, error) {
	ips := make(map[string]bool)

	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return nil, fmt.Errorf("error resolving host: %v", err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("error listening for ICMP packets: %v", err)
	}
	defer conn.Close()

	for ttl := 1; ttl <= 64; ttl++ {
		deadline := time.Now().Add(3 * time.Second)
		conn.SetDeadline(deadline)

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   12345,
				Seq:  1,
				Data: []byte(""),
			},
		}
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			return nil, fmt.Errorf("error marshaling ICMP message: %v", err)
		}

		if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
			return nil, fmt.Errorf("error setting TTL: %v", err)
		}

		if _, err := conn.WriteTo(msgBytes, ipAddr); err != nil {
			return nil, fmt.Errorf("error sending ICMP message: %v", err)
		}

		recvBuf := make([]byte, 1500)
		n, addr, err := conn.ReadFrom(recvBuf)
		if err != nil {
			continue
		}

		ips[addr.String()] = true

		msg1, err1 := icmp.ParseMessage(icmpProtocol, recvBuf[:n])
		if err1 != nil {
			return nil, fmt.Errorf("error parsing ICMP message: %v", err1)
		}

		switch msg1.Type {
		case ipv4.ICMPTypeTimeExceeded:
		case ipv4.ICMPTypeEchoReply:
			return ips, nil
		default:
		}
	}
	return ips, nil
}
