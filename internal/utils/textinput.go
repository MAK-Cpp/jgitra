package utils

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
)

func NewTextInput() textinput.Model {
	t := textinput.New()
	t.KeyMap.CharacterBackward = key.NewBinding()
	t.KeyMap.CharacterForward = key.NewBinding()
	t.KeyMap.WordBackward = key.NewBinding()
	t.KeyMap.WordForward = key.NewBinding()
	return t
}
