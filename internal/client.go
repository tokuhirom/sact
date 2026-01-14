package internal

import (
	"context"
	"fmt"
	"log/slog"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"

	apprun "github.com/tokuhirom/sact/pkg/openapi/apprun_dedicated"
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
	ResourceTypeLoadBalancer
	ResourceTypeNFS
	ResourceTypeSSHKey
	ResourceTypeAutoBackup
	ResourceTypeSimpleMonitor
	ResourceTypeBridge
	ResourceTypeContainerRegistry
	ResourceTypeAppRunCluster
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
	ResourceTypeLoadBalancer,
	ResourceTypeNFS,
	ResourceTypeSSHKey,
	ResourceTypeAutoBackup,
	ResourceTypeSimpleMonitor,
	ResourceTypeBridge,
	ResourceTypeContainerRegistry,
	ResourceTypeAppRunCluster,
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
	case ResourceTypeLoadBalancer:
		return "LoadBalancer"
	case ResourceTypeNFS:
		return "NFS"
	case ResourceTypeSSHKey:
		return "SSHKey"
	case ResourceTypeAutoBackup:
		return "AutoBackup"
	case ResourceTypeSimpleMonitor:
		return "SimpleMonitor"
	case ResourceTypeBridge:
		return "Bridge"
	case ResourceTypeContainerRegistry:
		return "ContainerRegistry"
	case ResourceTypeAppRunCluster:
		return "AppRunCluster"
	default:
		return "Unknown"
	}
}

type SakuraClient struct {
	caller       iaas.APICaller
	zone         string
	apprunClient *apprun.Client
}

// apprunSecuritySource implements apprun.SecuritySource for BasicAuth
type apprunSecuritySource struct {
	username string
	password string
}

func (s *apprunSecuritySource) BasicAuth(ctx context.Context, operationName apprun.OperationName) (apprun.BasicAuth, error) {
	return apprun.BasicAuth{
		Username: s.username,
		Password: s.password,
	}, nil
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

// GetAppRunClient returns the AppRun Dedicated API client (lazy initialization)
func (c *SakuraClient) GetAppRunClient() (*apprun.Client, error) {
	if c.apprunClient != nil {
		return c.apprunClient, nil
	}

	// Use api-client-go to get credentials from profile or environment variables
	clientOpts, err := client.DefaultOption()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	token := clientOpts.AccessToken
	secret := clientOpts.AccessTokenSecret

	if token == "" || secret == "" {
		return nil, fmt.Errorf("credentials not found in profile or environment variables")
	}

	secSource := &apprunSecuritySource{
		username: token,
		password: secret,
	}

	apprunClient, err := apprun.NewClient(
		"https://secure.sakura.ad.jp/cloud/api/apprun-dedicated/1.0",
		secSource,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create AppRun client: %w", err)
	}

	c.apprunClient = apprunClient
	return c.apprunClient, nil
}
