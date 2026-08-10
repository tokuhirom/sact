package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthStatus(t *testing.T) {
	opts, _, err := LoadProfileAndZone()
	if err != nil {
		t.Skipf("Skipping test: no profile or environment variables configured: %v", err)
		return
	}

	client, err := NewSakuraClient(opts, "tk1b")
	assert.NoError(t, err)

	ctx := context.Background()
	authStatus, err := client.GetAuthStatus(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, authStatus)

	t.Logf("AccountName: %s", authStatus.AccountName)
	t.Logf("AccountID: %s", authStatus.AccountID)
	t.Logf("AccountCode: %s", authStatus.AccountCode)
}
