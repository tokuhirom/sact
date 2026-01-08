package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/types"
)

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
	Description    string
	Tags           []string
	CPU            int
	MemoryGB       int
	InterfaceCount int
	Disks          []DiskInfo
	IPAddresses    []string
	CreatedAt      string
}

type DiskInfo struct {
	Name   string
	SizeGB int
}

// Implement list.Item interface
func (s Server) FilterValue() string {
	return s.Name
}

func (s Server) Title() string {
	return s.Name
}

func (s Server) Description() string {
	return fmt.Sprintf("ID: %s | Status: %s", s.ID, s.InstanceStatus)
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

	searched, err := serverOp.Find(ctx, c.zone, &iaas.FindCondition{})
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
