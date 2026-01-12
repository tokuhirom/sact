package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// Internet represents a router (Internet) resource
type Internet struct {
	ID            string
	Name          string
	Desc          string
	Zone          string
	BandWidthMbps int
	SwitchID      string
	CreatedAt     string
}

type InternetDetail struct {
	Internet
	Tags           []string
	NetworkMaskLen int
	CreatedAt      string
}

// Implement list.Item interface for Internet
func (i Internet) FilterValue() string {
	return i.Name
}

func (i Internet) Title() string {
	return i.Name
}

func (i Internet) Description() string {
	desc := fmt.Sprintf("ID: %s", i.ID)
	if i.Desc != "" {
		desc += " | " + i.Desc
	}
	return desc
}

func (c *SakuraClient) ListInternet(ctx context.Context) ([]Internet, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching internet (routers) from Sakura Cloud",
		slog.String("zone", c.zone))

	internetOp := iaas.NewInternetOp(c.caller)

	searched, err := internetOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch internet",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	internetList := make([]Internet, 0, len(searched.Internet))
	for _, i := range searched.Internet {
		// Format created at
		createdAt := ""
		if !i.CreatedAt.IsZero() {
			createdAt = i.CreatedAt.Format("2006-01-02")
		}

		// Get switch ID
		switchID := ""
		if i.Switch != nil && !i.Switch.ID.IsEmpty() {
			switchID = i.Switch.ID.String()
		}

		internetList = append(internetList, Internet{
			ID:            i.ID.String(),
			Name:          i.Name,
			Desc:          i.Description,
			Zone:          c.zone,
			BandWidthMbps: i.BandWidthMbps,
			SwitchID:      switchID,
			CreatedAt:     createdAt,
		})
	}

	slog.Info("Successfully fetched internet",
		slog.String("zone", c.zone),
		slog.Int("count", len(internetList)))

	return internetList, nil
}

func (c *SakuraClient) GetInternetDetail(ctx context.Context, internetID string) (*InternetDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching internet detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("internetID", internetID))

	internetOp := iaas.NewInternetOp(c.caller)

	id := types.StringID(internetID)

	i, err := internetOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch internet detail",
			slog.String("zone", c.zone),
			slog.String("internetID", internetID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !i.CreatedAt.IsZero() {
		createdAt = i.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Get switch ID
	switchID := ""
	if i.Switch != nil && !i.Switch.ID.IsEmpty() {
		switchID = i.Switch.ID.String()
	}

	detail := &InternetDetail{
		Internet: Internet{
			ID:            i.ID.String(),
			Name:          i.Name,
			Desc:          i.Description,
			Zone:          c.zone,
			BandWidthMbps: i.BandWidthMbps,
			SwitchID:      switchID,
			CreatedAt:     createdAt,
		},
		Tags:           i.Tags,
		NetworkMaskLen: i.NetworkMaskLen,
		CreatedAt:      createdAt,
	}

	slog.Info("Successfully fetched internet detail",
		slog.String("zone", c.zone),
		slog.String("internetID", internetID))

	return detail, nil
}
