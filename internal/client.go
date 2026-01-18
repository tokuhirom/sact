package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"

	apprun "github.com/tokuhirom/sact/pkg/openapi/apprun_dedicated"
)

// debugTransport logs HTTP request/response for debugging
type debugTransport struct {
	base http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	slog.Debug("AppRun API request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()))

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Log response for non-2xx status codes
	if resp.StatusCode >= 400 {
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			slog.Error("Failed to read error response body", slog.Any("error", readErr))
			return nil, readErr
		}

		slog.Error("AppRun API error response",
			slog.Int("status", resp.StatusCode),
			slog.String("body", string(body)))

		// Restore the body for further processing
		resp.Body = io.NopCloser(bytes.NewReader(body))
	}

	return resp, nil
}

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
	ResourceTypeAppRunDedicated
	ResourceTypeMonitoringLogStorage
	ResourceTypeMonitoringMetricsStorage
	ResourceTypeMonitoringTraceStorage
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
	ResourceTypeAppRunDedicated,
	ResourceTypeMonitoringLogStorage,
	ResourceTypeMonitoringMetricsStorage,
	ResourceTypeMonitoringTraceStorage,
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
	case ResourceTypeAppRunDedicated:
		return "AppRun Dedicated"
	case ResourceTypeMonitoringLogStorage:
		return "Monitoring Suite - Log Storage"
	case ResourceTypeMonitoringMetricsStorage:
		return "Monitoring Suite - Metrics Storage"
	case ResourceTypeMonitoringTraceStorage:
		return "Monitoring Suite - Trace Storage"
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

func NewSakuraClient(opts *api.CallerOptions, zone string) (*SakuraClient, error) {
	if zone == "" {
		slog.Error("Zone is empty")
		return nil, fmt.Errorf("zone must be specified")
	}

	slog.Info("Creating Sakura Cloud API caller", slog.String("zone", zone))

	caller := api.NewCallerWithOptions(opts)

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

	// Use debug transport to log error responses
	httpClient := &http.Client{
		Transport: &debugTransport{base: http.DefaultTransport},
	}

	apprunClient, err := apprun.NewClient(
		"https://secure.sakura.ad.jp/cloud/api/apprun-dedicated/1.0",
		secSource,
		apprun.WithClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create AppRun client: %w", err)
	}

	c.apprunClient = apprunClient
	return c.apprunClient, nil
}
