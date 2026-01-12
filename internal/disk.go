package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// Disk represents a disk resource
type Disk struct {
	ID         string
	Name       string
	Desc       string
	Zone       string
	SizeGB     int
	Connection string
	ServerID   string
	ServerName string
	CreatedAt  string
}

type DiskDetail struct {
	Disk
	Tags            []string
	DiskPlanName    string
	SourceDiskID    string
	SourceArchiveID string
	Availability    string
	EncryptionAlgo  string
	CreatedAt       string
	ModifiedAt      string
}

// Implement list.Item interface for Disk
func (d Disk) FilterValue() string {
	return d.Name
}

func (d Disk) Title() string {
	return d.Name
}

func (d Disk) Description() string {
	desc := fmt.Sprintf("ID: %s", d.ID)
	if d.Desc != "" {
		desc += " | " + d.Desc
	}
	return desc
}

func (c *SakuraClient) ListDisks(ctx context.Context) ([]Disk, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching disks from Sakura Cloud",
		slog.String("zone", c.zone))

	diskOp := iaas.NewDiskOp(c.caller)

	searched, err := diskOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch disks",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	disks := make([]Disk, 0, len(searched.Disks))
	for _, d := range searched.Disks {
		// Get connection type
		connection := string(d.Connection)

		// Get server info
		serverID := ""
		serverName := ""
		if !d.ServerID.IsEmpty() {
			serverID = d.ServerID.String()
			serverName = d.ServerName
		}

		// Format created at
		createdAt := ""
		if !d.CreatedAt.IsZero() {
			createdAt = d.CreatedAt.Format("2006-01-02")
		}

		disks = append(disks, Disk{
			ID:         d.ID.String(),
			Name:       d.Name,
			Desc:       d.Description,
			Zone:       c.zone,
			SizeGB:     d.SizeMB / 1024,
			Connection: connection,
			ServerID:   serverID,
			ServerName: serverName,
			CreatedAt:  createdAt,
		})
	}

	slog.Info("Successfully fetched disks",
		slog.String("zone", c.zone),
		slog.Int("count", len(disks)))

	return disks, nil
}

func (c *SakuraClient) GetDiskDetail(ctx context.Context, diskID string) (*DiskDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching disk detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("diskID", diskID))

	diskOp := iaas.NewDiskOp(c.caller)

	id := types.StringID(diskID)

	d, err := diskOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch disk detail",
			slog.String("zone", c.zone),
			slog.String("diskID", diskID),
			slog.Any("error", err))
		return nil, err
	}

	// Get connection type
	connection := string(d.Connection)

	// Get server info
	serverID := ""
	serverName := ""
	if !d.ServerID.IsEmpty() {
		serverID = d.ServerID.String()
		serverName = d.ServerName
	}

	// Format created at
	createdAt := ""
	if !d.CreatedAt.IsZero() {
		createdAt = d.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Format modified at
	modifiedAt := ""
	if !d.ModifiedAt.IsZero() {
		modifiedAt = d.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Get source info
	sourceDiskID := ""
	if !d.SourceDiskID.IsEmpty() {
		sourceDiskID = d.SourceDiskID.String()
	}
	sourceArchiveID := ""
	if !d.SourceArchiveID.IsEmpty() {
		sourceArchiveID = d.SourceArchiveID.String()
	}

	detail := &DiskDetail{
		Disk: Disk{
			ID:         d.ID.String(),
			Name:       d.Name,
			Desc:       d.Description,
			Zone:       c.zone,
			SizeGB:     d.SizeMB / 1024,
			Connection: connection,
			ServerID:   serverID,
			ServerName: serverName,
			CreatedAt:  createdAt,
		},
		Tags:            d.Tags,
		DiskPlanName:    d.DiskPlanName,
		SourceDiskID:    sourceDiskID,
		SourceArchiveID: sourceArchiveID,
		Availability:    string(d.Availability),
		EncryptionAlgo:  string(d.EncryptionAlgorithm),
		CreatedAt:       createdAt,
		ModifiedAt:      modifiedAt,
	}

	slog.Info("Successfully fetched disk detail",
		slog.String("zone", c.zone),
		slog.String("diskID", diskID))

	return detail, nil
}
