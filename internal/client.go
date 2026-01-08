package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// ResourceType represents different cloud resource types
type ResourceType int

const (
	ResourceTypeServer ResourceType = iota
	ResourceTypeSwitch
	ResourceTypeDNS
	// Future: ResourceTypeDB, ResourceTypeELB, ResourceTypeAppRun
)

func (r ResourceType) String() string {
	switch r {
	case ResourceTypeServer:
		return "Server"
	case ResourceTypeSwitch:
		return "Switch"
	case ResourceTypeDNS:
		return "DNS"
	default:
		return "Unknown"
	}
}

type SakuraClient struct {
	caller iaas.APICaller
	zone   string
}

func NewSakuraClient(zone string) (*SakuraClient, error) {
	if zone == "" {
		slog.Error("Zone is empty")
		return nil, fmt.Errorf("zone must be specified")
	}

	slog.Info("Creating Sakura Cloud API caller", slog.String("zone", zone))

	caller, err := api.NewCaller()
	if err != nil {
		return nil, fmt.Errorf("failed to create Sakura Cloud API caller: %w", err)
	}

	slog.Info("Sakura Cloud API caller created successfully!", slog.String("zone", zone))

	return &SakuraClient{
		caller: caller,
		zone:   zone,
	}, nil
}

type Server struct {
	ID             string
	Name           string
	InstanceStatus string
	Zone           string
}

type ServerDetail struct {
	Server
	Description     string
	Tags            []string
	CPU             int
	MemoryGB        int
	InterfaceCount  int
	Disks           []DiskInfo
	IPAddresses     []string
	UserIPAddresses []string
	CreatedAt       string
}

type DiskInfo struct {
	Name   string
	SizeGB int
}

// Switch represents a switch resource
type Switch struct {
	ID           string
	Name         string
	Desc         string
	Zone         string
	DefaultRoute string
	ServerCount  int
	CreatedAt    string
}

type SwitchDetail struct {
	Switch
	Tags        []string
	SubnetCount int
	ServerCount int
	CreatedAt   string
}

// DNS represents a DNS zone resource
type DNS struct {
	ID          string
	Name        string
	Desc        string
	Zone        string
	RecordCount int
	CreatedAt   string
}

type DNSDetail struct {
	DNS
	Tags         []string
	RecordCount  int
	NameServers  []string
	IconID       string
	CreatedAt    string
	ModifiedAt   string
}

// Implement list.Item interface for Server
func (s Server) FilterValue() string {
	return s.Name
}

func (s Server) Title() string {
	return s.Name
}

func (s Server) Description() string {
	return fmt.Sprintf("ID: %s | Status: %s", s.ID, s.InstanceStatus)
}

// Implement list.Item interface for Switch
func (s Switch) FilterValue() string {
	return s.Name
}

func (s Switch) Title() string {
	return s.Name
}

func (s Switch) Description() string {
	desc := fmt.Sprintf("ID: %s", s.ID)
	if s.Desc != "" {
		desc += " | " + s.Desc
	}
	return desc
}

// Implement list.Item interface for DNS
func (d DNS) FilterValue() string {
	return d.Name
}

func (d DNS) Title() string {
	return d.Name
}

func (d DNS) Description() string {
	desc := fmt.Sprintf("ID: %s", d.ID)
	if d.Desc != "" {
		desc += " | " + d.Desc
	}
	return desc
}

func (c *SakuraClient) GetAuthStatus(ctx context.Context) (*iaas.AuthStatus, error) {
	op := iaas.NewAuthStatusOp(c.caller)
	return op.Read(ctx)
}

func (c *SakuraClient) ListServers(ctx context.Context) ([]Server, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching servers from Sakura Cloud",
		slog.String("zone", c.zone))

	serverOp := iaas.NewServerOp(c.caller)

	searched, err := serverOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch servers",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	servers := make([]Server, 0, len(searched.Servers))
	for _, s := range searched.Servers {
		status := "UNKNOWN"
		if s.InstanceStatus != "" {
			status = string(s.InstanceStatus)
		}
		servers = append(servers, Server{
			ID:             s.ID.String(),
			Name:           s.Name,
			InstanceStatus: status,
			Zone:           c.zone,
		})
	}

	slog.Info("Successfully fetched servers",
		slog.String("zone", c.zone),
		slog.Int("count", len(servers)))

	return servers, nil
}

