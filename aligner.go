package smithwatr

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/lib/pq"
)

func Align(db *sql.DB, genomeTarget int, limit int, b62 Blosum62, conf Env) {
	mChan := make(chan Alignment)
	resChan := make(chan Alignment)
	var mWG sync.WaitGroup

	genesTarget := GetGenome(db, genomeTarget, limit)

	for i := 1; i <= conf.WorkersNum; i++ {
		mWG.Add(1)
		go matcherWorker(db, mWG, mChan, resChan, b62, conf)
	}

	go saveResults(db, resChan)

	count := 0
	for {
		count += 1
		gene := getAJob(db)
		if gene.ID > 0 {
			log.Printf("Alignment %d for %s, size %d", count, gene.Gene, gene.SeqLen)
			for _, g := range genesTarget {
				mChan <- Alignment{Gene1: gene, Gene2: g}
			}
			err := db.QueryRow(`UPDATE jobs
			                      SET status = 'finished'
														WHERE gene_id = $1`, gene.ID).Scan()
			if err != sql.ErrNoRows {
				Check(err)
			}
		} else {
			close(mChan)
			break
		}
	}
	mWG.Wait()
	close(resChan)
}

func getAJob(db *sql.DB) Gene {
	var id, genomeID int
	var gene, sequence string

	q1 := `SELECT id, genome_id, gene, sequence 
	        FROM genes g 
					  JOIN jobs j on j.gene_id = g.id
					WHERE j.status = 'pending'
					LIMIT 1`

	q2 := `UPDATE jobs
	         SET status = 'started'
					   WHERE gene_id = $1`

	err := db.QueryRow(q1).Scan(&id, &genomeID, &gene, &sequence)
	Check(err)

	seqRunes := []rune(sequence)

	err = db.QueryRow(q2, id).Scan()
	if err != sql.ErrNoRows {
		Check(err)
	}

	return Gene{ID: id, GenomeID: genomeID, Gene: gene, Seq: seqRunes,
		SeqLen: len(seqRunes)}
}

func saveResults(db *sql.DB, resChan <-chan Alignment) {
	res := make([]Alignment, 1000)
	i := 0
	k := 1
	for gm := range resChan {
		res[i] = gm
		i++
		if i%1000 == 0 {
			bulkSave(db, res)
			// log.Printf("%d: saved", i*k)
			i = 0
			k++
		}
	}
	bulkSave(db, res[0:i-1])
}

func bulkSave(db *sql.DB, gms []Alignment) {
	batch := gms
	columns := []string{"gene_id", "match_gene_id", "score", "identical_num",
		"similar_num", "ident_percent", "sim_percent"}
	transaction, err := db.Begin()
	Check(err)

	stmt, err := transaction.Prepare(pq.CopyIn("genes_matches", columns...))
	Check(err)

	for _, gm := range batch {
		ident, sim := gm.IdentitySimilarity()
		_, err = stmt.Exec(gm.Gene1.ID, gm.Gene2.ID, gm.Score, gm.Identical,
			gm.Similar, ident, sim)
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

func matcherWorker(db *sql.DB, mWG sync.WaitGroup, mChan <-chan Alignment,
	resChan chan<- Alignment, b62 Blosum62, conf Env) {
	defer mWG.Done()
	for g := range mChan {
		resChan <- SmithWaterman(g.Gene1, g.Gene2, b62, conf)
	}
}

func GetGenome(db *sql.DB, genome int, num int) []Gene {
	var ID, genomeID int
	var gene, sequence string
	var res []Gene
	q := `SELECT id, genome_id, gene, sequence
	       FROM genes where genome_id = $1`
	if num > 0 {
		q = fmt.Sprintf("%s LIMIT %d", q, num)
	}

	rows, err := db.Query(q, genome)
	Check(err)
	for rows.Next() {
		err := rows.Scan(&ID, &genomeID, &gene, &sequence)
		Check(err)
		seqRunes := []rune(sequence)
		gene := Gene{ID: ID, GenomeID: genomeID, Gene: gene, Seq: seqRunes,
			SeqLen: len(seqRunes)}
		res = append(res, gene)
	}
	err = rows.Close()
	Check(err)
	return res
}
