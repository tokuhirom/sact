package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthStatus(t *testing.T) {
	client, err := NewSakuraClient("tk1b")
	assert.NoError(t, err)

	ctx := context.Background()
	authStatus, err := client.GetAuthStatus(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, authStatus)

	t.Logf("AccountName: %s", authStatus.AccountName)
	t.Logf("AccountID: %s", authStatus.AccountID)
	t.Logf("AccountCode: %s", authStatus.AccountCode)
}
