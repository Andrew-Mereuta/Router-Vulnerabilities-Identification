from scapy.all import *

class ITTLProbe:

    @staticmethod
    def get_tcp_ittl_probe(ip, port=33533):
        syn_packet = IP(dst=ip)/TCP(dport=port, flags='S')
        syn_ack_packet = sr1(syn_packet, timeout=1, verbose=0)

        if syn_ack_packet:
            ttl = syn_ack_packet[IP].ttl
            window_size = syn_ack_packet[TCP].window

            print(f"Target: {ip}, Port: {port}, TTL: {ttl}, Window Size: {window_size}")
        else:
            print(f"No response from {ip}, port {port}.")


    @staticmethod
    def get_icmp_ittl_probe(ip, port=33533):
        packet = IP(dst=ip)/ICMP()
        
        response = sr1(packet, timeout=1, verbose=0)
        
        if response is None:
            print(f"No response from {ip}")
        else:
            # Retrieve and print the TTL value from the response
            ttl = response[IP].ttl
            print(f"Response from {ip} with TTL: {ttl}")

ITTLProbe.get_icmp_ittl_probe('193.239.116.205')

# Example usage

# ITTLProbe.get_tcp_ittl_probe('193.239.116.205')
