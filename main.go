package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	cli "github.com/daptheHuman/multiport-listener/cli"
)

func main() {
	p := tea.NewProgram(cli.InitalModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	select {}
}
