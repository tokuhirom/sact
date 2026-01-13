package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// Bridge represents a bridge resource
type Bridge struct {
	ID          string
	Name        string
	Desc        string
	Zone        string
	Region      string
	SwitchCount int
}

type BridgeSwitchInfo struct {
	ID             string
	Name           string
	ZoneName       string
	Scope          string
	ServerCount    int
	ApplianceCount int
}

type BridgeDetail struct {
	Bridge
	Switches     []BridgeSwitchInfo
	SwitchInZone *BridgeSwitchInfo
	CreatedAt    string
}

// Implement list.Item interface for Bridge
func (b Bridge) FilterValue() string {
	return b.Name
}

func (b Bridge) Title() string {
	return b.Name
}

func (b Bridge) Description() string {
	desc := fmt.Sprintf("ID: %s", b.ID)
	if b.Desc != "" {
		desc += " | " + b.Desc
	}
	return desc
}

func (c *SakuraClient) ListBridges(ctx context.Context) ([]Bridge, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching bridges from Sakura Cloud",
		slog.String("zone", c.zone))

	bridgeOp := iaas.NewBridgeOp(c.caller)

	searched, err := bridgeOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch bridges",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	bridges := make([]Bridge, 0, len(searched.Bridges))
	for _, b := range searched.Bridges {
		// Get region name
		region := ""
		if b.Region != nil {
			region = b.Region.Name
		}

		// Count switches
		switchCount := len(b.BridgeInfo)

		bridges = append(bridges, Bridge{
			ID:          b.ID.String(),
			Name:        b.Name,
			Desc:        b.Description,
			Zone:        c.zone,
			Region:      region,
			SwitchCount: switchCount,
		})
	}

	slog.Info("Successfully fetched bridges",
		slog.String("zone", c.zone),
		slog.Int("count", len(bridges)))

	return bridges, nil
}

func (c *SakuraClient) GetBridgeDetail(ctx context.Context, bridgeID string) (*BridgeDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching bridge detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("bridgeID", bridgeID))

	bridgeOp := iaas.NewBridgeOp(c.caller)

	id := types.StringID(bridgeID)

	b, err := bridgeOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch bridge detail",
			slog.String("zone", c.zone),
			slog.String("bridgeID", bridgeID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !b.CreatedAt.IsZero() {
		createdAt = b.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Get region name
	region := ""
	if b.Region != nil {
		region = b.Region.Name
	}

	// Convert bridge info (switches)
	switches := make([]BridgeSwitchInfo, 0, len(b.BridgeInfo))
	for _, sw := range b.BridgeInfo {
		switches = append(switches, BridgeSwitchInfo{
			ID:       sw.ID.String(),
			Name:     sw.Name,
			ZoneName: sw.ZoneName,
		})
	}

	// Get switch in zone info
	var switchInZone *BridgeSwitchInfo
	if b.SwitchInZone != nil {
		switchInZone = &BridgeSwitchInfo{
			ID:             b.SwitchInZone.ID.String(),
			Name:           b.SwitchInZone.Name,
			Scope:          string(b.SwitchInZone.Scope),
			ServerCount:    b.SwitchInZone.ServerCount,
			ApplianceCount: b.SwitchInZone.ApplianceCount,
		}
	}

	detail := &BridgeDetail{
		Bridge: Bridge{
			ID:          b.ID.String(),
			Name:        b.Name,
			Desc:        b.Description,
			Zone:        c.zone,
			Region:      region,
			SwitchCount: len(b.BridgeInfo),
		},
		Switches:     switches,
		SwitchInZone: switchInZone,
		CreatedAt:    createdAt,
	}

	slog.Info("Successfully fetched bridge detail",
		slog.String("zone", c.zone),
		slog.String("bridgeID", bridgeID))

	return detail, nil
}
