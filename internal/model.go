package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/table"
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

// Helper functions to create table columns for each resource type
func getServerColumns(width int) []table.Column {
	nameWidth := 35
	idWidth := 20
	statusWidth := 10

	if width < 100 {
		nameWidth = 30
		idWidth = 15
	} else if width > 140 {
		nameWidth = 45
	}

	return []table.Column{
		{Title: "Name", Width: nameWidth},
		{Title: "ID", Width: idWidth},
		{Title: "Status", Width: statusWidth},
	}
}

func getSwitchColumns(width int) []table.Column {
	nameWidth := 35
	idWidth := 20
	serversWidth := 10

	if width < 100 {
		return []table.Column{
			{Title: "Name", Width: 30},
			{Title: "ID", Width: 20},
		}
	} else if width < 140 {
		return []table.Column{
			{Title: "Name", Width: nameWidth},
			{Title: "ID", Width: idWidth},
			{Title: "Servers", Width: serversWidth},
		}
	} else {
		return []table.Column{
			{Title: "Name", Width: 40},
			{Title: "ID", Width: idWidth},
			{Title: "Servers", Width: serversWidth},
			{Title: "Route", Width: 15},
			{Title: "Created", Width: 12},
		}
	}
}

func getDNSColumns(width int) []table.Column {
	nameWidth := 35
	idWidth := 20
	recordsWidth := 10

	if width < 100 {
		return []table.Column{
			{Title: "Name", Width: 30},
			{Title: "ID", Width: 20},
		}
	} else if width < 140 {
		return []table.Column{
			{Title: "Name", Width: nameWidth},
			{Title: "ID", Width: idWidth},
			{Title: "Records", Width: recordsWidth},
		}
	} else {
		return []table.Column{
			{Title: "Name", Width: 40},
			{Title: "ID", Width: idWidth},
			{Title: "Records", Width: recordsWidth},
			{Title: "Created", Width: 12},
		}
	}
}

func getELBColumns(width int) []table.Column {
	if width < 100 {
		return []table.Column{
			{Title: "Name", Width: 30},
			{Title: "ID", Width: 20},
			{Title: "VIP", Width: 15},
		}
	} else if width < 140 {
		return []table.Column{
			{Title: "Name", Width: 35},
			{Title: "ID", Width: 20},
			{Title: "VIP", Width: 15},
			{Title: "Servers", Width: 10},
		}
	} else {
		return []table.Column{
			{Title: "Name", Width: 40},
			{Title: "ID", Width: 20},
			{Title: "VIP", Width: 15},
			{Title: "Servers", Width: 10},
			{Title: "Plan", Width: 15},
		}
	}
}

func getGSLBColumns(width int) []table.Column {
	if width < 100 {
		return []table.Column{
			{Title: "Name", Width: 30},
			{Title: "ID", Width: 20},
		}
	} else if width < 140 {
		return []table.Column{
			{Title: "Name", Width: 35},
			{Title: "ID", Width: 20},
			{Title: "FQDN", Width: 30},
		}
	} else {
		return []table.Column{
			{Title: "Name", Width: 40},
			{Title: "ID", Width: 20},
			{Title: "FQDN", Width: 40},
			{Title: "Servers", Width: 10},
		}
	}
}

func getDBColumns(width int) []table.Column {
	if width < 100 {
		return []table.Column{
			{Title: "Name", Width: 30},
			{Title: "ID", Width: 20},
			{Title: "Status", Width: 10},
		}
	} else if width < 140 {
		return []table.Column{
			{Title: "Name", Width: 35},
			{Title: "ID", Width: 20},
			{Title: "Type", Width: 10},
			{Title: "Status", Width: 10},
		}
	} else {
		return []table.Column{
			{Title: "Name", Width: 40},
			{Title: "ID", Width: 20},
			{Title: "Type", Width: 10},
			{Title: "Plan", Width: 20},
			{Title: "Status", Width: 10},
		}
	}
}

