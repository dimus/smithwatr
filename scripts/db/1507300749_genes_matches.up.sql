CREATE TABLE genes_matches (
    gene_id int NOT NULL,
    match_gene_id int NOT NULL,
    score int NOT NULL,
    percentage float NOT NULL,
    CONSTRAINT genes_matches_pkey PRIMARY KEY (gene_id, match_gene_id)
);
