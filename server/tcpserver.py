import socket

def start_server():
    # Create a TCP socket
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Bind the socket to the port
    server_socket.bind(('localhost', 8000))

    # Enable the server to accept connections
    server_socket.listen(5)
    print("Server is listening on port 8000...")

    while True:
        # Wait for a client to connect
        client_socket, addr = server_socket.accept()
        print(f"Connection from {addr} has been established!")

        # Receive data from the client
        data = client_socket.recv(1024).decode('utf-8')
        print(f"Received: {data}")

        # Optionally, send a response back to the client
        response = "Hello from the server!"
        client_socket.sendall(response.encode('utf-8'))

        # Close the client socket
        client_socket.close()

if __name__ == "__main__":
    start_server()
