from scapy.all import *
import collections
from enum import Enum, auto

class SequenceNumberBehavior(Enum):
    INCREMENTAL = auto()
    RANDOM = auto()
    ZERO = auto()
    DUPLICATES = auto()
    UNDETERMINED = auto()

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
    def get_icmp_ittl_probe(target_ip, count=3):
        ttl_values = []
        sequence_numbers = []
        response_sizes = []
        for i in range(count):
            packet = IP(dst=target_ip)/ICMP(seq=i)
            
            response = sr1(packet, timeout=1, verbose=0)
            
            if response is not None:
                ttl = response[IP].ttl
                seq = response[ICMP].seq
                size = len(response)
                ttl_values.append(ttl)
                sequence_numbers.append(seq)
                response_sizes.append(size)
                print(f"Response from {target_ip}: TTL = {ttl}, Seq = {seq}, Size: {size}")
            else:
                print(f"No response for packet {i}")

        return ttl_values, sequence_numbers, response_sizes

    @staticmethod
    def analyze_sequence_numbers(sequence_numbers):
        if not sequence_numbers:
            print("No sequence numbers to analyze.")
            return SequenceNumberBehavior.UNDETERMINED
        
        counter = collections.Counter(sequence_numbers)
        if len(counter) == 1 and sequence_numbers[0] == 0:
            print("All sequence numbers are zero.")
            return SequenceNumberBehavior.ZERO
        elif len(counter) == len(sequence_numbers):
            if all(x < y for x, y in zip(sequence_numbers, sequence_numbers[1:])):
                print("Sequence numbers are incremental.")
                return SequenceNumberBehavior.INCREMENTAL
            else:
                print("Sequence numbers are random.")
                return SequenceNumberBehavior.RANDOM
        else:
            print("Sequence numbers contain duplicates.")
        
            return SequenceNumberBehavior.DUPLICATES
        return SequenceNumberBehavior.UNDETERMINED

ttl_values, sequence_numbers, response_sizes = ITTLProbe.get_icmp_ittl_probe('193.239.116.205')
ITTLProbe.analyze_sequence_numbers(sequence_numbers)


# ITTLProbe.get_icmp_ittl_probe()

# Example usage

# ITTLProbe.get_tcp_ittl_probe('193.239.116.205')
