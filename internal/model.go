package internal

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("5")).
			MarginBottom(1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Italic(true)

	zoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("10"))
	upStatusStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	downStatusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	otherStatusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

// Custom delegate for single-line resource display (handles Server and Switch)
type resourceDelegate struct{}

func (d resourceDelegate) Height() int                             { return 1 }
func (d resourceDelegate) Spacing() int                            { return 0 }
func (d resourceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d resourceDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var str string

	// Handle Server
	if server, ok := item.(Server); ok {
		statusStyle := otherStatusStyle
		switch server.InstanceStatus {
		case "UP":
			statusStyle = upStatusStyle
		case "DOWN":
			statusStyle = downStatusStyle
		}

		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("▸ %-40s ID: %-20s Status: %s",
				server.Name,
				server.ID,
				statusStyle.Render(server.InstanceStatus)))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s ID: %-20s Status: %s",
				server.Name,
				server.ID,
				statusStyle.Render(server.InstanceStatus)))
		}
	} else if sw, ok := item.(Switch); ok {
		// Handle Switch
		width := m.Width()

		// Build info string based on available width
		var info string
		if width < 100 {
			// Compact mode: just ID
			info = fmt.Sprintf("ID: %-20s", sw.ID)
		} else if width < 140 {
			// Medium mode: ID + ServerCount
			info = fmt.Sprintf("ID: %-20s Servers: %d", sw.ID, sw.ServerCount)
		} else {
			// Full mode: ID + ServerCount + DefaultRoute + CreatedAt
			routeInfo := ""
			if sw.DefaultRoute != "" {
				routeInfo = fmt.Sprintf("Route: %-15s", sw.DefaultRoute)
			}
			dateInfo := ""
			if sw.CreatedAt != "" {
				dateInfo = fmt.Sprintf("Created: %s", sw.CreatedAt)
			}

			info = fmt.Sprintf("ID: %-20s Servers: %-2d %s %s",
				sw.ID, sw.ServerCount, routeInfo, dateInfo)
		}

		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("▸ %-40s %s",
				sw.Name,
				info))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %s",
				sw.Name,
				info))
		}
	} else if dns, ok := item.(DNS); ok {
		// Handle DNS
		width := m.Width()

		// Build info string based on available width
		var info string
		if width < 100 {
			// Compact mode: just ID
			info = fmt.Sprintf("ID: %-20s", dns.ID)
		} else if width < 140 {
			// Medium mode: ID + RecordCount
			info = fmt.Sprintf("ID: %-20s Records: %d", dns.ID, dns.RecordCount)
		} else {
			// Full mode: ID + RecordCount + CreatedAt
			dateInfo := ""
			if dns.CreatedAt != "" {
				dateInfo = fmt.Sprintf("Created: %s", dns.CreatedAt)
			}

			info = fmt.Sprintf("ID: %-20s Records: %-3d %s",
				dns.ID, dns.RecordCount, dateInfo)
		}

		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("▸ %-40s %s",
				dns.Name,
				info))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %s",
				dns.Name,
				info))
		}
	} else {
		return
	}

	fmt.Fprint(w, str)
}

type model struct {
	client        *SakuraClient
	list          list.Model
	zones         []string
	currentZone   string
	cursor        int
	err           error
	loading       bool
	quitting      bool
	accountName   string
	windowHeight  int
	windowWidth   int
	searchMode    bool
	searchInput   textinput.Model
	searchQuery   string
	searchMatches []int // Indices of matching items
	currentMatch  int   // Current match index in searchMatches
	detailMode    bool
	serverDetail  *ServerDetail
	switchDetail  *SwitchDetail
	dnsDetail     *DNSDetail
	detailLoading bool
	resourceType  ResourceType
	detailViewport viewport.Model
}

type serversLoadedMsg struct {
	servers []Server
	err     error
}

type switchesLoadedMsg struct {
	switches []Switch
	err      error
}

type authStatusLoadedMsg struct {
	accountName string
	err         error
}

type serverDetailLoadedMsg struct {
	detail *ServerDetail
	err    error
}

