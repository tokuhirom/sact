package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	monitoringsuite "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

// MonitoringLogStorage wraps the generated LogStorage type for TUI display
type MonitoringLogStorage struct {
	v1.LogStorage
}

// MonitoringMetricsStorage wraps the generated MetricsStorage type for TUI display
type MonitoringMetricsStorage struct {
	v1.MetricsStorage
}

// MonitoringTraceStorage wraps the generated TraceStorage type for TUI display
type MonitoringTraceStorage struct {
	v1.TraceStorage
}

// MonitoringLogRouting wraps the generated LogRouting type for TUI display
type MonitoringLogRouting struct {
	v1.LogRouting
}

// MonitoringMetricsRouting wraps the generated MetricsRouting type for TUI display
type MonitoringMetricsRouting struct {
	v1.MetricsRouting
}

// Helper functions to safely get values from optional types
func getOptString(s v1.OptString) string {
	if v, ok := s.Get(); ok {
		return v
	}
	return ""
}

func getNilInt64AsString(n v1.NilInt64) string {
	if v, ok := n.Get(); ok {
		return strconv.FormatInt(v, 10)
	}
	return ""
}

func getOptNilInt64AsString(n v1.OptNilInt64) string {
	if v, ok := n.Get(); ok {
		return strconv.FormatInt(v, 10)
	}
	return ""
}

// Implement list.Item interface for MonitoringLogStorage
func (s MonitoringLogStorage) FilterValue() string {
	return getOptString(s.Name)
}

func (s MonitoringLogStorage) Title() string {
	return getOptString(s.Name)
}

func (s MonitoringLogStorage) Description() string {
	resourceID := getNilInt64AsString(s.ResourceID)
	return fmt.Sprintf("ID: %s | Expire: %d days | Routings: %d", resourceID, s.ExpireDay.Or(0), s.Usage.LogRoutings)
}

// Implement list.Item interface for MonitoringMetricsStorage
func (s MonitoringMetricsStorage) FilterValue() string {
	return getOptString(s.Name)
}

func (s MonitoringMetricsStorage) Title() string {
	return getOptString(s.Name)
}

func (s MonitoringMetricsStorage) Description() string {
	resourceID := getNilInt64AsString(s.ResourceID)
	return fmt.Sprintf("ID: %s | Routings: %d | Alert Rules: %d", resourceID, s.Usage.MetricsRoutings, s.Usage.AlertRules)
}

// Implement list.Item interface for MonitoringTraceStorage
func (s MonitoringTraceStorage) FilterValue() string {
	return getOptString(s.Name)
}

func (s MonitoringTraceStorage) Title() string {
	return getOptString(s.Name)
}

func (s MonitoringTraceStorage) Description() string {
	return fmt.Sprintf("ID: %d | Retention: %d days", s.ResourceID, s.RetentionPeriodDays)
}

// Implement list.Item interface for MonitoringLogRouting
func (r MonitoringLogRouting) FilterValue() string {
	return r.UID.String()
}

func (r MonitoringLogRouting) Title() string {
	return r.UID.String()[:8]
}

func (r MonitoringLogRouting) Description() string {
	resourceID := getOptNilInt64AsString(r.ResourceID)
	logStorageID := getOptNilInt64AsString(r.LogStorageID)
	return fmt.Sprintf("ResourceID: %s -> LogStorage: %s | Variant: %s", resourceID, logStorageID, r.Variant)
}

// Implement list.Item interface for MonitoringMetricsRouting
func (r MonitoringMetricsRouting) FilterValue() string {
	return r.UID.String()
}

func (r MonitoringMetricsRouting) Title() string {
	return r.UID.String()[:8]
}

func (r MonitoringMetricsRouting) Description() string {
	resourceID := getOptNilInt64AsString(r.ResourceID)
	metricsStorageID := getOptNilInt64AsString(r.MetricsStorageID)
	return fmt.Sprintf("ResourceID: %s -> MetricsStorage: %s | Variant: %s", resourceID, metricsStorageID, r.Variant)
}

// MonitoringLogStorageDetail contains detailed info about a log storage
type MonitoringLogStorageDetail struct {
	v1.LogStorage
	Routings []MonitoringLogRouting
}

// MonitoringMetricsStorageDetail contains detailed info about a metrics storage
type MonitoringMetricsStorageDetail struct {
	v1.MetricsStorage
	Routings []MonitoringMetricsRouting
}

// MonitoringTraceStorageDetail contains detailed info about a trace storage
type MonitoringTraceStorageDetail struct {
	v1.TraceStorage
}

// getMonitoringClient creates a monitoring suite API client
func (c *SakuraClient) getMonitoringClient() (*v1.Client, error) {
	return monitoringsuite.NewClient()
}

