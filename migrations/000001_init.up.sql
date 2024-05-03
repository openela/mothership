-- Copyright 2024 The Mothership Authors
-- SPDX-License-Identifier: Apache-2.0

CREATE TABLE workers
(
    name              VARCHAR(255) PRIMARY KEY,
    worker_id         VARCHAR(255) UNIQUE NOT NULL,
    create_time       TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    last_checkin_time TIMESTAMPTZ,
    api_secret        VARCHAR(255)        NOT NULL
);

CREATE TABLE entries
(
    name            VARCHAR(255) PRIMARY KEY,
    entry_id        VARCHAR(255) NOT NULL,
    create_time     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    os_release      TEXT         NOT NULL,
    sha256_sum      VARCHAR(255) NOT NULL,
    repository_name VARCHAR(255) NOT NULL,
    worker_id       VARCHAR(255) REFERENCES workers (worker_id),
    batch_name      VARCHAR(255),
    user_email      TEXT,
    commit_uri      TEXT         NOT NULL,
    commit_hash     TEXT         NOT NULL,
    commit_branch   TEXT         NOT NULL,
    commit_tag      TEXT         NOT NULL,
    state           NUMERIC      NOT NULL,
    package_name    TEXT         NOT NULL
);

CREATE TABLE batches
(
    name           VARCHAR(255) PRIMARY KEY,
    batch_id       VARCHAR(255) UNIQUE,
    create_time    TIMESTAMPTZ                         NOT NULL DEFAULT NOW(),
    update_time    TIMESTAMPTZ                         NOT NULL DEFAULT NOW(),
    seal_time      TIMESTAMPTZ,
    worker_id      TEXT REFERENCES workers (worker_id) NOT NULL,
    bugtracker_uri TEXT
);

CREATE TABLE bugtracker_configs
(
    name        VARCHAR(255) PRIMARY KEY,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    config      JSONB       NOT NULL
);