package smithwatr

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/lib/pq"
)

type GeneMatcher struct {
	Gene1        Gene
	Gene2        Gene
	Score        int
	Identical    int
	Similar      int
	IdentPercent float32
	SimPercent   float32
}

func Align(db *sql.DB, genome1 int, genome2 int, num int,
	b62 Blosum62, conf Env) {
	g1 := GetGenome(db, genome1, num)
	g2 := GetGenome(db, genome2, num)
	mChan := make(chan GeneMatcher)
	resChan := make(chan GeneMatcher)
	var mWG sync.WaitGroup

	for i := 1; i <= conf.WorkersNum; i++ {
		mWG.Add(1)
		go matcherWorker(db, mWG, mChan, resChan, b62, conf)
	}

	go saveResults(db, resChan)

	for _, gene1 := range g1 {
		for _, gene2 := range g2 {
			mChan <- GeneMatcher{Gene1: gene1, Gene2: gene2}
		}
	}
	close(mChan)

	mWG.Wait()
	close(resChan)
}

func saveResults(db *sql.DB, resChan <-chan GeneMatcher) {
	res := make([]GeneMatcher, 1000)
	i := 0
	k := 1
	for gm := range resChan {
		res[i] = gm
		i++
		if i%1000 == 0 {
			bulkSave(db, res)
			log.Printf("%d: saved", i*k)
			i = 0
			k++
		}
	}
	bulkSave(db, res[0:i-1])
}

func bulkSave(db *sql.DB, gms []GeneMatcher) {
	batch := gms
	columns := []string{"gene_id", "match_gene_id", "score", "identical_num",
		"similar_num", "ident_percent", "sim_percent"}
	transaction, err := db.Begin()
	Check(err)

	stmt, err := transaction.Prepare(pq.CopyIn("genes_matches", columns...))
	Check(err)

	for _, gm := range batch {
		_, err = stmt.Exec(gm.Gene1.ID, gm.Gene2.ID, gm.Score, gm.Identical,
			gm.Similar, gm.IdentPercent, gm.SimPercent)
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

func matcherWorker(db *sql.DB, mWG sync.WaitGroup, mChan <-chan GeneMatcher,
	resChan chan<- GeneMatcher, b62 Blosum62, conf Env) {
	defer mWG.Done()
	for g := range mChan {
		res := SmithWaterman([]rune(g.Gene1.Seq), []rune(g.Gene2.Seq), b62, conf)
		g.Score = res.Score
		g.Identical = res.Identical
		g.Similar = res.Similar
		g.IdentPercent, g.SimPercent = res.IdentitySimilarity()
		resChan <- g
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

		gene := Gene{ID: ID, GenomeID: genomeID, Gene: gene, Seq: sequence}
		res = append(res, gene)
	}
	err = rows.Close()
	Check(err)
	return res
}