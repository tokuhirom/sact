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
	IPAddresses    []string
	Port           int
	DefaultUser    string
	NetworkMaskLen int
	DefaultRoute   string
	SwitchID       string
	SourceNetworks []string
	WebUI          string
	// Backup settings
	BackupTime     string
	BackupWeekdays []string
	BackupRotate   int
	// Replication
	ReplicationModel string
	ReplicationIP    string
	// Host info
	HostName     string
	Availability string
	ModifiedAt   string
	CreatedAt    string
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

	// Format modified at
	modifiedAt := ""
	if !db.ModifiedAt.IsZero() {
		modifiedAt = db.ModifiedAt.Format("2006-01-02 15:04:05")
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

	// Get network info from API (not hardcoded)
	ipAddresses := db.IPAddresses
	networkMaskLen := db.NetworkMaskLen
	defaultRoute := db.DefaultRoute
	switchID := ""
	if !db.SwitchID.IsEmpty() {
		switchID = db.SwitchID.String()
	}

	// Get port and user from CommonSetting
	port := 0
	defaultUser := ""
	var sourceNetworks []string
	webUI := ""

	if db.CommonSetting != nil {
		port = db.CommonSetting.ServicePort
		defaultUser = db.CommonSetting.DefaultUser
		sourceNetworks = db.CommonSetting.SourceNetwork
		webUI = string(db.CommonSetting.WebUI)
	}

	// Fallback port if not set
	if port == 0 && db.Conf != nil {
		if db.Conf.DatabaseName == "postgres" {
			port = 5432
		} else {
			port = 3306
		}
	}

	// Fallback user if not set
	if defaultUser == "" && db.Conf != nil && db.Conf.DefaultUser != "" {
		defaultUser = db.Conf.DefaultUser
	}

	// Backup settings
	backupTime := ""
	var backupWeekdays []string
	backupRotate := 0
	if db.BackupSetting != nil {
		backupTime = db.BackupSetting.Time
		backupRotate = db.BackupSetting.Rotate
		for _, day := range db.BackupSetting.DayOfWeek {
			backupWeekdays = append(backupWeekdays, string(day))
		}
	}

	// Replication settings
	replicationModel := ""
	replicationIP := ""
	if db.ReplicationSetting != nil {
		replicationModel = string(db.ReplicationSetting.Model)
		replicationIP = db.ReplicationSetting.IPAddress
	}

	// Host info
	hostName := db.InstanceHostName

	// Availability
	availability := string(db.Availability)

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
		Tags:             db.Tags,
		CPU:              cpu,
		MemoryGB:         memoryGB,
		DiskSizeGB:       diskSizeGB,
		IPAddresses:      ipAddresses,
		Port:             port,
		DefaultUser:      defaultUser,
		NetworkMaskLen:   networkMaskLen,
		DefaultRoute:     defaultRoute,
		SwitchID:         switchID,
		SourceNetworks:   sourceNetworks,
		WebUI:            webUI,
		BackupTime:       backupTime,
		BackupWeekdays:   backupWeekdays,
		BackupRotate:     backupRotate,
		ReplicationModel: replicationModel,
		ReplicationIP:    replicationIP,
		HostName:         hostName,
		Availability:     availability,
		ModifiedAt:       modifiedAt,
		CreatedAt:        createdAt,
	}

	slog.Info("Successfully fetched DB detail",
		slog.String("dbID", dbID))

	return detail, nil
}
