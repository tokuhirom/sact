package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// testNewClient creates a client for testing, skipping if no profile is configured
func testNewClient(t *testing.T, zone string) *SakuraClient {
	opts, _, err := LoadProfileAndZone()
	if err != nil {
		t.Skipf("Skipping test: no profile or environment variables configured: %v", err)
		return nil
	}

	client, err := NewSakuraClient(opts, zone)
	require.NoError(t, err)
	return client
}

func TestServerList(t *testing.T) {
	client := testNewClient(t, "tk1b")
	if client == nil {
		return
	}

	servers, err := client.ListServers(t.Context())
	require.NoError(t, err)
	require.IsType(t, []Server{}, servers)
	for i, server := range servers {
		t.Logf("server %d: %+v", i, server)
	}
}

func TestProxyLBList(t *testing.T) {
	client := testNewClient(t, "tk1b")
	if client == nil {
		return
	}

	elbs, err := client.ListELB(t.Context())
	require.NoError(t, err)
	require.IsType(t, []ELB{}, elbs)
	for i, elb := range elbs {
		t.Logf("elb %d: %+v", i, elb)
	}
}