func (c *SakuraClient) SetZone(zone string) {
	slog.Info("Switching zone",
		slog.String("from", c.zone),
		slog.String("to", zone))
	c.zone = zone
}

func (c *SakuraClient) GetZone() string {
	return c.zone
}

func (c *SakuraClient) GetServerDetail(ctx context.Context, serverID string) (*ServerDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching server detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("serverID", serverID))

	serverOp := iaas.NewServerOp(c.caller)

	// Parse server ID
	id := types.StringID(serverID)

	server, err := serverOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch server detail",
			slog.String("zone", c.zone),
			slog.String("serverID", serverID),
			slog.Any("error", err))
		return nil, err
	}

	detail := &ServerDetail{
		Server: Server{
			ID:             server.ID.String(),
			Name:           server.Name,
			InstanceStatus: string(server.InstanceStatus),
			Zone:           c.zone,
		},
		Description: server.Description,
		Tags:        server.Tags,
		CPU:         server.CPU,
		MemoryGB:    server.GetMemoryGB(),
	}

	// Get IP addresses
	if len(server.Interfaces) > 0 {
		detail.InterfaceCount = len(server.Interfaces)
		for _, iface := range server.Interfaces {
			if iface.IPAddress != "" {
				detail.IPAddresses = append(detail.IPAddresses, iface.IPAddress)
			}
			if iface.UserIPAddress != "" {
				detail.UserIPAddresses = append(detail.UserIPAddresses, iface.UserIPAddress)
			}
		}
	}

	// Get disk info
	for _, disk := range server.Disks {
		if disk != nil {
			detail.Disks = append(detail.Disks, DiskInfo{
				Name:   disk.Name,
				SizeGB: disk.GetSizeGB(),
			})
		}
	}

	if !server.CreatedAt.IsZero() {
		detail.CreatedAt = server.CreatedAt.Format("2006-01-02 15:04:05")
	}

	slog.Info("Successfully fetched server detail",
		slog.String("zone", c.zone),
		slog.String("serverID", serverID))

	return detail, nil
}

func (c *SakuraClient) ListSwitches(ctx context.Context) ([]Switch, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching switches from Sakura Cloud",
		slog.String("zone", c.zone))

	switchOp := iaas.NewSwitchOp(c.caller)

	searched, err := switchOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch switches",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	// Get server list once for counting connected servers
	serverOp := iaas.NewServerOp(c.caller)
	servers, err := serverOp.Find(ctx, c.zone, &iaas.FindCondition{})
	var serverList []*iaas.Server
	if err == nil {
		serverList = servers.Servers
	}

	switches := make([]Switch, 0, len(searched.Switches))
	for _, sw := range searched.Switches {
		// Count connected servers
		serverCount := 0
		for _, server := range serverList {
			for _, iface := range server.Interfaces {
				if iface.SwitchID == sw.ID {
					serverCount++
					break
				}
			}
		}

		// Get default route
		defaultRoute := ""
		if len(sw.Subnets) > 0 && sw.Subnets[0].DefaultRoute != "" {
			defaultRoute = sw.Subnets[0].DefaultRoute
		}

		// Format created at
		createdAt := ""
		if !sw.CreatedAt.IsZero() {
			createdAt = sw.CreatedAt.Format("2006-01-02")
		}

		switches = append(switches, Switch{
			ID:           sw.ID.String(),
			Name:         sw.Name,
			Desc:         sw.Description,
			Zone:         c.zone,
			DefaultRoute: defaultRoute,
			ServerCount:  serverCount,
			CreatedAt:    createdAt,
		})
	}

	slog.Info("Successfully fetched switches",
		slog.String("zone", c.zone),
		slog.Int("count", len(switches)))

	return switches, nil
}

