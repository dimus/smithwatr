CREATE TABLE genomes (
    id serial NOT NULL,
    file_name character varying(255) NOT NULL,
    species character varying(255) NOT NULL,
    CONSTRAINT titles_pkey PRIMARY KEY (id)
);

INSERT INTO genomes (file_name, genome) VALUES
('Araport11_genes.201606.pep.fasta.gz', 'Arabidopsis thaliana (L.) Heynh.'),
('Caenorhabditis_elegans.WBcel235.pep.all.fa.gz', 'Caenorhabditis elegans (Maupas, 1900)');
