CREATE TABLE genes (
    id serial NOT NULL,
    genome_id int NOT NULL,
    gene character varying(255) NOT NULL,
    description text NOT NULL,
    sequence text NOT NULL,
    CONSTRAINT genes_pkey PRIMARY KEY (id)
);

CREATE INDEX gene_index ON genes USING btree (gene);
CREATE INDEX genome_id_index ON genes USING btree (genome_id, gene);
