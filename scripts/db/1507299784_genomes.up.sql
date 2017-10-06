CREATE TABLE genomes (
    id serial NOT NULL,
    file_name character varying(255) NOT NULL,
    species character varying(255) NOT NULL,
    species_short character varying(255) NOT NULL,
    CONSTRAINT genomes_pkey PRIMARY KEY (id)
);

INSERT INTO genomes (file_name, species, species_short) VALUES
('Araport11_genes.201606.pep.fasta.gz', 'Arabidopsis thaliana (L.) Heynh.', 'Arabidopsis'),
('Caenorhabditis_elegans.WBcel235.pep.all.fa.gz', 'Caenorhabditis elegans (Maupas, 1900)', 'C. elegans'),
('Homo_sapiens.GRCh38.pep.all.fa.gz', 'Homo sapiens L. 1958', 'Homo sapiens');
