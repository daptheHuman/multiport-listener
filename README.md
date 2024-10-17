## **Multiport Traffic Listener (Golang Project)**

### **Overview**

This project is a **Multiport Traffic Listener** written in Golang. It captures and logs **network packets** on a range of ports using **pcap** and **gopacket** libraries. Each port is listened to concurrently using **goroutines**, and relevant packet information, including HTTP request methods, is logged.

This tool can be helpful for:

- **Network monitoring**.
- **Debugging services** running on multiple ports.
- **Learning about network protocols** and **application payloads**.

---

### **Project Structure**

```
.
├── go.sum             # Go module checksums
├── listener           # Contains listener logic
│   └── listener.go    # Main code for packet listening and parsing
├── main.go            # Entry point of the application
└── server       # Sample Python TCP server (for testing)
```

---

### **Features**

- **Concurrent listening** on multiple ports using goroutines.
- Logs **source and destination IPs** and **HTTP methods** (e.g., `GET`, `POST`).
- Automatically defaults to port 80 if no input is provided.
- Uses **pcap** filters to capture only relevant packets for specified ports.

---

### **Dependencies**

1. **Gopacket**: For packet parsing and inspection.

---

### **Usage**

#### 1. Run a Test TCP Server on Port 8000

You can use the provided `/server` to simulate traffic.

#### 2. Run the Multiport Listener in Golang

Compile and run the listener:

```bash
sudo go run . "80,8000,8080"

```

The application will start listening on ports `80`, `8000`, and `8080` concurrently.
If you don't provide any input, the listener defaults to **port 80**:

```bash
sudo go run main.go ""
```

### **Potential Improvements**

1. **TLS/SSL Support**: Capture encrypted traffic.
2. **Docker Compatibility**: Extend support for containerized environments.
3. **Bubble Tea Dashboard**: Add a terminal dashboard for real-time packet logging.
