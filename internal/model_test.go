package internal

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	client, err := NewSakuraClient("tk1b")
	assert.NoError(t, err)

	m := InitialModel(client, "tk1b")

	assert.Equal(t, "tk1b", m.currentZone)
	assert.Equal(t, 1, m.cursor) // tk1b is at index 1
	assert.True(t, m.loading)
	assert.False(t, m.quitting)
	assert.False(t, m.searchMode)
	assert.Equal(t, "", m.searchQuery)
}

func TestInitialModelWithDifferentZone(t *testing.T) {
	client, err := NewSakuraClient("is1a")
	assert.NoError(t, err)

	m := InitialModel(client, "is1a")

	assert.Equal(t, "is1a", m.currentZone)
	assert.Equal(t, 2, m.cursor) // is1a is at index 2
}

func TestUpdateWindowSize(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.Equal(t, 100, m.windowWidth)
	assert.Equal(t, 50, m.windowHeight)
}

func TestUpdateQuit(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updated, cmd := m.Update(msg)
	m = updated.(model)

	assert.True(t, m.quitting)
	assert.NotNil(t, cmd)
}

func TestEscInListViewDoesNotQuit(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = false

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.quitting)
}

func TestUpdateZoneSwitch(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = false // Simulate finished loading

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}}
	updated, cmd := m.Update(msg)
	m = updated.(model)

	assert.Equal(t, "is1a", m.currentZone) // Should cycle to next zone
	assert.Equal(t, 2, m.cursor)
	assert.True(t, m.loading) // Should start loading
	assert.NotNil(t, cmd)

	// Search should be cleared on zone switch
	assert.Equal(t, "", m.searchQuery)
	assert.Empty(t, m.searchMatches)
}

func TestUpdateEnterSearchMode(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.True(t, m.searchMode)
}

func TestSearchModeEscape(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.searchMode = true
	m.searchInput.SetValue("test")

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.searchMode)
	assert.Equal(t, "", m.searchInput.Value())
}

func TestSearchModeEnter(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.searchMode = true
	m.searchInput.SetValue("web")

	// Add some test servers
	servers := []Server{
		{ID: "1", Name: "web-server-1", InstanceStatus: "UP"},
		{ID: "2", Name: "db-server-1", InstanceStatus: "UP"},
		{ID: "3", Name: "web-server-2", InstanceStatus: "DOWN"},
	}
	items := make([]list.Item, len(servers))
	for i, s := range servers {
		items[i] = s
	}
	m.list.SetItems(items)

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.searchMode)
	assert.Equal(t, "web", m.searchQuery)
	assert.Equal(t, 2, len(m.searchMatches)) // Should find 2 web servers
	assert.Equal(t, 0, m.currentMatch)
	assert.Equal(t, 0, m.list.Index()) // Should jump to first match
}

func TestNextMatch(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	servers := []Server{
		{ID: "1", Name: "web-server-1", InstanceStatus: "UP"},
		{ID: "2", Name: "db-server-1", InstanceStatus: "UP"},
		{ID: "3", Name: "web-server-2", InstanceStatus: "DOWN"},
	}
	items := make([]list.Item, len(servers))
	for i, s := range servers {
		items[i] = s
	}
	m.list.SetItems(items)

	m.searchQuery = "web"
	m.searchMatches = []int{0, 2}
	m.currentMatch = 0

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.Equal(t, 1, m.currentMatch)
	assert.Equal(t, 2, m.list.Index()) // Should jump to second match

	// Test wrapping
	updated, _ = m.Update(msg)
	m = updated.(model)
	assert.Equal(t, 0, m.currentMatch) // Should wrap to first
	assert.Equal(t, 0, m.list.Index())
}

func TestPrevMatch(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	servers := []Server{
		{ID: "1", Name: "web-server-1", InstanceStatus: "UP"},
		{ID: "2", Name: "db-server-1", InstanceStatus: "UP"},
		{ID: "3", Name: "web-server-2", InstanceStatus: "DOWN"},
	}
	items := make([]list.Item, len(servers))
	for i, s := range servers {
		items[i] = s
	}
	m.list.SetItems(items)

	m.searchQuery = "web"
	m.searchMatches = []int{0, 2}
	m.currentMatch = 1

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.Equal(t, 0, m.currentMatch)
	assert.Equal(t, 0, m.list.Index())

	// Test wrapping backwards
	updated, _ = m.Update(msg)
	m = updated.(model)
	assert.Equal(t, 1, m.currentMatch) // Should wrap to last
	assert.Equal(t, 2, m.list.Index())
}

