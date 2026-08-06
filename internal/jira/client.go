package jira

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
)

var (
	UnauthorizedError = errors.New("jira rejected authentication: check email and API token")
	ForbiddenError    = errors.New("authentication succeeded, but the user has no permission")
	WrongPathError    = errors.New("jira API was not found: check the Jira URL")
)

type Client struct {
	config *config.Config
	api    string
	client *http.Client
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
		api:    config.Jira.BaseURL + "/rest/api/3",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) execute(method string, path string, body io.Reader, v any) error {
	fullPath, err := url.JoinPath(c.api, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.config.Jira.UserEmail, c.config.Jira.APIToken)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		if v == nil {
			return nil
		}
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil { // TODO: json??
			return fmt.Errorf("decode Jira response: %w", err)
		}

	case http.StatusUnauthorized:
		return UnauthorizedError

	case http.StatusForbidden:
		return ForbiddenError

	case http.StatusNotFound:
		return WrongPathError

	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf(
			"jira returned %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}
	return nil
}

func (c *Client) get(path string, v any) error {
	return c.execute(http.MethodGet, path, nil, v)
}

type UserResponse struct {
	AccountID    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

func (c *Client) Myself() (UserResponse, error) {
	var user UserResponse
	if err := c.get("/myself", &user); err != nil {
		return UserResponse{}, err
	}
	return user, nil
}

type ProjectsResponse struct {
	Projects []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"values"`
}

func (c *Client) ProjectSearch() (ProjectsResponse, error) {
	var projects ProjectsResponse
	if err := c.get("/project/search", &projects); err != nil {
		return ProjectsResponse{}, err
	}
	return projects, nil
}
