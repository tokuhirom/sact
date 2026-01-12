package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

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
