package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	client "github.com/sacloud/api-client-go"
)

// MonitoringLogStorage represents a log storage
type MonitoringLogStorage struct {
	ID         int64  `json:"id"`
	ResourceID int64  `json:"resource_id"`
	Name       string `json:"name"`
	Desc       string `json:"description"`
	ExpireDay  int    `json:"expire_day"`
	IsSystem   bool   `json:"is_system"`
	CreatedAt  string `json:"created_at"`
	Endpoints  struct {
		Ingester struct {
			Address  string `json:"address"`
			Insecure bool   `json:"insecure"`
		} `json:"ingester"`
	} `json:"endpoints"`
	Usage struct {
		LogRoutings     int `json:"log_routings"`
		LogMeasureRules int `json:"log_measure_rules"`
	} `json:"usage"`
}

// MonitoringMetricsStorage represents a metrics storage
type MonitoringMetricsStorage struct {
	ID         int64  `json:"id"`
	ResourceID int64  `json:"resource_id"`
	Name       string `json:"name"`
	Desc       string `json:"description"`
	IsSystem   bool   `json:"is_system"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Endpoints  struct {
		Address string `json:"address"`
	} `json:"endpoints"`
	Usage struct {
		MetricsRoutings int `json:"metrics_routings"`
		AlertRules      int `json:"alert_rules"`
		LogMeasureRules int `json:"log_measure_rules"`
	} `json:"usage"`
}

// MonitoringLogRouting represents a log routing rule
type MonitoringLogRouting struct {
	UID              string `json:"uid"`
	Name             string `json:"name"`
	Enabled          bool   `json:"enabled"`
	SourceResourceID int64  `json:"source_resource_id"`
	DestResourceID   int64  `json:"dest_resource_id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// MonitoringMetricsRouting represents a metrics routing rule
type MonitoringMetricsRouting struct {
	UID              string `json:"uid"`
	Name             string `json:"name"`
	Enabled          bool   `json:"enabled"`
	SourceResourceID int64  `json:"source_resource_id"`
	DestResourceID   int64  `json:"dest_resource_id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Implement list.Item interface for MonitoringLogStorage
func (s MonitoringLogStorage) FilterValue() string {
	return s.Name
}

func (s MonitoringLogStorage) Title() string {
	return s.Name
}

func (s MonitoringLogStorage) Description() string {
	return fmt.Sprintf("ID: %d | Expire: %d days | Routings: %d", s.ResourceID, s.ExpireDay, s.Usage.LogRoutings)
}

// Implement list.Item interface for MonitoringMetricsStorage
func (s MonitoringMetricsStorage) FilterValue() string {
	return s.Name
}

func (s MonitoringMetricsStorage) Title() string {
	return s.Name
}

func (s MonitoringMetricsStorage) Description() string {
	return fmt.Sprintf("ID: %d | Routings: %d | Alert Rules: %d", s.ResourceID, s.Usage.MetricsRoutings, s.Usage.AlertRules)
}

// Implement list.Item interface for MonitoringLogRouting
func (r MonitoringLogRouting) FilterValue() string {
	return r.Name
}

func (r MonitoringLogRouting) Title() string {
	return r.Name
}

func (r MonitoringLogRouting) Description() string {
	enabled := "enabled"
	if !r.Enabled {
		enabled = "disabled"
	}
	return fmt.Sprintf("UID: %s | %s | Src: %d -> Dest: %d", r.UID[:8], enabled, r.SourceResourceID, r.DestResourceID)
}

// Implement list.Item interface for MonitoringMetricsRouting
func (r MonitoringMetricsRouting) FilterValue() string {
	return r.Name
}

func (r MonitoringMetricsRouting) Title() string {
	return r.Name
}

func (r MonitoringMetricsRouting) Description() string {
	enabled := "enabled"
	if !r.Enabled {
		enabled = "disabled"
	}
	return fmt.Sprintf("UID: %s | %s | Src: %d -> Dest: %d", r.UID[:8], enabled, r.SourceResourceID, r.DestResourceID)
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

// monitoringHTTPClient creates an HTTP client for monitoring suite API
func (c *SakuraClient) monitoringHTTPClient() (*http.Client, string, string, error) {
	clientOpts, err := client.DefaultOption()
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get credentials: %w", err)
	}

	token := clientOpts.AccessToken
	secret := clientOpts.AccessTokenSecret

	if token == "" || secret == "" {
		return nil, "", "", fmt.Errorf("credentials not found")
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return httpClient, token, secret, nil
}

// monitoringRequest makes a request to the monitoring suite API
func (c *SakuraClient) monitoringRequest(ctx context.Context, method, path string, result interface{}) error {
	httpClient, token, secret, err := c.monitoringHTTPClient()
	if err != nil {
		return err
	}

	// Use is1a zone for monitoring suite API
	baseURL := "https://secure.sakura.ad.jp/cloud/zone/is1a/api/monitoring/1.0"
	url := baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(token, secret)
	req.Header.Set("Accept", "application/json")

	slog.Debug("Monitoring API request", slog.String("method", method), slog.String("url", url))

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		slog.Error("Monitoring API error", slog.Int("status", resp.StatusCode), slog.String("body", string(body)))
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// ListMonitoringLogStorages fetches all log storages
func (c *SakuraClient) ListMonitoringLogStorages(ctx context.Context) ([]MonitoringLogStorage, error) {
	slog.Info("Fetching Monitoring Log Storages")

	var response struct {
		Count   int                    `json:"count"`
		Results []MonitoringLogStorage `json:"results"`
	}

	if err := c.monitoringRequest(ctx, "GET", "/logs/storages/", &response); err != nil {
		slog.Error("Failed to fetch log storages", slog.Any("error", err))
		return nil, err
	}

	slog.Info("Successfully fetched log storages", slog.Int("count", len(response.Results)))
	return response.Results, nil
}

// ListMonitoringMetricsStorages fetches all metrics storages
func (c *SakuraClient) ListMonitoringMetricsStorages(ctx context.Context) ([]MonitoringMetricsStorage, error) {
	slog.Info("Fetching Monitoring Metrics Storages")

	var response struct {
		Count   int                        `json:"count"`
		Results []MonitoringMetricsStorage `json:"results"`
	}

	if err := c.monitoringRequest(ctx, "GET", "/metrics/storages/", &response); err != nil {
		slog.Error("Failed to fetch metrics storages", slog.Any("error", err))
		return nil, err
	}

	slog.Info("Successfully fetched metrics storages", slog.Int("count", len(response.Results)))
	return response.Results, nil
}

// ListMonitoringLogRoutings fetches all log routing rules
func (c *SakuraClient) ListMonitoringLogRoutings(ctx context.Context) ([]MonitoringLogRouting, error) {
	slog.Info("Fetching Monitoring Log Routings")

	var response struct {
		Count   int                    `json:"count"`
		Results []MonitoringLogRouting `json:"results"`
	}

	if err := c.monitoringRequest(ctx, "GET", "/logs/routings/", &response); err != nil {
		slog.Error("Failed to fetch log routings", slog.Any("error", err))
		return nil, err
	}

	slog.Info("Successfully fetched log routings", slog.Int("count", len(response.Results)))
	return response.Results, nil
}

// ListMonitoringMetricsRoutings fetches all metrics routing rules
func (c *SakuraClient) ListMonitoringMetricsRoutings(ctx context.Context) ([]MonitoringMetricsRouting, error) {
	slog.Info("Fetching Monitoring Metrics Routings")

	var response struct {
		Count   int                        `json:"count"`
		Results []MonitoringMetricsRouting `json:"results"`
	}

	if err := c.monitoringRequest(ctx, "GET", "/metrics/routings/", &response); err != nil {
		slog.Error("Failed to fetch metrics routings", slog.Any("error", err))
		return nil, err
	}

	slog.Info("Successfully fetched metrics routings", slog.Int("count", len(response.Results)))
	return response.Results, nil
}

// GetMonitoringLogStorageDetail fetches detail for a log storage
func (c *SakuraClient) GetMonitoringLogStorageDetail(ctx context.Context, resourceID int64) (*MonitoringLogStorageDetail, error) {
	slog.Info("Fetching Log Storage detail", slog.Int64("resourceID", resourceID))

	var storage MonitoringLogStorage
	if err := c.monitoringRequest(ctx, "GET", fmt.Sprintf("/logs/storages/%d/", resourceID), &storage); err != nil {
		slog.Error("Failed to fetch log storage detail", slog.Any("error", err))
		return nil, err
	}

	// Also fetch routings
	routings, _ := c.ListMonitoringLogRoutings(ctx)

	// Filter routings for this storage
	var filteredRoutings []MonitoringLogRouting
	for _, r := range routings {
		if r.SourceResourceID == resourceID || r.DestResourceID == resourceID {
			filteredRoutings = append(filteredRoutings, r)
		}
	}

	detail := &MonitoringLogStorageDetail{
		MonitoringLogStorage: storage,
		Routings:             filteredRoutings,
	}

	slog.Info("Successfully fetched log storage detail", slog.Int64("resourceID", resourceID))
	return detail, nil
}

// GetMonitoringMetricsStorageDetail fetches detail for a metrics storage
func (c *SakuraClient) GetMonitoringMetricsStorageDetail(ctx context.Context, resourceID int64) (*MonitoringMetricsStorageDetail, error) {
	slog.Info("Fetching Metrics Storage detail", slog.Int64("resourceID", resourceID))

	var storage MonitoringMetricsStorage
	if err := c.monitoringRequest(ctx, "GET", fmt.Sprintf("/metrics/storages/%d/", resourceID), &storage); err != nil {
		slog.Error("Failed to fetch metrics storage detail", slog.Any("error", err))
		return nil, err
	}

	// Also fetch routings
	routings, _ := c.ListMonitoringMetricsRoutings(ctx)

	// Filter routings for this storage
	var filteredRoutings []MonitoringMetricsRouting
	for _, r := range routings {
		if r.SourceResourceID == resourceID || r.DestResourceID == resourceID {
			filteredRoutings = append(filteredRoutings, r)
		}
	}

	detail := &MonitoringMetricsStorageDetail{
		MonitoringMetricsStorage: storage,
		Routings:                 filteredRoutings,
	}

	slog.Info("Successfully fetched metrics storage detail", slog.Int64("resourceID", resourceID))
	return detail, nil
}
