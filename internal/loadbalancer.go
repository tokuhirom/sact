package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// LoadBalancer represents a load balancer resource
type LoadBalancer struct {
	ID             string
	Name           string
	Desc           string
	Zone           string
	InstanceStatus string
	VIPCount       int
	CreatedAt      string
}

// LoadBalancerServer represents a real server in a VIP
type LoadBalancerServer struct {
	IPAddress string
	Port      int
	Enabled   bool
}

// LoadBalancerVIP represents a virtual IP configuration
type LoadBalancerVIP struct {
	VirtualIPAddress string
	Port             int
	DelayLoop        int
	SorryServer      string
	Description      string
	Servers          []LoadBalancerServer
}

type LoadBalancerDetail struct {
	LoadBalancer
	Tags           []string
	PlanID         string
	SwitchID       string
	DefaultRoute   string
	NetworkMaskLen int
	IPAddresses    []string
	VRID           int
	VIPs           []LoadBalancerVIP
	CreatedAt      string
}

// Implement list.Item interface for LoadBalancer
func (lb LoadBalancer) FilterValue() string {
	return lb.Name
}

func (lb LoadBalancer) Title() string {
	return lb.Name
}

func (lb LoadBalancer) Description() string {
	desc := fmt.Sprintf("ID: %s", lb.ID)
	if lb.Desc != "" {
		desc += " | " + lb.Desc
	}
	return desc
}

func (c *SakuraClient) ListLoadBalancers(ctx context.Context) ([]LoadBalancer, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching load balancers from Sakura Cloud",
		slog.String("zone", c.zone))

	lbOp := iaas.NewLoadBalancerOp(c.caller)

	searched, err := lbOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch load balancers",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	loadBalancers := make([]LoadBalancer, 0, len(searched.LoadBalancers))
	for _, lb := range searched.LoadBalancers {
		// Format created at
		createdAt := ""
		if !lb.CreatedAt.IsZero() {
			createdAt = lb.CreatedAt.Format("2006-01-02")
		}

		// Count VIPs
		vipCount := len(lb.VirtualIPAddresses)

		loadBalancers = append(loadBalancers, LoadBalancer{
			ID:             lb.ID.String(),
			Name:           lb.Name,
			Desc:           lb.Description,
			Zone:           c.zone,
			InstanceStatus: string(lb.InstanceStatus),
			VIPCount:       vipCount,
			CreatedAt:      createdAt,
		})
	}

	slog.Info("Successfully fetched load balancers",
		slog.String("zone", c.zone),
		slog.Int("count", len(loadBalancers)))

	return loadBalancers, nil
}

func (c *SakuraClient) GetLoadBalancerDetail(ctx context.Context, loadBalancerID string) (*LoadBalancerDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching load balancer detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("loadBalancerID", loadBalancerID))

	lbOp := iaas.NewLoadBalancerOp(c.caller)

	id := types.StringID(loadBalancerID)

	lb, err := lbOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch load balancer detail",
			slog.String("zone", c.zone),
			slog.String("loadBalancerID", loadBalancerID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !lb.CreatedAt.IsZero() {
		createdAt = lb.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Convert VIPs
	vips := make([]LoadBalancerVIP, 0, len(lb.VirtualIPAddresses))
	for _, vip := range lb.VirtualIPAddresses {
		servers := make([]LoadBalancerServer, 0, len(vip.Servers))
		for _, srv := range vip.Servers {
			servers = append(servers, LoadBalancerServer{
				IPAddress: srv.IPAddress,
				Port:      int(srv.Port),
				Enabled:   bool(srv.Enabled),
			})
		}
		vips = append(vips, LoadBalancerVIP{
			VirtualIPAddress: vip.VirtualIPAddress,
			Port:             int(vip.Port),
			DelayLoop:        int(vip.DelayLoop),
			SorryServer:      vip.SorryServer,
			Description:      vip.Description,
			Servers:          servers,
		})
	}

	// Convert tags
	tags := make([]string, 0, len(lb.Tags))
	tags = append(tags, lb.Tags...)

	detail := &LoadBalancerDetail{
		LoadBalancer: LoadBalancer{
			ID:             lb.ID.String(),
			Name:           lb.Name,
			Desc:           lb.Description,
			Zone:           c.zone,
			InstanceStatus: string(lb.InstanceStatus),
			VIPCount:       len(lb.VirtualIPAddresses),
			CreatedAt:      createdAt,
		},
		Tags:           tags,
		PlanID:         lb.PlanID.String(),
		SwitchID:       lb.SwitchID.String(),
		DefaultRoute:   lb.DefaultRoute,
		NetworkMaskLen: lb.NetworkMaskLen,
		IPAddresses:    lb.IPAddresses,
		VRID:           lb.VRID,
		VIPs:           vips,
		CreatedAt:      createdAt,
	}

	slog.Info("Successfully fetched load balancer detail",
		slog.String("zone", c.zone),
		slog.String("loadBalancerID", loadBalancerID))

	return detail, nil
}
