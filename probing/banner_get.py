import socket

ip_address = '193.239.116.205' 
port = 50 

try:
    connection = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    connection.connect((ip_address, port))

    connection.send(b"GET / HTTP/1.1\r\nHost: " + ip_address.encode() + b"\r\n\r\n")
    response = connection.recv(4096)
    print(response.decode())
except Exception as e:
    print(f"Error: {e}")
finally:
    connection.close()
