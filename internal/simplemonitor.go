package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// SimpleMonitor represents a simple monitor resource
type SimpleMonitor struct {
	ID           string
	Name         string
	Desc         string
	Target       string
	Protocol     string
	Enabled      bool
	Health       string
	Availability string
}

type SimpleMonitorDetail struct {
	SimpleMonitor
	Tags                []string
	DelayLoop           int
	MaxCheckAttempts    int
	RetryInterval       int
	Timeout             int
	NotifyEmailEnabled  bool
	NotifySlackEnabled  bool
	NotifyInterval      int
	Port                int
	Path                string
	Host                string
	ContainsString      string
	CreatedAt           string
	ModifiedAt          string
	LastCheckedAt       string
	LastHealthChangedAt string
	LatestLogs          []string
}

// Implement list.Item interface for SimpleMonitor
func (s SimpleMonitor) FilterValue() string {
	return s.Name
}

func (s SimpleMonitor) Title() string {
	return s.Name
}

func (s SimpleMonitor) Description() string {
	desc := fmt.Sprintf("Target: %s", s.Target)
	if s.Desc != "" {
		desc += " | " + s.Desc
	}
	return desc
}

func (c *SakuraClient) ListSimpleMonitors(ctx context.Context) ([]SimpleMonitor, error) {
	slog.Info("Fetching simple monitors from Sakura Cloud")

	simpleMonitorOp := iaas.NewSimpleMonitorOp(c.caller)

	searched, err := simpleMonitorOp.Find(ctx, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch simple monitors",
			slog.Any("error", err))
		return nil, err
	}

	simpleMonitors := make([]SimpleMonitor, 0, len(searched.SimpleMonitors))
	for _, sm := range searched.SimpleMonitors {
		// Get protocol
		protocol := ""
		if sm.HealthCheck != nil {
			protocol = string(sm.HealthCheck.Protocol)
		}

		simpleMonitors = append(simpleMonitors, SimpleMonitor{
			ID:           sm.ID.String(),
			Name:         sm.Name,
			Desc:         sm.Description,
			Target:       sm.Target,
			Protocol:     protocol,
			Enabled:      bool(sm.Enabled),
			Health:       "",
			Availability: string(sm.Availability),
		})
	}

	slog.Info("Successfully fetched simple monitors",
		slog.Int("count", len(simpleMonitors)))

	return simpleMonitors, nil
}

func (c *SakuraClient) GetSimpleMonitorDetail(ctx context.Context, simpleMonitorID string) (*SimpleMonitorDetail, error) {
	slog.Info("Fetching simple monitor detail from Sakura Cloud",
		slog.String("simpleMonitorID", simpleMonitorID))

	simpleMonitorOp := iaas.NewSimpleMonitorOp(c.caller)

	id := types.StringID(simpleMonitorID)

	sm, err := simpleMonitorOp.Read(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch simple monitor detail",
			slog.String("simpleMonitorID", simpleMonitorID),
			slog.Any("error", err))
		return nil, err
	}

	// Get health status
	healthStatus, err := simpleMonitorOp.HealthStatus(ctx, id)
	if err != nil {
		slog.Warn("Failed to fetch simple monitor health status",
			slog.String("simpleMonitorID", simpleMonitorID),
			slog.Any("error", err))
		// Continue without health status
	}

	// Format created at
	createdAt := ""
	if !sm.CreatedAt.IsZero() {
		createdAt = sm.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Format modified at
	modifiedAt := ""
	if !sm.ModifiedAt.IsZero() {
		modifiedAt = sm.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Get protocol and health check details
	protocol := ""
	port := 0
	path := ""
	host := ""
	containsString := ""
	if sm.HealthCheck != nil {
		protocol = string(sm.HealthCheck.Protocol)
		port = int(sm.HealthCheck.Port)
		path = sm.HealthCheck.Path
		host = sm.HealthCheck.Host
		containsString = sm.HealthCheck.ContainsString
	}

	// Convert tags
	tags := make([]string, 0, len(sm.Tags))
	tags = append(tags, sm.Tags...)

	detail := &SimpleMonitorDetail{
		SimpleMonitor: SimpleMonitor{
			ID:           sm.ID.String(),
			Name:         sm.Name,
			Desc:         sm.Description,
			Target:       sm.Target,
			Protocol:     protocol,
			Enabled:      bool(sm.Enabled),
			Health:       "",
			Availability: string(sm.Availability),
		},
		Tags:               tags,
		DelayLoop:          sm.DelayLoop,
		MaxCheckAttempts:   sm.MaxCheckAttempts,
		RetryInterval:      sm.RetryInterval,
		Timeout:            sm.Timeout,
		NotifyEmailEnabled: bool(sm.NotifyEmailEnabled),
		NotifySlackEnabled: bool(sm.NotifySlackEnabled),
		NotifyInterval:     sm.NotifyInterval,
		Port:               port,
		Path:               path,
		Host:               host,
		ContainsString:     containsString,
		CreatedAt:          createdAt,
		ModifiedAt:         modifiedAt,
	}

	// Add health status if available
	if healthStatus != nil {
		detail.Health = string(healthStatus.Health)
		if !healthStatus.LastCheckedAt.IsZero() {
			detail.LastCheckedAt = healthStatus.LastCheckedAt.Format("2006-01-02 15:04:05")
		}
		if !healthStatus.LastHealthChangedAt.IsZero() {
			detail.LastHealthChangedAt = healthStatus.LastHealthChangedAt.Format("2006-01-02 15:04:05")
		}
		detail.LatestLogs = healthStatus.LatestLogs
	}

	slog.Info("Successfully fetched simple monitor detail",
		slog.String("simpleMonitorID", simpleMonitorID))

	return detail, nil
}
