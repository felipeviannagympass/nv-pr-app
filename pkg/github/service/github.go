package github

import (
	"fmt"

	"nv-pr-app/pkg/github/model"
	"nv-pr-app/pkg/nvhttp"

	"nv-pr-app/internal/config"
)

const (
	URL      = "api.github.com"
	PROTOCOL = "https"
)

type Service struct {
	gc *config.Github
}

func New(githubConfig *config.Github) *Service {
	return &Service{
		gc: githubConfig,
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

	pullRequests := []model.PullRequest{}
	err = nvhttp.UnmarshalBody(resp, &pullRequests)
	if err != nil {
		return nil, err
	}
	return pullRequests, nil
}