type switchDetailLoadedMsg struct {
	detail *SwitchDetail
	err    error
}

type dnsLoadedMsg struct {
	dnsList []DNS
	err     error
}

type dnsDetailLoadedMsg struct {
	detail *DNSDetail
	err    error
}

func loadServers(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		servers, err := client.ListServers(ctx)
		return serversLoadedMsg{servers: servers, err: err}
	}
}

func loadSwitches(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		switches, err := client.ListSwitches(ctx)
		return switchesLoadedMsg{switches: switches, err: err}
	}
}

func loadAuthStatus(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		authStatus, err := client.GetAuthStatus(ctx)
		if err != nil {
			slog.Error("Failed to load auth status", slog.Any("error", err))
			return authStatusLoadedMsg{err: err}
		}
		slog.Info("Auth status loaded successfully", slog.String("account", authStatus.AccountName))
		return authStatusLoadedMsg{accountName: authStatus.AccountName}
	}
}

func loadServerDetail(client *SakuraClient, serverID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetServerDetail(ctx, serverID)
		if err != nil {
			slog.Error("Failed to load server detail", slog.Any("error", err))
			return serverDetailLoadedMsg{err: err}
		}
		slog.Info("Server detail loaded successfully", slog.String("serverID", serverID))
		return serverDetailLoadedMsg{detail: detail}
	}
}

func loadSwitchDetail(client *SakuraClient, switchID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetSwitchDetail(ctx, switchID)
		if err != nil {
			slog.Error("Failed to load switch detail", slog.Any("error", err))
			return switchDetailLoadedMsg{err: err}
		}
		slog.Info("Switch detail loaded successfully", slog.String("switchID", switchID))
		return switchDetailLoadedMsg{detail: detail}
	}
}

func loadDNS(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		dnsList, err := client.ListDNS(ctx)
		return dnsLoadedMsg{dnsList: dnsList, err: err}
	}
}

func loadDNSDetail(client *SakuraClient, dnsID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetDNSDetail(ctx, dnsID)
		if err != nil {
			slog.Error("Failed to load DNS detail", slog.Any("error", err))
			return dnsDetailLoadedMsg{err: err}
		}
		slog.Info("DNS detail loaded successfully", slog.String("dnsID", dnsID))
		return dnsDetailLoadedMsg{detail: detail}
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

	// Create list with custom delegate
	delegate := resourceDelegate{}
	resourceList := list.New([]list.Item{}, delegate, 0, 0)
	resourceList.Title = "Servers"
	resourceList.SetShowStatusBar(false)
	resourceList.SetFilteringEnabled(false) // Disable built-in filtering
	resourceList.Styles.Title = titleStyle

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 50

	return model{
		client:       client,
		list:         resourceList,
		zones:        zones,
		currentZone:  defaultZone,
		cursor:       cursor,
		loading:      true,
		searchInput:  ti,
		resourceType: ResourceTypeServer,
	}
}

func (m model) Init() tea.Cmd {
	slog.Info("Initializing TUI model", slog.String("zone", m.currentZone))
	return tea.Batch(
		loadServers(m.client),
		loadAuthStatus(m.client),
	)
}

