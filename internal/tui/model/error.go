package model

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type Error struct {
	Err error
}

func (m Error) Init() tea.Cmd {
	panic(fmt.Sprintf("Cannot use model.Error, it's needs for returning errors from models. Error: %s", m.Err))
}

func (m Error) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	panic(fmt.Sprintf("Cannot use model.Error, it's needs for returning errors from models. Error: %s", m.Err))
}

func (m Error) View() tea.View {
	panic(fmt.Sprintf("Cannot use model.Error, it's needs for returning errors from models. Error: %s", m.Err))
}
