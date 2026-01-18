package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAndSetDefaultZone(t *testing.T) {
	tests := []struct {
		name     string
		zone     string
		expected string
	}{
		{
			name:     "valid zone tk1a",
			zone:     "tk1a",
			expected: "tk1a",
		},
		{
			name:     "valid zone tk1b",
			zone:     "tk1b",
			expected: "tk1b",
		},
		{
			name:     "valid zone is1a",
			zone:     "is1a",
			expected: "is1a",
		},
		{
			name:     "valid zone is1b",
			zone:     "is1b",
			expected: "is1b",
		},
		{
			name:     "valid zone is1c",
			zone:     "is1c",
			expected: "is1c",
		},
		{
			name:     "empty zone falls back to tk1b",
			zone:     "",
			expected: "tk1b",
		},
		{
			name:     "invalid zone falls back to tk1b",
			zone:     "invalid",
			expected: "tk1b",
		},
		{
			name:     "unknown region falls back to tk1b",
			zone:     "us1a",
			expected: "tk1b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateAndSetDefaultZone(tt.zone)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadProfileAndZone(t *testing.T) {
	// This test requires either a usacloud profile or environment variables to be set
	// Skip if not configured
	opts, zone, err := LoadProfileAndZone()

	if err != nil {
		t.Skipf("Skipping test: no profile or environment variables configured: %v", err)
		return
	}

	// If we get here, credentials should be loaded
	assert.NotNil(t, opts)
	assert.NotEmpty(t, opts.AccessToken)
	assert.NotEmpty(t, opts.AccessTokenSecret)

	// Zone should be one of the valid zones
	validZone := false
	for _, v := range validZones {
		if zone == v {
			validZone = true
			break
		}
	}
	assert.True(t, validZone, "zone should be valid: %s", zone)
}
