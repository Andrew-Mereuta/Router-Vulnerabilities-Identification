from scapy.all import *
import sys

def send_udp_probes(target_ip, target_port=33533):
    """
    Sends three UDP packets with a 12-byte payload of all zeroes to the target IP and port.
    """
    payload = Raw(load=b'\x00' * 12)  # 12 bytes of all zero payload
    for _ in range(3):  # Send 3 packets
        packet = IP(dst=target_ip)/UDP(dport=target_port)/payload
        send(packet, verbose=False)

def capture_and_analyze(target_ip, timeout=10):
    """
    Captures ICMP port unreachable responses from the target IP for a given timeout and analyzes the IPID sequence.
    """
    # Filter for ICMP type 3 (destination unreachable) code 3 (port unreachable)
    packets = sniff(filter=f"icmp and src host {target_ip}", timeout=timeout)
    
    # Extract IPID values from the ICMP responses
    ipids = [packet[IP].id for packet in packets if packet.haslayer(ICMP) and packet[ICMP].type == 3 and packet[ICMP].code == 3]
    
    # Analyze IPID sequence
    return analyze_ipid_sequence(ipids)

def analyze_ipid_sequence(ipids):
    """
    Analyzes the sequence of IPID values to determine if they are incremental, random, static, zero, or contain duplicates.
    """
    if not ipids:
        return "No ICMP port unreachable responses captured."
    
    analysis_result = "Analysis Result: "
    if all(ipid == 0 for ipid in ipids):
        analysis_result += "Zero"
    elif all(ipid == ipids[0] for ipid in ipids):
        analysis_result += "Static"
    elif len(set(ipids)) == len(ipids):
        analysis_result += "Random or Incremental"
    else:
        analysis_result += "Duplicate"

    return analysis_result

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <target_ip>")
        sys.exit(1)
    
    target_ip = sys.argv[1]

    print("Sending UDP probes to", target_ip)
    send_udp_probes(target_ip)

    print("Capturing and analyzing ICMP responses...")
    result = capture_and_analyze(target_ip)
    print(result)
