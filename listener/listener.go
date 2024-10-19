package listener

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gopacket "github.com/google/gopacket"
	pcap "github.com/google/gopacket/pcap"
)

func ParseInput(input string) ([]int, error) {
	var ports []int
	if input == "" {
		ports = append(ports, 80) // Default port
	} else {
		for _, port := range strings.Split(input, ",") {
			if p, err := strconv.ParseInt(port, 10, 64); err == nil {
				ports = append(ports, int(p))
			} else {
				return nil, fmt.Errorf("invalid port: %s", port)
			}
		}
	}
	return ports, nil
}

func handlePacket(packet gopacket.Packet, log chan string) {
	now := time.Now().Format("2006-01-02 15:04:05")

	// Extract network and transport layer details
	networkLayer := packet.NetworkLayer()     // ipv4 or ipv6
	transportLayer := packet.TransportLayer() // tcp or udp
	// linkLayer := packet.LinkLayer()           // ethernet

	if networkLayer != nil && transportLayer != nil {
		src, dst := networkLayer.NetworkFlow().Endpoints()

		// Optionally, print the application payload
		if appLayer := packet.ApplicationLayer(); appLayer != nil && len(appLayer.Payload()) > 0 {
			payload := string(appLayer.Payload())
			log <- fmt.Sprintf("%s: %s -> %s", now, src, dst)

			httpMethod := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE"}
			if slices.Contains(httpMethod, strings.Split(payload, " ")[0]) {
				method := strings.Split(payload, " ")[0]
				log <- fmt.Sprintf("METHOD: [%s] -> Payload: [%s]\n\n", method, strings.TrimSpace(payload))
			} else {
				log <- fmt.Sprintf("Payload: [%q]\n\n", payload)
			}

		}
	}
}

func ListenPortRange(selectedDevice string, portInput string, log chan string) tea.Cmd {
	var ports []int

	ports, err := ParseInput(portInput)
	if err != nil {
		log <- fmt.Sprintf("Error parsing input: %v", err)
	}
	log <- fmt.Sprintf("Listening on ports: %v", ports)
	for _, port := range ports {
		go func(p int) {
			log <- fmt.Sprintf("Listening on port %d\n", port)
			ListenPort(selectedDevice, p, log) // Start a goroutine for each port
		}(port)
	}

	return nil

}

func AllDevices() ([]string, error) {
	devices := []string{}
	if allDevices, err := pcap.FindAllDevs(); err != nil {
		return nil, err
	} else {
		for _, device := range allDevices {
			devices = append(devices, device.Name)
		}
	}
	return devices, nil
}

func ListenPort(selectedDevice string, port int, log chan string) {
	filter := fmt.Sprintf("port %d", port)
	if handle, err := pcap.OpenLive(selectedDevice, 1600, true, pcap.BlockForever); err != nil {
		log <- err.Error()
	} else if err := handle.SetBPFFilter(filter); err != nil {
		log <- err.Error()
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			handlePacket(packet, log)
		}
	}
}
