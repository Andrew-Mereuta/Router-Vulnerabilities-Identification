from scapy.all import *

# This script reads a list of domains from a file, performs a traceroute to each domain, 
# and writes the IPs of the routers to a file.
# 
# To correctly work, it must be executed with root privileges.

# Function to read domains from a file
def read_domains():
    with open("input/domains.txt") as f:
        # Read each line, strip whitespace, and return as a list
        return [line.strip() for line in f.readlines()]

# Function to write router IPs to a file
def write_router_ips(router_ips):
    with open("input/router_ips3.txt", "w") as f:
        for ip in router_ips:
            f.write(ip + "\n")

# Function to perform a traceroute to a domain
def traceroute_domain(domain):
    # Perform the traceroute
    result, unans = traceroute([domain])
    # Get the route packets
    route_packets = list(result.get_trace().values())[0]
    router_ips = set()
    # Loop over the route packets
    for (ip, value_bool) in route_packets.values():
        router_ips.add(ip)
    return router_ips

# Function to find all router IPs
def find_router_ips():
    # Set to store all IPs
    all_ips = set()
    # Loop over all domains
    for domain in read_domains():
        # Update the set with the IPs from the traceroute
        all_ips.update(traceroute_domain(domain))
    return all_ips

if __name__ == "__main__":
    all_ips = find_router_ips()
    write_router_ips(all_ips)