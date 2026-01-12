package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// VPCRouter represents a VPC Router resource
type VPCRouter struct {
	ID             string
	Name           string
	Desc           string
	Zone           string
	Plan           string
	Version        int
	InstanceStatus string
	CreatedAt      string
}

type VPCRouterDetail struct {
	VPCRouter
	Tags              []string
	PublicIPAddresses []string
	NICs              []VPCRouterNIC
	CreatedAt         string
}

type VPCRouterNIC struct {
	Index     int
	SwitchID  string
	IPAddress string
}

// Implement list.Item interface for VPCRouter
func (v VPCRouter) FilterValue() string {
	return v.Name
}

func (v VPCRouter) Title() string {
	return v.Name
}

func (v VPCRouter) Description() string {
	desc := fmt.Sprintf("ID: %s", v.ID)
	if v.Desc != "" {
		desc += " | " + v.Desc
	}
	return desc
}

func (c *SakuraClient) ListVPCRouters(ctx context.Context) ([]VPCRouter, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching VPC routers from Sakura Cloud",
		slog.String("zone", c.zone))

	vpcRouterOp := iaas.NewVPCRouterOp(c.caller)

	searched, err := vpcRouterOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch VPC routers",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	vpcRouterList := make([]VPCRouter, 0, len(searched.VPCRouters))
	for _, v := range searched.VPCRouters {
		// Format created at
		createdAt := ""
		if !v.CreatedAt.IsZero() {
			createdAt = v.CreatedAt.Format("2006-01-02")
		}

		// Get plan name
		planName := v.PlanID.String()

		vpcRouterList = append(vpcRouterList, VPCRouter{
			ID:             v.ID.String(),
			Name:           v.Name,
			Desc:           v.Description,
			Zone:           c.zone,
			Plan:           planName,
			Version:        v.Version,
			InstanceStatus: string(v.InstanceStatus),
			CreatedAt:      createdAt,
		})
	}

	slog.Info("Successfully fetched VPC routers",
		slog.String("zone", c.zone),
		slog.Int("count", len(vpcRouterList)))

	return vpcRouterList, nil
}

func (c *SakuraClient) GetVPCRouterDetail(ctx context.Context, vpcRouterID string) (*VPCRouterDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching VPC router detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("vpcRouterID", vpcRouterID))

	vpcRouterOp := iaas.NewVPCRouterOp(c.caller)

	id := types.StringID(vpcRouterID)

	v, err := vpcRouterOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch VPC router detail",
			slog.String("zone", c.zone),
			slog.String("vpcRouterID", vpcRouterID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !v.CreatedAt.IsZero() {
		createdAt = v.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Get plan name
	planName := v.PlanID.String()

	// Get public IP addresses
	publicIPs := make([]string, 0)
	for _, iface := range v.Interfaces {
		if iface.IPAddress != "" {
			publicIPs = append(publicIPs, iface.IPAddress)
		}
	}

	// Get NICs
	nics := make([]VPCRouterNIC, 0)
	for i, iface := range v.Interfaces {
		switchID := ""
		if iface.SwitchID != 0 {
			switchID = fmt.Sprintf("%d", iface.SwitchID)
		}
		nics = append(nics, VPCRouterNIC{
			Index:     i,
			SwitchID:  switchID,
			IPAddress: iface.IPAddress,
		})
	}

	detail := &VPCRouterDetail{
		VPCRouter: VPCRouter{
			ID:             v.ID.String(),
			Name:           v.Name,
			Desc:           v.Description,
			Zone:           c.zone,
			Plan:           planName,
			Version:        v.Version,
			InstanceStatus: string(v.InstanceStatus),
			CreatedAt:      createdAt,
		},
		Tags:              v.Tags,
		PublicIPAddresses: publicIPs,
		NICs:              nics,
		CreatedAt:         createdAt,
	}

	slog.Info("Successfully fetched VPC router detail",
		slog.String("zone", c.zone),
		slog.String("vpcRouterID", vpcRouterID))

	return detail, nil
}
