DROP TYPE IF EXISTS job_status;

CREATE TYPE job_status AS ENUM ('pending', 'started', 'finished') ;

CREATE TABLE jobs (
    gene_id int NOT NULL,
    status job_status NOT NULL DEFAULT 'pending',
    CONSTRAINT jobs_pkey PRIMARY KEY (gene_id)
);

CREATE INDEX status_index ON jobs USING btree (status);
