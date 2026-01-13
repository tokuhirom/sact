package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
)

// ResourceType represents different cloud resource types
type ResourceType int

const (
	ResourceTypeServer ResourceType = iota
	ResourceTypeSwitch
	ResourceTypeDNS
	ResourceTypeELB
	ResourceTypeGSLB
	ResourceTypeDB
	ResourceTypeDisk
	ResourceTypeArchive
	ResourceTypeInternet
	ResourceTypeVPCRouter
	ResourceTypePacketFilter
)

// AllResourceTypes returns all available resource types
var AllResourceTypes = []ResourceType{
	ResourceTypeServer,
	ResourceTypeSwitch,
	ResourceTypeDNS,
	ResourceTypeELB,
	ResourceTypeGSLB,
	ResourceTypeDB,
	ResourceTypeDisk,
	ResourceTypeArchive,
	ResourceTypeInternet,
	ResourceTypeVPCRouter,
	ResourceTypePacketFilter,
}

func (r ResourceType) String() string {
	switch r {
	case ResourceTypeServer:
		return "Server"
	case ResourceTypeSwitch:
		return "Switch"
	case ResourceTypeDNS:
		return "DNS"
	case ResourceTypeELB:
		return "ELB"
	case ResourceTypeGSLB:
		return "GSLB"
	case ResourceTypeDB:
		return "DB"
	case ResourceTypeDisk:
		return "Disk"
	case ResourceTypeArchive:
		return "Archive"
	case ResourceTypeInternet:
		return "Internet"
	case ResourceTypeVPCRouter:
		return "VPCRouter"
	case ResourceTypePacketFilter:
		return "PacketFilter"
	default:
		return "Unknown"
	}
}

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

func (c *SakuraClient) SetZone(zone string) {
	slog.Info("Switching zone",
		slog.String("from", c.zone),
		slog.String("to", zone))
	c.zone = zone
}

func (c *SakuraClient) GetZone() string {
	return c.zone
}

func (c *SakuraClient) GetAuthStatus(ctx context.Context) (*iaas.AuthStatus, error) {
	op := iaas.NewAuthStatusOp(c.caller)
	return op.Read(ctx)
}
