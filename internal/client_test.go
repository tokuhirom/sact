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
	for i, server := range servers {
		t.Logf("server %d: %+v", i, server)
	}
}

func TestProxyLBList(t *testing.T) {
	client, err := NewSakuraClient("tk1b")
	require.NoError(t, err)

	elbs, err := client.ListELB(t.Context())
	require.NoError(t, err)
	require.IsType(t, []ELB{}, elbs)
	for i, elb := range elbs {
		t.Logf("elb %d: %+v", i, elb)
	}
}