// Helper functions to convert resources to table rows
func serverToRow(server Server, width int) table.Row {
	statusStr := server.InstanceStatus
	switch server.InstanceStatus {
	case "UP":
		statusStr = upStatusStyle.Render(server.InstanceStatus)
	case "DOWN":
		statusStr = downStatusStyle.Render(server.InstanceStatus)
	default:
		statusStr = otherStatusStyle.Render(server.InstanceStatus)
	}

	return table.Row{
		server.Name,
		server.ID,
		statusStr,
	}
}

func switchToRow(sw Switch, width int) table.Row {
	if width < 100 {
		return table.Row{
			sw.Name,
			sw.ID,
		}
	} else if width < 140 {
		return table.Row{
			sw.Name,
			sw.ID,
			fmt.Sprintf("%d", sw.ServerCount),
		}
	} else {
		createdAt := ""
		if sw.CreatedAt != "" && len(sw.CreatedAt) >= 10 {
			createdAt = sw.CreatedAt[:10]
		}
		return table.Row{
			sw.Name,
			sw.ID,
			fmt.Sprintf("%d", sw.ServerCount),
			sw.DefaultRoute,
			createdAt,
		}
	}
}

func dnsToRow(dns DNS, width int) table.Row {
	if width < 100 {
		return table.Row{
			dns.Name,
			dns.ID,
		}
	} else if width < 140 {
		return table.Row{
			dns.Name,
			dns.ID,
			fmt.Sprintf("%d", dns.RecordCount),
		}
	} else {
		createdAt := ""
		if dns.CreatedAt != "" && len(dns.CreatedAt) >= 10 {
			createdAt = dns.CreatedAt[:10]
		}
		return table.Row{
			dns.Name,
			dns.ID,
			fmt.Sprintf("%d", dns.RecordCount),
			createdAt,
		}
	}
}

func elbToRow(elb ELB, width int) table.Row {
	if width < 100 {
		return table.Row{
			elb.Name,
			elb.ID,
			elb.VIP,
		}
	} else if width < 140 {
		return table.Row{
			elb.Name,
			elb.ID,
			elb.VIP,
			fmt.Sprintf("%d", elb.ServerCount),
		}
	} else {
		return table.Row{
			elb.Name,
			elb.ID,
			elb.VIP,
			fmt.Sprintf("%d", elb.ServerCount),
			elb.Plan,
		}
	}
}

func gslbToRow(gslb GSLB, width int) table.Row {
	if width < 100 {
		return table.Row{
			gslb.Name,
			gslb.ID,
		}
	} else if width < 140 {
		fqdn := gslb.FQDN
		if len(fqdn) > 30 {
			fqdn = fqdn[:27] + "..."
		}
		return table.Row{
			gslb.Name,
			gslb.ID,
			fqdn,
		}
	} else {
		fqdn := gslb.FQDN
		if len(fqdn) > 40 {
			fqdn = fqdn[:37] + "..."
		}
		return table.Row{
			gslb.Name,
			gslb.ID,
			fqdn,
			fmt.Sprintf("%d", gslb.ServerCount),
		}
	}
}

func dbToRow(db DB, width int) table.Row {
	statusStr := db.InstanceStatus
	switch db.InstanceStatus {
	case "UP":
		statusStr = upStatusStyle.Render(db.InstanceStatus)
	case "DOWN":
		statusStr = downStatusStyle.Render(db.InstanceStatus)
	default:
		statusStr = otherStatusStyle.Render(db.InstanceStatus)
	}

	if width < 100 {
		return table.Row{
			db.Name,
			db.ID,
			statusStr,
		}
	} else if width < 140 {
		return table.Row{
			db.Name,
			db.ID,
			db.DBType,
			statusStr,
		}
	} else {
		return table.Row{
			db.Name,
			db.ID,
			db.DBType,
			db.Plan,
			statusStr,
		}
	}
}

