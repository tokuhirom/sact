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
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %s",
				server.Name,
				server.ID,
				statusStyle.Render(server.InstanceStatus)))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %s",
				server.Name,
				server.ID,
				statusStyle.Render(server.InstanceStatus)))
		}
	} else if sw, ok := item.(Switch); ok {
		// Handle Switch
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %d",
				sw.Name,
				sw.ID,
				sw.ServerCount))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %d",
				sw.Name,
				sw.ID,
				sw.ServerCount))
		}
	} else if dns, ok := item.(DNS); ok {
		// Handle DNS
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %d",
				dns.Name,
				dns.ID,
				dns.RecordCount))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %d",
				dns.Name,
				dns.ID,
				dns.RecordCount))
		}
	} else if elb, ok := item.(ELB); ok {
		// Handle ELB
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %-15s %d",
				elb.Name,
				elb.ID,
				elb.VIP,
				elb.ServerCount))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %-15s %d",
				elb.Name,
				elb.ID,
				elb.VIP,
				elb.ServerCount))
		}
	} else if gslb, ok := item.(GSLB); ok {
		// Handle GSLB
		fqdn := gslb.FQDN
		if len(fqdn) > 30 {
			fqdn = fqdn[:27] + "..."
		}
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %s",
				gslb.Name,
				gslb.ID,
				fqdn))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %s",
				gslb.Name,
				gslb.ID,
				fqdn))
		}
	} else if db, ok := item.(DB); ok {
		// Handle DB
		statusStyle := otherStatusStyle
		switch db.InstanceStatus {
		case "UP":
			statusStyle = upStatusStyle
		case "DOWN":
			statusStyle = downStatusStyle
		}

		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %-10s %s",
				db.Name,
				db.ID,
				db.DBType,
				statusStyle.Render(db.InstanceStatus)))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %-10s %s",
				db.Name,
				db.ID,
				db.DBType,
				statusStyle.Render(db.InstanceStatus)))
		}
	} else if disk, ok := item.(Disk); ok {
		// Handle Disk
		serverInfo := "-"
		if disk.ServerName != "" {
			serverInfo = disk.ServerName
		}
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %4dGB %-10s %s",
				disk.Name,
				disk.ID,
				disk.SizeGB,
				disk.Connection,
				serverInfo))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %4dGB %-10s %s",
				disk.Name,
				disk.ID,
				disk.SizeGB,
				disk.Connection,
				serverInfo))
		}
	} else if archive, ok := item.(Archive); ok {
		// Handle Archive
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %4dGB %-10s %s",
				archive.Name,
				archive.ID,
				archive.SizeGB,
				archive.Scope,
				archive.Availability))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %4dGB %-10s %s",
				archive.Name,
				archive.ID,
				archive.SizeGB,
				archive.Scope,
				archive.Availability))
		}
	} else if internet, ok := item.(Internet); ok {
		// Handle Internet
		switchID := "-"
		if internet.SwitchID != "" {
			switchID = internet.SwitchID
		}
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %6dMbps %s",
				internet.Name,
				internet.ID,
				internet.BandWidthMbps,
				switchID))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %6dMbps %s",
				internet.Name,
				internet.ID,
				internet.BandWidthMbps,
				switchID))
		}
	} else if vpcRouter, ok := item.(VPCRouter); ok {
		// Handle VPCRouter
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %-10s v%-7d %s",
				vpcRouter.Name,
				vpcRouter.ID,
				vpcRouter.Plan,
				vpcRouter.Version,
				vpcRouter.InstanceStatus))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %-10s v%-7d %s",
				vpcRouter.Name,
				vpcRouter.ID,
				vpcRouter.Plan,
				vpcRouter.Version,
				vpcRouter.InstanceStatus))
		}
	} else if pf, ok := item.(PacketFilter); ok {
		// Handle PacketFilter
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %d rules",
				pf.Name,
				pf.ID,
				pf.RuleCount))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %d rules",
				pf.Name,
				pf.ID,
				pf.RuleCount))
		}
	} else if lb, ok := item.(LoadBalancer); ok {
		// Handle LoadBalancer
		statusStyle := otherStatusStyle
		switch lb.InstanceStatus {
		case "up":
			statusStyle = upStatusStyle
		case "down":
			statusStyle = downStatusStyle
		}
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %d VIPs  %s",
				lb.Name,
				lb.ID,
				lb.VIPCount,
				statusStyle.Render(lb.InstanceStatus)))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %d VIPs  %s",
				lb.Name,
				lb.ID,
				lb.VIPCount,
				statusStyle.Render(lb.InstanceStatus)))
		}
	} else if nfs, ok := item.(NFS); ok {
		// Handle NFS
		statusStyle := otherStatusStyle
		switch nfs.InstanceStatus {
		case "up":
			statusStyle = upStatusStyle
		case "down":
			statusStyle = downStatusStyle
		}
		if index == m.Index() {
			str = selectedItemStyle.Render(fmt.Sprintf("> %-40s %-20s %-15s %s",
				nfs.Name,
				nfs.ID,
				nfs.SwitchName,
				statusStyle.Render(nfs.InstanceStatus)))
		} else {
			str = itemStyle.Render(fmt.Sprintf("  %-40s %-20s %-15s %s",
				nfs.Name,
				nfs.ID,
				nfs.SwitchName,
				statusStyle.Render(nfs.InstanceStatus)))
		}
	} else {
		return
	}

	fmt.Fprint(w, str)
}

