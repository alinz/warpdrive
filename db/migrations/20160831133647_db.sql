
-- +goose Up

-- users table

CREATE SEQUENCE user_id_seq
    START WITH 2
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE users (
    id bigint DEFAULT nextval('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name varchar(128) NOT NULL,
    email varchar(256) NOT NULL,
    password varchar(128) NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE users ADD UNIQUE (email);

-- apps table

CREATE SEQUENCE app_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE apps (
    id bigint DEFAULT nextval('app_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name varchar(256) NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE apps ADD UNIQUE (name);

-- cycles table

CREATE SEQUENCE cycle_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE cycles (
    id bigint DEFAULT nextval('cycle_id_seq'::regclass) NOT NULL PRIMARY KEY,
    app_id bigint NOT NULL,
    name varchar(128) NOT NULL,
    public_key text NOT NULL,
    private_key text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE cycles ADD UNIQUE (app_id, name);
ALTER TABLE cycles ADD FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- releases table

CREATE SEQUENCE release_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE releases (
    id bigint DEFAULT nextval('release_id_seq'::regclass) NOT NULL PRIMARY KEY,
    cycle_id bigint NOT NULL,
    platform int NOT NULL,
    version varchar(256) NOT NULL,
    major bigint NOT NULL,
    minor bigint NOT NULL,
    patch bigint NOT NULL,
    build varchar(128) NOT NULL,
    note text DEFAULT '',
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    locked boolean DEFAULT FALSE NOT NULL
);

ALTER TABLE releases ADD FOREIGN KEY (cycle_id) REFERENCES cycles(id) ON DELETE CASCADE ON UPDATE CASCADE;
-- # each platform can have their own versions.
ALTER TABLE releases ADD UNIQUE (cycle_id, platform, version);

-- bundles table

CREATE SEQUENCE bundle_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE bundles (
    id bigint DEFAULT nextval('bundle_id_seq'::regclass) NOT NULL PRIMARY KEY,
    release_id bigint NOT NULL,
    hash varchar(128) NOT NULL,
    name text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE bundles ADD FOREIGN KEY ("release_id") REFERENCES releases("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- permissions table

CREATE SEQUENCE permission_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE permissions (
    id bigint DEFAULT nextval('permission_id_seq'::regclass) NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    app_id bigint NOT NULL
);

ALTER TABLE permissions ADD FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE permissions ADD FOREIGN KEY ("app_id") REFERENCES apps("id") ON DELETE CASCADE ON UPDATE CASCADE;
-- # each user must have only one permission with one app.
ALTER TABLE permissions ADD UNIQUE (user_id, app_id);

-- the password hash represents `root`
INSERT INTO users (id, name, email, password) VALUES (1, 'Mr. Robot', 'admin@pressly.com', '$2a$10$aWHlz4foaCzhIxRPwz.PeuA1c328upMUkc6iJqx5h4ggly1hY0DMS');

-- +goose Down

DROP TABLE IF EXISTS permissions;
DROP SEQUENCE IF EXISTS permission_id_seq;
DROP TABLE IF EXISTS bundles;
DROP SEQUENCE IF EXISTS bundle_id_seq;
DROP TABLE IF EXISTS releases;
DROP SEQUENCE IF EXISTS release_id_seq;
DROP TABLE IF EXISTS cycles;
DROP SEQUENCE IF EXISTS cycle_id_seq;
DROP TABLE IF EXISTS apps;
DROP SEQUENCE IF EXISTS app_id_seq;
DROP TABLE IF EXISTS users;
DROP SEQUENCE IF EXISTS user_id_seq;




