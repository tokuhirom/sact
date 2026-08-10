package internal

import (
	"fmt"
	"log/slog"

	"github.com/sacloud/iaas-api-go/helper/api"
)

var validZones = []string{"tk1a", "tk1b", "is1a", "is1b", "is1c"}

// LoadProfileAndZone loads credentials and default zone from usacloud profile or environment variables.
// It returns:
// - opts: Caller options containing credentials from profile or environment variables
// - defaultZone: The zone to use (from profile's Zone field, or "tk1b" as fallback)
// - err: Error if credentials cannot be loaded
func LoadProfileAndZone() (opts *api.CallerOptions, defaultZone string, err error) {
	// Load credentials from profile or environment variables
	// This automatically:
	// - Reads ~/.usacloud/current to find active profile
	// - Falls back to environment variables (SAKURACLOUD_ACCESS_TOKEN, etc.)
	// - Provides ProfileConfigValue() to access profile's Zone field
	opts, err = api.DefaultOption()
	if err != nil {
		return nil, "", fmt.Errorf("failed to load credentials: %w\n\nPlease configure credentials using one of:\n  1. Run 'usacloud config' to set up a profile\n  2. Set environment variables: SAKURACLOUD_ACCESS_TOKEN and SAKURACLOUD_ACCESS_TOKEN_SECRET", err)
	}

	// Check if credentials are actually present
	if opts.AccessToken == "" || opts.AccessTokenSecret == "" {
		return nil, "", fmt.Errorf("credentials not found\n\nPlease configure credentials using one of:\n  1. Run 'usacloud config' to set up a profile\n  2. Set environment variables: SAKURACLOUD_ACCESS_TOKEN and SAKURACLOUD_ACCESS_TOKEN_SECRET")
	}

	// Extract zone from profile (if available)
	zone := ""
	if profileConfig := opts.ProfileConfigValue(); profileConfig != nil {
		zone = profileConfig.Zone
		slog.Info("Loaded zone from profile", slog.String("zone", zone))
	}

	// Validate and set default zone
	defaultZone = validateAndSetDefaultZone(zone)

	slog.Info("Profile loaded successfully",
		slog.String("default_zone", defaultZone),
		slog.Bool("from_profile", zone != ""))

	return opts, defaultZone, nil
}

// validateAndSetDefaultZone validates the zone and returns a valid zone or falls back to tk1b
func validateAndSetDefaultZone(zone string) string {
	if zone == "" {
		slog.Info("No zone specified in profile, using default", slog.String("default", "tk1b"))
		return "tk1b"
	}

	// Check if zone is valid
	for _, validZone := range validZones {
		if zone == validZone {
			return zone
		}
	}

	// Invalid zone, fall back to default
	slog.Warn("Invalid zone in profile, using default",
		slog.String("invalid_zone", zone),
		slog.String("default", "tk1b"))
	return "tk1b"
}
