package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	apprun "github.com/tokuhirom/sact/pkg/openapi/apprun_dedicated"
)

// AppRunCluster represents an AppRun Dedicated cluster resource
type AppRunCluster struct {
	ID               string
	Name             string
	ASGCount         int
	HasLetsEncrypt   bool
	ServicePrincipal string
	CreatedAt        string
}

// AppRunClusterDetail contains detailed information about an AppRun cluster
type AppRunClusterDetail struct {
	AppRunCluster
	Ports             []AppRunPort
	AutoScalingGroups []AppRunASG
}

// AppRunPort represents a load balancer port
type AppRunPort struct {
	Port     uint16
	Protocol string
}

// AppRunASG represents an auto scaling group summary
type AppRunASG struct {
	ID   string
	Name string
}

// Implement list.Item interface for AppRunCluster
func (c AppRunCluster) FilterValue() string {
	return c.Name
}

func (c AppRunCluster) Title() string {
	return c.Name
}

func (c AppRunCluster) Description() string {
	return fmt.Sprintf("ASG: %d | Created: %s", c.ASGCount, c.CreatedAt)
}

// ListAppRunClusters fetches all AppRun Dedicated clusters
func (c *SakuraClient) ListAppRunClusters(ctx context.Context) ([]AppRunCluster, error) {
	slog.Info("Fetching AppRun Dedicated clusters")

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	var allClusters []AppRunCluster
	var cursor apprun.OptClusterID

	for {
		params := apprun.ListClustersParams{
			MaxItems: 100,
			Cursor:   cursor,
		}

		resp, err := client.ListClusters(ctx, params)
		if err != nil {
			slog.Error("Failed to fetch AppRun clusters", slog.Any("error", err))
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			createdAt := ""
			if cluster.Created > 0 {
				createdAt = time.Unix(int64(cluster.Created), 0).Format("2006-01-02 15:04:05")
			}

			allClusters = append(allClusters, AppRunCluster{
				ID:               uuid.UUID(cluster.ClusterID).String(),
				Name:             cluster.Name,
				ASGCount:         len(cluster.AutoScalingGroups),
				HasLetsEncrypt:   cluster.HasLetsEncryptEmail,
				ServicePrincipal: cluster.ServicePrincipalID,
				CreatedAt:        createdAt,
			})
		}

		if !resp.NextCursor.Set {
			break
		}
		cursor = resp.NextCursor
	}

	slog.Info("Successfully fetched AppRun clusters", slog.Int("count", len(allClusters)))
	return allClusters, nil
}

// GetAppRunClusterDetail fetches detailed information about a specific cluster
func (c *SakuraClient) GetAppRunClusterDetail(ctx context.Context, clusterID string) (*AppRunClusterDetail, error) {
	slog.Info("Fetching AppRun cluster detail", slog.String("clusterID", clusterID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsed, err := uuid.Parse(clusterID)
	if err != nil {
		slog.Error("Invalid cluster ID", slog.String("clusterID", clusterID), slog.Any("error", err))
		return nil, fmt.Errorf("invalid cluster ID: %w", err)
	}

	params := apprun.GetClusterParams{
		ClusterID: apprun.ClusterID(parsed),
	}

	resp, err := client.GetCluster(ctx, params)
	if err != nil {
		slog.Error("Failed to fetch AppRun cluster detail", slog.String("clusterID", clusterID), slog.Any("error", err))
		return nil, err
	}

	cluster := resp.Cluster

	createdAt := ""
	if cluster.Created > 0 {
		createdAt = time.Unix(int64(cluster.Created), 0).Format("2006-01-02 15:04:05")
	}

	// Convert ports
	ports := make([]AppRunPort, 0, len(cluster.Ports))
	for _, p := range cluster.Ports {
		ports = append(ports, AppRunPort{
			Port:     p.Port,
			Protocol: string(p.Protocol),
		})
	}

	// Convert ASGs
	asgs := make([]AppRunASG, 0, len(cluster.AutoScalingGroups))
	for _, asg := range cluster.AutoScalingGroups {
		asgs = append(asgs, AppRunASG{
			ID:   uuid.UUID(asg.AutoScalingGroupID).String(),
			Name: asg.Name,
		})
	}

	detail := &AppRunClusterDetail{
		AppRunCluster: AppRunCluster{
			ID:               uuid.UUID(cluster.ClusterID).String(),
			Name:             cluster.Name,
			ASGCount:         len(cluster.AutoScalingGroups),
			HasLetsEncrypt:   cluster.HasLetsEncryptEmail,
			ServicePrincipal: cluster.ServicePrincipalID,
			CreatedAt:        createdAt,
		},
		Ports:             ports,
		AutoScalingGroups: asgs,
	}

	slog.Info("Successfully fetched AppRun cluster detail", slog.String("clusterID", clusterID))
	return detail, nil
}
