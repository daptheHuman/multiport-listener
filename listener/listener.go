package listener

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	gopacket "github.com/google/gopacket"
	pcap "github.com/google/gopacket/pcap"
)

func parseInput(input string) ([]int, error) {
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

func handlePacket(packet gopacket.Packet) {
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
			log.Printf("%s: %s -> %s\n", now, src, dst)

			httpMethod := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE"}
			if slices.Contains(httpMethod, strings.Split(payload, " ")[0]) {
				method := strings.Split(payload, " ")[0]
				log.Printf("METHOD: [%s] -> Payload: [%s]\n\n", method, payload)
			} else {
				log.Printf("Payload: [%q]\n\n", payload)
			}

		}
	}
}

func ListenPortRange(portInput string) {
	var ports []int

	ports, err := parseInput(portInput)
	if err != nil {
		panic(err)
	}

	for _, port := range ports {
		go func(p int) {
			if err := ListenPort(p); err != nil { // Start a goroutine for each port
				log.Printf("Error listening on port %d: %v", p, err)
			}
		}(port)
	}
}

func ListenPort(port int) error {
	devices, err := pcap.FindAllDevs()
	if devices == nil || err != nil {
		return fmt.Errorf("no devices found")
	}
	filter := fmt.Sprintf("port %d", port)
	if handle, err := pcap.OpenLive("lo0", 1600, true, pcap.BlockForever); err != nil {
		return err
	} else if err := handle.SetBPFFilter(filter); err != nil {
		return err
	} else {
		fmt.Printf("Listening on port %d...\n", port)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			handlePacket(packet)
		}

		return nil
	}
}