// ListMonitoringLogStorages fetches all log storages
func (c *SakuraClient) ListMonitoringLogStorages(ctx context.Context) ([]MonitoringLogStorage, error) {
	slog.Info("Fetching Monitoring Log Storages")

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	op := monitoringsuite.NewLogsStorageOp(monClient)
	storages, err := op.List(ctx, monitoringsuite.LogsStoragesListParams{})
	if err != nil {
		slog.Error("Failed to fetch log storages", slog.Any("error", err))
		return nil, err
	}

	result := make([]MonitoringLogStorage, len(storages))
	for i, storage := range storages {
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

	op := monitoringsuite.NewMetricsStorageOp(monClient)
	storages, err := op.List(ctx, monitoringsuite.MetricsStorageListParams{})
	if err != nil {
		slog.Error("Failed to fetch metrics storages", slog.Any("error", err))
		return nil, err
	}

	result := make([]MonitoringMetricsStorage, len(storages))
	for i, storage := range storages {
		result[i] = MonitoringMetricsStorage{MetricsStorage: storage}
	}

	slog.Info("Successfully fetched metrics storages", slog.Int("count", len(result)))
	return result, nil
}

// ListMonitoringTraceStorages fetches all trace storages
func (c *SakuraClient) ListMonitoringTraceStorages(ctx context.Context) ([]MonitoringTraceStorage, error) {
	slog.Info("Fetching Monitoring Trace Storages")

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	op := monitoringsuite.NewTracesStorageOp(monClient)
	storages, err := op.List(ctx, monitoringsuite.TracesStorageListParams{})
	if err != nil {
		slog.Error("Failed to fetch trace storages", slog.Any("error", err))
		return nil, err
	}

	result := make([]MonitoringTraceStorage, len(storages))
	for i, storage := range storages {
		result[i] = MonitoringTraceStorage{TraceStorage: storage}
	}

	slog.Info("Successfully fetched trace storages", slog.Int("count", len(result)))
	return result, nil
}

// ListMonitoringLogRoutings fetches all log routing rules
func (c *SakuraClient) ListMonitoringLogRoutings(ctx context.Context) ([]MonitoringLogRouting, error) {
	slog.Info("Fetching Monitoring Log Routings")

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	op := monitoringsuite.NewLogRoutingOp(monClient)
	routings, err := op.List(ctx, monitoringsuite.LogsRoutingsListParams{})
	if err != nil {
		slog.Error("Failed to fetch log routings", slog.Any("error", err))
		return nil, err
	}

	result := make([]MonitoringLogRouting, len(routings))
	for i, routing := range routings {
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

	op := monitoringsuite.NewMetricsRoutingOp(monClient)
	routings, err := op.List(ctx, monitoringsuite.MetricsRoutingsListParams{})
	if err != nil {
		slog.Error("Failed to fetch metrics routings", slog.Any("error", err))
		return nil, err
	}

	result := make([]MonitoringMetricsRouting, len(routings))
	for i, routing := range routings {
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

	op := monitoringsuite.NewLogsStorageOp(monClient)
	storage, err := op.Read(ctx, resourceID)
	if err != nil {
		slog.Error("Failed to fetch log storage detail", slog.Any("error", err))
		return nil, err
	}

	// Also fetch routings
	routings, _ := c.ListMonitoringLogRoutings(ctx)

	// Filter routings for this storage
	var filteredRoutings []MonitoringLogRouting
	for _, r := range routings {
		// Use embedded LogStorage's ResourceID
		storageResourceID := getNilInt64AsString(r.LogStorage.ResourceID)
		if storageResourceID == resourceID {
			filteredRoutings = append(filteredRoutings, r)
		}
	}

	detail := &MonitoringLogStorageDetail{
		LogStorage: *storage,
		Routings:   filteredRoutings,
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

	op := monitoringsuite.NewMetricsStorageOp(monClient)
	storage, err := op.Read(ctx, resourceID)
	if err != nil {
		slog.Error("Failed to fetch metrics storage detail", slog.Any("error", err))
		return nil, err
	}

	// Also fetch routings
	routings, _ := c.ListMonitoringMetricsRoutings(ctx)

	// Filter routings for this storage
	var filteredRoutings []MonitoringMetricsRouting
	for _, r := range routings {
		// Use embedded MetricsStorage's ResourceID
		storageResourceID := getNilInt64AsString(r.MetricsStorage.ResourceID)
		if storageResourceID == resourceID {
			filteredRoutings = append(filteredRoutings, r)
		}
	}

	detail := &MonitoringMetricsStorageDetail{
		MetricsStorage: *storage,
		Routings:       filteredRoutings,
	}

	slog.Info("Successfully fetched metrics storage detail", slog.String("resourceID", resourceID))
	return detail, nil
}

// GetMonitoringTraceStorageDetail fetches detail for a trace storage
func (c *SakuraClient) GetMonitoringTraceStorageDetail(ctx context.Context, resourceID string) (*MonitoringTraceStorageDetail, error) {
	slog.Info("Fetching Trace Storage detail", slog.String("resourceID", resourceID))

	monClient, err := c.getMonitoringClient()
	if err != nil {
		return nil, err
	}

	op := monitoringsuite.NewTracesStorageOp(monClient)
	storage, err := op.Read(ctx, resourceID)
	if err != nil {
		slog.Error("Failed to fetch trace storage detail", slog.Any("error", err))
		return nil, err
	}

	detail := &MonitoringTraceStorageDetail{
		TraceStorage: *storage,
	}

	slog.Info("Successfully fetched trace storage detail", slog.String("resourceID", resourceID))
	return detail, nil
}
