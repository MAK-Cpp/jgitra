package model

import tea "charm.land/bubbletea/v2"

type ErrorModel struct {
	Error error
}

func Error(err error) ErrorModel {
	return ErrorModel{
		Error: err,
	}
}

func (m ErrorModel) Init() tea.Cmd {
	panic("Error model")
}

func (m ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	panic("Error model")
}

func (m ErrorModel) View() tea.View {
	panic("Error model")
}
