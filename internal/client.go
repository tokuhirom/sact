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
	ResourceTypeELB
	ResourceTypeGSLB
	ResourceTypeDB
	ResourceTypeDisk
)

// AllResourceTypes returns all available resource types
var AllResourceTypes = []ResourceType{
	ResourceTypeServer,
	ResourceTypeSwitch,
	ResourceTypeDNS,
	ResourceTypeELB,
	ResourceTypeGSLB,
	ResourceTypeDB,
	ResourceTypeDisk,
}

func (r ResourceType) String() string {
	switch r {
	case ResourceTypeServer:
		return "Server"
	case ResourceTypeSwitch:
		return "Switch"
	case ResourceTypeDNS:
		return "DNS"
	case ResourceTypeELB:
		return "ELB"
	case ResourceTypeGSLB:
		return "GSLB"
	case ResourceTypeDB:
		return "DB"
	case ResourceTypeDisk:
		return "Disk"
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

type DNSRecord struct {
	Name  string
	Type  string
	RData string
	TTL   int
}

type DNSDetail struct {
	DNS
	Tags        []string
	RecordCount int
	Records     []DNSRecord
	NameServers []string
	IconID      string
	CreatedAt   string
	ModifiedAt  string
}

// ELB represents an Enhanced Load Balancer resource
type ELB struct {
	ID          string
	Name        string
	Desc        string
	Zone        string
	VIP         string
	Plan        string
	ServerCount int
	CreatedAt   string
}

type ELBServer struct {
	IPAddress string
	Port      int
	Enabled   bool
}

type ELBDetail struct {
	ELB
	Tags      []string
	Servers   []ELBServer
	FQDN      string
	CreatedAt string
}

// GSLB represents a Global Server Load Balancer resource
type GSLB struct {
	ID          string
	Name        string
	Desc        string
	FQDN        string
	ServerCount int
	CreatedAt   string
}

type GSLBServer struct {
	IPAddress string
	Enabled   bool
	Weight    int
}

type GSLBDetail struct {
	GSLB
	Tags       []string
	Servers    []GSLBServer
	HealthPath string
	DelayLoop  int
	Weighted   bool
	CreatedAt  string
}

// DB represents a Database Appliance resource
type DB struct {
	ID             string
	Name           string
	Desc           string
	Zone           string
	DBType         string
	Plan           string
	InstanceStatus string
	CreatedAt      string
}

type DBDetail struct {
	DB
	Tags           []string
	CPU            int
	MemoryGB       int
	DiskSizeGB     int
	IPAddress      string
	Port           int
	DefaultUser    string
	NetworkMaskLen int
	DefaultRoute   string
	CreatedAt      string
}

// Disk represents a disk resource
type Disk struct {
	ID         string
	Name       string
	Desc       string
	Zone       string
	SizeGB     int
	Connection string
	ServerID   string
	ServerName string
	CreatedAt  string
}

type DiskDetail struct {
	Disk
	Tags             []string
	DiskPlanName     string
	SourceDiskID     string
	SourceArchiveID  string
	Availability     string
	EncryptionAlgo   string
	CreatedAt        string
	ModifiedAt       string
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

// Implement list.Item interface for ELB
func (e ELB) FilterValue() string {
	return e.Name
}

func (e ELB) Title() string {
	return e.Name
}

func (e ELB) Description() string {
	desc := fmt.Sprintf("ID: %s", e.ID)
	if e.Desc != "" {
		desc += " | " + e.Desc
	}
	return desc
}

// Implement list.Item interface for GSLB
func (g GSLB) FilterValue() string {
	return g.Name
}

func (g GSLB) Title() string {
	return g.Name
}

func (g GSLB) Description() string {
	desc := fmt.Sprintf("ID: %s", g.ID)
	if g.Desc != "" {
		desc += " | " + g.Desc
	}
	return desc
}

// Implement list.Item interface for DB
func (d DB) FilterValue() string {
	return d.Name
}

func (d DB) Title() string {
	return d.Name
}

func (d DB) Description() string {
	desc := fmt.Sprintf("ID: %s", d.ID)
	if d.Desc != "" {
		desc += " | " + d.Desc
	}
	return desc
}

// Implement list.Item interface for Disk
func (d Disk) FilterValue() string {
	return d.Name
}

func (d Disk) Title() string {
	return d.Name
}

func (d Disk) Description() string {
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

	// Get DNS records
	records := []DNSRecord{}
	for _, rec := range d.Records {
		records = append(records, DNSRecord{
			Name:  rec.Name,
			Type:  string(rec.Type),
			RData: rec.RData,
			TTL:   rec.TTL,
		})
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
		Records:     records,
		NameServers: nameServers,
		IconID:      d.IconID.String(),
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
	}

	slog.Info("Successfully fetched DNS detail",
		slog.String("dnsID", dnsID))

	return detail, nil
}

func (c *SakuraClient) ListELB(ctx context.Context) ([]ELB, error) {
	slog.Info("Fetching ELBs from Sakura Cloud")

	elbOp := iaas.NewProxyLBOp(c.caller)

	searched, err := elbOp.Find(ctx, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch ELBs",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	elbList := make([]ELB, 0, len(searched.ProxyLBs))
	for _, elb := range searched.ProxyLBs {
		// Get VIP
		vip := ""
		if elb.VirtualIPAddress != "" {
			vip = elb.VirtualIPAddress
		}

		// Get plan
		plan := elb.Plan.String()

		// Count servers
		serverCount := len(elb.Servers)

		// Format created at
		createdAt := ""
		if !elb.CreatedAt.IsZero() {
			createdAt = elb.CreatedAt.Format("2006-01-02")
		}

		elbList = append(elbList, ELB{
			ID:          elb.ID.String(),
			Name:        elb.Name,
			Desc:        elb.Description,
			Zone:        "global", // ELB is a global resource
			VIP:         vip,
			Plan:        plan,
			ServerCount: serverCount,
			CreatedAt:   createdAt,
		})
	}

	slog.Info("Successfully fetched ELBs",
		slog.String("zone", c.zone),
		slog.Int("count", len(elbList)))

	return elbList, nil
}

func (c *SakuraClient) GetELBDetail(ctx context.Context, elbID string) (*ELBDetail, error) {
	slog.Info("Fetching ELB detail from Sakura Cloud",
		slog.String("elbID", elbID))

	elbOp := iaas.NewProxyLBOp(c.caller)

	id := types.StringID(elbID)

	elb, err := elbOp.Read(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch ELB detail",
			slog.String("elbID", elbID),
			slog.Any("error", err))
		return nil, err
	}

	// Get VIP
	vip := ""
	if elb.VirtualIPAddress != "" {
		vip = elb.VirtualIPAddress
	}

	// Get plan
	plan := elb.Plan.String()

	// Format created at
	createdAt := ""
	if !elb.CreatedAt.IsZero() {
		createdAt = elb.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Get servers
	servers := []ELBServer{}
	for _, srv := range elb.Servers {
		servers = append(servers, ELBServer{
			IPAddress: srv.IPAddress,
			Port:      srv.Port,
			Enabled:   srv.Enabled,
		})
	}

	detail := &ELBDetail{
		ELB: ELB{
			ID:          elb.ID.String(),
			Name:        elb.Name,
			Desc:        elb.Description,
			Zone:        "global",
			VIP:         vip,
			Plan:        plan,
			ServerCount: len(elb.Servers),
			CreatedAt:   createdAt,
		},
		Tags:      elb.Tags,
		Servers:   servers,
		FQDN:      elb.FQDN,
		CreatedAt: createdAt,
	}

	slog.Info("Successfully fetched ELB detail",
		slog.String("elbID", elbID))

	return detail, nil
}

func (c *SakuraClient) ListGSLB(ctx context.Context) ([]GSLB, error) {
	slog.Info("Fetching GSLBs from Sakura Cloud")

	gslbOp := iaas.NewGSLBOp(c.caller)

	searched, err := gslbOp.Find(ctx, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch GSLBs",
			slog.Any("error", err))
		return nil, err
	}

	gslbList := make([]GSLB, 0, len(searched.GSLBs))
	for _, gslb := range searched.GSLBs {
		// Get FQDN
		fqdn := gslb.FQDN

		// Count servers
		serverCount := len(gslb.DestinationServers)

		// Format created at
		createdAt := ""
		if !gslb.CreatedAt.IsZero() {
			createdAt = gslb.CreatedAt.Format("2006-01-02")
		}

		gslbList = append(gslbList, GSLB{
			ID:          gslb.ID.String(),
			Name:        gslb.Name,
			Desc:        gslb.Description,
			FQDN:        fqdn,
			ServerCount: serverCount,
			CreatedAt:   createdAt,
		})
	}

	slog.Info("Successfully fetched GSLBs",
		slog.Int("count", len(gslbList)))

	return gslbList, nil
}

func (c *SakuraClient) GetGSLBDetail(ctx context.Context, gslbID string) (*GSLBDetail, error) {
	slog.Info("Fetching GSLB detail from Sakura Cloud",
		slog.String("gslbID", gslbID))

	gslbOp := iaas.NewGSLBOp(c.caller)

	id := types.StringID(gslbID)

	gslb, err := gslbOp.Read(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch GSLB detail",
			slog.String("gslbID", gslbID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !gslb.CreatedAt.IsZero() {
		createdAt = gslb.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Get servers
	servers := []GSLBServer{}
	for _, srv := range gslb.DestinationServers {
		servers = append(servers, GSLBServer{
			IPAddress: srv.IPAddress,
			Enabled:   srv.Enabled.Bool(),
			Weight:    int(srv.Weight),
		})
	}

	// Get health check path
	healthPath := ""
	if gslb.HealthCheck.Path != "" {
		healthPath = gslb.HealthCheck.Path
	}

	// Get delay loop
	delayLoop := gslb.DelayLoop

	// Check if weighted
	weighted := gslb.Weighted.Bool()

	detail := &GSLBDetail{
		GSLB: GSLB{
			ID:          gslb.ID.String(),
			Name:        gslb.Name,
			Desc:        gslb.Description,
			FQDN:        gslb.FQDN,
			ServerCount: len(gslb.DestinationServers),
			CreatedAt:   createdAt,
		},
		Tags:       gslb.Tags,
		Servers:    servers,
		HealthPath: healthPath,
		DelayLoop:  delayLoop,
		Weighted:   weighted,
		CreatedAt:  createdAt,
	}

	slog.Info("Successfully fetched GSLB detail",
		slog.String("gslbID", gslbID))

	return detail, nil
}

func (c *SakuraClient) ListDB(ctx context.Context) ([]DB, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching DBs from Sakura Cloud",
		slog.String("zone", c.zone))

	dbOp := iaas.NewDatabaseOp(c.caller)

	searched, err := dbOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch DBs",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	dbList := make([]DB, 0, len(searched.Databases))
	for _, db := range searched.Databases {
		// Get DB type
		dbType := "Unknown"
		if db.Conf != nil && db.Conf.DatabaseName != "" {
			dbType = db.Conf.DatabaseName
		}

		// Get plan from PlanID
		plan := db.PlanID.String()

		// Format created at
		createdAt := ""
		if !db.CreatedAt.IsZero() {
			createdAt = db.CreatedAt.Format("2006-01-02")
		}

		// Get instance status
		instanceStatus := string(db.InstanceStatus)

		dbList = append(dbList, DB{
			ID:             db.ID.String(),
			Name:           db.Name,
			Desc:           db.Description,
			Zone:           c.zone,
			DBType:         dbType,
			Plan:           plan,
			InstanceStatus: instanceStatus,
			CreatedAt:      createdAt,
		})
	}

	slog.Info("Successfully fetched DBs",
		slog.String("zone", c.zone),
		slog.Int("count", len(dbList)))

	return dbList, nil
}

func (c *SakuraClient) GetDBDetail(ctx context.Context, dbID string) (*DBDetail, error) {
	slog.Info("Fetching DB detail from Sakura Cloud",
		slog.String("dbID", dbID),
		slog.String("zone", c.zone))

	dbOp := iaas.NewDatabaseOp(c.caller)

	id := types.StringID(dbID)

	db, err := dbOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch DB detail",
			slog.String("dbID", dbID),
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	// Get DB type
	dbType := "Unknown"
	if db.Conf != nil && db.Conf.DatabaseName != "" {
		dbType = db.Conf.DatabaseName
	}

	// Get plan from PlanID
	plan := db.PlanID.String()

	// Format created at
	createdAt := ""
	if !db.CreatedAt.IsZero() {
		createdAt = db.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Get instance status
	instanceStatus := string(db.InstanceStatus)

	// Parse plan - simplified
	cpu := 0
	memoryGB := 0
	diskSizeGB := 0

	// Try to extract from plan name
	planStr := plan
	if planStr == "10" {
		cpu = 1
		memoryGB = 2
		diskSizeGB = 20
	} else if planStr == "30" {
		cpu = 2
		memoryGB = 4
		diskSizeGB = 30
	} else if planStr == "50" {
		cpu = 4
		memoryGB = 8
		diskSizeGB = 50
	} else {
		cpu = 1
		memoryGB = 2
		diskSizeGB = 20
	}

	// Get network info
	ipAddress := ""
	port := 0
	defaultUser := ""
	networkMaskLen := 0
	defaultRoute := ""

	if len(db.Interfaces) > 0 {
		ipAddress = db.Interfaces[0].IPAddress
		networkMaskLen = 24 // Default
		defaultRoute = ""
	}

	if db.Conf != nil {
		if db.Conf.DatabaseName == "postgres" {
			port = 5432
		} else {
			port = 3306
		}
		if db.Conf.DefaultUser != "" {
			defaultUser = db.Conf.DefaultUser
		}
	}

	detail := &DBDetail{
		DB: DB{
			ID:             db.ID.String(),
			Name:           db.Name,
			Desc:           db.Description,
			Zone:           c.zone,
			DBType:         dbType,
			Plan:           plan,
			InstanceStatus: instanceStatus,
			CreatedAt:      createdAt,
		},
		Tags:           db.Tags,
		CPU:            cpu,
		MemoryGB:       memoryGB,
		DiskSizeGB:     diskSizeGB,
		IPAddress:      ipAddress,
		Port:           port,
		DefaultUser:    defaultUser,
		NetworkMaskLen: networkMaskLen,
		DefaultRoute:   defaultRoute,
		CreatedAt:      createdAt,
	}

	slog.Info("Successfully fetched DB detail",
		slog.String("dbID", dbID))

	return detail, nil
}

func (c *SakuraClient) ListDisks(ctx context.Context) ([]Disk, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching disks from Sakura Cloud",
		slog.String("zone", c.zone))

	diskOp := iaas.NewDiskOp(c.caller)

	searched, err := diskOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch disks",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	disks := make([]Disk, 0, len(searched.Disks))
	for _, d := range searched.Disks {
		// Get connection type
		connection := string(d.Connection)

		// Get server info
		serverID := ""
		serverName := ""
		if !d.ServerID.IsEmpty() {
			serverID = d.ServerID.String()
			serverName = d.ServerName
		}

		// Format created at
		createdAt := ""
		if !d.CreatedAt.IsZero() {
			createdAt = d.CreatedAt.Format("2006-01-02")
		}

		disks = append(disks, Disk{
			ID:         d.ID.String(),
			Name:       d.Name,
			Desc:       d.Description,
			Zone:       c.zone,
			SizeGB:     d.SizeMB / 1024,
			Connection: connection,
			ServerID:   serverID,
			ServerName: serverName,
			CreatedAt:  createdAt,
		})
	}

	slog.Info("Successfully fetched disks",
		slog.String("zone", c.zone),
		slog.Int("count", len(disks)))

	return disks, nil
}

func (c *SakuraClient) GetDiskDetail(ctx context.Context, diskID string) (*DiskDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching disk detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("diskID", diskID))

	diskOp := iaas.NewDiskOp(c.caller)

	id := types.StringID(diskID)

	d, err := diskOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch disk detail",
			slog.String("zone", c.zone),
			slog.String("diskID", diskID),
			slog.Any("error", err))
		return nil, err
	}

	// Get connection type
	connection := string(d.Connection)

	// Get server info
	serverID := ""
	serverName := ""
	if !d.ServerID.IsEmpty() {
		serverID = d.ServerID.String()
		serverName = d.ServerName
	}

	// Format created at
	createdAt := ""
	if !d.CreatedAt.IsZero() {
		createdAt = d.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Format modified at
	modifiedAt := ""
	if !d.ModifiedAt.IsZero() {
		modifiedAt = d.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Get source info
	sourceDiskID := ""
	if !d.SourceDiskID.IsEmpty() {
		sourceDiskID = d.SourceDiskID.String()
	}
	sourceArchiveID := ""
	if !d.SourceArchiveID.IsEmpty() {
		sourceArchiveID = d.SourceArchiveID.String()
	}

	detail := &DiskDetail{
		Disk: Disk{
			ID:         d.ID.String(),
			Name:       d.Name,
			Desc:       d.Description,
			Zone:       c.zone,
			SizeGB:     d.SizeMB / 1024,
			Connection: connection,
			ServerID:   serverID,
			ServerName: serverName,
			CreatedAt:  createdAt,
		},
		Tags:            d.Tags,
		DiskPlanName:    d.DiskPlanName,
		SourceDiskID:    sourceDiskID,
		SourceArchiveID: sourceArchiveID,
		Availability:    string(d.Availability),
		EncryptionAlgo:  string(d.EncryptionAlgorithm),
		CreatedAt:       createdAt,
		ModifiedAt:      modifiedAt,
	}

	slog.Info("Successfully fetched disk detail",
		slog.String("zone", c.zone),
		slog.String("diskID", diskID))

	return detail, nil
}
