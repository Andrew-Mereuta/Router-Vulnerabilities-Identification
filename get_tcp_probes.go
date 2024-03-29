package main

import (
    "fmt"
    "log"
    "net"
    "time"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
)

func main() {

    targetIP := "193.239.116.205"
    localIP := "192.168.1.162"  // Set your local IP here
    iface := "wlo1"   // Set your network interface here
    targetPort := 33533              // Target port for the TCP probes
    fmt.Printf("start")

    // Open up a pcap handle for packet reads/writes.
    handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
    if err != nil {
        log.Fatal(err)
    }
    defer handle.Close()

    // Send TCP probes
    if err := sendTCPProbes(handle, localIP, targetIP, targetPort); err != nil {
        log.Fatalf("Error sending TCP probes: %v", err)
    }

    // Analyze the TCP responses to determine the nature of the IPID counter
    analysisDuration := 10 * time.Second // Duration to capture and analyze packets
    ipidClassification := analyzeTCPIPIDs(handle, targetIP, analysisDuration)
    fmt.Printf("TCP IPID sequence classification: %s\n", ipidClassification)
}

// sendTCPProbes sends TCP probes to the target
func sendTCPProbes(handle *pcap.Handle, localIP, targetIP string, targetPort int) error {
    packets := []*layers.TCP{
        {SrcPort: 12345, DstPort: layers.TCPPort(targetPort), Seq: 100, ACK: false, SYN: true},
        {SrcPort: 12345, DstPort: layers.TCPPort(targetPort), Seq: 200, ACK: true},
        {SrcPort: 12345, DstPort: layers.TCPPort(targetPort), Seq: 300, ACK: true},
    }

    for _, tcp := range packets {
        if err := sendTCPPacket(handle, localIP, targetIP, tcp); err != nil {
            return err
        }
    }

    return nil
}

// sendTCPPacket crafts and sends a single TCP packet
func sendTCPPacket(handle *pcap.Handle, srcIP, dstIP string, tcp *layers.TCP) error {
    ip := &layers.IPv4{
        Version:  4,
        TTL:      64,
        SrcIP:    net.ParseIP(srcIP),
        DstIP:    net.ParseIP(dstIP),
        Protocol: layers.IPProtocolTCP,
    }
    tcp.SetNetworkLayerForChecksum(ip)
    buf := gopacket.NewSerializeBuffer()
    opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
    gopacket.SerializeLayers(buf, opts, ip, tcp)
    return handle.WritePacketData(buf.Bytes())
}

// analyzeTCPIPIDs captures and analyzes TCP response IPID values
func analyzeTCPIPIDs(handle *pcap.Handle, targetIP string, duration time.Duration) string {
    start := time.Now()
    var ipids []uint16

    // Set BPF filter
    filter := fmt.Sprintf("ip src %s and tcp", targetIP)
    if err := handle.SetBPFFilter(filter); err != nil {
        log.Fatalf("Could not set BPF filter: %v", err)
    }

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        if time.Since(start) > duration {
            break
        }
        if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
            ip, _ := ipLayer.(*layers.IPv4)
            ipids = append(ipids, ip.Id)
        }
    }

    return classifyIPIDSequence(ipids)
}

// classifyIPIDSequence analyzes the collected IPID values and classifies them.
func classifyIPIDSequence(ipids []uint16) string {
	if len(ipids) == 0 {
		return "no data"
	}

	fmt.Printf("GOT IPDS")

	// Initialize variables to track sequence characteristics
	isIncremental := true
	isRandom := false
	isStatic := true
	isZero := true
	duplicates := make(map[uint16]int)

	previousIPID := ipids[0]
	for i, ipid := range ipids {
		if i > 0 {
			if ipid != previousIPID {
				isStatic = false
			}
			if ipid != previousIPID+1 && isIncremental {
				isIncremental = false
			}
			if ipid != 0 {
				isZero = false
			}
		}

		duplicates[ipid]++

		previousIPID = ipid
	}

	for _, count := range duplicates {
		if count > 1 {
			isRandom = true
			break
		}
	}

	// Determine classification based on observed characteristics
	switch {
	case isZero:
		return "zero"
	case isStatic:
		return "static"
	case isIncremental:
		return "incremental"
	case isRandom:
		return "random or duplicate"
	default:
		return "unclassified"
	}
}