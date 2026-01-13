package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// AutoBackup represents an auto backup resource
type AutoBackup struct {
	ID           string
	Name         string
	Desc         string
	Zone         string
	DiskID       string
	MaxBackups   int
	Weekdays     string
	Availability string
	CreatedAt    string
}

type AutoBackupDetail struct {
	AutoBackup
	Tags           []string
	BackupWeekdays []string
	AccountID      string
	ZoneName       string
	CreatedAt      string
	ModifiedAt     string
}

// Implement list.Item interface for AutoBackup
func (a AutoBackup) FilterValue() string {
	return a.Name
}

func (a AutoBackup) Title() string {
	return a.Name
}

func (a AutoBackup) Description() string {
	desc := fmt.Sprintf("ID: %s", a.ID)
	if a.Desc != "" {
		desc += " | " + a.Desc
	}
	return desc
}

func (c *SakuraClient) ListAutoBackups(ctx context.Context) ([]AutoBackup, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching auto backups from Sakura Cloud",
		slog.String("zone", c.zone))

	autoBackupOp := iaas.NewAutoBackupOp(c.caller)

	searched, err := autoBackupOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch auto backups",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	autoBackups := make([]AutoBackup, 0, len(searched.AutoBackups))
	for _, ab := range searched.AutoBackups {
		// Format created at
		createdAt := ""
		if !ab.CreatedAt.IsZero() {
			createdAt = ab.CreatedAt.Format("2006-01-02")
		}

		// Get backup weekdays
		weekdays := formatWeekdays(ab.BackupSpanWeekdays)

		// Get disk ID
		diskID := ""
		if ab.DiskID != 0 {
			diskID = ab.DiskID.String()
		}

		autoBackups = append(autoBackups, AutoBackup{
			ID:           ab.ID.String(),
			Name:         ab.Name,
			Desc:         ab.Description,
			Zone:         c.zone,
			DiskID:       diskID,
			MaxBackups:   ab.MaximumNumberOfArchives,
			Weekdays:     weekdays,
			Availability: string(ab.Availability),
			CreatedAt:    createdAt,
		})
	}

	slog.Info("Successfully fetched auto backups",
		slog.String("zone", c.zone),
		slog.Int("count", len(autoBackups)))

	return autoBackups, nil
}

func (c *SakuraClient) GetAutoBackupDetail(ctx context.Context, autoBackupID string) (*AutoBackupDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching auto backup detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("autoBackupID", autoBackupID))

	autoBackupOp := iaas.NewAutoBackupOp(c.caller)

	id := types.StringID(autoBackupID)

	ab, err := autoBackupOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch auto backup detail",
			slog.String("zone", c.zone),
			slog.String("autoBackupID", autoBackupID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !ab.CreatedAt.IsZero() {
		createdAt = ab.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Format modified at
	modifiedAt := ""
	if !ab.ModifiedAt.IsZero() {
		modifiedAt = ab.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Get backup weekdays as strings
	weekdayStrs := make([]string, 0, len(ab.BackupSpanWeekdays))
	for _, day := range ab.BackupSpanWeekdays {
		weekdayStrs = append(weekdayStrs, string(day))
	}

	// Convert tags
	tags := make([]string, 0, len(ab.Tags))
	tags = append(tags, ab.Tags...)

	detail := &AutoBackupDetail{
		AutoBackup: AutoBackup{
			ID:           ab.ID.String(),
			Name:         ab.Name,
			Desc:         ab.Description,
			Zone:         c.zone,
			DiskID:       ab.DiskID.String(),
			MaxBackups:   ab.MaximumNumberOfArchives,
			Weekdays:     formatWeekdays(ab.BackupSpanWeekdays),
			Availability: string(ab.Availability),
			CreatedAt:    createdAt,
		},
		Tags:           tags,
		BackupWeekdays: weekdayStrs,
		AccountID:      ab.AccountID.String(),
		ZoneName:       ab.ZoneName,
		CreatedAt:      createdAt,
		ModifiedAt:     modifiedAt,
	}

	slog.Info("Successfully fetched auto backup detail",
		slog.String("zone", c.zone),
		slog.String("autoBackupID", autoBackupID))

	return detail, nil
}

func formatWeekdays(weekdays []types.EDayOfTheWeek) string {
	if len(weekdays) == 0 {
		return "-"
	}
	strs := make([]string, 0, len(weekdays))
	for _, day := range weekdays {
		strs = append(strs, string(day))
	}
	return strings.Join(strs, ", ")
}
