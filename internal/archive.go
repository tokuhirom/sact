package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// Archive represents an archive resource
type Archive struct {
	ID           string
	Name         string
	Desc         string
	Zone         string
	SizeGB       int
	Scope        string
	Availability string
	CreatedAt    string
}

type ArchiveDetail struct {
	Archive
	Tags            []string
	SourceDiskID    string
	SourceArchiveID string
	BundleInfo      string
	ModifiedAt      string
}

// Implement list.Item interface for Archive
func (a Archive) FilterValue() string {
	return a.Name
}

func (a Archive) Title() string {
	return a.Name
}

func (a Archive) Description() string {
	desc := fmt.Sprintf("ID: %s", a.ID)
	if a.Desc != "" {
		desc += " | " + a.Desc
	}
	return desc
}

func (c *SakuraClient) ListArchives(ctx context.Context) ([]Archive, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching archives from Sakura Cloud",
		slog.String("zone", c.zone))

	archiveOp := iaas.NewArchiveOp(c.caller)

	searched, err := archiveOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch archives",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	archives := make([]Archive, 0, len(searched.Archives))
	for _, a := range searched.Archives {
		// Format created at
		createdAt := ""
		if !a.CreatedAt.IsZero() {
			createdAt = a.CreatedAt.Format("2006-01-02")
		}

		archives = append(archives, Archive{
			ID:           a.ID.String(),
			Name:         a.Name,
			Desc:         a.Description,
			Zone:         c.zone,
			SizeGB:       a.SizeMB / 1024,
			Scope:        string(a.Scope),
			Availability: string(a.Availability),
			CreatedAt:    createdAt,
		})
	}

	slog.Info("Successfully fetched archives",
		slog.String("zone", c.zone),
		slog.Int("count", len(archives)))

	return archives, nil
}

func (c *SakuraClient) GetArchiveDetail(ctx context.Context, archiveID string) (*ArchiveDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching archive detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("archiveID", archiveID))

	archiveOp := iaas.NewArchiveOp(c.caller)

	id := types.StringID(archiveID)

	a, err := archiveOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch archive detail",
			slog.String("zone", c.zone),
			slog.String("archiveID", archiveID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !a.CreatedAt.IsZero() {
		createdAt = a.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Format modified at
	modifiedAt := ""
	if !a.ModifiedAt.IsZero() {
		modifiedAt = a.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Get source info
	sourceDiskID := ""
	if !a.SourceDiskID.IsEmpty() {
		sourceDiskID = a.SourceDiskID.String()
	}
	sourceArchiveID := ""
	if !a.SourceArchiveID.IsEmpty() {
		sourceArchiveID = a.SourceArchiveID.String()
	}

	// Get bundle info
	bundleInfo := ""
	if a.BundleInfo != nil && a.BundleInfo.HostClass != "" {
		bundleInfo = a.BundleInfo.HostClass
	}

	detail := &ArchiveDetail{
		Archive: Archive{
			ID:           a.ID.String(),
			Name:         a.Name,
			Desc:         a.Description,
			Zone:         c.zone,
			SizeGB:       a.SizeMB / 1024,
			Scope:        string(a.Scope),
			Availability: string(a.Availability),
			CreatedAt:    createdAt,
		},
		Tags:            a.Tags,
		SourceDiskID:    sourceDiskID,
		SourceArchiveID: sourceArchiveID,
		BundleInfo:      bundleInfo,
		ModifiedAt:      modifiedAt,
	}

	slog.Info("Successfully fetched archive detail",
		slog.String("zone", c.zone),
		slog.String("archiveID", archiveID))

	return detail, nil
}
