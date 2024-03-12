package ip_extractor

import (
	"bufio"
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

func ExtractIPs() {

	// Create file to store all unique ip addresses
	file, err := os.Create("./input/ips.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	// ipsSet := make(map[string]bool)
	// getDomains reads the domain.txt file and returns all unique website domains
	hosts := getDomains()

	for _, host := range hosts {
		// traceroute(host) performs a traceroute operation in golang and returns set of ips, that were travelled through
		ips, err := traceroute(host)
		if err != nil {
			fmt.Println("Traceroute error:", err)
			return
		}

		// Write the domain to the file
        _, fErr := writer.WriteString(host)
        if fErr != nil {
            fmt.Println("Error writing to file:", fErr)
            return
        }

		for ip := range ips {
            _, fErr := writer.WriteString(", " + ip)
            if fErr != nil {
                fmt.Println("Error writing to file:", fErr)
                return
            }
        }

		_, fErr = writer.WriteString("\n")
        if fErr != nil {
            fmt.Println("Error writing to file:", fErr)
            return
        }

		if fErr := writer.Flush(); fErr != nil {
			fmt.Println("Error flushing writer:", fErr)
			return
		}

		// for ip := range ips {
		// 	ipsSet[ip] = true
		// }
	}

	// Save all ip addresses to the file.
	// for ip := range ipsSet {
	// 	_, fErr := writer.WriteString(ip + "\n")
	// 	if fErr != nil {
	// 		fmt.Println("Error writing to file:", fErr)
	// 		return
	// 	}
	// }

	// if fErr := writer.Flush(); fErr != nil {
	// 	fmt.Println("Error flushing writer:", fErr)
	// 	return
	// }
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
