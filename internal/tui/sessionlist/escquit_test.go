package sessionlist

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eleonorayaya/utena/internal/claude"
	"github.com/eleonorayaya/utena/internal/git"
	"github.com/eleonorayaya/utena/internal/session"
	"github.com/eleonorayaya/utena/internal/tui/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var escKey = tea.KeyMsg{Type: tea.KeyEsc}

func escQuitList(t *testing.T, ss ...session.Session) Model {
	t.Helper()
	m, _ := New(WithEscQuits()).Update(tea.WindowSizeMsg{Width: 120, Height: 20})
	m, _ = m.Update(provider.SessionsStateUpdatedMsg{Sessions: ss})
	return m
}

func TestEsc_DoesNotQuitByDefault(t *testing.T) {
	m := listWith(t, withClaude(1, "a", claude.StatusIdle))

	_, cmd, _ := m.OnKeyMsg(escKey)

	assert.Nil(t, cmd)
}

func TestEsc_QuitsWithEscQuits(t *testing.T) {
	m := escQuitList(t, withClaude(1, "a", claude.StatusIdle))

	_, cmd, handled := m.OnKeyMsg(escKey)

	assert.True(t, handled)
	require.NotNil(t, cmd)
	assert.Equal(t, tea.QuitMsg{}, cmd())
}

func TestEsc_ClosesPRPickerBeforeQuitting(t *testing.T) {
	m := escQuitList(t, withPRs(withClaude(1, "a", claude.StatusIdle),
		git.PullRequest{Number: 1, State: git.PRStateOpen, HTMLURL: "u1"},
		git.PullRequest{Number: 2, State: git.PRStateOpen, HTMLURL: "u2"}))
	m, _, _ = pressO(m)
	require.True(t, m.picker.IsOpen())

	m, cmd, _ := m.OnKeyMsg(escKey)

	assert.Nil(t, cmd)
	assert.False(t, m.picker.IsOpen())
}

func TestEscQuits_HelpShowsEsc(t *testing.T) {
	assert.Equal(t, "q/esc", New(WithEscQuits()).keys().Quit.Help().Key)
	assert.Equal(t, "q", New().keys().Quit.Help().Key)
}
