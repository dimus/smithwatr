package smithwatr

import (
	"fmt"
)

const different = 0
const similar = 1
const identical = 2

const substitution = 0
const deletion = 1
const insertion = 2

type Match struct {
	L1    rune
	L2    rune
	I     int
	J     int
	Score int
	Type  int
	Subst int
}

type Alignment struct {
	Gene1     Gene
	Gene2     Gene
	Score     int
	Identical int
	Similar   int
	Path      []Match
}

type MaxScore struct {
	Score int
	I     int
	J     int
}

type ScoreMatrix [][]int

// SmithWaterman calculates results of alignment for two peptide sequences
// according to Smith-Waterman algorithm. It takes aminoacid sequences of
// two proteins and returns result of calculation in a structure
func SmithWaterman(g1 Gene, g2 Gene, b62 Blosum62, conf Env) Alignment {
	var res Alignment
	res.Gene1 = g1
	res.Gene2 = g2
	matrix, max := res.calculateScoreMatrix(conf, b62)
	res.Score = max.Score
	res.calculatePath(matrix, max)
	return res
}

func (a *Alignment) Show(cols int) string {
	res := []rune{'\n'}
	res = append(res, []rune(a.Gene1.Gene)...)
	res = append(res, []rune(fmt.Sprintf(" length: %d\n", a.Gene1.SeqLen))...)
	res = append(res, formatSeq(a.Gene1.Seq, cols)...)
	res = append(res, []rune(a.Gene2.Gene)...)
	res = append(res, []rune(fmt.Sprintf(" length: %d\n", a.Gene2.SeqLen))...)
	res = append(res, formatSeq(a.Gene2.Seq, cols)...)
	res = append(res, '\n')
	res = append(res, []rune(fmt.Sprintf("Score: %d, ", a.Score))...)
	ident, sim := a.IdentitySimilarity()
	res = append(res, []rune(fmt.Sprintf("Identical: %0.1f, Similar: %0.1f",
		ident, sim))...)
	res = append(res, '\n', '\n')

	i := 0
	offset := 0
	for {
		if i == cols {
			res = append(res, showRow(a.Path[offset:offset+i])...)
			res = append(res, '\n')
			i = 0
			offset += cols
		} else if offset+i == len(a.Path) {
			break
		}
		i++
	}
	res = append(res, showRow(a.Path[offset:offset+i])...)
	res = append(res, '\n')
	return string(res)
}

func formatSeq(seq []rune, cols int) []rune {
	l := len(seq)
	var res []rune
	for i := 0; i < l; i += cols {
		end := i + cols
		if end > l {
			end = l
		}
		res = append(res, seq[i:end]...)
		res = append(res, '\n')
	}
	return res
}

func showRow(row []Match) []rune {
	res := []rune{}
	l := len(row)
	line1 := make([]rune, l)
	line2 := make([]rune, l)
	line3 := make([]rune, l)

	for i, v := range row {
		l1 := v.L1
		if v.Type == insertion {
			l1 = '-'
		}
		l2 := v.L2
		if v.Type == deletion {
			l2 = '-'
		}

		line1[i] = l1
		line2[i] = alignmentRune(v)
		line3[i] = l2
	}
	res = append(res, line1...)
	res = append(res, '\n')
	res = append(res, line2...)
	res = append(res, '\n')
	res = append(res, line3...)
	res = append(res, '\n')
	return res
}

func alignmentRune(m Match) rune {
	r := ' '
	if m.Subst == identical {
		r = '|'
	} else if m.Subst == similar {
		r = ':'
	}
	return r
}

func (a *Alignment) IdentitySimilarity() (float32, float32) {
	var length int
	if length = a.Gene1.SeqLen; length < a.Gene2.SeqLen {
		length = a.Gene2.SeqLen
	}
	identity := 100 * float32(a.Identical) / float32(length)
	similarity := 100 * float32(a.Identical+a.Similar) / float32(length)
	return identity, similarity
}

func (a *Alignment) calculateScoreMatrix(conf Env,
	b62 Blosum62) (ScoreMatrix, MaxScore) {
	var max MaxScore
	var gain, score int
	matrix := make(ScoreMatrix, a.Gene1.SeqLen+1)
	for i := range matrix {
		matrix[i] = make([]int, a.Gene2.SeqLen+1)
	}

	for j := 1; j <= a.Gene2.SeqLen; j++ {
		for i := 1; i <= a.Gene1.SeqLen; i++ {
			if gain = b62[a.Gene1.Seq[i-1]][a.Gene2.Seq[j-1]]; gain > 0 {
				score = matrix[i-1][j-1] + gain
			} else {
				cells := []int{matrix[i-1][j-1] + gain, del(matrix, i, j, conf),
					ins(matrix, i, j, conf)}
				score = maximum(cells)
			}
			if score < 0 {
				score = 0
			}
			matrix[i][j] = score
			if score > max.Score {
				max = MaxScore{score, i, j}
			}
		}
	}
	return matrix, max
}

func (a *Alignment) calculatePath(matrix ScoreMatrix, max MaxScore) {
	var path []Match
	var match Match
	i := max.I
	j := max.J
	score := matrix[i][j]
	for {
		if score > 0 {
			match, i, j = previous(matrix, i, j, a)
			path = append(path, match)
			score = matrix[i][j]
		} else {
			break
		}
	}
	a.Path = Reverse(path)
	// a.Path = path
}

func Reverse(path []Match) []Match {
	last := len(path) - 1
	for i := 0; i < len(path)/2; i++ {
		path[i], path[last-i] = path[last-i], path[i]
	}
	return path
}

func previous(matrix ScoreMatrix, i int, j int,
	a *Alignment) (Match, int, int) {
	match := Match{L1: a.Gene1.Seq[i-1], L2: a.Gene2.Seq[j-1], I: i, J: j,
		Score: matrix[i][j]}
	candidates := [][]int{{matrix[i-1][j-1], substitution},
		{matrix[i][j-1], insertion},
		{matrix[i-1][j], deletion}}
	prev := candidates[0]
	for _, v := range candidates[1:len(candidates)] {
		if prev[0] < v[0] {
			prev = v
		}
	}

	match.Type = prev[1]
	if match.Type == substitution {
		if a.Gene1.Seq[i-1] == a.Gene2.Seq[j-1] {
			match.Subst = identical
			a.Identical += 1
		} else if matrix[i][j]-prev[0] > 0 {
			match.Subst = similar
			a.Similar += 1
		} else {
			match.Subst = different
		}
	}
	switch match.Type {
	case deletion:
		i = i - 1
	case insertion:
		j = j - 1
	default:
		i = i - 1
		j = j - 1
	}

	return match, i, j
}

func del(matrix ScoreMatrix, i int, j int, conf Env) int {
	return matrix[i-1][j] - conf.GapExtends
}

func ins(matrix ScoreMatrix, i int, j int, conf Env) int {
	return matrix[i][j-1] - conf.GapExtends
}

func maximum(ary []int) int {
	max := ary[0]
	for _, v := range ary[1 : len(ary)-1] {
		if max < v {
			max = v
		}
	}
	return max
}
