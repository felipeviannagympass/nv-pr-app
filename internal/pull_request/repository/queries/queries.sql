-- name: ListPullRequests :many
SELECT * FROM pull_requests
ORDER BY owner, repository, number;

-- name: ListPullRequestsNotNotified :many
SELECT * FROM pull_requests
where notified = false
ORDER BY owner, repository, number;

-- name: GetPullRequestByOwnerAndRepositoryAndNumber :one
SELECT * FROM pull_requests
where owner = ?
and repository = ?
and number = ?;

-- name: CreatePullRequest :exec
INSERT INTO pull_requests (
  owner, repository, number, notified, created_at, updated_at 
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: SetNotified :exec
UPDATE pull_requests
set notified = true
, updated_at = ?
where id = ?;

