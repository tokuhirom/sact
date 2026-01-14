package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	client "github.com/sacloud/api-client-go"
	monitoring "github.com/tokuhirom/sact/pkg/openapi/monitoring_suite"
)

// MonitoringLogStorage wraps the generated LogStorage type for TUI display
type MonitoringLogStorage struct {
	monitoring.LogStorage
}

// MonitoringMetricsStorage wraps the generated MetricsStorage type for TUI display
type MonitoringMetricsStorage struct {
	monitoring.MetricsStorage
}

// MonitoringLogRouting wraps the generated LogRouting type for TUI display
type MonitoringLogRouting struct {
	monitoring.LogRouting
}

// MonitoringMetricsRouting wraps the generated MetricsRouting type for TUI display
type MonitoringMetricsRouting struct {
	monitoring.MetricsRouting
}

// Helper functions to safely get values from pointers
func getStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getInt64Ptr(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func getIntPtr(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func getBoolPtr(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Implement list.Item interface for MonitoringLogStorage
func (s MonitoringLogStorage) FilterValue() string {
	return getStringPtr(s.Name)
}

func (s MonitoringLogStorage) Title() string {
	return getStringPtr(s.Name)
}

func (s MonitoringLogStorage) Description() string {
	routings := 0
	if s.Usage != nil {
		routings = s.Usage.LogRoutings
	}
	return fmt.Sprintf("ID: %s | Expire: %d days | Routings: %d", getStringPtr(s.ResourceId), getIntPtr(s.ExpireDay), routings)
}

// Implement list.Item interface for MonitoringMetricsStorage
func (s MonitoringMetricsStorage) FilterValue() string {
	return getStringPtr(s.Name)
}

func (s MonitoringMetricsStorage) Title() string {
	return getStringPtr(s.Name)
}

func (s MonitoringMetricsStorage) Description() string {
	routings := 0
	alertRules := 0
	if s.Usage != nil {
		routings = s.Usage.MetricsRoutings
		alertRules = s.Usage.AlertRules
	}
	return fmt.Sprintf("ID: %s | Routings: %d | Alert Rules: %d", getStringPtr(s.ResourceId), routings, alertRules)
}

// Implement list.Item interface for MonitoringLogRouting
func (r MonitoringLogRouting) FilterValue() string {
	if r.Uid != nil {
		return r.Uid.String()
	}
	return ""
}

func (r MonitoringLogRouting) Title() string {
	if r.Uid != nil {
		return r.Uid.String()[:8]
	}
	return ""
}

func (r MonitoringLogRouting) Description() string {
	return fmt.Sprintf("ResourceID: %s -> LogStorage: %d | Variant: %s",
		getStringPtr(r.ResourceId),
		getInt64Ptr(r.LogStorageId),
		r.Variant)
}

// Implement list.Item interface for MonitoringMetricsRouting
func (r MonitoringMetricsRouting) FilterValue() string {
	if r.Uid != nil {
		return r.Uid.String()
	}
	return ""
}

func (r MonitoringMetricsRouting) Title() string {
	if r.Uid != nil {
		return r.Uid.String()[:8]
	}
	return ""
}

func (r MonitoringMetricsRouting) Description() string {
	return fmt.Sprintf("ResourceID: %s -> MetricsStorage: %d | Variant: %s",
		getStringPtr(r.ResourceId),
		getInt64Ptr(r.MetricsStorageId),
		r.Variant)
}

// MonitoringLogStorageDetail contains detailed info about a log storage
type MonitoringLogStorageDetail struct {
	MonitoringLogStorage
	Tags     []string
	Routings []MonitoringLogRouting
}

// MonitoringMetricsStorageDetail contains detailed info about a metrics storage
type MonitoringMetricsStorageDetail struct {
	MonitoringMetricsStorage
	Tags     []string
	Routings []MonitoringMetricsRouting
}

// getMonitoringClient creates a monitoring suite API client
func (c *SakuraClient) getMonitoringClient() (*monitoring.ClientWithResponses, error) {
	clientOpts, err := client.DefaultOption()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	token := clientOpts.AccessToken
	secret := clientOpts.AccessTokenSecret

	if token == "" || secret == "" {
		return nil, fmt.Errorf("credentials not found")
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create client with basic auth
	monClient, err := monitoring.NewClientWithResponses(
		"https://secure.sakura.ad.jp/cloud/zone/is1a/api/monitoring/1.0",
		monitoring.WithHTTPClient(httpClient),
		monitoring.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.SetBasicAuth(token, secret)
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitoring client: %w", err)
	}

	return monClient, nil
}

// ListMonitoringLogStorages fetches all log storages
func (c *SakuraClient) ListMonitoringLogStorages(ctx context.Context) ([]MonitoringLogStorage, error) {
	slog.Info("Fetching Monitoring Log Storages")

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	resp, err := monClient.LogsStoragesListWithResponse(ctx, nil)
	if err != nil {
		slog.Error("Failed to fetch log storages", slog.Any("error", err))
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status())
	}

	result := make([]MonitoringLogStorage, len(resp.JSON200.Results))
	for i, storage := range resp.JSON200.Results {
		result[i] = MonitoringLogStorage{LogStorage: storage}
	}

	slog.Info("Successfully fetched log storages", slog.Int("count", len(result)))
	return result, nil
}

// ListMonitoringMetricsStorages fetches all metrics storages
func (c *SakuraClient) ListMonitoringMetricsStorages(ctx context.Context) ([]MonitoringMetricsStorage, error) {
	slog.Info("Fetching Monitoring Metrics Storages")

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	resp, err := monClient.MetricsStoragesListWithResponse(ctx, nil)
	if err != nil {
		slog.Error("Failed to fetch metrics storages", slog.Any("error", err))
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status())
	}

	result := make([]MonitoringMetricsStorage, len(resp.JSON200.Results))
	for i, storage := range resp.JSON200.Results {
		result[i] = MonitoringMetricsStorage{MetricsStorage: storage}
	}

	slog.Info("Successfully fetched metrics storages", slog.Int("count", len(result)))
	return result, nil
}

// ListMonitoringLogRoutings fetches all log routing rules
func (c *SakuraClient) ListMonitoringLogRoutings(ctx context.Context) ([]MonitoringLogRouting, error) {
	slog.Info("Fetching Monitoring Log Routings")

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	resp, err := monClient.LogsRoutingsListWithResponse(ctx, nil)
	if err != nil {
		slog.Error("Failed to fetch log routings", slog.Any("error", err))
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status())
	}

	result := make([]MonitoringLogRouting, len(resp.JSON200.Results))
	for i, routing := range resp.JSON200.Results {
		result[i] = MonitoringLogRouting{LogRouting: routing}
	}

	slog.Info("Successfully fetched log routings", slog.Int("count", len(result)))
	return result, nil
}

// ListMonitoringMetricsRoutings fetches all metrics routing rules
func (c *SakuraClient) ListMonitoringMetricsRoutings(ctx context.Context) ([]MonitoringMetricsRouting, error) {
	slog.Info("Fetching Monitoring Metrics Routings")

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	resp, err := monClient.MetricsRoutingsListWithResponse(ctx, nil)
	if err != nil {
		slog.Error("Failed to fetch metrics routings", slog.Any("error", err))
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status())
	}

	result := make([]MonitoringMetricsRouting, len(resp.JSON200.Results))
	for i, routing := range resp.JSON200.Results {
		result[i] = MonitoringMetricsRouting{MetricsRouting: routing}
	}

	slog.Info("Successfully fetched metrics routings", slog.Int("count", len(result)))
	return result, nil
}

