package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"jgitra/internal/config"
	"jgitra/internal/utils"
)

type ValidateCmd struct {
}

type user struct {
	AccountID    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

func (v *ValidateCmd) Run(c *config.Config) error {
	baseURL, err := url.ParseRequestURI(c.Jira.BaseURL)
	if err != nil {
		return err
	}
	myselfEndpoint := baseURL.ResolveReference(&url.URL{Path: "/rest/api/3/myself"})
	req, err := http.NewRequest(http.MethodGet, myselfEndpoint.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Jira.UserEmail, c.Jira.APIToken)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var u user
	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
			return fmt.Errorf("decode Jira response: %w", err)
		}

	case http.StatusUnauthorized:
		return errors.New("jira rejected authentication: check email and API token")

	case http.StatusForbidden:
		return errors.New("authentication succeeded, but the user has no permission")

	case http.StatusNotFound:
		return errors.New("jira API was not found: check the Jira URL")

	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf(
			"jira returned %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}
	fmt.Printf("valid config\njira response:\n%s\n", utils.PrettyJSON(u))
	return nil
}
