package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/gen2brain/beeep"

	"nv-pr-app/pkg/github/model"
	githubService "nv-pr-app/pkg/github/service"

	"nv-pr-app/internal/config"
	"nv-pr-app/internal/pull_request/repository"
)

type PullRequestService struct {
	querier  repository.Querier
	projects []config.Project
	github   *githubService.Service
	now      func() time.Time
}

func New(querier repository.Querier, projects []config.Project, github *githubService.Service) *PullRequestService {
	now := time.Now
	return &PullRequestService{
		querier:  querier,
		projects: projects,
		github:   github,
		now:      now,
	}
}

func (p *PullRequestService) NotifyPullRequests(ctx context.Context) error {
	pullRequests, err := p.querier.ListPullRequestsNotNotified(ctx)
	if err != nil {
		return err
	}

	for _, pr := range pullRequests {
		err = beeep.Notify("Alerta Pull Request", pr.Repository, "ic_launcher_foreground.png")
		if err != nil {
			return err
		}

		err = beeep.Alert("Alerta Pull Request", pr.Repository, "ic_launcher_foreground.png")
		if err != nil {
			return err
		}

		err = p.querier.SetNotified(ctx, repository.SetNotifiedParams{
			ID:        pr.ID,
			UpdatedAt: p.now(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PullRequestService) AddNonexistentPullRequests(ctx context.Context) error {
	for _, project := range p.projects {
		for _, repo := range project.Repos {
			pullRequests, err := p.github.ListPr(project.Owner, repo)
			if err != nil {
				return err
			}

			err = p.addPullRequestsByGithub(pullRequests, ctx, project, repo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PullRequestService) addPullRequestsByGithub(pullRequests []model.PullRequest, ctx context.Context, project config.Project, repo string) error {
	for _, pr := range pullRequests {
		pullRequest, err := p.getPullRequestFromDB(ctx, project, repo, pr)
		if err != nil {
			return err
		}
		if pullRequest == nil {
			err = p.createPullRequest(ctx, project.Owner, repo, int64(pr.Number))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PullRequestService) getPullRequestFromDB(ctx context.Context, project config.Project, repo string, pr model.PullRequest) (*repository.PullRequest, error) {
	pullRequest, err := p.querier.GetPullRequestByOwnerAndRepositoryAndNumber(ctx, repository.GetPullRequestByOwnerAndRepositoryAndNumberParams{
		Owner:      project.Owner,
		Repository: repo,
		Number:     int64(pr.Number),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pullRequest, nil
}

func (p *PullRequestService) createPullRequest(ctx context.Context, owner, repo string, number int64) error {
	err := p.querier.CreatePullRequest(ctx, repository.CreatePullRequestParams{
		Owner:      owner,
		Repository: repo,
		Number:     number,
		Notified:   false,
		CreatedAt:  p.now(),
		UpdatedAt:  p.now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PullRequestService) ListPullRequests(ctx context.Context, owner, repo string) ([]model.PullRequest, error) {
	pullRequests, err := p.github.ListPr(owner, repo)
	if err != nil {
		return nil, err
	}
	return pullRequests, nil
}
