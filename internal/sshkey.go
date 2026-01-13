package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// SSHKey represents a SSH public key resource
type SSHKey struct {
	ID          string
	Name        string
	Desc        string
	Fingerprint string
	CreatedAt   string
}

type SSHKeyDetail struct {
	SSHKey
	PublicKey string
	CreatedAt string
}

// Implement list.Item interface for SSHKey
func (s SSHKey) FilterValue() string {
	return s.Name
}

func (s SSHKey) Title() string {
	return s.Name
}

func (s SSHKey) Description() string {
	desc := fmt.Sprintf("ID: %s", s.ID)
	if s.Desc != "" {
		desc += " | " + s.Desc
	}
	return desc
}

func (c *SakuraClient) ListSSHKeys(ctx context.Context) ([]SSHKey, error) {
	slog.Info("Fetching SSH keys from Sakura Cloud")

	sshKeyOp := iaas.NewSSHKeyOp(c.caller)

	searched, err := sshKeyOp.Find(ctx, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch SSH keys",
			slog.Any("error", err))
		return nil, err
	}

	sshKeys := make([]SSHKey, 0, len(searched.SSHKeys))
	for _, key := range searched.SSHKeys {
		// Format created at
		createdAt := ""
		if !key.CreatedAt.IsZero() {
			createdAt = key.CreatedAt.Format("2006-01-02")
		}

		// Truncate fingerprint for display
		fingerprint := key.Fingerprint
		if len(fingerprint) > 30 {
			fingerprint = fingerprint[:27] + "..."
		}

		sshKeys = append(sshKeys, SSHKey{
			ID:          key.ID.String(),
			Name:        key.Name,
			Desc:        key.Description,
			Fingerprint: fingerprint,
			CreatedAt:   createdAt,
		})
	}

	slog.Info("Successfully fetched SSH keys",
		slog.Int("count", len(sshKeys)))

	return sshKeys, nil
}

func (c *SakuraClient) GetSSHKeyDetail(ctx context.Context, sshKeyID string) (*SSHKeyDetail, error) {
	slog.Info("Fetching SSH key detail from Sakura Cloud",
		slog.String("sshKeyID", sshKeyID))

	sshKeyOp := iaas.NewSSHKeyOp(c.caller)

	id := types.StringID(sshKeyID)

	key, err := sshKeyOp.Read(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch SSH key detail",
			slog.String("sshKeyID", sshKeyID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !key.CreatedAt.IsZero() {
		createdAt = key.CreatedAt.Format("2006-01-02 15:04:05")
	}

	detail := &SSHKeyDetail{
		SSHKey: SSHKey{
			ID:          key.ID.String(),
			Name:        key.Name,
			Desc:        key.Description,
			Fingerprint: key.Fingerprint,
			CreatedAt:   createdAt,
		},
		PublicKey: key.PublicKey,
		CreatedAt: createdAt,
	}

	slog.Info("Successfully fetched SSH key detail",
		slog.String("sshKeyID", sshKeyID))

	return detail, nil
}
