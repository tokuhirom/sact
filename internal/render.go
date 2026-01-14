package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func renderServerDetail(detail *ServerDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Server: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Status:      %s\n", detail.InstanceStatus))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Description != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Description))
	}

	b.WriteString(fmt.Sprintf("CPU:         %d Core(s)\n", detail.CPU))
	b.WriteString(fmt.Sprintf("Memory:      %d GB\n", detail.MemoryGB))

	if len(detail.IPAddresses) > 0 {
		b.WriteString(fmt.Sprintf("IP Address:  %s\n", strings.Join(detail.IPAddresses, ", ")))
	}

	if len(detail.UserIPAddresses) > 0 {
		b.WriteString(fmt.Sprintf("User IP:     %s\n", strings.Join(detail.UserIPAddresses, ", ")))
	}

	if len(detail.Disks) > 0 {
		b.WriteString("\nDisks:\n")
		for _, disk := range detail.Disks {
			b.WriteString(fmt.Sprintf("  - %s (%d GB)\n", disk.Name, disk.SizeGB))
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderSwitchDetail(detail *SwitchDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Switch: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Subnets:     %d\n", detail.SubnetCount))
	b.WriteString(fmt.Sprintf("Servers:     %d connected\n", detail.ServerCount))

	if detail.DefaultRoute != "" {
		b.WriteString(fmt.Sprintf("Route:       %s\n", detail.DefaultRoute))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderDNSDetail(detail *DNSDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("DNS: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Records:     %d\n", detail.RecordCount))

	// Display DNS records in table format
	if len(detail.Records) > 0 {
		b.WriteString("\nDNS Records:\n")
		b.WriteString(fmt.Sprintf("  %-8s %-30s %-8s %s\n", "Type", "Name", "TTL", "Data"))
		b.WriteString(fmt.Sprintf("  %-8s %-30s %-8s %s\n", "----", "----", "---", "----"))
		for _, rec := range detail.Records {
			// Truncate long RData
			rdata := rec.RData
			if len(rdata) > 60 {
				rdata = rdata[:57] + "..."
			}
			b.WriteString(fmt.Sprintf("  %-8s %-30s %-8d %s\n",
				rec.Type,
				rec.Name,
				rec.TTL,
				rdata))
		}
	}

	if len(detail.NameServers) > 0 {
		b.WriteString("\nName Servers:\n")
		for _, ns := range detail.NameServers {
			b.WriteString(fmt.Sprintf("  - %s\n", ns))
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	if detail.ModifiedAt != "" {
		b.WriteString(fmt.Sprintf("Modified:    %s\n", detail.ModifiedAt))
	}

	return b.String()
}

func renderELBDetail(detail *ELBDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("ELB: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("VIP:         %s\n", detail.VIP))
	b.WriteString(fmt.Sprintf("Plan:        %s\n", detail.Plan))

	if detail.FQDN != "" {
		b.WriteString(fmt.Sprintf("FQDN:        %s\n", detail.FQDN))
	}

	b.WriteString(fmt.Sprintf("Servers:     %d\n", detail.ServerCount))

	// Display server list in table format
	if len(detail.Servers) > 0 {
		b.WriteString("\nServers:\n")
		b.WriteString(fmt.Sprintf("  %-20s %-8s %s\n", "IP Address", "Port", "Status"))
		b.WriteString(fmt.Sprintf("  %-20s %-8s %s\n", "----------", "----", "------"))
		for _, server := range detail.Servers {
			status := "Disabled"
			if server.Enabled {
				status = "Enabled"
			}
			b.WriteString(fmt.Sprintf("  %-20s %-8d %s\n",
				server.IPAddress,
				server.Port,
				status))
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderGSLBDetail(detail *GSLBDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("GSLB: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("FQDN:        %s\n", detail.FQDN))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Servers:     %d\n", detail.ServerCount))

	// Display health check settings
	if detail.HealthPath != "" {
		b.WriteString(fmt.Sprintf("Health Path: %s\n", detail.HealthPath))
	}
	b.WriteString(fmt.Sprintf("Delay Loop:  %d sec\n", detail.DelayLoop))
	if detail.Weighted {
		b.WriteString("Weighted:    Yes\n")
	} else {
		b.WriteString("Weighted:    No\n")
	}

	// Display server list in table format using bubbles table
	if len(detail.Servers) > 0 {
		b.WriteString("\nServers:\n")

		var columns []table.Column
		var rows []table.Row

		if detail.Weighted {
			columns = []table.Column{
				{Title: "IP Address", Width: 20},
				{Title: "Weight", Width: 8},
				{Title: "Status", Width: 10},
			}

			for _, server := range detail.Servers {
				status := "Disabled"
				if server.Enabled {
					status = "Enabled"
				}
				rows = append(rows, table.Row{
					server.IPAddress,
					fmt.Sprintf("%d", server.Weight),
					status,
				})
			}
		} else {
			columns = []table.Column{
				{Title: "IP Address", Width: 20},
				{Title: "Status", Width: 10},
			}

			for _, server := range detail.Servers {
				status := "Disabled"
				if server.Enabled {
					status = "Enabled"
				}
				rows = append(rows, table.Row{
					server.IPAddress,
					status,
				})
			}
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(false),
			table.WithHeight(len(rows)),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		b.WriteString(t.View())
		b.WriteString("\n")
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderDBDetail(detail *DBDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("DB: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("DB Type:     %s\n", detail.DBType))
	b.WriteString(fmt.Sprintf("Status:      %s\n", detail.InstanceStatus))
	b.WriteString(fmt.Sprintf("Plan:        %s\n", detail.Plan))
	b.WriteString(fmt.Sprintf("CPU:         %d Core(s)\n", detail.CPU))
	b.WriteString(fmt.Sprintf("Memory:      %d GB\n", detail.MemoryGB))
	b.WriteString(fmt.Sprintf("Disk Size:   %d GB\n", detail.DiskSizeGB))

	if detail.IPAddress != "" {
		b.WriteString(fmt.Sprintf("IP Address:  %s\n", detail.IPAddress))
		if detail.NetworkMaskLen > 0 {
			b.WriteString(fmt.Sprintf("Netmask:     /%d\n", detail.NetworkMaskLen))
		}
		if detail.DefaultRoute != "" {
			b.WriteString(fmt.Sprintf("Gateway:     %s\n", detail.DefaultRoute))
		}
	}

	if detail.Port > 0 {
		b.WriteString(fmt.Sprintf("Port:        %d\n", detail.Port))
	}

	if detail.DefaultUser != "" {
		b.WriteString(fmt.Sprintf("User:        %s\n", detail.DefaultUser))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderDiskDetail(detail *DiskDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Disk: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Size:        %d GB\n", detail.SizeGB))
	b.WriteString(fmt.Sprintf("Connection:  %s\n", detail.Connection))
	b.WriteString(fmt.Sprintf("Plan:        %s\n", detail.DiskPlanName))
	b.WriteString(fmt.Sprintf("Availability: %s\n", detail.Availability))

	if detail.EncryptionAlgo != "" {
		b.WriteString(fmt.Sprintf("Encryption:  %s\n", detail.EncryptionAlgo))
	}

	if detail.ServerID != "" {
		b.WriteString(fmt.Sprintf("\nServer ID:   %s\n", detail.ServerID))
		b.WriteString(fmt.Sprintf("Server Name: %s\n", detail.ServerName))
	} else {
		b.WriteString("\nServer:      (not attached)\n")
	}

	if detail.SourceDiskID != "" {
		b.WriteString(fmt.Sprintf("\nSource Disk: %s\n", detail.SourceDiskID))
	}
	if detail.SourceArchiveID != "" {
		b.WriteString(fmt.Sprintf("Source Archive: %s\n", detail.SourceArchiveID))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}
	if detail.ModifiedAt != "" {
		b.WriteString(fmt.Sprintf("Modified:    %s\n", detail.ModifiedAt))
	}

	return b.String()
}

func renderArchiveDetail(detail *ArchiveDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Archive: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Size:        %d GB\n", detail.SizeGB))
	b.WriteString(fmt.Sprintf("Scope:       %s\n", detail.Scope))
	b.WriteString(fmt.Sprintf("Availability: %s\n", detail.Availability))

	if detail.BundleInfo != "" {
		b.WriteString(fmt.Sprintf("Bundle:      %s\n", detail.BundleInfo))
	}

	if detail.SourceDiskID != "" {
		b.WriteString(fmt.Sprintf("\nSource Disk: %s\n", detail.SourceDiskID))
	}
	if detail.SourceArchiveID != "" {
		b.WriteString(fmt.Sprintf("Source Archive: %s\n", detail.SourceArchiveID))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}
	if detail.ModifiedAt != "" {
		b.WriteString(fmt.Sprintf("Modified:    %s\n", detail.ModifiedAt))
	}

	return b.String()
}

func renderInternetDetail(detail *InternetDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Internet: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Bandwidth:   %d Mbps\n", detail.BandWidthMbps))
	b.WriteString(fmt.Sprintf("Netmask:     /%d\n", detail.NetworkMaskLen))

	if detail.SwitchID != "" {
		b.WriteString(fmt.Sprintf("Switch ID:   %s\n", detail.SwitchID))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderVPCRouterDetail(detail *VPCRouterDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("VPC Router: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Plan:        %s\n", detail.Plan))
	b.WriteString(fmt.Sprintf("Version:     %d\n", detail.Version))
	b.WriteString(fmt.Sprintf("Status:      %s\n", detail.InstanceStatus))

	if len(detail.PublicIPAddresses) > 0 {
		b.WriteString(fmt.Sprintf("Public IPs:  %s\n", strings.Join(detail.PublicIPAddresses, ", ")))
	}

	if len(detail.NICs) > 0 {
		b.WriteString("\nNetwork Interfaces:\n")
		for _, nic := range detail.NICs {
			if nic.SwitchID != "" || nic.IPAddress != "" {
				b.WriteString(fmt.Sprintf("  NIC%d: Switch=%s, IP=%s\n", nic.Index, nic.SwitchID, nic.IPAddress))
			}
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderPacketFilterDetail(detail *PacketFilterDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Packet Filter: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Rules:       %d\n", detail.RuleCount))

	if detail.ExpressionHash != "" {
		b.WriteString(fmt.Sprintf("Hash:        %s\n", detail.ExpressionHash))
	}

	// Display rules in table format
	if len(detail.Rules) > 0 {
		b.WriteString("\nFilter Rules:\n")
		b.WriteString(fmt.Sprintf("  %-8s %-18s %-12s %-12s %-8s %s\n",
			"Protocol", "Source", "SrcPort", "DstPort", "Action", "Description"))
		b.WriteString(fmt.Sprintf("  %-8s %-18s %-12s %-12s %-8s %s\n",
			"--------", "------", "-------", "-------", "------", "-----------"))
		for _, rule := range detail.Rules {
			srcNet := rule.SourceNetwork
			if srcNet == "" {
				srcNet = "*"
			}
			srcPort := rule.SourcePort
			if srcPort == "" {
				srcPort = "*"
			}
			dstPort := rule.DestinationPort
			if dstPort == "" {
				dstPort = "*"
			}
			b.WriteString(fmt.Sprintf("  %-8s %-18s %-12s %-12s %-8s %s\n",
				rule.Protocol,
				srcNet,
				srcPort,
				dstPort,
				rule.Action,
				rule.Description))
		}
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderLoadBalancerDetail(detail *LoadBalancerDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Load Balancer: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Status:      %s\n", detail.InstanceStatus))
	b.WriteString(fmt.Sprintf("VRID:        %d\n", detail.VRID))

	if len(detail.IPAddresses) > 0 {
		b.WriteString(fmt.Sprintf("IP Addresses: %s\n", strings.Join(detail.IPAddresses, ", ")))
	}

	b.WriteString(fmt.Sprintf("Network:     /%d (Default Route: %s)\n", detail.NetworkMaskLen, detail.DefaultRoute))

	if detail.SwitchID != "" && detail.SwitchID != "0" {
		b.WriteString(fmt.Sprintf("Switch ID:   %s\n", detail.SwitchID))
	}

	// Display VIPs
	if len(detail.VIPs) > 0 {
		b.WriteString("\nVirtual IPs:\n")
		for i, vip := range detail.VIPs {
			b.WriteString(fmt.Sprintf("\n  VIP %d: %s:%d\n", i+1, vip.VirtualIPAddress, vip.Port))
			if vip.Description != "" {
				b.WriteString(fmt.Sprintf("    Description: %s\n", vip.Description))
			}
			b.WriteString(fmt.Sprintf("    Delay Loop:  %d sec\n", vip.DelayLoop))
			if vip.SorryServer != "" {
				b.WriteString(fmt.Sprintf("    Sorry Server: %s\n", vip.SorryServer))
			}

			// Display real servers
			if len(vip.Servers) > 0 {
				b.WriteString("    Real Servers:\n")
				for _, srv := range vip.Servers {
					enabled := "enabled"
					if !srv.Enabled {
						enabled = "disabled"
					}
					b.WriteString(fmt.Sprintf("      - %s:%d (%s)\n", srv.IPAddress, srv.Port, enabled))
				}
			}
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderNFSDetail(detail *NFSDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("NFS: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Status:      %s\n", detail.InstanceStatus))

	if detail.PlanID != "" && detail.PlanID != "0" {
		b.WriteString(fmt.Sprintf("Plan ID:     %s\n", detail.PlanID))
	}

	if len(detail.IPAddresses) > 0 {
		b.WriteString(fmt.Sprintf("IP Addresses: %s\n", strings.Join(detail.IPAddresses, ", ")))
	}

	if detail.NetworkMaskLen > 0 {
		b.WriteString(fmt.Sprintf("Network:     /%d\n", detail.NetworkMaskLen))
	}

	if detail.DefaultRoute != "" {
		b.WriteString(fmt.Sprintf("Gateway:     %s\n", detail.DefaultRoute))
	}

	if detail.SwitchID != "" && detail.SwitchID != "0" {
		b.WriteString(fmt.Sprintf("Switch ID:   %s\n", detail.SwitchID))
	}

	if detail.SwitchName != "" {
		b.WriteString(fmt.Sprintf("Switch Name: %s\n", detail.SwitchName))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderSSHKeyDetail(detail *SSHKeyDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("SSH Key: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Fingerprint: %s\n", detail.Fingerprint))

	if detail.PublicKey != "" {
		b.WriteString("\nPublic Key:\n")
		// Split long public key for better display
		pubKey := detail.PublicKey
		if len(pubKey) > 80 {
			for i := 0; i < len(pubKey); i += 80 {
				end := i + 80
				if end > len(pubKey) {
					end = len(pubKey)
				}
				b.WriteString(fmt.Sprintf("  %s\n", pubKey[i:end]))
			}
		} else {
			b.WriteString(fmt.Sprintf("  %s\n", pubKey))
		}
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderAutoBackupDetail(detail *AutoBackupDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Auto Backup: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Availability: %s\n", detail.Availability))

	if detail.DiskID != "" && detail.DiskID != "0" {
		b.WriteString(fmt.Sprintf("Disk ID:     %s\n", detail.DiskID))
	}

	b.WriteString(fmt.Sprintf("Max Backups: %d\n", detail.MaxBackups))

	if len(detail.BackupWeekdays) > 0 {
		b.WriteString(fmt.Sprintf("Weekdays:    %s\n", strings.Join(detail.BackupWeekdays, ", ")))
	}

	if detail.ZoneName != "" {
		b.WriteString(fmt.Sprintf("Zone Name:   %s\n", detail.ZoneName))
	}

	if detail.AccountID != "" && detail.AccountID != "0" {
		b.WriteString(fmt.Sprintf("Account ID:  %s\n", detail.AccountID))
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	if detail.ModifiedAt != "" {
		b.WriteString(fmt.Sprintf("Modified:    %s\n", detail.ModifiedAt))
	}

	return b.String()
}

func renderSimpleMonitorDetail(detail *SimpleMonitorDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Simple Monitor: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Target:      %s\n", detail.Target))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	b.WriteString(fmt.Sprintf("Protocol:    %s\n", detail.Protocol))
	if detail.Port > 0 {
		b.WriteString(fmt.Sprintf("Port:        %d\n", detail.Port))
	}
	if detail.Path != "" {
		b.WriteString(fmt.Sprintf("Path:        %s\n", detail.Path))
	}
	if detail.Host != "" {
		b.WriteString(fmt.Sprintf("Host:        %s\n", detail.Host))
	}

	enabledStr := "No"
	if detail.Enabled {
		enabledStr = "Yes"
	}
	b.WriteString(fmt.Sprintf("Enabled:     %s\n", enabledStr))

	if detail.Health != "" {
		b.WriteString(fmt.Sprintf("Health:      %s\n", detail.Health))
	}

	b.WriteString(fmt.Sprintf("\nMonitoring Settings:\n"))
	b.WriteString(fmt.Sprintf("  Delay Loop:       %d sec\n", detail.DelayLoop))
	b.WriteString(fmt.Sprintf("  Max Check Attempts: %d\n", detail.MaxCheckAttempts))
	b.WriteString(fmt.Sprintf("  Retry Interval:   %d sec\n", detail.RetryInterval))
	b.WriteString(fmt.Sprintf("  Timeout:          %d sec\n", detail.Timeout))

	b.WriteString(fmt.Sprintf("\nNotification Settings:\n"))
	emailEnabled := "No"
	if detail.NotifyEmailEnabled {
		emailEnabled = "Yes"
	}
	b.WriteString(fmt.Sprintf("  Email:            %s\n", emailEnabled))
	slackEnabled := "No"
	if detail.NotifySlackEnabled {
		slackEnabled = "Yes"
	}
	b.WriteString(fmt.Sprintf("  Slack:            %s\n", slackEnabled))
	if detail.NotifyInterval > 0 {
		b.WriteString(fmt.Sprintf("  Notify Interval:  %d hours\n", detail.NotifyInterval))
	}

	if detail.LastCheckedAt != "" {
		b.WriteString(fmt.Sprintf("\nLast Checked:   %s\n", detail.LastCheckedAt))
	}
	if detail.LastHealthChangedAt != "" {
		b.WriteString(fmt.Sprintf("Health Changed: %s\n", detail.LastHealthChangedAt))
	}

	if len(detail.LatestLogs) > 0 {
		b.WriteString("\nLatest Logs:\n")
		for _, log := range detail.LatestLogs {
			b.WriteString(fmt.Sprintf("  %s\n", log))
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:        %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	if detail.ModifiedAt != "" {
		b.WriteString(fmt.Sprintf("Modified:    %s\n", detail.ModifiedAt))
	}

	return b.String()
}

func renderBridgeDetail(detail *BridgeDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Bridge: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:          %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("Zone:        %s\n", detail.Zone))

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", detail.Desc))
	}

	if detail.Region != "" {
		b.WriteString(fmt.Sprintf("Region:      %s\n", detail.Region))
	}

	b.WriteString(fmt.Sprintf("Switches:    %d connected\n", detail.SwitchCount))

	// Display switch in zone info
	if detail.SwitchInZone != nil {
		b.WriteString("\nSwitch in Zone:\n")
		b.WriteString(fmt.Sprintf("  ID:             %s\n", detail.SwitchInZone.ID))
		b.WriteString(fmt.Sprintf("  Name:           %s\n", detail.SwitchInZone.Name))
		if detail.SwitchInZone.Scope != "" {
			b.WriteString(fmt.Sprintf("  Scope:          %s\n", detail.SwitchInZone.Scope))
		}
		b.WriteString(fmt.Sprintf("  Servers:        %d\n", detail.SwitchInZone.ServerCount))
		b.WriteString(fmt.Sprintf("  Appliances:     %d\n", detail.SwitchInZone.ApplianceCount))
	}

	// Display connected switches
	if len(detail.Switches) > 0 {
		b.WriteString("\nConnected Switches:\n")
		for _, sw := range detail.Switches {
			b.WriteString(fmt.Sprintf("  - %s (ID: %s, Zone: %s)\n", sw.Name, sw.ID, sw.ZoneName))
		}
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:     %s\n", detail.CreatedAt))
	}

	return b.String()
}

func renderContainerRegistryDetail(detail *ContainerRegistryDetail) string {
	var b strings.Builder

	b.WriteString(selectedStyle.Render(fmt.Sprintf("Container Registry: %s", detail.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("ID:            %s\n", detail.ID))
	b.WriteString(fmt.Sprintf("FQDN:          %s\n", detail.FQDN))
	b.WriteString(fmt.Sprintf("Access Level:  %s\n", detail.AccessLevel))
	b.WriteString(fmt.Sprintf("Availability:  %s\n", detail.Availability))

	if detail.SubDomainLabel != "" {
		b.WriteString(fmt.Sprintf("Subdomain:     %s\n", detail.SubDomainLabel))
	}

	if detail.VirtualDomain != "" {
		b.WriteString(fmt.Sprintf("Virtual Domain: %s\n", detail.VirtualDomain))
	}

	if detail.Desc != "" {
		b.WriteString(fmt.Sprintf("Description:   %s\n", detail.Desc))
	}

	// Display users
	b.WriteString(fmt.Sprintf("\nUsers:         %d\n", detail.UserCount))
	if len(detail.Users) > 0 {
		for _, user := range detail.Users {
			b.WriteString(fmt.Sprintf("  - %s (%s)\n", user.UserName, user.Permission))
		}
	}

	// Display images
	b.WriteString("\nImages:\n")
	if detail.ImagesError != "" {
		b.WriteString(fmt.Sprintf("  Error: %s\n", detail.ImagesError))
	} else if len(detail.Images) == 0 {
		b.WriteString("  (no images)\n")
	} else {
		for _, img := range detail.Images {
			sizeStr := "-"
			if img.Size > 0 {
				sizeStr = formatBytes(img.Size)
			}
			createdStr := "-"
			if img.CreatedAt != "" {
				createdStr = img.CreatedAt
			}
			b.WriteString(fmt.Sprintf("  - %s:%s  [%s]  %s\n", img.Repository, img.Tag, sizeStr, createdStr))
		}
	}

	if len(detail.Tags) > 0 {
		b.WriteString(fmt.Sprintf("\nTags:          %s\n", strings.Join(detail.Tags, ", ")))
	}

	if detail.CreatedAt != "" {
		b.WriteString(fmt.Sprintf("\nCreated:       %s\n", detail.CreatedAt))
	}

	if detail.ModifiedAt != "" {
		b.WriteString(fmt.Sprintf("Modified:      %s\n", detail.ModifiedAt))
	}

	return b.String()
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
