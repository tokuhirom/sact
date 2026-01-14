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

// AppRunASG represents an auto scaling group for list display
type AppRunASG struct {
	ID              string
	Name            string
	Zone            string
	MinNodes        int32
	MaxNodes        int32
	WorkerNodeCount int32
	ServiceClass    string
	ClusterID       string // parent cluster ID for navigation
}

// AppRunASGInterface represents a network interface for ASG
type AppRunASGInterface struct {
	Index          int16
	Upstream       string
	NetmaskLen     int16
	DefaultGateway string
	PacketFilterID string
	ConnectsToLB   bool
	IPPool         []string
}

// AppRunASGDetail contains detailed information about an ASG
type AppRunASGDetail struct {
	AppRunASG
	NameServers []string
	Interfaces  []AppRunASGInterface
	Deleting    bool
}

// AppRunWorkerNode represents a worker node
type AppRunWorkerNode struct {
	ID             string
	ResourceID     string
	Status         string
	Draining       bool
	ArchiveVersion string
	CreatedAt      string
	ErrorMessage   string
	IPAddresses    []string // flattened from network interfaces
	ClusterID      string   // parent cluster ID
	ASGID          string   // parent ASG ID
}

// Implement list.Item interface for AppRunWorkerNode
func (w AppRunWorkerNode) FilterValue() string {
	return w.ResourceID
}

func (w AppRunWorkerNode) Title() string {
	return w.ResourceID
}

func (w AppRunWorkerNode) Description() string {
	status := w.Status
	if w.Draining {
		status += " (draining)"
	}
	return fmt.Sprintf("Status: %s | IPs: %s", status, joinStrings(w.IPAddresses, ", "))
}

func joinStrings(s []string, sep string) string {
	if len(s) == 0 {
		return "-"
	}
	result := s[0]
	for i := 1; i < len(s); i++ {
		result += sep + s[i]
	}
	return result
}

// AppRunLB represents a load balancer for list display
type AppRunLB struct {
	ID           string
	Name         string
	ServiceClass string
	CreatedAt    string
	ClusterID    string // parent cluster ID
	ASGID        string // parent ASG ID
}

// AppRunLBInterface represents a load balancer network interface
type AppRunLBInterface struct {
	Index          int16
	Upstream       string
	VIP            string
	DefaultGateway string
	NetmaskLen     int16
	IPPool         []string
}

// AppRunLBDetail contains detailed information about a load balancer
type AppRunLBDetail struct {
	AppRunLB
	NameServers []string
	Interfaces  []AppRunLBInterface
	Deleting    bool
}

// AppRunApplication represents an application for list display
type AppRunApplication struct {
	ID            string
	Name          string
	ClusterID     string
	ClusterName   string
	ActiveVersion int32
	DesiredCount  int32
}