type model struct {
	client             *SakuraClient
	list               list.Model
	zones              []string
	currentZone        string
	cursor             int
	err                error
	loading            bool
	quitting           bool
	accountName        string
	windowHeight       int
	windowWidth        int
	searchMode         bool
	searchInput        textinput.Model
	searchQuery        string
	searchMatches      []int // Indices of matching items
	currentMatch       int   // Current match index in searchMatches
	detailMode         bool
	serverDetail       *ServerDetail
	switchDetail       *SwitchDetail
	dnsDetail          *DNSDetail
	elbDetail          *ELBDetail
	gslbDetail         *GSLBDetail
	dbDetail           *DBDetail
	diskDetail         *DiskDetail
	archiveDetail      *ArchiveDetail
	internetDetail     *InternetDetail
	vpcRouterDetail    *VPCRouterDetail
	packetFilterDetail *PacketFilterDetail
	loadBalancerDetail *LoadBalancerDetail
	nfsDetail          *NFSDetail
	detailLoading      bool
	resourceType       ResourceType
	detailViewport     viewport.Model
	// Resource type selector
	resourceSelectMode   bool
	resourceSelectCursor int
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

type disksLoadedMsg struct {
	disks []Disk
	err   error
}

type diskDetailLoadedMsg struct {
	detail *DiskDetail
	err    error
}

type archivesLoadedMsg struct {
	archives []Archive
	err      error
}

type archiveDetailLoadedMsg struct {
	detail *ArchiveDetail
	err    error
}

type internetLoadedMsg struct {
	internet []Internet
	err      error
}

type internetDetailLoadedMsg struct {
	detail *InternetDetail
	err    error
}

type vpcRoutersLoadedMsg struct {
	vpcRouters []VPCRouter
	err        error
}

type vpcRouterDetailLoadedMsg struct {
	detail *VPCRouterDetail
	err    error
}

type packetFiltersLoadedMsg struct {
	packetFilters []PacketFilter
	err           error
}

type packetFilterDetailLoadedMsg struct {
	detail *PacketFilterDetail
	err    error
}

type loadBalancersLoadedMsg struct {
	loadBalancers []LoadBalancer
	err           error
}

type loadBalancerDetailLoadedMsg struct {
	detail *LoadBalancerDetail
	err    error
}

type nfsLoadedMsg struct {
	nfsList []NFS
	err     error
}

type nfsDetailLoadedMsg struct {
	detail *NFSDetail
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

func loadDisks(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		disks, err := client.ListDisks(ctx)
		return disksLoadedMsg{disks: disks, err: err}
	}
}

func loadDiskDetail(client *SakuraClient, diskID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetDiskDetail(ctx, diskID)
		if err != nil {
			slog.Error("Failed to load disk detail", slog.Any("error", err))
			return diskDetailLoadedMsg{err: err}
		}
		slog.Info("Disk detail loaded successfully", slog.String("diskID", diskID))
		return diskDetailLoadedMsg{detail: detail}
	}
}

