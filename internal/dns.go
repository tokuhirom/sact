package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// DNS represents a DNS zone resource
type DNS struct {
	ID          string
	Name        string
	Desc        string
	Zone        string
	RecordCount int
	CreatedAt   string
}

type DNSRecord struct {
	Name  string
	Type  string
	RData string
	TTL   int
}

type DNSDetail struct {
	DNS
	Tags        []string
	RecordCount int
	Records     []DNSRecord
	NameServers []string
	IconID      string
	CreatedAt   string
	ModifiedAt  string
}

// Implement list.Item interface for DNS
func (d DNS) FilterValue() string {
	return d.Name
}

func (d DNS) Title() string {
	return d.Name
}

func (d DNS) Description() string {
	desc := fmt.Sprintf("ID: %s", d.ID)
	if d.Desc != "" {
		desc += " | " + d.Desc
	}
	return desc
}

func (c *SakuraClient) ListDNS(ctx context.Context) ([]DNS, error) {
	slog.Info("Fetching DNS zones from Sakura Cloud")

	dnsOp := iaas.NewDNSOp(c.caller)

	searched, err := dnsOp.Find(ctx, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch DNS zones",
			slog.Any("error", err))
		return nil, err
	}

	dnsList := make([]DNS, 0, len(searched.DNS))
	for _, d := range searched.DNS {
		// Count records
		recordCount := len(d.Records)

		// Format created at
		createdAt := ""
		if !d.CreatedAt.IsZero() {
			createdAt = d.CreatedAt.Format("2006-01-02")
		}

		dnsList = append(dnsList, DNS{
			ID:          d.ID.String(),
			Name:        d.Name,
			Desc:        d.Description,
			Zone:        "global", // DNS is a global resource
			RecordCount: recordCount,
			CreatedAt:   createdAt,
		})
	}

	slog.Info("Successfully fetched DNS zones",
		slog.Int("count", len(dnsList)))

	return dnsList, nil
}

func (c *SakuraClient) GetDNSDetail(ctx context.Context, dnsID string) (*DNSDetail, error) {
	slog.Info("Fetching DNS detail from Sakura Cloud",
		slog.String("dnsID", dnsID))

	dnsOp := iaas.NewDNSOp(c.caller)

	id := types.StringID(dnsID)

	d, err := dnsOp.Read(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch DNS detail",
			slog.String("dnsID", dnsID),
			slog.Any("error", err))
		return nil, err
	}

	// Format created at
	createdAt := ""
	if !d.CreatedAt.IsZero() {
		createdAt = d.CreatedAt.Format("2006-01-02")
	}

	// Format modified at
	modifiedAt := ""
	if !d.ModifiedAt.IsZero() {
		modifiedAt = d.ModifiedAt.Format("2006-01-02 15:04:05")
	}

	// Get name servers
	nameServers := []string{}
	for _, ns := range d.DNSNameServers {
		nameServers = append(nameServers, ns)
	}

	// Get DNS records
	records := []DNSRecord{}
	for _, rec := range d.Records {
		records = append(records, DNSRecord{
			Name:  rec.Name,
			Type:  string(rec.Type),
			RData: rec.RData,
			TTL:   rec.TTL,
		})
	}

	detail := &DNSDetail{
		DNS: DNS{
			ID:          d.ID.String(),
			Name:        d.Name,
			Desc:        d.Description,
			Zone:        "global",
			RecordCount: len(d.Records),
			CreatedAt:   createdAt,
		},
		Tags:        d.Tags,
		RecordCount: len(d.Records),
		Records:     records,
		NameServers: nameServers,
		IconID:      d.IconID.String(),
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
	}

	slog.Info("Successfully fetched DNS detail",
		slog.String("dnsID", dnsID))

	return detail, nil
}
