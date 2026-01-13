package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// PacketFilter represents a packet filter resource
type PacketFilter struct {
	ID        string
	Name      string
	Desc      string
	Zone      string
	RuleCount int
	CreatedAt string
}

// PacketFilterRule represents a single rule in a packet filter
type PacketFilterRule struct {
	Protocol        string
	SourceNetwork   string
	SourcePort      string
	DestinationPort string
	Action          string
	Description     string
}

type PacketFilterDetail struct {
	PacketFilter
	Rules          []PacketFilterRule
	ExpressionHash string
	CreatedAt      string
}

// Implement list.Item interface for PacketFilter
func (p PacketFilter) FilterValue() string {
	return p.Name
}

func (p PacketFilter) Title() string {
	return p.Name
}

func (p PacketFilter) Description() string {
	desc := fmt.Sprintf("ID: %s", p.ID)
	if p.Desc != "" {
		desc += " | " + p.Desc
	}
	return desc
}

func (c *SakuraClient) ListPacketFilters(ctx context.Context) ([]PacketFilter, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching packet filters from Sakura Cloud",
		slog.String("zone", c.zone))

	pfOp := iaas.NewPacketFilterOp(c.caller)

	searched, err := pfOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch packet filters",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	packetFilters := make([]PacketFilter, 0, len(searched.PacketFilters))
	for _, pf := range searched.PacketFilters {
		// Format created at
		createdAt := ""
		if !pf.CreatedAt.IsZero() {
			createdAt = pf.CreatedAt.Format("2006-01-02")
		}

		// Count rules
		ruleCount := len(pf.Expression)

		packetFilters = append(packetFilters, PacketFilter{
			ID:        pf.ID.String(),
			Name:      pf.Name,
			Desc:      pf.Description,
			Zone:      c.zone,
			RuleCount: ruleCount,
			CreatedAt: createdAt,
		})
	}

	slog.Info("Successfully fetched packet filters",
		slog.String("zone", c.zone),
		slog.Int("count", len(packetFilters)))

	return packetFilters, nil
}

func (c *SakuraClient) GetPacketFilterDetail(ctx context.Context, packetFilterID string) (*PacketFilterDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching packet filter detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("packetFilterID", packetFilterID))

	pfOp := iaas.NewPacketFilterOp(c.caller)

	id := types.StringID(packetFilterID)

	pf, err := pfOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch packet filter detail",
			slog.String("zone", c.zone),
			slog.String("packetFilterID", packetFilterID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !pf.CreatedAt.IsZero() {
		createdAt = pf.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Convert rules
	rules := make([]PacketFilterRule, 0, len(pf.Expression))
	for _, expr := range pf.Expression {
		rules = append(rules, PacketFilterRule{
			Protocol:        string(expr.Protocol),
			SourceNetwork:   string(expr.SourceNetwork),
			SourcePort:      string(expr.SourcePort),
			DestinationPort: string(expr.DestinationPort),
			Action:          string(expr.Action),
			Description:     expr.Description,
		})
	}

	detail := &PacketFilterDetail{
		PacketFilter: PacketFilter{
			ID:        pf.ID.String(),
			Name:      pf.Name,
			Desc:      pf.Description,
			Zone:      c.zone,
			RuleCount: len(pf.Expression),
			CreatedAt: createdAt,
		},
		Rules:          rules,
		ExpressionHash: pf.ExpressionHash,
		CreatedAt:      createdAt,
	}

	slog.Info("Successfully fetched packet filter detail",
		slog.String("zone", c.zone),
		slog.String("packetFilterID", packetFilterID))

	return detail, nil
}
