package cli

import (
	"fmt"

	"jgitra/internal/jira"
	"jgitra/internal/tui/message"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

type ValidateCmd struct {
}

type respUserMsg struct {
	user jira.UserResponse
}

func makeRequest(jira *jira.Client) tea.Cmd {
	return func() tea.Msg {
		u, err := jira.Myself()
		if err != nil {
			return message.ResponseError{Error: err}
		}
		return respUserMsg{u}
	}
}

type validateModel struct {
	sp      spinner.Model
	loading bool

	client *jira.Client
	user   jira.UserResponse
	err    error
}

func newValidateModel(client *jira.Client) validateModel {
	sp := spinner.New()
	sp.Spinner = spinner.MiniDot

	return validateModel{
		client:  client,
		sp:      sp,
		loading: true,
	}
}

func (m validateModel) Init() tea.Cmd {
	return tea.Batch(
		m.sp.Tick,
		makeRequest(m.client),
	)
}

func (m validateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		}

	case respUserMsg:
		m.loading = false
		m.user = msg.user
		return m, tea.Quit

	case message.ResponseError:
		m.loading = false
		m.err = msg.Error
		return m, tea.Quit
	}

	if m.loading {
		sp, spCmd := m.sp.Update(msg)
		m.sp = sp
		return m, spCmd
	}
	return m, nil
}

func (m validateModel) View() tea.View {
	view := tea.NewView(m.sp.View() + " validating config")
	return view
}

func (v *ValidateCmd) Run(jira *jira.Client) error {
	m, err := tea.NewProgram(newValidateModel(jira)).Run()
	if err == nil {
		vm := m.(validateModel)
		if vm.err != nil {
			return fmt.Errorf("validation error: %s", vm.err)
		}
		fmt.Printf("Validation succeeded!\nHello, %s\n", vm.user.DisplayName)
	}
	return err
}