type model struct {
	client         *SakuraClient
	resourceTable  table.Model
	zones          []string
	currentZone    string
	cursor         int
	err            error
	loading        bool
	quitting       bool
	accountName    string
	windowHeight   int
	windowWidth    int
	searchMode     bool
	searchInput    textinput.Model
	searchQuery    string
	searchMatches  []int // Indices of matching items
	currentMatch   int   // Current match index in searchMatches
	detailMode     bool
	serverDetail   *ServerDetail
	switchDetail   *SwitchDetail
	dnsDetail      *DNSDetail
	elbDetail      *ELBDetail
	gslbDetail     *GSLBDetail
	dbDetail       *DBDetail
	detailLoading  bool
	resourceType   ResourceType
	detailViewport viewport.Model
	// Resource type selector
	resourceSelectMode   bool
	resourceSelectCursor int
	// Raw resource data (for search and detail view)
	servers  []Server
	switches []Switch
	dnsList  []DNS
	elbList  []ELB
	gslbList []GSLB
	dbList   []DB
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

type elbLoadedMsg struct {
	elbList []ELB
	err     error
}

type elbDetailLoadedMsg struct {
	detail *ELBDetail
	err    error
}

type gslbLoadedMsg struct {
	gslbList []GSLB
	err      error
}

type gslbDetailLoadedMsg struct {
	detail *GSLBDetail
	err    error
}

type dbLoadedMsg struct {
	dbList []DB
	err    error
}

type dbDetailLoadedMsg struct {
	detail *DBDetail
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

func loadELB(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		elbList, err := client.ListELB(ctx)
		return elbLoadedMsg{elbList: elbList, err: err}
	}
}

func loadELBDetail(client *SakuraClient, elbID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetELBDetail(ctx, elbID)
		if err != nil {
			slog.Error("Failed to load ELB detail", slog.Any("error", err))
			return elbDetailLoadedMsg{err: err}
		}
		slog.Info("ELB detail loaded successfully", slog.String("elbID", elbID))
		return elbDetailLoadedMsg{detail: detail}
	}
}

func loadGSLB(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		gslbList, err := client.ListGSLB(ctx)
		return gslbLoadedMsg{gslbList: gslbList, err: err}
	}
}

func loadGSLBDetail(client *SakuraClient, gslbID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetGSLBDetail(ctx, gslbID)
		if err != nil {
			slog.Error("Failed to load GSLB detail", slog.Any("error", err))
			return gslbDetailLoadedMsg{err: err}
		}
		slog.Info("GSLB detail loaded successfully", slog.String("gslbID", gslbID))
		return gslbDetailLoadedMsg{detail: detail}
	}
}

func loadDB(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		dbList, err := client.ListDB(ctx)
		return dbLoadedMsg{dbList: dbList, err: err}
	}
}

func loadDBDetail(client *SakuraClient, dbID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetDBDetail(ctx, dbID)
		if err != nil {
			slog.Error("Failed to load DB detail", slog.Any("error", err))
			return dbDetailLoadedMsg{err: err}
		}
		slog.Info("DB detail loaded successfully", slog.String("dbID", dbID))
		return dbDetailLoadedMsg{detail: detail}
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

	// Create table with server columns
	columns := getServerColumns(100) // Default width
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 50

	return model{
		client:        client,
		resourceTable: t,
		zones:         zones,
		currentZone:   defaultZone,
		cursor:        cursor,
		loading:       true,
		searchInput:   ti,
		resourceType:  ResourceTypeServer,
	}
}

// Rebuild table from current resource data
func (m *model) rebuildTable() {
	var rows []table.Row
	var columns []table.Column

	switch m.resourceType {
	case ResourceTypeServer:
		columns = getServerColumns(m.windowWidth)
		for _, server := range m.servers {
			rows = append(rows, serverToRow(server, m.windowWidth))
		}
	case ResourceTypeSwitch:
		columns = getSwitchColumns(m.windowWidth)
		for _, sw := range m.switches {
			rows = append(rows, switchToRow(sw, m.windowWidth))
		}
	case ResourceTypeDNS:
		columns = getDNSColumns(m.windowWidth)
		for _, dns := range m.dnsList {
			rows = append(rows, dnsToRow(dns, m.windowWidth))
		}
	case ResourceTypeELB:
		columns = getELBColumns(m.windowWidth)
		for _, elb := range m.elbList {
			rows = append(rows, elbToRow(elb, m.windowWidth))
		}
	case ResourceTypeGSLB:
		columns = getGSLBColumns(m.windowWidth)
		for _, gslb := range m.gslbList {
			rows = append(rows, gslbToRow(gslb, m.windowWidth))
		}
	case ResourceTypeDB:
		columns = getDBColumns(m.windowWidth)
		for _, db := range m.dbList {
			rows = append(rows, dbToRow(db, m.windowWidth))
		}
	}

	m.resourceTable.SetColumns(columns)
	m.resourceTable.SetRows(rows)
}

