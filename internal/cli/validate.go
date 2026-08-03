package cli

import (
	"fmt"

	"jgitra/internal/client"
)

type ValidateCmd struct {
}

func (v *ValidateCmd) Run(jira *client.JiraClient) error {
	u, err := jira.Myself()
	if err == nil {
		fmt.Printf("valid config\nHello, %s\n", u.DisplayName)
	}
	return err
}