// AppRunVersion represents an application version for list display
type AppRunVersion struct {
	Version         int32
	Image           string
	ActiveNodeCount int64
	CreatedAt       string
	ApplicationID   string // parent application ID
	ClusterID       string // parent cluster ID
	IsActive        bool   // true if this is the active version
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

// Implement list.Item interface for AppRunASG
func (a AppRunASG) FilterValue() string {
	return a.Name
}

func (a AppRunASG) Title() string {
	return a.Name
}

func (a AppRunASG) Description() string {
	return fmt.Sprintf("Zone: %s | Nodes: %d/%d-%d", a.Zone, a.WorkerNodeCount, a.MinNodes, a.MaxNodes)
}

// Implement list.Item interface for AppRunLB
func (l AppRunLB) FilterValue() string {
	return l.Name
}

func (l AppRunLB) Title() string {
	return l.Name
}

func (l AppRunLB) Description() string {
	return fmt.Sprintf("Service: %s | Created: %s", l.ServiceClass, l.CreatedAt)
}

// Implement list.Item interface for AppRunApplication
func (a AppRunApplication) FilterValue() string {
	return a.Name
}

func (a AppRunApplication) Title() string {
	return a.Name
}

func (a AppRunApplication) Description() string {
	versionStr := "-"
	if a.ActiveVersion > 0 {
		versionStr = fmt.Sprintf("v%d", a.ActiveVersion)
	}
	return fmt.Sprintf("Version: %s | DesiredCount: %d", versionStr, a.DesiredCount)
}

// Implement list.Item interface for AppRunVersion
func (v AppRunVersion) FilterValue() string {
	return fmt.Sprintf("v%d", v.Version)
}

func (v AppRunVersion) Title() string {
	return fmt.Sprintf("v%d", v.Version)
}

func (v AppRunVersion) Description() string {
	return fmt.Sprintf("Nodes: %d | %s", v.ActiveNodeCount, v.Image)
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
			MaxItems: 30, // API limit is 30
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

// ListAppRunASGs fetches all ASGs for a specific cluster
func (c *SakuraClient) ListAppRunASGs(ctx context.Context, clusterID string) ([]AppRunASG, error) {
	slog.Info("Fetching AppRun ASGs", slog.String("clusterID", clusterID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsed, err := uuid.Parse(clusterID)
	if err != nil {
		return nil, fmt.Errorf("invalid cluster ID: %w", err)
	}

	var allASGs []AppRunASG
	var cursor apprun.OptAutoScalingGroupID

	for {
		params := apprun.ListAutoScalingGroupsParams{
			ClusterID: apprun.ClusterID(parsed),
			MaxItems:  30,
			Cursor:    cursor,
		}

		resp, err := client.ListAutoScalingGroups(ctx, params)
		if err != nil {
			slog.Error("Failed to fetch AppRun ASGs", slog.Any("error", err))
			return nil, err
		}

		for _, asg := range resp.AutoScalingGroups {
			allASGs = append(allASGs, AppRunASG{
				ID:              uuid.UUID(asg.AutoScalingGroupID).String(),
				Name:            asg.Name,
				Zone:            asg.Zone,
				MinNodes:        asg.MinNodes,
				MaxNodes:        asg.MaxNodes,
				WorkerNodeCount: asg.WorkerNodeCount,
				ServiceClass:    asg.WorkerServiceClassPath,
				ClusterID:       clusterID,
			})
		}

		if !resp.NextCursor.Set {
			break
		}
		cursor = resp.NextCursor
	}

	slog.Info("Successfully fetched AppRun ASGs", slog.Int("count", len(allASGs)))
	return allASGs, nil
}

// GetAppRunASGDetail fetches detailed information about a specific ASG
func (c *SakuraClient) GetAppRunASGDetail(ctx context.Context, clusterID, asgID string) (*AppRunASGDetail, error) {
	slog.Info("Fetching AppRun ASG detail", slog.String("clusterID", clusterID), slog.String("asgID", asgID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsedCluster, err := uuid.Parse(clusterID)
	if err != nil {
		return nil, fmt.Errorf("invalid cluster ID: %w", err)
	}
	parsedASG, err := uuid.Parse(asgID)
	if err != nil {
		return nil, fmt.Errorf("invalid ASG ID: %w", err)
	}

	params := apprun.GetAutoScalingGroupParams{
		ClusterID:          apprun.ClusterID(parsedCluster),
		AutoScalingGroupID: apprun.AutoScalingGroupID(parsedASG),
	}

	resp, err := client.GetAutoScalingGroup(ctx, params)
	if err != nil {
		slog.Error("Failed to fetch AppRun ASG detail", slog.String("asgID", asgID), slog.Any("error", err))
		return nil, err
	}

	asg := resp.AutoScalingGroup

	// Convert name servers
	nameServers := make([]string, len(asg.NameServers))
	for i, ns := range asg.NameServers {
		nameServers[i] = string(ns)
	}

	// Convert interfaces
	interfaces := make([]AppRunASGInterface, len(asg.Interfaces))
	for i, iface := range asg.Interfaces {
		ipPool := make([]string, len(iface.IpPool))
		for j, ip := range iface.IpPool {
			ipPool[j] = fmt.Sprintf("%s-%s", string(ip.Start), string(ip.End))
		}

		netmaskLen := int16(0)
		if iface.NetmaskLen.Set {
			netmaskLen = iface.NetmaskLen.Value
		}
		defaultGateway := ""
		if iface.DefaultGateway.Set {
			defaultGateway = iface.DefaultGateway.Value
		}
		packetFilterID := ""
		if iface.PacketFilterID.Set {
			packetFilterID = iface.PacketFilterID.Value
		}

		interfaces[i] = AppRunASGInterface{
			Index:          iface.InterfaceIndex,
			Upstream:       iface.Upstream,
			NetmaskLen:     netmaskLen,
			DefaultGateway: defaultGateway,
			PacketFilterID: packetFilterID,
			ConnectsToLB:   iface.ConnectsToLB,
			IPPool:         ipPool,
		}
	}

	detail := &AppRunASGDetail{
		AppRunASG: AppRunASG{
			ID:              uuid.UUID(asg.AutoScalingGroupID).String(),
			Name:            asg.Name,
			Zone:            asg.Zone,
			MinNodes:        asg.MinNodes,
			MaxNodes:        asg.MaxNodes,
			WorkerNodeCount: asg.WorkerNodeCount,
			ServiceClass:    asg.WorkerServiceClassPath,
			ClusterID:       clusterID,
		},
		NameServers: nameServers,
		Interfaces:  interfaces,
		Deleting:    asg.Deleting,
	}

	slog.Info("Successfully fetched AppRun ASG detail", slog.String("asgID", asgID))
	return detail, nil
}

// ListAppRunWorkerNodes fetches all worker nodes for a specific ASG
func (c *SakuraClient) ListAppRunWorkerNodes(ctx context.Context, clusterID, asgID string) ([]AppRunWorkerNode, error) {
	slog.Info("Fetching AppRun Worker Nodes", slog.String("clusterID", clusterID), slog.String("asgID", asgID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsedCluster, err := uuid.Parse(clusterID)
	if err != nil {
		return nil, fmt.Errorf("invalid cluster ID: %w", err)
	}
	parsedASG, err := uuid.Parse(asgID)
	if err != nil {
		return nil, fmt.Errorf("invalid ASG ID: %w", err)
	}

	var allNodes []AppRunWorkerNode
	var cursor apprun.OptWorkerNodeID

	for {
		params := apprun.ListWorkerNodesParams{
			ClusterID:          apprun.ClusterID(parsedCluster),
			AutoScalingGroupID: apprun.AutoScalingGroupID(parsedASG),
			MaxItems:           30,
			Cursor:             cursor,
		}

		resp, err := client.ListWorkerNodes(ctx, params)
		if err != nil {
			slog.Error("Failed to fetch AppRun Worker Nodes", slog.Any("error", err))
			return nil, err
		}

		for _, node := range resp.WorkerNodes {
			createdAt := ""
			if node.Created > 0 {
				createdAt = time.Unix(int64(node.Created), 0).Format("2006-01-02 15:04:05")
			}

			// Extract IP addresses from network interfaces
			var ipAddresses []string
			for _, iface := range node.NetworkInterfaces {
				for _, addr := range iface.Addresses {
					ipAddresses = append(ipAddresses, addr.Address)
				}
			}

			resourceID := ""
			if !node.ResourceID.Null {
				resourceID = node.ResourceID.Value
			}

			errorMsg := ""
			if node.CreateErrorMessage.Set {
				errorMsg = node.CreateErrorMessage.Value
			}

			archiveVersion := ""
			if node.ArchiveVersion.Set {
				archiveVersion = node.ArchiveVersion.Value
			}

			allNodes = append(allNodes, AppRunWorkerNode{
				ID:             uuid.UUID(node.WorkerNodeID).String(),
				ResourceID:     resourceID,
				Status:         string(node.Status),
				Draining:       node.Draining,
				ArchiveVersion: archiveVersion,
				CreatedAt:      createdAt,
				ErrorMessage:   errorMsg,
				IPAddresses:    ipAddresses,
				ClusterID:      clusterID,
				ASGID:          asgID,
			})
		}

		if !resp.NextCursor.Set {
			break
		}
		cursor = resp.NextCursor
	}

	slog.Info("Successfully fetched AppRun Worker Nodes", slog.Int("count", len(allNodes)))
	return allNodes, nil
}

// ListAppRunLBs fetches all load balancers for a specific ASG
func (c *SakuraClient) ListAppRunLBs(ctx context.Context, clusterID, asgID string) ([]AppRunLB, error) {
	slog.Info("Fetching AppRun LBs", slog.String("clusterID", clusterID), slog.String("asgID", asgID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsedCluster, err := uuid.Parse(clusterID)
	if err != nil {
		return nil, fmt.Errorf("invalid cluster ID: %w", err)
	}

	parsedASG, err := uuid.Parse(asgID)
	if err != nil {
		return nil, fmt.Errorf("invalid ASG ID: %w", err)
	}

	var allLBs []AppRunLB
	var cursor apprun.OptLoadBalancerID

	for {
		params := apprun.ListLoadBalancersParams{
			ClusterID:          apprun.ClusterID(parsedCluster),
			AutoScalingGroupID: apprun.AutoScalingGroupID(parsedASG),
			MaxItems:           30,
			Cursor:             cursor,
		}

		resp, err := client.ListLoadBalancers(ctx, params)
		if err != nil {
			slog.Error("Failed to fetch AppRun LBs", slog.Any("error", err))
			return nil, err
		}

		for _, lb := range resp.LoadBalancers {
			createdAt := ""
			if lb.Created > 0 {
				createdAt = time.Unix(int64(lb.Created), 0).Format("2006-01-02 15:04:05")
			}
			allLBs = append(allLBs, AppRunLB{
				ID:           uuid.UUID(lb.LoadBalancerID).String(),
				Name:         lb.Name,
				ServiceClass: lb.ServiceClassPath,
				CreatedAt:    createdAt,
				ClusterID:    clusterID,
				ASGID:        asgID,
			})
		}

		if !resp.NextCursor.Set {
			break
		}
		cursor = resp.NextCursor
	}

	slog.Info("Successfully fetched AppRun LBs", slog.Int("count", len(allLBs)))
	return allLBs, nil
}

// ListAppRunApplications fetches all applications for a specific cluster
func (c *SakuraClient) ListAppRunApplications(ctx context.Context, clusterID string) ([]AppRunApplication, error) {
	slog.Info("Fetching AppRun Applications", slog.String("clusterID", clusterID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsed, err := uuid.Parse(clusterID)
	if err != nil {
		return nil, fmt.Errorf("invalid cluster ID: %w", err)
	}

	var allApps []AppRunApplication
	var cursor apprun.OptString

	for {
		params := apprun.ListApplicationsParams{
			ClusterID: apprun.OptClusterID{Value: apprun.ClusterID(parsed), Set: true},
			MaxItems:  30,
			Cursor:    cursor,
		}

		resp, err := client.ListApplications(ctx, params)
		if err != nil {
			slog.Error("Failed to fetch AppRun Applications", slog.Any("error", err))
			return nil, err
		}

		for _, app := range resp.Applications {
			activeVersion := int32(0)
			if !app.ActiveVersion.Null {
				activeVersion = app.ActiveVersion.Value
			}
			desiredCount := int32(0)
			if !app.DesiredCount.Null {
				desiredCount = app.DesiredCount.Value
			}
			allApps = append(allApps, AppRunApplication{
				ID:            uuid.UUID(app.ApplicationID).String(),
				Name:          app.Name,
				ClusterID:     clusterID,
				ClusterName:   app.ClusterName,
				ActiveVersion: activeVersion,
				DesiredCount:  desiredCount,
			})
		}

		if resp.NextCursor.Value == "" {
			break
		}
		cursor = resp.NextCursor
	}

	slog.Info("Successfully fetched AppRun Applications", slog.Int("count", len(allApps)))
	return allApps, nil
}

// ListAppRunVersions fetches all versions for a specific application
func (c *SakuraClient) ListAppRunVersions(ctx context.Context, applicationID, clusterID string, activeVersion int32) ([]AppRunVersion, error) {
	slog.Info("Fetching AppRun Versions", slog.String("applicationID", applicationID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsed, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application ID: %w", err)
	}

	var allVersions []AppRunVersion
	var cursor apprun.OptApplicationVersionNumber

	for {
		params := apprun.ListApplicationVersionsParams{
			ApplicationID: apprun.ApplicationID(parsed),
			MaxItems:      30,
			Cursor:        cursor,
		}

		resp, err := client.ListApplicationVersions(ctx, params)
		if err != nil {
			slog.Error("Failed to fetch AppRun Versions", slog.Any("error", err))
			return nil, err
		}

		for _, ver := range resp.Versions {
			createdAt := ""
			if ver.Created > 0 {
				createdAt = time.Unix(int64(ver.Created), 0).Format("2006-01-02 15:04:05")
			}
			allVersions = append(allVersions, AppRunVersion{
				Version:         int32(ver.Version),
				Image:           ver.Image,
				ActiveNodeCount: ver.ActiveNodeCount,
				CreatedAt:       createdAt,
				ApplicationID:   applicationID,
				ClusterID:       clusterID,
				IsActive:        int32(ver.Version) == activeVersion,
			})
		}

		if !resp.NextCursor.Set {
			break
		}
		cursor = resp.NextCursor
	}

	slog.Info("Successfully fetched AppRun Versions", slog.Int("count", len(allVersions)))
	return allVersions, nil
}

// GetAppRunLBDetail fetches detailed information about a specific load balancer
func (c *SakuraClient) GetAppRunLBDetail(ctx context.Context, clusterID, asgID, lbID string) (*AppRunLBDetail, error) {
	slog.Info("Fetching AppRun LB detail", slog.String("lbID", lbID))

	client, err := c.GetAppRunClient()
	if err != nil {
		slog.Error("Failed to get AppRun client", slog.Any("error", err))
		return nil, err
	}

	parsedCluster, err := uuid.Parse(clusterID)
	if err != nil {
		return nil, fmt.Errorf("invalid cluster ID: %w", err)
	}
	parsedASG, err := uuid.Parse(asgID)
	if err != nil {
		return nil, fmt.Errorf("invalid ASG ID: %w", err)
	}
	parsedLB, err := uuid.Parse(lbID)
	if err != nil {
		return nil, fmt.Errorf("invalid LB ID: %w", err)
	}

	params := apprun.GetLoadBalancerParams{
		ClusterID:          apprun.ClusterID(parsedCluster),
		AutoScalingGroupID: apprun.AutoScalingGroupID(parsedASG),
		LoadBalancerID:     apprun.LoadBalancerID(parsedLB),
	}

	resp, err := client.GetLoadBalancer(ctx, params)
	if err != nil {
		slog.Error("Failed to fetch AppRun LB detail", slog.String("lbID", lbID), slog.Any("error", err))
		return nil, err
	}

	lb := resp.LoadBalancer

	createdAt := ""
	if lb.Created > 0 {
		createdAt = time.Unix(int64(lb.Created), 0).Format("2006-01-02 15:04:05")
	}

	// Convert name servers
	nameServers := make([]string, len(lb.NameServers))
	for i, ns := range lb.NameServers {
		nameServers[i] = string(ns)
	}

	// Convert interfaces
	interfaces := make([]AppRunLBInterface, len(lb.Interfaces))
	for i, iface := range lb.Interfaces {
		ipPool := make([]string, len(iface.IpPool))
		for j, ip := range iface.IpPool {
			ipPool[j] = fmt.Sprintf("%s-%s", string(ip.Start), string(ip.End))
		}

		vip := ""
		if iface.Vip.Set {
			vip = iface.Vip.Value
		}
		gw := ""
		if iface.DefaultGateway.Set {
			gw = iface.DefaultGateway.Value
		}
		netmask := int16(0)
		if iface.NetmaskLen.Set {
			netmask = iface.NetmaskLen.Value
		}

		interfaces[i] = AppRunLBInterface{
			Index:          iface.InterfaceIndex,
			Upstream:       iface.Upstream,
			VIP:            vip,
			DefaultGateway: gw,
			NetmaskLen:     netmask,
			IPPool:         ipPool,
		}
	}

	detail := &AppRunLBDetail{
		AppRunLB: AppRunLB{
			ID:           uuid.UUID(lb.LoadBalancerID).String(),
			Name:         lb.Name,
			ServiceClass: lb.ServiceClassPath,
			CreatedAt:    createdAt,
			ClusterID:    clusterID,
			ASGID:        asgID,
		},
		NameServers: nameServers,
		Interfaces:  interfaces,
		Deleting:    lb.Deleting,
	}

	slog.Info("Successfully fetched AppRun LB detail", slog.String("lbID", lbID))
	return detail, nil
}
