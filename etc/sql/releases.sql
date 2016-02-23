DROP TABLE IF EXISTS releases;
DROP SEQUENCE IF EXISTS release_id_seq;

CREATE SEQUENCE release_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE releases (
    id bigint DEFAULT nextval('release_id_seq'::regclass) NOT NULL,
    cycle_id bigint NOT NULL,
    platform int NOT NULL,
    version bigint NOT NULL,
    note text DEFAULT '',
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    locked boolean DEFAULT FALSE NOT NULL
);

ALTER TABLE ONLY releases ADD CONSTRAINT releases_pkey PRIMARY KEY (id);
ALTER TABLE releases ADD FOREIGN KEY (cycle_id) REFERENCES cycles(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- # each platform can have their own versions.
ALTER TABLE releases ADD UNIQUE (cycle_id, platform, version);
