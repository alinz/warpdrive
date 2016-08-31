
-- +goose Up

-- users table

CREATE TABLE users (
    id bigint PRIMARY KEY,
    name varchar(128) NOT NULL,
    email varchar(256) NOT NULL,
    password varchar(128) NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE users ADD UNIQUE ("email");

-- apps table

CREATE TABLE apps (
    id bigint PRIMARY KEY,
    name varchar(256) NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE apps ADD UNIQUE ("name");

-- cycles table

CREATE TABLE cycles (
    id bigint PRIMARY KEY,
    app_id bigint NOT NULL,
    name varchar(128) NOT NULL,
    public_key text NOT NULL,
    private_key text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE cycles ADD UNIQUE ("name");
ALTER TABLE cycles ADD FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- releases table

CREATE TABLE releases (
    id bigint PRIMARY KEY,
    cycle_id bigint NOT NULL,
    platform int NOT NULL,
    version bigint NOT NULL,
    note text DEFAULT '',
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    locked boolean DEFAULT FALSE NOT NULL
);

ALTER TABLE releases ADD FOREIGN KEY (cycle_id) REFERENCES cycles(id) ON DELETE CASCADE ON UPDATE CASCADE;
-- # each platform can have their own versions.
ALTER TABLE releases ADD UNIQUE (cycle_id, platform, version);

-- bundles table

CREATE TABLE bundles (
    id bigint PRIMARY KEY,
    release_id bigint NOT NULL,
    hash varchar(128) NOT NULL,
    name varchar(1024) NOT NULL,
    type int NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE bundles ADD FOREIGN KEY ("release_id") REFERENCES releases("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- permissions table

CREATE TABLE permissions (
    id bigint PRIMARY KEY,
    user_id bigint NOT NULL,
    app_id bigint NOT NULL,
    permission int NOT NULL
);

ALTER TABLE permissions ADD FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE permissions ADD FOREIGN KEY ("app_id") REFERENCES apps("id") ON DELETE CASCADE ON UPDATE CASCADE;
-- # each user must have only one permission with one app.
ALTER TABLE permissions ADD UNIQUE (user_id, app_id);

-- the password hash represents `root`
INSERT INTO users (id, name, email, password) VALUES (1, 'Mr. Robot', 'root', '$2a$10$u0NH2a95xrx83EZQhZ5nNesbKrjAMW3GZRe3MZHjXlB.Hqca.nrca');

-- +goose Down

DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS bundles;
DROP TABLE IF EXISTS releases;
DROP TABLE IF EXISTS cycles;
DROP TABLE IF EXISTS apps;
DROP TABLE IF EXISTS users;