func (c *SakuraClient) GetSwitchDetail(ctx context.Context, switchID string) (*SwitchDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching switch detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("switchID", switchID))

	switchOp := iaas.NewSwitchOp(c.caller)

	id := types.StringID(switchID)

	sw, err := switchOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch switch detail",
			slog.String("zone", c.zone),
			slog.String("switchID", switchID),
			slog.Any("error", err))
		return nil, err
	}

	// Get default route
	defaultRoute := ""
	if len(sw.Subnets) > 0 && sw.Subnets[0].DefaultRoute != "" {
		defaultRoute = sw.Subnets[0].DefaultRoute
	}

	// Format created at
	createdAt := ""
	if !sw.CreatedAt.IsZero() {
		createdAt = sw.CreatedAt.Format("2006-01-02")
	}

	// Count connected servers
	serverOp := iaas.NewServerOp(c.caller)
	servers, err := serverOp.Find(ctx, c.zone, &iaas.FindCondition{})
	serverCount := 0
	if err == nil {
		for _, server := range servers.Servers {
			for _, iface := range server.Interfaces {
				if iface.SwitchID == sw.ID {
					serverCount++
					break
				}
			}
		}
	}

	detail := &SwitchDetail{
		Switch: Switch{
			ID:           sw.ID.String(),
			Name:         sw.Name,
			Desc:         sw.Description,
			Zone:         c.zone,
			DefaultRoute: defaultRoute,
			ServerCount:  serverCount,
			CreatedAt:    createdAt,
		},
		Tags:        sw.Tags,
		SubnetCount: len(sw.Subnets),
		ServerCount: serverCount,
	}

	if !sw.CreatedAt.IsZero() {
		detail.CreatedAt = sw.CreatedAt.Format("2006-01-02 15:04:05")
	}

	slog.Info("Successfully fetched switch detail",
		slog.String("zone", c.zone),
		slog.String("switchID", switchID))

	return detail, nil
}

func (c *SakuraClient) ListDNS(ctx context.Context) ([]DNS, error) {
	slog.Info("Fetching DNS zones from Sakura Cloud")

	dnsOp := iaas.NewDNSOp(c.caller)

	searched, err := dnsOp.Find(ctx, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch DNS zones",
			slog.Any("error", err))
		return nil, err
	}

	dnsList := make([]DNS, 0, len(searched.DNS))
	for _, d := range searched.DNS {
		// Count records
		recordCount := len(d.Records)

		// Format created at
		createdAt := ""
		if !d.CreatedAt.IsZero() {
			createdAt = d.CreatedAt.Format("2006-01-02")
		}

		dnsList = append(dnsList, DNS{
			ID:          d.ID.String(),
			Name:        d.Name,
			Desc:        d.Description,
			Zone:        "global", // DNS is a global resource
			RecordCount: recordCount,
			CreatedAt:   createdAt,
		})
	}

	slog.Info("Successfully fetched DNS zones",
		slog.Int("count", len(dnsList)))

	return dnsList, nil
}

func (c *SakuraClient) GetDNSDetail(ctx context.Context, dnsID string) (*DNSDetail, error) {
	slog.Info("Fetching DNS detail from Sakura Cloud",
		slog.String("dnsID", dnsID))

	dnsOp := iaas.NewDNSOp(c.caller)

	id := types.StringID(dnsID)

	d, err := dnsOp.Read(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch DNS detail",
			slog.String("dnsID", dnsID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !d.CreatedAt.IsZero() {
		createdAt = d.CreatedAt.Format("2006-01-02")
	}

	// Format modified at
	modifiedAt := ""
	if !d.ModifiedAt.IsZero() {
		modifiedAt = d.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Get name servers
	nameServers := []string{}
	for _, ns := range d.DNSNameServers {
		nameServers = append(nameServers, ns)
	}

	detail := &DNSDetail{
		DNS: DNS{
			ID:          d.ID.String(),
			Name:        d.Name,
			Desc:        d.Description,
			Zone:        "global",
			RecordCount: len(d.Records),
			CreatedAt:   createdAt,
		},
		Tags:        d.Tags,
		RecordCount: len(d.Records),
		NameServers: nameServers,
		IconID:      d.IconID.String(),
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
	}

	slog.Info("Successfully fetched DNS detail",
		slog.String("dnsID", dnsID))

	return detail, nil
}
