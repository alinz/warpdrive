DROP TABLE IF EXISTS cycles;
DROP SEQUENCE IF EXISTS cycle_id_seq;

CREATE SEQUENCE cycle_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE cycles (
    id bigint DEFAULT nextval('cycle_id_seq'::regclass) NOT NULL,
    app_id bigint NOT NULL,
    name varchar(128) NOT NULL,
    public_key text NOT NULL,
    private_key text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);

ALTER TABLE ONLY cycles ADD CONSTRAINT cycles_pkey PRIMARY KEY (id);
ALTER TABLE cycles ADD UNIQUE ("name");

ALTER TABLE cycles ADD FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE ON UPDATE CASCADE;
