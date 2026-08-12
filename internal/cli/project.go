package cli

import (
	"fmt"
	"strings"

	"jgitra/internal/config"
	"jgitra/internal/jira"
	"jgitra/internal/tui/message"
	"jgitra/internal/tui/model"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ProjectCmd struct {
}

type projectsLoadingModel struct {
	client        *jira.Client
	chosenProject string
	sp            spinner.Model
}

func newProjectLoadingModel(client *jira.Client, project string) projectsLoadingModel {
	sp := spinner.New()
	sp.Spinner = spinner.MiniDot

	return projectsLoadingModel{
		client:        client,
		chosenProject: project,
		sp:            sp,
	}
}

func (m projectsLoadingModel) getProjects() tea.Msg {
	projects, err := m.client.ProjectSearch()
	if err != nil {
		return message.ResponseError{Error: err}
	}
	return projects
}

func (m projectsLoadingModel) Init() tea.Cmd {
	return tea.Batch(
		m.sp.Tick,
		m.getProjects,
	)
}

func (m projectsLoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		}

	case message.ResponseError:
		return model.Error(msg.Error), tea.Quit

	case jira.ProjectsResponse:
		return newProjectsTableModel(msg, m.chosenProject), nil
	}

	sp, spCmd := m.sp.Update(msg)
	m.sp = sp
	return m, spCmd
}

func (m projectsLoadingModel) View() tea.View {
	view := tea.NewView(m.sp.View() + " loading projects")
	return view
}

type projectsTableModel struct {
	chosenProject string
	t             table.Model
}

func newProjectsTableModel(projects jira.ProjectsResponse, chosenProject string) projectsTableModel {
	rows := make([]table.Row, len(projects.Projects))
	cursor := -1

	for i, project := range projects.Projects {
		rows[i] = table.Row{project.Key, project.Name}
		if chosenProject == project.Key {
			cursor = i
		}
	}

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Key", Width: 15},
			{Title: "Name", Width: 15},
		}),
		table.WithRows(rows),
		table.WithHeight(2),
		table.WithWidth(50),
		table.WithFocused(true),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		Border(lipgloss.NormalBorder(), true, false)
	s.Selected = s.Selected.SetString(">")
	t.SetStyles(s)
	if cursor != -1 {
		t.SetCursor(cursor)
	}

	return projectsTableModel{
		chosenProject: chosenProject,
		t:             t,
	}
}

func (m projectsTableModel) Init() tea.Cmd {
	return nil
}

func (m projectsTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "enter":
			m.chosenProject = m.t.SelectedRow()[0]
			return m, tea.Quit
		}
	}
	sp, spCmd := m.t.Update(msg)
	m.t = sp
	return m, spCmd
}

func (m projectsTableModel) View() tea.View {
	var sb strings.Builder
	sb.WriteString(m.t.View())
	sb.WriteString("\n")
	sb.WriteString(m.t.HelpView())
	view := tea.NewView(sb.String())
	return view
}

func (p *ProjectCmd) Run(client *jira.Client, config *config.Config) error {
	m, err := tea.NewProgram(newProjectLoadingModel(client, config.Jira.Project)).Run()
	if err == nil {
		switch m := m.(type) {
		case model.ErrorModel:
			return fmt.Errorf("setting project error: %s", m.Error)
		case projectsTableModel:
			config.Jira.Project = m.chosenProject
			err = config.Save()
			if err == nil {
				fmt.Printf("Chosen project: %s\n", m.chosenProject)
			}
			return err
		default:
			return fmt.Errorf("unexpected model type: %T", m)
		}
	}
	return err
}
