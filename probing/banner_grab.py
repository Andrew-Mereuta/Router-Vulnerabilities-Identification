import socket

def banner_grabbing(ip_address, port):
    try:
        client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        
        client_socket.settimeout(2)
        
        client_socket.connect((ip_address, port))
        
        banner = client_socket.recv(1024)
        
        print(f"[+] Banner for {ip_address}:{port}")
        print(banner.decode().strip())
    
    except Exception as e:
        print(f"[-] Error: {e}")
    
    finally:
        client_socket.close()

def banner_grabbing_for_ips(file_path, ports):
    try:
        with open(file_path, 'r') as file:
            for line in file:
                ip_address = line.strip()  
                for port in ports:
                    print(f"SCAN HOST: {ip_address} on port: {port}")    
                    banner_grabbing(ip_address, port)
    
    except FileNotFoundError:
        print("File not found.")
    
    except Exception as e:
        print(f"[ERROR]: {e}")

"""
execute banner grabbing at HTTP port 80 on ips from a file.
"""
def main():
    file_path = '../input/router_ips3.txt'
    
    # HTTP, HTTPS, FTP
    ports = [80, 8080, 21, 22]
    
    banner_grabbing_for_ips(file_path, ports)

if __name__ == "__main__":
    main()
