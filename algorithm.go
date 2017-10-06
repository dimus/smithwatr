package smithwatr

import "log"

const different = 0
const similar = 1
const identical = 2

const diagonal = 0
const vertical = 1
const horisontal = 2

type Alignment struct {
	Seq1      []rune
	Seq2      []rune
	SeqLen1   int
	SeqLen2   int
	Score     int
	Identical int
	Similar   int
}

type MaxScore struct {
	Score int
	I     int
	J     int
}

type Score struct {
	Gap1    int
	Gap2    int
	Score   int
	Compass int
	Match   int
}

type ScoreMatrix [][]Score

// SmithWaterman calculates results of alignment for two peptide sequences
// according to Smith-Waterman algorithm. It takes aminoacid sequences of
// two proteins and returns result of calculation in a structure
func SmithWaterman(seq1 []rune, seq2 []rune, b62 Blosum62, conf Env) Alignment {
	var res Alignment
	res.Seq1 = seq1
	res.Seq2 = seq2
	res.SeqLen1 = len(seq1)
	res.SeqLen2 = len(seq2)
	matrix, max := res.calculateScoreMatrix(conf, b62)
	res.calculatePath(matrix, max)
	return res
}

func (a *Alignment) calculatePath(m ScoreMatrix, max MaxScore) {
	i := max.I
	j := max.J
	cell := m[i][j]
	score := max.Score
	for {
		if score > 0 {
			switch cell.Match {
			case identical:
				a.Identical += 1
			case similar:
				a.Similar += 1
			}

			switch cell.Compass {
			case diagonal:
				i -= 1
				j -= 1
			case vertical:
				j -= 1
			case horisontal:
				i -= 1
			}
			cell = m[i][j]
			score = cell.Score
		} else {
			break
		}
	}
}

func (a *Alignment) calculateScoreMatrix(conf Env,
	b62 Blosum62) (ScoreMatrix, MaxScore) {
	var max MaxScore
	matrix := make(ScoreMatrix, a.SeqLen1+1)
	for i := range matrix {
		matrix[i] = make([]Score, a.SeqLen2+1)
	}

	log.Println()
	for j := 1; j <= a.SeqLen2; j++ {
		for i := 1; i <= a.SeqLen1; i++ {
			log.Println(i, j)
			newScore, newCompass := calculateGap1(matrix, i, j, conf)
			newScore, newCompass = calculateGap2(matrix, i, j, newScore,
				newCompass, conf)
			score := calculateMatch(a, matrix, i, j, newScore, newCompass, b62)
			if max.Score < score {
				max.Score = score
				max.I = i
				max.J = j
			}

		}
	}
	return matrix, max
}

func calculateMatch(a *Alignment, matrix [][]Score, i int, j int, newScore int,
	newCompass int, b62 Blosum62) int {
	cellDiagonal := matrix[i-1][j-1]
	cell := matrix[i][j]
	r1 := a.Seq1[i-1]
	r2 := a.Seq2[j-1]
	gain := b62[r1][r2]
	if r1 == r2 {
		cell.Match = identical
	} else if gain > 0 {
		cell.Match = similar
	}

	if s := cellDiagonal.Score + gain; s > newScore {
		newScore = s
		newCompass = diagonal
	}

	if newScore < 0 {
		newScore = 0
	}

	cell.Score = newScore
	cell.Compass = newCompass
	return newScore
}

func calculateGap1(matrix [][]Score, i int, j int, conf Env) (int, int) {
	newCompass := 0
	newScore := 0
	log.Println("Matrix length", len(matrix[i]))
	cellUp := matrix[i][j-1]
	cell := matrix[i][j]
	if cellUp.Gap2 <= conf.GapExtends && cellUp.Score <= conf.GapOpens {
		cell.Gap2 = 0
	} else if cellUp.Gap2-conf.GapExtends > cellUp.Score-conf.GapOpens {
		newScore = cellUp.Gap2 - conf.GapExtends
		cell.Gap2 = newScore
		newCompass = horisontal
	} else {
		newScore = cellUp.Score - conf.GapOpens
		newCompass = horisontal
	}
	return newScore, newCompass
}

func calculateGap2(matrix [][]Score, i int, j int, newScore int,
	newCompass int, conf Env) (int, int) {
	cellLeft := matrix[i-1][j]
	cell := matrix[i][j]
	if cellLeft.Gap1 <= conf.GapExtends && cellLeft.Score <= conf.GapOpens {
		cell.Gap1 = 0
	} else if cellLeft.Gap1-conf.GapExtends > cellLeft.Score-conf.GapOpens {
		cell.Gap1 = cellLeft.Gap1 - conf.GapExtends
		if cell.Gap1 > newScore {
			newScore = cell.Gap1
			newCompass = vertical
		}
	} else {
		cell.Gap1 = cellLeft.Score - conf.GapOpens
		if cell.Gap1 > newScore {
			newScore = cell.Gap1
			newCompass = vertical
		}
	}
	return newScore, newCompass
}
