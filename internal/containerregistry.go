package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// ContainerRegistry represents a container registry resource
type ContainerRegistry struct {
	ID          string
	Name        string
	Desc        string
	FQDN        string
	AccessLevel string
	UserCount   int
}

type ContainerRegistryUser struct {
	UserName   string
	Permission string
}

type ContainerRegistryDetail struct {
	ContainerRegistry
	Tags           []string
	VirtualDomain  string
	SubDomainLabel string
	Availability   string
	Users          []ContainerRegistryUser
	CreatedAt      string
	ModifiedAt     string
}

// Implement list.Item interface for ContainerRegistry
func (cr ContainerRegistry) FilterValue() string {
	return cr.Name
}

func (cr ContainerRegistry) Title() string {
	return cr.Name
}

func (cr ContainerRegistry) Description() string {
	desc := fmt.Sprintf("FQDN: %s", cr.FQDN)
	if cr.Desc != "" {
		desc += " | " + cr.Desc
	}
	return desc
}

func (c *SakuraClient) ListContainerRegistries(ctx context.Context) ([]ContainerRegistry, error) {
	slog.Info("Fetching container registries from Sakura Cloud")

	containerRegistryOp := iaas.NewContainerRegistryOp(c.caller)

	searched, err := containerRegistryOp.Find(ctx, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch container registries",
			slog.Any("error", err))
		return nil, err
	}

	containerRegistries := make([]ContainerRegistry, 0, len(searched.ContainerRegistries))
	for _, cr := range searched.ContainerRegistries {
		// Get user count
		userCount := 0
		users, err := containerRegistryOp.ListUsers(ctx, cr.ID)
		if err == nil && users != nil {
			userCount = len(users.Users)
		}

		containerRegistries = append(containerRegistries, ContainerRegistry{
			ID:          cr.ID.String(),
			Name:        cr.Name,
			Desc:        cr.Description,
			FQDN:        cr.FQDN,
			AccessLevel: string(cr.AccessLevel),
			UserCount:   userCount,
		})
	}

	slog.Info("Successfully fetched container registries",
		slog.Int("count", len(containerRegistries)))

	return containerRegistries, nil
}

func (c *SakuraClient) GetContainerRegistryDetail(ctx context.Context, containerRegistryID string) (*ContainerRegistryDetail, error) {
	slog.Info("Fetching container registry detail from Sakura Cloud",
		slog.String("containerRegistryID", containerRegistryID))

	containerRegistryOp := iaas.NewContainerRegistryOp(c.caller)

	id := types.StringID(containerRegistryID)

	cr, err := containerRegistryOp.Read(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch container registry detail",
			slog.String("containerRegistryID", containerRegistryID),
			slog.Any("error", err))
		return nil, err
	}

	// Get users
	users := make([]ContainerRegistryUser, 0)
	userList, err := containerRegistryOp.ListUsers(ctx, id)
	if err != nil {
		slog.Warn("Failed to fetch container registry users",
			slog.String("containerRegistryID", containerRegistryID),
			slog.Any("error", err))
	} else if userList != nil {
		for _, u := range userList.Users {
			users = append(users, ContainerRegistryUser{
				UserName:   u.UserName,
				Permission: string(u.Permission),
			})
		}
	}

	// Format created at
	createdAt := ""
	if !cr.CreatedAt.IsZero() {
		createdAt = cr.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Format modified at
	modifiedAt := ""
	if !cr.ModifiedAt.IsZero() {
		modifiedAt = cr.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Convert tags
	tags := make([]string, 0, len(cr.Tags))
	tags = append(tags, cr.Tags...)

	detail := &ContainerRegistryDetail{
		ContainerRegistry: ContainerRegistry{
			ID:          cr.ID.String(),
			Name:        cr.Name,
			Desc:        cr.Description,
			FQDN:        cr.FQDN,
			AccessLevel: string(cr.AccessLevel),
			UserCount:   len(users),
		},
		Tags:           tags,
		VirtualDomain:  cr.VirtualDomain,
		SubDomainLabel: cr.SubDomainLabel,
		Availability:   string(cr.Availability),
		Users:          users,
		CreatedAt:      createdAt,
		ModifiedAt:     modifiedAt,
	}

	slog.Info("Successfully fetched container registry detail",
		slog.String("containerRegistryID", containerRegistryID))

	return detail, nil
}
