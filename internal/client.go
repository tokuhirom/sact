package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
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