// GetMonitoringLogStorageDetail fetches detail for a log storage
func (c *SakuraClient) GetMonitoringLogStorageDetail(ctx context.Context, resourceID string) (*MonitoringLogStorageDetail, error) {
	slog.Info("Fetching Log Storage detail", slog.String("resourceID", resourceID))

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	resp, err := monClient.LogsStoragesRetrieveWithResponse(ctx, resourceID)
	if err != nil {
		slog.Error("Failed to fetch log storage detail", slog.Any("error", err))
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status())
	}

	// Also fetch routings
	routings, _ := c.ListMonitoringLogRoutings(ctx)

	// Filter routings for this storage (by LogStorageId)
	var filteredRoutings []MonitoringLogRouting
	for _, r := range routings {
		logStorageID := getInt64Ptr(r.LogStorageId)
		if logStorageID == 0 {
			continue
		}
		if fmt.Sprintf("%d", logStorageID) == resourceID {
			filteredRoutings = append(filteredRoutings, r)
		}
	}

	var tags []string
	if resp.JSON200.Tags != nil {
		tags = *resp.JSON200.Tags
	}

	// Convert WrappedLogStorage to LogStorage
	storage := monitoring.LogStorage{
		AccountId:   resp.JSON200.AccountId,
		CreatedAt:   resp.JSON200.CreatedAt,
		Description: resp.JSON200.Description,
		Endpoints:   resp.JSON200.Endpoints,
		ExpireDay:   resp.JSON200.ExpireDay,
		Icon:        resp.JSON200.Icon,
		Id:          resp.JSON200.Id,
		IsSystem:    resp.JSON200.IsSystem,
		Name:        resp.JSON200.Name,
		ResourceId:  resp.JSON200.ResourceId,
		Tags:        resp.JSON200.Tags,
		Usage:       resp.JSON200.Usage,
	}

	detail := &MonitoringLogStorageDetail{
		MonitoringLogStorage: MonitoringLogStorage{LogStorage: storage},
		Tags:                 tags,
		Routings:             filteredRoutings,
	}

	slog.Info("Successfully fetched log storage detail", slog.String("resourceID", resourceID))
	return detail, nil
}

