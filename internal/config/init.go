package config

import (
	"strings"

	"jgitra/internal/utils"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbletea/v2"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
}

func (km keyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Up, km.Down, km.Enter}
}

func (km keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Up, km.Down},
		{km.Enter},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎", "save config"),
	),
}

type initModel struct {
	conf   *Config
	inputs []textinput.Model
	help   help.Model
	pos    int
}

func newInitModel(conf *Config) initModel {
	url := utils.NewTextInput()
	url.Placeholder = "https://example.atlassian.net"
	url.SetWidth(1000)
	url.CharLimit = 1000
	url.Prompt = "Base url of your jira: "
	url.SetValue(conf.Jira.BaseURL)
	url.Focus()

	email := utils.NewTextInput()
	email.Placeholder = "example@mail.com"
	email.SetWidth(1000)
	email.CharLimit = 1000
	email.Prompt = "Email address: "
	email.SetValue(conf.Jira.UserEmail)

	token := utils.NewTextInput()
	token.EchoMode = textinput.EchoPassword
	token.EchoCharacter = '•'
	token.Prompt = "API token: "
	token.SetValue(conf.Jira.APIToken)

	help := help.New()

	return initModel{
		conf: conf,
		inputs: []textinput.Model{
			url,
			email,
			token,
		},
		help: help,
	}
}

func (m initModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "up":
			m.inputs[m.pos].Blur()
			m.pos = max(0, m.pos-1)
			m.inputs[m.pos].Focus()
		case "down":
			m.inputs[m.pos].Blur()
			m.pos = min(len(m.inputs)-1, m.pos+1)
			m.inputs[m.pos].Focus()
		case "enter":
			return m, tea.Quit
		}
	}

	inputs := make([]textinput.Model, len(m.inputs))
	cmds := make([]tea.Cmd, len(m.inputs))

	for i, input := range m.inputs {
		inputs[i], cmds[i] = input.Update(msg)
		if inputs[i].Focused() {
			switch i {
			case 0:
				m.conf.Jira.BaseURL = inputs[i].Value()
			case 1:
				m.conf.Jira.UserEmail = inputs[i].Value()
			case 2:
				m.conf.Jira.APIToken = inputs[i].Value()
			}
		}
	}
	m.inputs = inputs

	return m, tea.Batch(cmds...)
}

func (m initModel) View() tea.View {
	var content strings.Builder
	content.WriteString("Initialization configuration\n\n")

	for _, input := range m.inputs {
		content.WriteString(input.View())
		content.WriteRune('\n')
	}

	content.WriteRune('\n')
	content.WriteString(m.help.View(keys))

	view := tea.NewView(content.String())
	view.AltScreen = true
	view.WindowTitle = "Init config"
	return view
}

func (c *Config) Init() error {
	_, err := tea.NewProgram(newInitModel(c)).Run()
	return err
}
