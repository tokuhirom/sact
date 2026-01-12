package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

// DB represents a Database Appliance resource
type DB struct {
	ID             string
	Name           string
	Desc           string
	Zone           string
	DBType         string
	Plan           string
	InstanceStatus string
	CreatedAt      string
}

type DBDetail struct {
	DB
	Tags           []string
	CPU            int
	MemoryGB       int
	DiskSizeGB     int
	IPAddress      string
	Port           int
	DefaultUser    string
	NetworkMaskLen int
	DefaultRoute   string
	CreatedAt      string
}

// Implement list.Item interface for DB
func (d DB) FilterValue() string {
	return d.Name
}

func (d DB) Title() string {
	return d.Name
}

func (d DB) Description() string {
	desc := fmt.Sprintf("ID: %s", d.ID)
	if d.Desc != "" {
		desc += " | " + d.Desc
	}
	return desc
}

func (c *SakuraClient) ListDB(ctx context.Context) ([]DB, error) {
	if c.zone == "" {
		slog.Error("Zone is not set in client")
		return nil, fmt.Errorf("zone is not set")
	}

	slog.Info("Fetching DBs from Sakura Cloud",
		slog.String("zone", c.zone))

	dbOp := iaas.NewDatabaseOp(c.caller)

	searched, err := dbOp.Find(ctx, c.zone, &iaas.FindCondition{
		Sort: search.SortKeys{
			search.SortKeyAsc("Name"),
		},
	})
	if err != nil {
		slog.Error("Failed to fetch DBs",
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	dbList := make([]DB, 0, len(searched.Databases))
	for _, db := range searched.Databases {
		// Get DB type
		dbType := "Unknown"
		if db.Conf != nil && db.Conf.DatabaseName != "" {
			dbType = db.Conf.DatabaseName
		}

		// Get plan from PlanID
		plan := db.PlanID.String()

		// Format created at
		createdAt := ""
		if !db.CreatedAt.IsZero() {
			createdAt = db.CreatedAt.Format("2006-01-02")
		}

		// Get instance status
		instanceStatus := string(db.InstanceStatus)

		dbList = append(dbList, DB{
			ID:             db.ID.String(),
			Name:           db.Name,
			Desc:           db.Description,
			Zone:           c.zone,
			DBType:         dbType,
			Plan:           plan,
			InstanceStatus: instanceStatus,
			CreatedAt:      createdAt,
		})
	}

	slog.Info("Successfully fetched DBs",
		slog.String("zone", c.zone),
		slog.Int("count", len(dbList)))

	return dbList, nil
}

func (c *SakuraClient) GetDBDetail(ctx context.Context, dbID string) (*DBDetail, error) {
	slog.Info("Fetching DB detail from Sakura Cloud",
		slog.String("dbID", dbID),
		slog.String("zone", c.zone))

	dbOp := iaas.NewDatabaseOp(c.caller)

	id := types.StringID(dbID)

	db, err := dbOp.Read(ctx, c.zone, id)
	if err != nil {
		slog.Error("Failed to fetch DB detail",
			slog.String("dbID", dbID),
			slog.String("zone", c.zone),
			slog.Any("error", err))
		return nil, err
	}

	// Get DB type
	dbType := "Unknown"
	if db.Conf != nil && db.Conf.DatabaseName != "" {
		dbType = db.Conf.DatabaseName
	}

	// Get plan from PlanID
	plan := db.PlanID.String()

	// Format created at
	createdAt := ""
	if !db.CreatedAt.IsZero() {
		createdAt = db.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Get instance status
	instanceStatus := string(db.InstanceStatus)

	// Parse plan - simplified
	cpu := 0
	memoryGB := 0
	diskSizeGB := 0

	// Try to extract from plan name
	planStr := plan
	switch planStr {
	case "10":
		cpu = 1
		memoryGB = 2
		diskSizeGB = 20
	case "30":
		cpu = 2
		memoryGB = 4
		diskSizeGB = 30
	case "50":
		cpu = 4
		memoryGB = 8
		diskSizeGB = 50
	default:
		cpu = 1
		memoryGB = 2
		diskSizeGB = 20
	}

	// Get network info
	ipAddress := ""
	port := 0
	defaultUser := ""
	networkMaskLen := 0
	defaultRoute := ""

	if len(db.Interfaces) > 0 {
		ipAddress = db.Interfaces[0].IPAddress
		networkMaskLen = 24 // Default
		defaultRoute = ""
	}

	if db.Conf != nil {
		if db.Conf.DatabaseName == "postgres" {
			port = 5432
		} else {
			port = 3306
		}
		if db.Conf.DefaultUser != "" {
			defaultUser = db.Conf.DefaultUser
		}
	}

	detail := &DBDetail{
		DB: DB{
			ID:             db.ID.String(),
			Name:           db.Name,
			Desc:           db.Description,
			Zone:           c.zone,
			DBType:         dbType,
			Plan:           plan,
			InstanceStatus: instanceStatus,
			CreatedAt:      createdAt,
		},
		Tags:           db.Tags,
		CPU:            cpu,
		MemoryGB:       memoryGB,
		DiskSizeGB:     diskSizeGB,
		IPAddress:      ipAddress,
		Port:           port,
		DefaultUser:    defaultUser,
		NetworkMaskLen: networkMaskLen,
		DefaultRoute:   defaultRoute,
		CreatedAt:      createdAt,
	}

	slog.Info("Successfully fetched DB detail",
		slog.String("dbID", dbID))

	return detail, nil
}
