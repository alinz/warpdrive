DROP TABLE IF EXISTS permissions;
DROP SEQUENCE IF EXISTS permission_id_seq;

CREATE SEQUENCE permission_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE permissions (
    id bigint DEFAULT nextval('permission_id_seq'::regclass) NOT NULL,
    user_id bigint NOT NULL,
    app_id bigint NOT NULL,
    permission int NOT NULL
);

ALTER TABLE ONLY permissions ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);
ALTER TABLE permissions ADD FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE permissions ADD FOREIGN KEY ("app_id") REFERENCES apps("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- # each user must have only one permission with one app.
ALTER TABLE permissions ADD UNIQUE (user_id, app_id);
