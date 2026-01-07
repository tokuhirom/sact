package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("5")).
			MarginBottom(1)

	zoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(1)
)

type model struct {
	client      *SakuraClient
	servers     []Server
	zones       []string
	currentZone string
	cursor      int
	err         error
	loading     bool
	quitting    bool
}

type serversLoadedMsg struct {
	servers []Server
	err     error
}

func loadServers(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		servers, err := client.ListServers(ctx)
		return serversLoadedMsg{servers: servers, err: err}
	}
}

func InitialModel(client *SakuraClient, defaultZone string) model {
	zones := []string{"tk1a", "tk1b", "is1a", "is1b", "is1c"}

	cursor := 0
	for i, zone := range zones {
		if zone == defaultZone {
			cursor = i
			break
		}
	}

	return model{
		client:      client,
		servers:     []Server{},
		zones:       zones,
		currentZone: defaultZone,
		cursor:      cursor,
		loading:     true,
	}
}

func (m model) Init() tea.Cmd {
	slog.Info("Initializing TUI model", slog.String("zone", m.currentZone))
	return loadServers(m.client)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			slog.Info("User requested quit")
			m.quitting = true
			return m, tea.Quit

		case "z":
			oldZone := m.currentZone
			m.cursor = (m.cursor + 1) % len(m.zones)
			m.currentZone = m.zones[m.cursor]
			slog.Info("User switched zone via keyboard",
				slog.String("from", oldZone),
				slog.String("to", m.currentZone))
			m.client.SetZone(m.currentZone)
			m.loading = true
			return m, loadServers(m.client)

		case "r":
			slog.Info("User requested refresh", slog.String("zone", m.currentZone))
			m.loading = true
			return m, loadServers(m.client)
		}

	case serversLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load servers", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Servers loaded successfully", slog.Int("count", len(msg.servers)))
		m.servers = msg.servers
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("sact - さくらのクラウド TUI"))
	b.WriteString("\n\n")

	b.WriteString("Zone: ")
	for i, zone := range m.zones {
		if i == m.cursor {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("[%s]", zone)))
		} else {
			b.WriteString(zoneStyle.Render(fmt.Sprintf(" %s ", zone)))
		}
		b.WriteString(" ")
	}
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("Loading servers...\n")
	} else if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else {
		b.WriteString(fmt.Sprintf("Servers (%d):\n\n", len(m.servers)))

		if len(m.servers) == 0 {
			b.WriteString("No servers found in this zone.\n")
		} else {
			b.WriteString(fmt.Sprintf("%-20s %-40s %-15s\n", "ID", "Name", "Status"))
			b.WriteString(strings.Repeat("-", 80))
			b.WriteString("\n")

			for _, server := range m.servers {
				status := server.InstanceStatus
				statusStyle := lipgloss.NewStyle()
				switch status {
				case "UP":
					statusStyle = statusStyle.Foreground(lipgloss.Color("10"))
				case "DOWN":
					statusStyle = statusStyle.Foreground(lipgloss.Color("8"))
				default:
					statusStyle = statusStyle.Foreground(lipgloss.Color("11"))
				}

				b.WriteString(fmt.Sprintf("%-20s %-40s %s\n",
					server.ID,
					server.Name,
					statusStyle.Render(status),
				))
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("z: switch zone | r: refresh | q: quit"))

	return b.String()
}
