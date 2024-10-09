/*
  sql for postgresql
*/

/* CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; */

CREATE TABLE request
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    payload json NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    job_id uuid NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE schedule
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    payload json NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    job_id uuid NOT NULL,
    PRIMARY KEY (id)
);
