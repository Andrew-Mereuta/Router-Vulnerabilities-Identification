from scapy.all import *
import sys

def send_tcp_probes(target_ip, target_port=33533):
    """
    Sends two ACK and one SYN packets to the target IP and port.
    """
    # Prepare the packet templates
    print(target_ip)
    ip_layer = IP(dst=target_ip)
    ack_packet = ip_layer/TCP(dport=target_port, flags="A")
    syn_packet = ip_layer/TCP(dport=target_port, flags="S")
    
    send(ack_packet, verbose=False)
    send(ack_packet, verbose=False)  
    send(syn_packet, verbose=False)

def capture_and_analyze(target_ip, timeout=10):
    """
    Captures the responses from the target IP for a given timeout and analyzes the IPID sequence.
    """
    packets = sniff(filter=f"tcp and dport 33533", timeout=timeout)
    if packets:
        for packet in packets:
            packet.show()  # Display packet details for debugging
    else:
        print("No packets captured.")


    
    # Extract IPID values
    ipids = [packet[IP].id for packet in packets if packet.haslayer(IP)]
    
    # Analyze IPID sequence
    analysis_result = analyze_ipid_sequence(ipids)
    return analysis_result

def analyze_ipid_sequence(ipids):
    """
    Analyzes the sequence of IPID values to determine if they are incremental, random, static, zero, or contain duplicates.
    """
    if not ipids:
        return "No packets captured."
    
    is_incremental = all(x+1 == y for x, y in zip(ipids, ipids[1:]))
    is_static = all(x == ipids[0] for x in ipids)
    is_zero = all(x == 0 for x in ipids)
    has_duplicates = len(ipids) != len(set(ipids))
    
    if is_zero:
        return "Zero"
    elif is_static:
        return "Static"
    elif is_incremental:
        return "Incremental"
    elif has_duplicates:
        return "Duplicate"
    else:
        return "Random"

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <target_ip>")
        sys.exit(1)
    
    target_ip = sys.argv[1]
    target_port = 33533  # Example port, adjust as needed
    
    print("Sending TCP probes...")
    send_tcp_probes(target_ip, target_port)
    
    print("Capturing and analyzing responses...")
    result = capture_and_analyze(target_ip)
    print(f"Analysis Result: {result}")