func TestPerformSearch(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	servers := []Server{
		{ID: "123", Name: "web-server-1", InstanceStatus: "UP"},
		{ID: "456", Name: "db-server-1", InstanceStatus: "UP"},
		{ID: "789", Name: "web-server-2", InstanceStatus: "DOWN"},
		{ID: "999", Name: "api-server", InstanceStatus: "UP"},
	}
	items := make([]list.Item, len(servers))
	for i, s := range servers {
		items[i] = s
	}
	m.list.SetItems(items)

	// Search by name
	m.searchQuery = "web"
	m.performSearch()
	assert.Equal(t, 2, len(m.searchMatches))
	assert.Contains(t, m.searchMatches, 0)
	assert.Contains(t, m.searchMatches, 2)

	// Search by ID
	m.searchQuery = "456"
	m.performSearch()
	assert.Equal(t, 1, len(m.searchMatches))
	assert.Equal(t, 1, m.searchMatches[0])

	// Case insensitive search
	m.searchQuery = "WEB"
	m.performSearch()
	assert.Equal(t, 2, len(m.searchMatches))
}

func TestViewQuitting(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.quitting = true

	output := m.View()
	assert.Equal(t, "Bye!\n", output)
}

func TestViewLoading(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = true

	output := m.View()
	assert.Contains(t, output, "sact")
	assert.Contains(t, output, "Loading servers...")
	assert.Contains(t, output, "tk1b") // Zone should be shown
}

func TestViewWithAccountName(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.accountName = "test-account"
	m.loading = false

	output := m.View()
	assert.Contains(t, output, "Account: test-account")
}

func TestViewSearchMode(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.searchMode = true
	m.loading = false

	output := m.View()
	assert.Contains(t, output, "/") // Search prompt
}

func TestViewWithSearchResults(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.searchQuery = "web"
	m.searchMatches = []int{0, 2, 5}
	m.currentMatch = 1
	m.loading = false

	output := m.View()
	assert.Contains(t, output, "Search: web (2/3)") // Should show current position
}

func TestViewZoneDisplay(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = false

	output := m.View()

	// Should show all zones
	assert.Contains(t, output, "tk1a")
	assert.Contains(t, output, "tk1b")
	assert.Contains(t, output, "is1a")
	assert.Contains(t, output, "is1b")
	assert.Contains(t, output, "is1c")

	// tk1b should be selected (indicated by brackets or highlighting)
	assert.Contains(t, output, "tk1b")
}

func TestServersLoadedMsg(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	servers := []Server{
		{ID: "1", Name: "server-1", InstanceStatus: "UP"},
		{ID: "2", Name: "server-2", InstanceStatus: "DOWN"},
	}

	msg := serversLoadedMsg{servers: servers, err: nil}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.loading)
	assert.Equal(t, 2, len(m.list.Items()))
}

func TestServersLoadedMsgWithError(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	msg := serversLoadedMsg{servers: nil, err: assert.AnError}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.loading)
	assert.NotNil(t, m.err)
}

func TestAuthStatusLoadedMsg(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	msg := authStatusLoadedMsg{accountName: "my-account", err: nil}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.Equal(t, "my-account", m.accountName)
}

func TestServerDelegateRender(t *testing.T) {
	// Test that delegate renders properly
	delegate := serverDelegate{}

	assert.Equal(t, 1, delegate.Height())
	assert.Equal(t, 0, delegate.Spacing())
}

func TestServerListItem(t *testing.T) {
	server := Server{
		ID:             "123456",
		Name:           "test-server",
		InstanceStatus: "UP",
		Zone:           "tk1b",
	}

	// Test list.Item interface implementation
	assert.Equal(t, "test-server", server.FilterValue())
	assert.Equal(t, "test-server", server.Title())
	assert.Contains(t, server.Description(), "123456")
	assert.Contains(t, server.Description(), "UP")
}

func TestNavigationKeysWork(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	servers := []Server{
		{ID: "1", Name: "server-1", InstanceStatus: "UP"},
		{ID: "2", Name: "server-2", InstanceStatus: "UP"},
		{ID: "3", Name: "server-3", InstanceStatus: "UP"},
	}
	items := make([]list.Item, len(servers))
	for i, s := range servers {
		items[i] = s
	}
	m.list.SetItems(items)
	m.loading = false

	// Test that j/k keys are delegated to list
	// The list component should handle navigation
	initialIndex := m.list.Index()

	// Simulate down key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updated, _ := m.Update(msg)
	m = updated.(model)

	// List should still be functional (index might change depending on list implementation)
	assert.NotNil(t, m.list)
	assert.GreaterOrEqual(t, m.list.Index(), initialIndex)
}