func (m *model) performSearch() {
	m.searchMatches = []int{}
	m.currentMatch = -1

	query := strings.ToLower(m.searchQuery)
	if query == "" {
		return
	}

	items := m.list.Items()
	for i, item := range items {
		// Handle Server
		if server, ok := item.(Server); ok {
			if strings.Contains(strings.ToLower(server.Name), query) ||
				strings.Contains(strings.ToLower(server.ID), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
			continue
		}
		// Handle Switch
		if sw, ok := item.(Switch); ok {
			if strings.Contains(strings.ToLower(sw.Name), query) ||
				strings.Contains(strings.ToLower(sw.ID), query) ||
				strings.Contains(strings.ToLower(sw.Desc), query) ||
				strings.Contains(strings.ToLower(sw.DefaultRoute), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
			continue
		}
		// Handle DNS
		if dns, ok := item.(DNS); ok {
			if strings.Contains(strings.ToLower(dns.Name), query) ||
				strings.Contains(strings.ToLower(dns.ID), query) ||
				strings.Contains(strings.ToLower(dns.Desc), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
		}
	}

	slog.Info("Search performed", slog.String("query", query), slog.Int("matches", len(m.searchMatches)))

	// Jump to first match
	if len(m.searchMatches) > 0 {
		m.currentMatch = 0
		m.list.Select(m.searchMatches[0])
	}
}

func (m *model) nextMatch() {
	if len(m.searchMatches) == 0 {
		return
	}
	m.currentMatch = (m.currentMatch + 1) % len(m.searchMatches)
	m.list.Select(m.searchMatches[m.currentMatch])
	slog.Info("Next match", slog.Int("match", m.currentMatch+1), slog.Int("total", len(m.searchMatches)))
}

func (m *model) prevMatch() {
	if len(m.searchMatches) == 0 {
		return
	}
	m.currentMatch--
	if m.currentMatch < 0 {
		m.currentMatch = len(m.searchMatches) - 1
	}
	m.list.Select(m.searchMatches[m.currentMatch])
	slog.Info("Previous match", slog.Int("match", m.currentMatch+1), slog.Int("total", len(m.searchMatches)))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
		slog.Debug("Window size updated", slog.Int("height", msg.Height), slog.Int("width", msg.Width))

		// Update list size - account for header area
		headerHeight := 8 // Title + Account + Zone + spacing
		if m.accountName == "" {
			headerHeight--
		}
		if m.searchMode {
			headerHeight++ // Add line for search input
		}
		m.list.SetSize(msg.Width, msg.Height-headerHeight)

		// Update detail viewport size if in detail mode
		if m.detailMode {
			m.detailViewport.Width = msg.Width
			m.detailViewport.Height = msg.Height - 10
		}
		return m, nil

	case tea.KeyMsg:
		// Handle detail mode
		if m.detailMode {
			switch msg.String() {
			case "esc", "q":
				m.detailMode = false
				m.serverDetail = nil
				m.switchDetail = nil
				m.dnsDetail = nil
				return m, nil
			default:
				// Pass other keys to viewport for scrolling
				var cmd tea.Cmd
				m.detailViewport, cmd = m.detailViewport.Update(msg)
				return m, cmd
			}
		}

		// Handle search mode
		if m.searchMode {
			switch msg.String() {
			case "enter":
				m.searchQuery = m.searchInput.Value()
				m.searchMode = false
				m.performSearch()
				return m, nil
			case "esc":
				m.searchMode = false
				m.searchInput.Reset()
				return m, nil
			default:
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				return m, cmd
			}
		}

		// Normal mode
		switch msg.String() {
		case "ctrl+c", "q":
			slog.Info("User requested quit")
			m.quitting = true
			return m, tea.Quit

		case "esc":
			// Ignore ESC in list view to prevent accidental exit
			return m, nil

		case "enter":
			// Show detail based on resource type
			if len(m.list.Items()) > 0 {
				selectedItem := m.list.SelectedItem()
				if server, ok := selectedItem.(Server); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadServerDetail(m.client, server.ID)
				}
				if sw, ok := selectedItem.(Switch); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadSwitchDetail(m.client, sw.ID)
				}
				if dns, ok := selectedItem.(DNS); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadDNSDetail(m.client, dns.ID)
				}
			}
			return m, nil

		case "/":
			m.searchMode = true
			m.searchInput.Focus()
			m.searchInput.Reset()
			return m, textinput.Blink

		case "n":
			m.nextMatch()
			return m, nil

		case "N":
			m.prevMatch()
			return m, nil

		case "t":
			// Cycle through resource types: Server -> Switch -> DNS -> Server
			switch m.resourceType {
			case ResourceTypeServer:
				m.resourceType = ResourceTypeSwitch
				m.list.Title = "Switches"
			case ResourceTypeSwitch:
				m.resourceType = ResourceTypeDNS
				m.list.Title = "DNS"
			case ResourceTypeDNS:
				m.resourceType = ResourceTypeServer
				m.list.Title = "Servers"
			}
			slog.Info("User switched resource type",
				slog.String("type", m.resourceType.String()))
			m.loading = true
			// Clear search when switching resource types
			m.searchQuery = ""
			m.searchMatches = []int{}
			m.currentMatch = -1
			// Load appropriate resources
			switch m.resourceType {
			case ResourceTypeServer:
				return m, loadServers(m.client)
			case ResourceTypeSwitch:
				return m, loadSwitches(m.client)
			case ResourceTypeDNS:
				return m, loadDNS(m.client)
			}
			return m, nil

		case "z":
			// Zone switching only affects Server and Switch (DNS is global)
			if m.resourceType == ResourceTypeDNS {
				return m, nil
			}
			oldZone := m.currentZone
			m.cursor = (m.cursor + 1) % len(m.zones)
			m.currentZone = m.zones[m.cursor]
			slog.Info("User switched zone via keyboard",
				slog.String("from", oldZone),
				slog.String("to", m.currentZone))
			m.client.SetZone(m.currentZone)
			m.loading = true
			// Clear search when switching zones
			m.searchQuery = ""
			m.searchMatches = []int{}
			m.currentMatch = -1
			// Load appropriate resources based on current type
			switch m.resourceType {
			case ResourceTypeServer:
				return m, loadServers(m.client)
			case ResourceTypeSwitch:
				return m, loadSwitches(m.client)
			}
			return m, nil

		case "r":
			slog.Info("User requested refresh", slog.String("zone", m.currentZone))
			m.loading = true
			// Refresh appropriate resources based on current type
			switch m.resourceType {
			case ResourceTypeServer:
				return m, loadServers(m.client)
			case ResourceTypeSwitch:
				return m, loadSwitches(m.client)
			case ResourceTypeDNS:
				return m, loadDNS(m.client)
			}
			return m, nil
		}

	case serversLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load servers", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Servers loaded successfully", slog.Int("count", len(msg.servers)))

		// Convert servers to list items
		items := make([]list.Item, len(msg.servers))
		for i, server := range msg.servers {
			items[i] = server
		}
		m.list.SetItems(items)
		return m, nil

	case switchesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load switches", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Switches loaded successfully", slog.Int("count", len(msg.switches)))

		// Convert switches to list items
		items := make([]list.Item, len(msg.switches))
		for i, sw := range msg.switches {
			items[i] = sw
		}
		m.list.SetItems(items)
		return m, nil

	case authStatusLoadedMsg:
		if msg.err != nil {
			slog.Error("Failed to load auth status", slog.Any("error", msg.err))
			return m, nil
		}
		slog.Info("Setting account name in model", slog.String("accountName", msg.accountName))
		m.accountName = msg.accountName
		return m, nil

	case serverDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load server detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.serverDetail = msg.detail
		// Setup viewport for detail view
		content := renderServerDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case switchDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load switch detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.switchDetail = msg.detail
		// Setup viewport for detail view
		content := renderSwitchDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case dnsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load DNS zones", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("DNS zones loaded successfully", slog.Int("count", len(msg.dnsList)))

		// Convert DNS to list items
		items := make([]list.Item, len(msg.dnsList))
		for i, dns := range msg.dnsList {
			items[i] = dns
		}
		m.list.SetItems(items)
		return m, nil

	case dnsDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load DNS detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.dnsDetail = msg.detail
		// Setup viewport for detail view
		content := renderDNSDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil
	}

	// Delegate to list for navigation
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var b strings.Builder

	// Header
	if m.accountName != "" {
		b.WriteString(statusBarStyle.Render(fmt.Sprintf("Account: %s", m.accountName)))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Detail mode view
	if m.detailMode {
		if m.detailLoading {
			b.WriteString("Loading details...\n")
		} else if m.serverDetail != nil || m.switchDetail != nil || m.dnsDetail != nil {
			b.WriteString(m.detailViewport.View())
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("↑/↓/j/k: scroll | ESC/q: back"))
		}
		return b.String()
	}

	// Zone selector and resource type
	b.WriteString("Zone: ")
	for i, zone := range m.zones {
		if i == m.cursor {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("[%s]", zone)))
		} else {
			b.WriteString(zoneStyle.Render(fmt.Sprintf(" %s ", zone)))
		}
		b.WriteString(" ")
	}
	b.WriteString(" | Type: ")
	b.WriteString(selectedStyle.Render(m.resourceType.String()))
	b.WriteString("\n")

	// Search mode or search status
	if m.searchMode {
		b.WriteString("/" + m.searchInput.View())
		b.WriteString("\n")
	} else if m.searchQuery != "" {
		matchInfo := fmt.Sprintf("Search: %s (%d matches)", m.searchQuery, len(m.searchMatches))
		if len(m.searchMatches) > 0 {
			matchInfo = fmt.Sprintf("Search: %s (%d/%d)", m.searchQuery, m.currentMatch+1, len(m.searchMatches))
		}
		b.WriteString(helpStyle.Render(matchInfo))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	// Resource list or loading/error
	if m.loading {
		b.WriteString(fmt.Sprintf("Loading %s...\n", strings.ToLower(m.resourceType.String())))
	} else if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else {
		b.WriteString(m.list.View())
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Enter: details | /: search | n/N: next/prev | t: type | z: zone | r: refresh | q: quit"))
	}

	return b.String()
}