// GetMonitoringMetricsStorageDetail fetches detail for a metrics storage
func (c *SakuraClient) GetMonitoringMetricsStorageDetail(ctx context.Context, resourceID string) (*MonitoringMetricsStorageDetail, error) {
	slog.Info("Fetching Metrics Storage detail", slog.String("resourceID", resourceID))

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	resp, err := monClient.MetricsStoragesRetrieveWithResponse(ctx, resourceID)
	if err != nil {
		slog.Error("Failed to fetch metrics storage detail", slog.Any("error", err))
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status())
	}

	// Also fetch routings
	routings, _ := c.ListMonitoringMetricsRoutings(ctx)

	// Filter routings for this storage (by MetricsStorageId)
	var filteredRoutings []MonitoringMetricsRouting
	for _, r := range routings {
		metricsStorageID := getInt64Ptr(r.MetricsStorageId)
		if metricsStorageID == 0 {
			continue
		}
		if fmt.Sprintf("%d", metricsStorageID) == resourceID {
			filteredRoutings = append(filteredRoutings, r)
		}
	}

	var tags []string
	if resp.JSON200.Tags != nil {
		tags = *resp.JSON200.Tags
	}

	// Convert WrappedMetricsStorage to MetricsStorage
	storage := monitoring.MetricsStorage{
		AccountId:   resp.JSON200.AccountId,
		CreatedAt:   resp.JSON200.CreatedAt,
		Description: resp.JSON200.Description,
		Endpoints:   resp.JSON200.Endpoints,
		Icon:        resp.JSON200.Icon,
		Id:          resp.JSON200.Id,
		IsSystem:    resp.JSON200.IsSystem,
		Name:        resp.JSON200.Name,
		ResourceId:  resp.JSON200.ResourceId,
		Tags:        resp.JSON200.Tags,
		UpdatedAt:   resp.JSON200.UpdatedAt,
		Usage:       resp.JSON200.Usage,
	}

	detail := &MonitoringMetricsStorageDetail{
		MonitoringMetricsStorage: MonitoringMetricsStorage{MetricsStorage: storage},
		Tags:                     tags,
		Routings:                 filteredRoutings,
	}

	slog.Info("Successfully fetched metrics storage detail", slog.String("resourceID", resourceID))
	return detail, nil
}
