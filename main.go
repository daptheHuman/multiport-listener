package main

import (
	"os"

	listener "github.com/daptheHuman/multiport-listener/listener"
)

func main() {
	portInput := os.Args[1]
	listener.ListenPortRange(portInput)

	// Prevent the program from exiting
	select {}
}
