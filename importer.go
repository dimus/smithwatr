package smithwatr

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lib/pq"
)

type Gene struct {
	ID       int
	GenomeID int
	Gene     string
	Desc     string
	Seq      string
}

func ImportData(db *sql.DB, conf Env) {
	d, err := os.Open(conf.DataDir)
	Check(err)
	defer func() {
		err := d.Close()
		Check(err)
	}()

	if GenesEmpty(db) {
		names, err := d.Readdirnames(-1)
		Check(err)
		for _, name := range names {
			l := len(name)
			if name[l-3:l] == ".gz" {
				genomeID := getGenomeID(db, name)
				path := filepath.Join(conf.DataDir, name)
				processFile(db, path, genomeID)
			}
		}
	}
}

func GenesEmpty(db *sql.DB) bool {
	var count int
	q := `SELECT count(*) from genes`
	err := db.QueryRow(q).Scan(&count)
	Check(err)
	if count == 0 {
		return true
	}
	return false
}

func getGenomeID(db *sql.DB, name string) int {
	var id int
	q := `SELECT id FROM genomes
	        WHERE file_name = $1`
	err := db.QueryRow(q, &name).Scan(&id)
	Check(err)
	return (id)
}

func processFile(db *sql.DB, path string, genomeID int) {
	f, err := os.Open(path)
	Check(err)
	gz, err := gzip.NewReader(f)
	Check(err)
	scanner := bufio.NewScanner(gz)
	genes := collectGenes(scanner, genomeID)
	saveGenes(db, genes)
}

func collectGenes(scanner *bufio.Scanner, genomeID int) []Gene {
	gene := Gene{}
	var seq []string
	res := []Gene{}
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '>' {
			if gene.Gene != "" {
				gene.Seq = joinSequence(seq)
				res = append(res, gene)
			}
			geneName, description := parseGeneHeader(line)
			gene = Gene{GenomeID: genomeID, Gene: geneName, Desc: description}
			seq = []string{}
		} else {
			seq = append(seq, strings.Trim(line, "\n\r"))
		}
	}
	err := scanner.Err()
	Check(err)
	return res
}

func joinSequence(seq []string) string {
	return strings.Join(seq, "")
}

func parseGeneHeader(line string) (string, string) {
	line = strings.Trim(line, "> \n\r")
	header := strings.SplitN(line, " ", 2)
	return header[0], header[1]
}

func saveGenes(db *sql.DB, genes []Gene) {
	batch := genes
	columns := []string{"genome_id", "gene", "description", "sequence"}
	transaction, err := db.Begin()
	Check(err)

	stmt, err := transaction.Prepare(pq.CopyIn("genes", columns...))
	Check(err)

	for _, p := range batch {
		_, err = stmt.Exec(p.GenomeID, p.Gene, p.Desc, p.Seq)
		Check(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(`
Bulk import of titles data failed, probably you need to empty all data
and start with an empty database.
`)
		log.Fatal(err)
	}

	err = stmt.Close()
	Check(err)

	err = transaction.Commit()
	Check(err)
}