// Get the currently selected resource ID (for detail view)
func (m model) getSelectedResourceID() string {
	selectedIdx := m.resourceTable.Cursor()

	switch m.resourceType {
	case ResourceTypeServer:
		if selectedIdx >= 0 && selectedIdx < len(m.servers) {
			return m.servers[selectedIdx].ID
		}
	case ResourceTypeSwitch:
		if selectedIdx >= 0 && selectedIdx < len(m.switches) {
			return m.switches[selectedIdx].ID
		}
	case ResourceTypeDNS:
		if selectedIdx >= 0 && selectedIdx < len(m.dnsList) {
			return m.dnsList[selectedIdx].ID
		}
	case ResourceTypeELB:
		if selectedIdx >= 0 && selectedIdx < len(m.elbList) {
			return m.elbList[selectedIdx].ID
		}
	case ResourceTypeGSLB:
		if selectedIdx >= 0 && selectedIdx < len(m.gslbList) {
			return m.gslbList[selectedIdx].ID
		}
	case ResourceTypeDB:
		if selectedIdx >= 0 && selectedIdx < len(m.dbList) {
			return m.dbList[selectedIdx].ID
		}
	}
	return ""
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

	switch m.resourceType {
	case ResourceTypeServer:
		for i, server := range m.servers {
			if strings.Contains(strings.ToLower(server.Name), query) ||
				strings.Contains(strings.ToLower(server.ID), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
		}
	case ResourceTypeSwitch:
		for i, sw := range m.switches {
			if strings.Contains(strings.ToLower(sw.Name), query) ||
				strings.Contains(strings.ToLower(sw.ID), query) ||
				strings.Contains(strings.ToLower(sw.Desc), query) ||
				strings.Contains(strings.ToLower(sw.DefaultRoute), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
		}
	case ResourceTypeDNS:
		for i, dns := range m.dnsList {
			if strings.Contains(strings.ToLower(dns.Name), query) ||
				strings.Contains(strings.ToLower(dns.ID), query) ||
				strings.Contains(strings.ToLower(dns.Desc), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
		}
	case ResourceTypeELB:
		for i, elb := range m.elbList {
			if strings.Contains(strings.ToLower(elb.Name), query) ||
				strings.Contains(strings.ToLower(elb.ID), query) ||
				strings.Contains(strings.ToLower(elb.Desc), query) ||
				strings.Contains(strings.ToLower(elb.VIP), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
		}
	case ResourceTypeGSLB:
		for i, gslb := range m.gslbList {
			if strings.Contains(strings.ToLower(gslb.Name), query) ||
				strings.Contains(strings.ToLower(gslb.ID), query) ||
				strings.Contains(strings.ToLower(gslb.Desc), query) ||
				strings.Contains(strings.ToLower(gslb.FQDN), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
		}
	case ResourceTypeDB:
		for i, db := range m.dbList {
			if strings.Contains(strings.ToLower(db.Name), query) ||
				strings.Contains(strings.ToLower(db.ID), query) ||
				strings.Contains(strings.ToLower(db.Desc), query) ||
				strings.Contains(strings.ToLower(db.DBType), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
		}
	}

	slog.Info("Search performed", slog.String("query", query), slog.Int("matches", len(m.searchMatches)))

	// Jump to first match
	if len(m.searchMatches) > 0 {
		m.currentMatch = 0
		m.resourceTable.SetCursor(m.searchMatches[0])
	}
}

func (m *model) nextMatch() {
	if len(m.searchMatches) == 0 {
		return
	}
	m.currentMatch = (m.currentMatch + 1) % len(m.searchMatches)
	m.resourceTable.SetCursor(m.searchMatches[m.currentMatch])
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
	m.resourceTable.SetCursor(m.searchMatches[m.currentMatch])
	slog.Info("Previous match", slog.Int("match", m.currentMatch+1), slog.Int("total", len(m.searchMatches)))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
		slog.Debug("Window size updated", slog.Int("height", msg.Height), slog.Int("width", msg.Width))

		// Update table size - account for header area
		headerHeight := 8 // Title + Account + Zone + spacing
		if m.accountName == "" {
			headerHeight--
		}
		if m.searchMode {
			headerHeight++ // Add line for search input
		}
		m.resourceTable.SetHeight(msg.Height - headerHeight)
		// Rebuild table with new width-dependent columns
		m.rebuildTable()

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
				m.elbDetail = nil
				m.gslbDetail = nil
				m.dbDetail = nil
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

		// Handle resource select mode
		if m.resourceSelectMode {
			switch msg.String() {
			case "j", "down":
				if m.resourceSelectCursor < len(AllResourceTypes)-1 {
					m.resourceSelectCursor++
				}
				return m, nil
			case "k", "up":
				if m.resourceSelectCursor > 0 {
					m.resourceSelectCursor--
				}
				return m, nil
			case "enter":
				selectedType := AllResourceTypes[m.resourceSelectCursor]
				if selectedType != m.resourceType {
					m.resourceType = selectedType
					m.list.Title = selectedType.String()
					if m.resourceType == ResourceTypeServer {
						m.list.Title = "Servers"
					} else if m.resourceType == ResourceTypeSwitch {
						m.list.Title = "Switches"
					}
					slog.Info("User switched resource type",
						slog.String("type", m.resourceType.String()))
					m.loading = true
					// Clear search when switching resource types
					m.searchQuery = ""
					m.searchMatches = []int{}
					m.currentMatch = -1
					m.resourceSelectMode = false
					// Load appropriate resources
					switch m.resourceType {
					case ResourceTypeServer:
						return m, loadServers(m.client)
					case ResourceTypeSwitch:
						return m, loadSwitches(m.client)
					case ResourceTypeDNS:
						return m, loadDNS(m.client)
					case ResourceTypeELB:
						return m, loadELB(m.client)
					case ResourceTypeGSLB:
						return m, loadGSLB(m.client)
					case ResourceTypeDB:
						return m, loadDB(m.client)
					}
				}
				m.resourceSelectMode = false
				return m, nil
			case "esc", "t", "q":
				m.resourceSelectMode = false
				return m, nil
			}
			return m, nil
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
			resourceID := m.getSelectedResourceID()
			if resourceID != "" {
				m.detailMode = true
				m.detailLoading = true
				switch m.resourceType {
				case ResourceTypeServer:
					return m, loadServerDetail(m.client, resourceID)
				case ResourceTypeSwitch:
					return m, loadSwitchDetail(m.client, resourceID)
				case ResourceTypeDNS:
					return m, loadDNSDetail(m.client, resourceID)
				case ResourceTypeELB:
					return m, loadELBDetail(m.client, resourceID)
				case ResourceTypeGSLB:
					return m, loadGSLBDetail(m.client, resourceID)
				case ResourceTypeDB:
					return m, loadDBDetail(m.client, resourceID)
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
			// Open resource type selector menu
			m.resourceSelectMode = true
			// Set cursor to current resource type
			for i, rt := range AllResourceTypes {
				if rt == m.resourceType {
					m.resourceSelectCursor = i
					break
				}
			}
			return m, nil

		case "z":
			// Zone switching only affects Server and Switch (DNS, ELB, and GSLB are global)
			if m.resourceType == ResourceTypeDNS || m.resourceType == ResourceTypeELB || m.resourceType == ResourceTypeGSLB {
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
			case ResourceTypeDB:
				return m, loadDB(m.client)
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
			case ResourceTypeELB:
				return m, loadELB(m.client)
			case ResourceTypeGSLB:
				return m, loadGSLB(m.client)
			case ResourceTypeDB:
				return m, loadDB(m.client)
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

		// Store server data and rebuild table
		m.servers = msg.servers
		m.rebuildTable()
		return m, nil

	case switchesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load switches", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Switches loaded successfully", slog.Int("count", len(msg.switches)))

		// Store switch data and rebuild table
		m.switches = msg.switches
		m.rebuildTable()
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

		// Store DNS data and rebuild table
		m.dnsList = msg.dnsList
		m.rebuildTable()
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

	case elbLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load ELBs", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("ELBs loaded successfully", slog.Int("count", len(msg.elbList)))

		// Store ELB data and rebuild table
		m.elbList = msg.elbList
		m.rebuildTable()
		return m, nil

	case elbDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load ELB detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.elbDetail = msg.detail
		// Setup viewport for detail view
		content := renderELBDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case gslbLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load GSLBs", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("GSLBs loaded successfully", slog.Int("count", len(msg.gslbList)))

		// Store GSLB data and rebuild table
		m.gslbList = msg.gslbList
		m.rebuildTable()
		return m, nil

	case gslbDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load GSLB detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.gslbDetail = msg.detail
		// Setup viewport for detail view
		content := renderGSLBDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case dbLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load DBs", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("DBs loaded successfully", slog.Int("count", len(msg.dbList)))

		// Store DB data and rebuild table
		m.dbList = msg.dbList
		m.rebuildTable()
		return m, nil

	case dbDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load DB detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.dbDetail = msg.detail
		// Setup viewport for detail view
		content := renderDBDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil
	}

	// Delegate to table for navigation
	var cmd tea.Cmd
	m.resourceTable, cmd = m.resourceTable.Update(msg)
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
		} else if m.serverDetail != nil || m.switchDetail != nil || m.dnsDetail != nil || m.elbDetail != nil || m.gslbDetail != nil || m.dbDetail != nil {
			b.WriteString(m.detailViewport.View())
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("↑/↓/j/k: scroll | ESC/q: back"))
		}
		return b.String()
	}

	// Resource type selector mode
	if m.resourceSelectMode {
		b.WriteString(titleStyle.Render("Select Resource Type"))
		b.WriteString("\n")
		for i, rt := range AllResourceTypes {
			if i == m.resourceSelectCursor {
				b.WriteString(selectedItemStyle.Render(fmt.Sprintf("▸ %s", rt.String())))
			} else {
				b.WriteString(itemStyle.Render(fmt.Sprintf("  %s", rt.String())))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("↑/↓/j/k: move | Enter: select | Esc: cancel"))
		return b.String()
	}

	// Zone selector and resource type
	// Show "global" for global resources (DNS, ELB, GSLB)
	if m.resourceType == ResourceTypeDNS || m.resourceType == ResourceTypeELB || m.resourceType == ResourceTypeGSLB {
		b.WriteString("Zone: ")
		b.WriteString(zoneStyle.Render("global"))
		b.WriteString(" | Type: ")
		b.WriteString(selectedStyle.Render(m.resourceType.String()))
		b.WriteString("\n")
	} else {
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
	}

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

	// Resource table or loading/error
	if m.loading {
		b.WriteString(fmt.Sprintf("Loading %s...\n", strings.ToLower(m.resourceType.String())))
	} else if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else {
		// Add title above the table
		var title string
		switch m.resourceType {
		case ResourceTypeServer:
			title = "Servers"
		case ResourceTypeSwitch:
			title = "Switches"
		case ResourceTypeDNS:
			title = "DNS Zones"
		case ResourceTypeELB:
			title = "Load Balancers"
		case ResourceTypeGSLB:
			title = "GSLB"
		case ResourceTypeDB:
			title = "Databases"
		}
		b.WriteString(titleStyle.Render(title))
		b.WriteString("\n")
		b.WriteString(m.resourceTable.View())
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

func renderELBDetail(detail *ELBDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("ELB: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("VIP:         %s\n", detail.VIP))
	b.WriteString(fmt.Sprintf("Plan:        %s\n", detail.Plan))

	if detail.FQDN != "" {
		b.WriteString(fmt.Sprintf("FQDN:        %s\n", detail.FQDN))
	}

	b.WriteString(fmt.Sprintf("Servers:     %d\n", detail.ServerCount))

	// Display server list in table format
	if len(detail.Servers) > 0 {
		b.WriteString("\nServers:\n")
		b.WriteString(fmt.Sprintf("  %-20s %-8s %s\n", "IP Address", "Port", "Status"))
		b.WriteString(fmt.Sprintf("  %-20s %-8s %s\n", "----------", "----", "------"))
		for _, server := range detail.Servers {
			status := "Disabled"
			if server.Enabled {
				status = "Enabled"
			}
			b.WriteString(fmt.Sprintf("  %-20s %-8d %s\n",
				server.IPAddress,
				server.Port,
				status))
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

func renderGSLBDetail(detail *GSLBDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("GSLB: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("FQDN:        %s\n", detail.FQDN))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Servers:     %d\n", detail.ServerCount))

	// Display health check settings
	if detail.HealthPath != "" {
		b.WriteString(fmt.Sprintf("Health Path: %s\n", detail.HealthPath))
	}
	b.WriteString(fmt.Sprintf("Delay Loop:  %d sec\n", detail.DelayLoop))
	if detail.Weighted {
		b.WriteString("Weighted:    Yes\n")
	} else {
		b.WriteString("Weighted:    No\n")
	}

	// Display server list in table format using bubbles table
	if len(detail.Servers) > 0 {
		b.WriteString("\nServers:\n")

		var columns []table.Column
		var rows []table.Row

		if detail.Weighted {
			columns = []table.Column{
				{Title: "IP Address", Width: 20},
				{Title: "Weight", Width: 8},
				{Title: "Status", Width: 10},
			}

			for _, server := range detail.Servers {
				status := "Disabled"
				if server.Enabled {
					status = "Enabled"
				}
				rows = append(rows, table.Row{
					server.IPAddress,
					fmt.Sprintf("%d", server.Weight),
					status,
				})
			}
		} else {
			columns = []table.Column{
				{Title: "IP Address", Width: 20},
				{Title: "Status", Width: 10},
			}

			for _, server := range detail.Servers {
				status := "Disabled"
				if server.Enabled {
					status = "Enabled"
				}
				rows = append(rows, table.Row{
					server.IPAddress,
					status,
				})
			}
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(false),
			table.WithHeight(len(rows)),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		b.WriteString(t.View())
		b.WriteString("\n")
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderDBDetail(detail *DBDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("DB: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("DB Type:     %s\n", detail.DBType))
	b.WriteString(fmt.Sprintf("Status:      %s\n", detail.InstanceStatus))
	b.WriteString(fmt.Sprintf("Plan:        %s\n", detail.Plan))
	b.WriteString(fmt.Sprintf("CPU:         %d Core(s)\n", detail.CPU))
	b.WriteString(fmt.Sprintf("Memory:      %d GB\n", detail.MemoryGB))
	b.WriteString(fmt.Sprintf("Disk Size:   %d GB\n", detail.DiskSizeGB))

	if detail.IPAddress != "" {
		b.WriteString(fmt.Sprintf("IP Address:  %s\n", detail.IPAddress))
		if detail.NetworkMaskLen > 0 {
			b.WriteString(fmt.Sprintf("Netmask:     /%d\n", detail.NetworkMaskLen))
		}
		if detail.DefaultRoute != "" {
			b.WriteString(fmt.Sprintf("Gateway:     %s\n", detail.DefaultRoute))
		}
	}

	if detail.Port > 0 {
		b.WriteString(fmt.Sprintf("Port:        %d\n", detail.Port))
	}

	if detail.DefaultUser != "" {
		b.WriteString(fmt.Sprintf("User:        %s\n", detail.DefaultUser))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}
