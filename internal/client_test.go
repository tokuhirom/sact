package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerList(t *testing.T) {
	client, err := NewSakuraClient("tk1b")
	require.NoError(t, err)

	servers, err := client.ListServers(t.Context())
	require.NoError(t, err)
	require.IsType(t, []Server{}, servers)
}
