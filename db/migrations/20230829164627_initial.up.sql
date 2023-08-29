CREATE TABLE IF NOT EXISTS pull_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner TEXT NOT NULL,
    repository TEXT NOT NULL,
    number int not null,
    notified boolean not null,
    created_at timestamp not null,
    updated_at timestamp not null
);