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
