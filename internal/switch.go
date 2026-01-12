package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

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