func loadArchives(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		archives, err := client.ListArchives(ctx)
		return archivesLoadedMsg{archives: archives, err: err}
	}
}

func loadArchiveDetail(client *SakuraClient, archiveID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetArchiveDetail(ctx, archiveID)
		if err != nil {
			slog.Error("Failed to load archive detail", slog.Any("error", err))
			return archiveDetailLoadedMsg{err: err}
		}
		slog.Info("Archive detail loaded successfully", slog.String("archiveID", archiveID))
		return archiveDetailLoadedMsg{detail: detail}
	}
}

func loadInternet(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		internet, err := client.ListInternet(ctx)
		return internetLoadedMsg{internet: internet, err: err}
	}
}

func loadInternetDetail(client *SakuraClient, internetID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetInternetDetail(ctx, internetID)
		if err != nil {
			slog.Error("Failed to load internet detail", slog.Any("error", err))
			return internetDetailLoadedMsg{err: err}
		}
		slog.Info("Internet detail loaded successfully", slog.String("internetID", internetID))
		return internetDetailLoadedMsg{detail: detail}
	}
}

func loadVPCRouters(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		vpcRouters, err := client.ListVPCRouters(ctx)
		return vpcRoutersLoadedMsg{vpcRouters: vpcRouters, err: err}
	}
}

func loadVPCRouterDetail(client *SakuraClient, vpcRouterID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetVPCRouterDetail(ctx, vpcRouterID)
		if err != nil {
			slog.Error("Failed to load VPC router detail", slog.Any("error", err))
			return vpcRouterDetailLoadedMsg{err: err}
		}
		slog.Info("VPC router detail loaded successfully", slog.String("vpcRouterID", vpcRouterID))
		return vpcRouterDetailLoadedMsg{detail: detail}
	}
}

func loadPacketFilters(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		packetFilters, err := client.ListPacketFilters(ctx)
		return packetFiltersLoadedMsg{packetFilters: packetFilters, err: err}
	}
}

func loadPacketFilterDetail(client *SakuraClient, packetFilterID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetPacketFilterDetail(ctx, packetFilterID)
		if err != nil {
			slog.Error("Failed to load packet filter detail", slog.Any("error", err))
			return packetFilterDetailLoadedMsg{err: err}
		}
		slog.Info("Packet filter detail loaded successfully", slog.String("packetFilterID", packetFilterID))
		return packetFilterDetailLoadedMsg{detail: detail}
	}
}

func loadLoadBalancers(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		loadBalancers, err := client.ListLoadBalancers(ctx)
		return loadBalancersLoadedMsg{loadBalancers: loadBalancers, err: err}
	}
}

func loadLoadBalancerDetail(client *SakuraClient, loadBalancerID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetLoadBalancerDetail(ctx, loadBalancerID)
		if err != nil {
			slog.Error("Failed to load load balancer detail", slog.Any("error", err))
			return loadBalancerDetailLoadedMsg{err: err}
		}
		slog.Info("Load balancer detail loaded successfully", slog.String("loadBalancerID", loadBalancerID))
		return loadBalancerDetailLoadedMsg{detail: detail}
	}
}

func loadNFS(client *SakuraClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		nfsList, err := client.ListNFS(ctx)
		return nfsLoadedMsg{nfsList: nfsList, err: err}
	}
}

