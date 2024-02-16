package github

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"nv-pr-app/pkg/github/model"
	"nv-pr-app/pkg/nvhttp"

	"nv-pr-app/internal/config"
)

const (
	URL        = "api.github.com"
	PROTOCOL   = "https"
	URL_TOKENS = "https://github.com/settings/tokens"
)

type Service struct {
	gc            *config.Github
	alreadyOpened bool
}

func New(githubConfig *config.Github) *Service {
	return &Service{
		gc:            githubConfig,
		alreadyOpened: false,
	}
}

func (g *Service) ListPr(owner, repo string) ([]model.PullRequest, error) {
	url := fmt.Sprintf("%s://%s/repos/%s/%s/pulls", PROTOCOL, URL, owner, repo)

	nvHttpClient := nvhttp.New(g.gc.AccessToken)

	resp, err := nvHttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	prs, err := g.CheckAuthorizeSSO(resp.StatusCode, owner, repo)
	if err != nil {
		return nil, err
	}
	if prs != nil {
		return prs, nil
	}

	pullRequests := []model.PullRequest{}
	err = nvhttp.UnmarshalBody(resp, &pullRequests)
	if err != nil {
		return nil, err
	}
	return pullRequests, nil
}

func (g *Service) CheckAuthorizeSSO(statusCode int, owner, repo string) ([]model.PullRequest, error) {
	if statusCode == 403 {
		if !g.alreadyOpened {
			cmd := exec.Command(g.getCommand(), URL_TOKENS)
			err := cmd.Start()
			fmt.Println(err)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		}
		g.alreadyOpened = true
		time.Sleep(10 * time.Second)
		return g.ListPr(owner, repo)
	}
	return nil, nil
}

func (g *Service) getCommand() string {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "start"
	default:
		cmd = "xdg-open"
	}
	return cmd
}
