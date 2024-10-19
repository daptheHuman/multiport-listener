package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	listener "github.com/daptheHuman/multiport-listener/listener"
)

// PacketMsg represents a message containing packet information.
// type PacketMsg string
type LogMsg string

// Model stores the application state.
type Model struct {
	stage     int
	textinput textinput.Model
	cursor    int

	availDevices   []string
	selectedDevice string
	ports          []int

	packets []string
	log     chan string
	err     error
}

func initTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "80,443,8080"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ti
}

func waitLogMessages(log chan string) tea.Cmd {
	return func() tea.Msg {
		return LogMsg(<-log)
	}
}

// NewModel initializes the Bubble Tea model.
func InitalModel() Model {
	ti := initTextInput()

	devices, err := listener.AllDevices()
	if err != nil {
		panic(err)
	}

	return Model{
		ports:        []int{},
		availDevices: devices,
		log:          make(chan string),
		textinput:    ti,
		err:          nil,
	}
}

// Init starts the packet listener.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
}

// Update processes incoming messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case LogMsg:
		m.packets = append(m.packets, string(msg)) // Add log message to packets
		return m, waitLogMessages(m.log)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyUp, tea.KeyCtrlK:
			if m.stage == 0 {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case tea.KeyDown, tea.KeyCtrlJ:
			if m.stage == 0 {
				if m.cursor < len(m.availDevices)-1 {
					m.cursor++
				}
			}

		case tea.KeyEnter:
			if m.stage == 0 {
				m.stage++
				m.selectedDevice = m.availDevices[m.cursor]
				return m, nil
			}

			ports, err := listener.ParseInput(m.textinput.Value())
			if err != nil {
				m.err = err
				return m, nil
			}
			m.ports = ports
			go func() {
				listener.ListenPortRange(m.selectedDevice, m.textinput.Value(), m.log)
			}()
			return m, waitLogMessages(m.log)

		}

		m.textinput, cmd = m.textinput.Update(msg)
	}

	return m, cmd
}

// View renders the UI.
func (m Model) View() string {
	var s string
	footer := "Press Ctrl+C to quit."

	if m.stage == 0 {
		s += "Please select the device to listen on:\n"

		for idx, device := range m.availDevices {
			if idx == m.cursor {
				s += fmt.Sprintf("[x] %s\n", device)
			} else {
				s += fmt.Sprintf("[ ] %s\n", device)
			}
		}
		s += "\n"
		s += "Use the arrow keys to navigate, press Enter to select a device\n"
		return fmt.Sprint(s, footer)
	}

	s = "Please enter a port to listen on:\n\n%s\n\n%s\n%s"
	var packetsView string
	for _, packet := range m.packets {
		packetsView += fmt.Sprintf("- %s\n", packet)
	}

	return fmt.Sprintf(
		s,
		m.textinput.View(),
		packetsView,
		footer,
	)

}
