package gitea

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	sdk "code.gitea.io/sdk/gitea"
)

type Client struct {
	client *sdk.Client
	owner  string
	repo   string
}

func NewClient(urlStr, token, owner, repo string) (*Client, error) {
	if token == "" {
		return nil, errors.New("token must not be empty")
	}
	client, err := sdk.NewClient(urlStr, sdk.SetToken(token))
	if err != nil {
		return nil, fmt.Errorf("creating gitea client: %w", err)
	}
	return &Client{
		client: client,
		owner:  owner,
		repo:   repo,
	}, nil
}

func NewClientFromOrigin() (*Client, error) {
	token := os.Getenv("GITEA_TOKEN")
	if token == "" {
		return nil, errors.New("GITEA_TOKEN environment variable not set")
	}

	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return nil, fmt.Errorf("getting origin remote: %w", err)
	}

	originURL := strings.TrimSpace(string(out))
	baseURL, owner, repo, err := parseOriginURL(originURL)
	if err != nil {
		return nil, fmt.Errorf("parsing origin URL: %w", err)
	}

	return NewClient(baseURL, token, owner, repo)
}

func parseOriginURL(originURL string) (baseURL, owner, repo string, err error) {
	originURL = strings.TrimSpace(originURL)

	if strings.HasPrefix(originURL, "git@") {
		parts := strings.SplitN(originURL, ":", 2)
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid SSH origin URL: %s", originURL)
		}
		host := strings.TrimPrefix(parts[0], "git@")
		baseURL = "https://" + host

		path := strings.TrimSuffix(parts[1], ".git")
		pathParts := strings.SplitN(path, "/", 2)
		if len(pathParts) != 2 {
			return "", "", "", fmt.Errorf("cannot parse owner/repo from URL: %s", originURL)
		}
		return baseURL, pathParts[0], pathParts[1], nil
	}

	u, err := url.Parse(originURL)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid origin URL: %w", err)
	}
	baseURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, ".git")
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) != 2 {
		return "", "", "", fmt.Errorf("cannot parse owner/repo from URL: %s", originURL)
	}
	return baseURL, pathParts[0], pathParts[1], nil
}

func (c *Client) RepoOwner() string {
	return c.owner + "/" + c.repo
}