func renderServerDetail(detail *ServerDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Server: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Status:      %s\n", detail.InstanceStatus))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Description != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Description))
	}

	b.WriteString(fmt.Sprintf("CPU:         %d Core(s)\n", detail.CPU))
	b.WriteString(fmt.Sprintf("Memory:      %d GB\n", detail.MemoryGB))

	if len(detail.IPAddresses) > 0 {
		b.WriteString(fmt.Sprintf("IP Address:  %s\n", strings.Join(detail.IPAddresses, ", ")))
	}

	if len(detail.UserIPAddresses) > 0 {
		b.WriteString(fmt.Sprintf("User IP:     %s\n", strings.Join(detail.UserIPAddresses, ", ")))
	}

	if len(detail.Disks) > 0 {
		b.WriteString("\nDisks:\n")
		for _, disk := range detail.Disks {
			b.WriteString(fmt.Sprintf("  - %s (%d GB)\n", disk.Name, disk.SizeGB))
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderSwitchDetail(detail *SwitchDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Switch: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Subnets:     %d\n", detail.SubnetCount))
	b.WriteString(fmt.Sprintf("Servers:     %d connected\n", detail.ServerCount))

	if detail.DefaultRoute != "" {
		b.WriteString(fmt.Sprintf("Route:       %s\n", detail.DefaultRoute))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderDNSDetail(detail *DNSDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("DNS: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Records:     %d\n", detail.RecordCount))

	// Display DNS records in table format
	if len(detail.Records) > 0 {
		b.WriteString("\nDNS Records:\n")
		b.WriteString(fmt.Sprintf("  %-8s %-30s %-8s %s\n", "Type", "Name", "TTL", "Data"))
		b.WriteString(fmt.Sprintf("  %-8s %-30s %-8s %s\n", "----", "----", "---", "----"))
		for _, rec := range detail.Records {
			// Truncate long RData
			rdata := rec.RData
			if len(rdata) > 60 {
				rdata = rdata[:57] + "..."
			}
			b.WriteString(fmt.Sprintf("  %-8s %-30s %-8d %s\n",
				rec.Type,
				rec.Name,
				rec.TTL,
				rdata))
		}
	}

	if len(detail.NameServers) > 0 {
		b.WriteString("\nName Servers:\n")
		for _, ns := range detail.NameServers {
			b.WriteString(fmt.Sprintf("  - %s\n", ns))
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	if detail.ModifiedAt != "" {
		b.WriteString(fmt.Sprintf("Modified:    %s\n", detail.ModifiedAt))
	}

	return b.String()
}
