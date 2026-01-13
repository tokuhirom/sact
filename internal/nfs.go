package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// NFS represents a NFS appliance resource
type NFS struct {
	ID             string
	Name           string
	Desc           string
	Zone           string
	InstanceStatus string
	SwitchName     string
	CreatedAt      string
}

type NFSDetail struct {
	NFS
	Tags           []string
	PlanID         string
	SwitchID       string
	DefaultRoute   string
	NetworkMaskLen int
	IPAddresses    []string
	CreatedAt      string
}

// Implement list.Item interface for NFS
func (n NFS) FilterValue() string {
	return n.Name
}

func (n NFS) Title() string {
	return n.Name
}

func (n NFS) Description() string {
	desc := fmt.Sprintf("ID: %s", n.ID)
	if n.Desc != "" {
		desc += " | " + n.Desc
	}
	return desc
}

func (c *SakuraClient) ListNFS(ctx context.Context) ([]NFS, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching NFS appliances from Sakura Cloud",
		slog.String("zone", c.zone))

	nfsOp := iaas.NewNFSOp(c.caller)

	searched, err := nfsOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch NFS appliances",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	nfsList := make([]NFS, 0, len(searched.NFS))
	for _, nfs := range searched.NFS {
		// Format created at
		createdAt := ""
		if !nfs.CreatedAt.IsZero() {
			createdAt = nfs.CreatedAt.Format("2006-01-02")
		}

		nfsList = append(nfsList, NFS{
			ID:             nfs.ID.String(),
			Name:           nfs.Name,
			Desc:           nfs.Description,
			Zone:           c.zone,
			InstanceStatus: string(nfs.InstanceStatus),
			SwitchName:     nfs.SwitchName,
			CreatedAt:      createdAt,
		})
	}

	slog.Info("Successfully fetched NFS appliances",
		slog.String("zone", c.zone),
		slog.Int("count", len(nfsList)))

	return nfsList, nil
}

func (c *SakuraClient) GetNFSDetail(ctx context.Context, nfsID string) (*NFSDetail, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching NFS detail from Sakura Cloud",
		slog.String("zone", c.zone),
		slog.String("nfsID", nfsID))

	nfsOp := iaas.NewNFSOp(c.caller)

	id := types.StringID(nfsID)

	nfs, err := nfsOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch NFS detail",
			slog.String("zone", c.zone),
			slog.String("nfsID", nfsID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !nfs.CreatedAt.IsZero() {
		createdAt = nfs.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Convert tags
	tags := make([]string, 0, len(nfs.Tags))
	tags = append(tags, nfs.Tags...)

	detail := &NFSDetail{
		NFS: NFS{
			ID:             nfs.ID.String(),
			Name:           nfs.Name,
			Desc:           nfs.Description,
			Zone:           c.zone,
			InstanceStatus: string(nfs.InstanceStatus),
			SwitchName:     nfs.SwitchName,
			CreatedAt:      createdAt,
		},
		Tags:           tags,
		PlanID:         nfs.PlanID.String(),
		SwitchID:       nfs.SwitchID.String(),
		DefaultRoute:   nfs.DefaultRoute,
		NetworkMaskLen: nfs.NetworkMaskLen,
		IPAddresses:    nfs.IPAddresses,
		CreatedAt:      createdAt,
	}

	slog.Info("Successfully fetched NFS detail",
		slog.String("zone", c.zone),
		slog.String("nfsID", nfsID))

	return detail, nil
}