func loadNFSDetail(client *SakuraClient, nfsID string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		detail, err := client.GetNFSDetail(ctx, nfsID)
		if err != nil {
			slog.Error("Failed to load NFS detail", slog.Any("error", err))
			return nfsDetailLoadedMsg{err: err}
		}
		slog.Info("NFS detail loaded successfully", slog.String("nfsID", nfsID))
		return nfsDetailLoadedMsg{detail: detail}
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
	resourceList.SetShowTitle(false)
	resourceList.SetShowStatusBar(false)
	resourceList.SetFilteringEnabled(false) // Disable built-in filtering

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
			continue
		}
		// Handle ELB
		if elb, ok := item.(ELB); ok {
			if strings.Contains(strings.ToLower(elb.Name), query) ||
				strings.Contains(strings.ToLower(elb.ID), query) ||
				strings.Contains(strings.ToLower(elb.Desc), query) ||
				strings.Contains(strings.ToLower(elb.VIP), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
			continue
		}
		// Handle GSLB
		if gslb, ok := item.(GSLB); ok {
			if strings.Contains(strings.ToLower(gslb.Name), query) ||
				strings.Contains(strings.ToLower(gslb.ID), query) ||
				strings.Contains(strings.ToLower(gslb.Desc), query) ||
				strings.Contains(strings.ToLower(gslb.FQDN), query) {
				m.searchMatches = append(m.searchMatches, i)
			}
			continue
		}
		// Handle DB
		if db, ok := item.(DB); ok {
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

func (m model) getTableHeader() string {
	switch m.resourceType {
	case ResourceTypeServer:
		return fmt.Sprintf("  %-40s %-20s %s", "Name", "ID", "Status")
	case ResourceTypeSwitch:
		return fmt.Sprintf("  %-40s %-20s %s", "Name", "ID", "Servers")
	case ResourceTypeDNS:
		return fmt.Sprintf("  %-40s %-20s %s", "Name", "ID", "Records")
	case ResourceTypeELB:
		return fmt.Sprintf("  %-40s %-20s %-15s %s", "Name", "ID", "VIP", "Servers")
	case ResourceTypeGSLB:
		return fmt.Sprintf("  %-40s %-20s %s", "Name", "ID", "FQDN")
	case ResourceTypeDB:
		return fmt.Sprintf("  %-40s %-20s %-10s %s", "Name", "ID", "Type", "Status")
	case ResourceTypeDisk:
		return fmt.Sprintf("  %-40s %-20s %6s %-10s %s", "Name", "ID", "Size", "Connection", "Server")
	case ResourceTypeArchive:
		return fmt.Sprintf("  %-40s %-20s %6s %-10s %s", "Name", "ID", "Size", "Scope", "Availability")
	case ResourceTypeInternet:
		return fmt.Sprintf("  %-40s %-20s %10s %s", "Name", "ID", "Bandwidth", "Switch ID")
	case ResourceTypeVPCRouter:
		return fmt.Sprintf("  %-40s %-20s %-10s %-8s %s", "Name", "ID", "Plan", "Version", "Status")
	case ResourceTypePacketFilter:
		return fmt.Sprintf("  %-40s %-20s %s", "Name", "ID", "Rules")
	case ResourceTypeLoadBalancer:
		return fmt.Sprintf("  %-40s %-20s %s  %s", "Name", "ID", "VIPs", "Status")
	case ResourceTypeNFS:
		return fmt.Sprintf("  %-40s %-20s %-15s %s", "Name", "ID", "Switch", "Status")
	default:
		return ""
	}
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
					case ResourceTypeDisk:
						return m, loadDisks(m.client)
					case ResourceTypeArchive:
						return m, loadArchives(m.client)
					case ResourceTypeInternet:
						return m, loadInternet(m.client)
					case ResourceTypeVPCRouter:
						return m, loadVPCRouters(m.client)
					case ResourceTypePacketFilter:
						return m, loadPacketFilters(m.client)
					case ResourceTypeLoadBalancer:
						return m, loadLoadBalancers(m.client)
					case ResourceTypeNFS:
						return m, loadNFS(m.client)
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
				if elb, ok := selectedItem.(ELB); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadELBDetail(m.client, elb.ID)
				}
				if gslb, ok := selectedItem.(GSLB); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadGSLBDetail(m.client, gslb.ID)
				}
				if db, ok := selectedItem.(DB); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadDBDetail(m.client, db.ID)
				}
				if disk, ok := selectedItem.(Disk); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadDiskDetail(m.client, disk.ID)
				}
				if archive, ok := selectedItem.(Archive); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadArchiveDetail(m.client, archive.ID)
				}
				if internet, ok := selectedItem.(Internet); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadInternetDetail(m.client, internet.ID)
				}
				if vpcRouter, ok := selectedItem.(VPCRouter); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadVPCRouterDetail(m.client, vpcRouter.ID)
				}
				if pf, ok := selectedItem.(PacketFilter); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadPacketFilterDetail(m.client, pf.ID)
				}
				if lb, ok := selectedItem.(LoadBalancer); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadLoadBalancerDetail(m.client, lb.ID)
				}
				if nfs, ok := selectedItem.(NFS); ok {
					m.detailMode = true
					m.detailLoading = true
					return m, loadNFSDetail(m.client, nfs.ID)
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
			case ResourceTypeDisk:
				return m, loadDisks(m.client)
			case ResourceTypeArchive:
				return m, loadArchives(m.client)
			case ResourceTypeInternet:
				return m, loadInternet(m.client)
			case ResourceTypeVPCRouter:
				return m, loadVPCRouters(m.client)
			case ResourceTypePacketFilter:
				return m, loadPacketFilters(m.client)
			case ResourceTypeLoadBalancer:
				return m, loadLoadBalancers(m.client)
			case ResourceTypeNFS:
				return m, loadNFS(m.client)
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
			case ResourceTypeDisk:
				return m, loadDisks(m.client)
			case ResourceTypeArchive:
				return m, loadArchives(m.client)
			case ResourceTypeInternet:
				return m, loadInternet(m.client)
			case ResourceTypeVPCRouter:
				return m, loadVPCRouters(m.client)
			case ResourceTypePacketFilter:
				return m, loadPacketFilters(m.client)
			case ResourceTypeLoadBalancer:
				return m, loadLoadBalancers(m.client)
			case ResourceTypeNFS:
				return m, loadNFS(m.client)
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

	case elbLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load ELBs", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("ELBs loaded successfully", slog.Int("count", len(msg.elbList)))

		// Convert ELBs to list items
		items := make([]list.Item, len(msg.elbList))
		for i, elb := range msg.elbList {
			items[i] = elb
		}
		m.list.SetItems(items)
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

		// Convert GSLBs to list items
		items := make([]list.Item, len(msg.gslbList))
		for i, gslb := range msg.gslbList {
			items[i] = gslb
		}
		m.list.SetItems(items)
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

		// Convert DBs to list items
		items := make([]list.Item, len(msg.dbList))
		for i, db := range msg.dbList {
			items[i] = db
		}
		m.list.SetItems(items)
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

	case disksLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load disks", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Disks loaded successfully", slog.Int("count", len(msg.disks)))

		// Convert disks to list items
		items := make([]list.Item, len(msg.disks))
		for i, disk := range msg.disks {
			items[i] = disk
		}
		m.list.SetItems(items)
		return m, nil

	case diskDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load disk detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.diskDetail = msg.detail
		// Setup viewport for detail view
		content := renderDiskDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case archivesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load archives", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Archives loaded successfully", slog.Int("count", len(msg.archives)))

		// Convert archives to list items
		items := make([]list.Item, len(msg.archives))
		for i, archive := range msg.archives {
			items[i] = archive
		}
		m.list.SetItems(items)
		return m, nil

	case archiveDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load archive detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.archiveDetail = msg.detail
		// Setup viewport for detail view
		content := renderArchiveDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case internetLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load internet", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Internet loaded successfully", slog.Int("count", len(msg.internet)))

		// Convert internet to list items
		items := make([]list.Item, len(msg.internet))
		for i, internet := range msg.internet {
			items[i] = internet
		}
		m.list.SetItems(items)
		return m, nil

	case internetDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load internet detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.internetDetail = msg.detail
		// Setup viewport for detail view
		content := renderInternetDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case vpcRoutersLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load VPC routers", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("VPC routers loaded successfully", slog.Int("count", len(msg.vpcRouters)))

		// Convert VPC routers to list items
		items := make([]list.Item, len(msg.vpcRouters))
		for i, vpcRouter := range msg.vpcRouters {
			items[i] = vpcRouter
		}
		m.list.SetItems(items)
		return m, nil

	case vpcRouterDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load VPC router detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.vpcRouterDetail = msg.detail
		// Setup viewport for detail view
		content := renderVPCRouterDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case packetFiltersLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load packet filters", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Packet filters loaded successfully", slog.Int("count", len(msg.packetFilters)))

		// Convert packet filters to list items
		items := make([]list.Item, len(msg.packetFilters))
		for i, pf := range msg.packetFilters {
			items[i] = pf
		}
		m.list.SetItems(items)
		return m, nil

	case packetFilterDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load packet filter detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.packetFilterDetail = msg.detail
		// Setup viewport for detail view
		content := renderPacketFilterDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case loadBalancersLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load load balancers", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("Load balancers loaded successfully", slog.Int("count", len(msg.loadBalancers)))

		// Convert load balancers to list items
		items := make([]list.Item, len(msg.loadBalancers))
		for i, lb := range msg.loadBalancers {
			items[i] = lb
		}
		m.list.SetItems(items)
		return m, nil

	case loadBalancerDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load load balancer detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.loadBalancerDetail = msg.detail
		// Setup viewport for detail view
		content := renderLoadBalancerDetail(msg.detail)
		m.detailViewport = viewport.New(m.windowWidth, m.windowHeight-10)
		m.detailViewport.SetContent(content)
		return m, nil

	case nfsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			slog.Error("Failed to load NFS appliances", slog.Any("error", msg.err))
			m.err = msg.err
			return m, nil
		}
		slog.Info("NFS appliances loaded successfully", slog.Int("count", len(msg.nfsList)))

		// Convert NFS to list items
		items := make([]list.Item, len(msg.nfsList))
		for i, nfs := range msg.nfsList {
			items[i] = nfs
		}
		m.list.SetItems(items)
		return m, nil

	case nfsDetailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			slog.Error("Failed to load NFS detail", slog.Any("error", msg.err))
			m.err = msg.err
			m.detailMode = false
			return m, nil
		}
		m.nfsDetail = msg.detail
		// Setup viewport for detail view
		content := renderNFSDetail(msg.detail)
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
		} else if m.serverDetail != nil || m.switchDetail != nil || m.dnsDetail != nil || m.elbDetail != nil || m.gslbDetail != nil || m.dbDetail != nil || m.diskDetail != nil || m.archiveDetail != nil || m.internetDetail != nil || m.vpcRouterDetail != nil || m.packetFilterDetail != nil || m.loadBalancerDetail != nil || m.nfsDetail != nil {
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
				b.WriteString(selectedItemStyle.Render(fmt.Sprintf("> %s", rt.String())))
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

	// Resource list or loading/error
	if m.loading {
		b.WriteString(fmt.Sprintf("Loading %s...\n", strings.ToLower(m.resourceType.String())))
	} else if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else {
		// Render table header based on resource type
		header := m.getTableHeader()
		if header != "" {
			b.WriteString(zoneStyle.Render(header))
			b.WriteString("\n")
		}
		b.WriteString(m.list.View())
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Enter: details | /: search | n/N: next/prev | t: type | z: zone | r: refresh | q: quit"))
	}

	return b.String()
}