func TestRefreshKey(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = false

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updated, cmd := m.Update(msg)
	m = updated.(model)

	assert.True(t, m.loading)
	assert.NotNil(t, cmd)
}

func TestViewDoesNotPanic(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")

	// Test various states don't cause panics
	assert.NotPanics(t, func() { m.View() })

	m.loading = true
	assert.NotPanics(t, func() { m.View() })

	m.loading = false
	m.searchMode = true
	assert.NotPanics(t, func() { m.View() })

	m.searchMode = false
	m.searchQuery = "test"
	m.searchMatches = []int{0, 1}
	assert.NotPanics(t, func() { m.View() })
}

func TestMultipleZoneSwitches(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = false

	zones := []string{"tk1b", "is1a", "is1b", "is1c", "tk1a"}

	for i := range zones {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}}
		updated, _ := m.Update(msg)
		m = updated.(model)

		nextExpectedZone := zones[(i+1)%len(zones)]
		assert.Equal(t, nextExpectedZone, m.currentZone, "After switch %d", i)
	}
}

func TestEmptyServerList(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = false

	// Set empty server list
	m.list.SetItems([]list.Item{})

	// Should not panic
	assert.NotPanics(t, func() { m.View() })

	// Search on empty list should not panic
	m.searchQuery = "test"
	m.performSearch()
	assert.Equal(t, 0, len(m.searchMatches))
}

func TestEnterDetailMode(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.loading = false

	// Add a test server
	servers := []Server{
		{ID: "123", Name: "test-server", InstanceStatus: "UP", Zone: "tk1b"},
	}
	items := make([]list.Item, len(servers))
	for i, s := range servers {
		items[i] = s
	}
	m.list.SetItems(items)

	// Press Enter to show details
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(msg)
	m = updated.(model)

	assert.True(t, m.detailMode)
	assert.True(t, m.detailLoading)
	assert.NotNil(t, cmd)
}

func TestExitDetailMode(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.detailMode = true
	m.serverDetail = &ServerDetail{
		Server: Server{ID: "123", Name: "test", InstanceStatus: "UP", Zone: "tk1b"},
	}

	// Press ESC to exit detail mode
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.detailMode)
	assert.Nil(t, m.serverDetail)
}

func TestDetailViewRendering(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.detailMode = true
	m.detailLoading = true

	// Test loading state
	output := m.View()
	assert.Contains(t, output, "Loading server details...")

	// Test detail display
	m.detailLoading = false
	m.serverDetail = &ServerDetail{
		Server: Server{
			ID:             "123456",
			Name:           "test-server",
			InstanceStatus: "UP",
			Zone:           "tk1b",
		},
		CPU:             4,
		MemoryGB:        8,
		IPAddresses:     []string{"192.168.1.1"},
		UserIPAddresses: []string{"10.0.0.1", "10.0.0.2"},
		Disks: []DiskInfo{
			{Name: "disk-1", SizeGB: 100},
		},
		Tags:      []string{"production", "web"},
		CreatedAt: "2024-01-01 12:00:00",
	}

	output = m.View()
	assert.Contains(t, output, "test-server")
	assert.Contains(t, output, "123456")
	assert.Contains(t, output, "UP")
	assert.Contains(t, output, "4 Core(s)")
	assert.Contains(t, output, "8 GB")
	assert.Contains(t, output, "192.168.1.1")
	assert.Contains(t, output, "User IP:")
	assert.Contains(t, output, "10.0.0.1")
	assert.Contains(t, output, "disk-1")
	assert.Contains(t, output, "100 GB")
	assert.Contains(t, output, "production")
	assert.Contains(t, output, "ESC or q to go back")
}

func TestServerDetailLoadedMsg(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.detailMode = true
	m.detailLoading = true

	detail := &ServerDetail{
		Server: Server{ID: "123", Name: "test", InstanceStatus: "UP", Zone: "tk1b"},
		CPU:    2,
	}

	msg := serverDetailLoadedMsg{detail: detail, err: nil}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.detailLoading)
	assert.NotNil(t, m.serverDetail)
	assert.Equal(t, "test", m.serverDetail.Name)
}

func TestServerDetailLoadedMsgWithError(t *testing.T) {
	client, _ := NewSakuraClient("tk1b")
	m := InitialModel(client, "tk1b")
	m.detailMode = true
	m.detailLoading = true

	msg := serverDetailLoadedMsg{detail: nil, err: assert.AnError}
	updated, _ := m.Update(msg)
	m = updated.(model)

	assert.False(t, m.detailLoading)
	assert.False(t, m.detailMode) // Should exit detail mode on error
	assert.NotNil(t, m.err)
}
