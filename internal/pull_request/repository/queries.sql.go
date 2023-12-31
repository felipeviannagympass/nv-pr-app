// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: queries.sql

package repository

import (
	"context"
	"time"
)

const createPullRequest = `-- name: CreatePullRequest :exec
INSERT INTO pull_requests (
  owner, repository, number, notified, created_at, updated_at 
) VALUES (
  ?, ?, ?, ?, ?, ?
)
`

type CreatePullRequestParams struct {
	Owner      string
	Repository string
	Number     int64
	Notified   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) CreatePullRequest(ctx context.Context, arg CreatePullRequestParams) error {
	_, err := q.exec(ctx, q.createPullRequestStmt, createPullRequest,
		arg.Owner,
		arg.Repository,
		arg.Number,
		arg.Notified,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getPullRequestByOwnerAndRepositoryAndNumber = `-- name: GetPullRequestByOwnerAndRepositoryAndNumber :one
SELECT id, owner, repository, number, notified, created_at, updated_at FROM pull_requests
where owner = ?
and repository = ?
and number = ?
`

type GetPullRequestByOwnerAndRepositoryAndNumberParams struct {
	Owner      string
	Repository string
	Number     int64
}

func (q *Queries) GetPullRequestByOwnerAndRepositoryAndNumber(ctx context.Context, arg GetPullRequestByOwnerAndRepositoryAndNumberParams) (PullRequest, error) {
	row := q.queryRow(ctx, q.getPullRequestByOwnerAndRepositoryAndNumberStmt, getPullRequestByOwnerAndRepositoryAndNumber, arg.Owner, arg.Repository, arg.Number)
	var i PullRequest
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Repository,
		&i.Number,
		&i.Notified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPullRequests = `-- name: ListPullRequests :many
SELECT id, owner, repository, number, notified, created_at, updated_at FROM pull_requests
ORDER BY owner, repository, number
`

func (q *Queries) ListPullRequests(ctx context.Context) ([]PullRequest, error) {
	rows, err := q.query(ctx, q.listPullRequestsStmt, listPullRequests)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PullRequest
	for rows.Next() {
		var i PullRequest
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Repository,
			&i.Number,
			&i.Notified,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPullRequestsNotNotified = `-- name: ListPullRequestsNotNotified :many
SELECT id, owner, repository, number, notified, created_at, updated_at FROM pull_requests
where notified = false
ORDER BY owner, repository, number
`

func (q *Queries) ListPullRequestsNotNotified(ctx context.Context) ([]PullRequest, error) {
	rows, err := q.query(ctx, q.listPullRequestsNotNotifiedStmt, listPullRequestsNotNotified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PullRequest
	for rows.Next() {
		var i PullRequest
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Repository,
			&i.Number,
			&i.Notified,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setNotified = `-- name: SetNotified :exec
UPDATE pull_requests
set notified = true
, updated_at = ?
where id = ?
`

type SetNotifiedParams struct {
	UpdatedAt time.Time
	ID        int64
}

func (q *Queries) SetNotified(ctx context.Context, arg SetNotifiedParams) error {
	_, err := q.exec(ctx, q.setNotifiedStmt, setNotified, arg.UpdatedAt, arg.ID)
	return err
}
