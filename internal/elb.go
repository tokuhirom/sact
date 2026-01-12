package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

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
