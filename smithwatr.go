package smithwatr

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Blosum62 map[rune]map[rune]int

// Env is a collection of environment variables.
type Env struct {
	DbHost     string
	DbUser     string
	Db         string
	DataDir    string
	GapOpens   int
	GapExtends int
	WorkersNum int
}

// Check handles error checking, and panicks if error is not nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// EnvVars imports all environment variables relevant for the data conversion.
func EnvVars() Env {
	emptyEnvs := make([]string, 0, 4)
	envVars := [7]string{"POSTGRES_HOST", "POSTGRES_USER", "POSTGRES_DB",
		"DATA_DIR", "GAP_OPEN_PENTALTY", "GAP_EXTENSION_PENALTY", "CPU_CAPACITY"}
	for i, v := range envVars {
		val, ok := os.LookupEnv(v)
		if ok {
			envVars[i] = val
		} else {
			emptyEnvs = append(emptyEnvs, v)
		}
	}
	if len(emptyEnvs) > 0 {
		envs := strings.Join(emptyEnvs, ", ")
		panic(fmt.Errorf("Environment variables %s are not defined", envs))
	}

	gopen, err := strconv.Atoi(envVars[4])
	Check(err)
	gext, err := strconv.Atoi(envVars[5])
	Check(err)

	return Env{DbHost: envVars[0], DbUser: envVars[1], Db: envVars[2],
		DataDir: envVars[3], GapOpens: gopen, GapExtends: gext,
		WorkersNum: calculateWorkersNum(envVars[6])}
}

func calculateWorkersNum(cpuLoad string) int {
	load, err := strconv.ParseFloat(cpuLoad, 64)
	Check(err)
	cpuNum := runtime.NumCPU()
	workersNum := int(math.Ceil(float64(cpuNum) * load))
	return workersNum
}

// InitBlosum62 creates a map with BLOSSUM62 weights values
func InitBlosum62() Blosum62 {
	b62 := make(Blosum62)

	b62['A'] = map[rune]int{'A': 4, 'R': -1, 'N': -2, 'D': -2, 'C': 0,
		'Q': -1, 'E': -1, 'G': 0, 'H': -2, 'I': -1, 'L': -1, 'K': -1, 'M': -1,
		'F': -2, 'P': -1, 'S': 1, 'T': 0, 'W': -3, 'Y': -2, 'V': 0, 'B': -2,
		'Z': -1, 'X': 0, '*': -4}
	b62['R'] = map[rune]int{'A': -1, 'R': 5, 'N': 0, 'D': -2, 'C': -3,
		'Q': 1, 'E': 0, 'G': -2, 'H': 0, 'I': -3, 'L': -2, 'K': 2, 'M': -1,
		'F': -3, 'P': -2, 'S': -1, 'T': -1, 'W': -3, 'Y': -2, 'V': -3, 'B': -1,
		'Z': 0, 'X': -1, '*': -4}
	b62['N'] = map[rune]int{'A': -2, 'R': 0, 'N': 6, 'D': 1, 'C': -3,
		'Q': 0, 'E': 0, 'G': 0, 'H': 1, 'I': -3, 'L': -3, 'K': 0, 'M': -2,
		'F': -3, 'P': -2, 'S': 1, 'T': 0, 'W': -4, 'Y': -2, 'V': -3, 'B': 3,
		'Z': 0, 'X': -1, '*': -4}
	b62['D'] = map[rune]int{'A': -2, 'R': -2, 'N': 1, 'D': 6, 'C': -3,
		'Q': 0, 'E': 2, 'G': -1, 'H': -1, 'I': -3, 'L': -4, 'K': -1, 'M': -3,
		'F': -3, 'P': -1, 'S': 0, 'T': -1, 'W': -4, 'Y': -3, 'V': -3, 'B': 4,
		'Z': 1, 'X': -1, '*': -4}
	b62['C'] = map[rune]int{'A': 0, 'R': -3, 'N': -3, 'D': -3, 'C': 9,
		'Q': -3, 'E': -4, 'G': -3, 'H': -3, 'I': -1, 'L': -1, 'K': -3, 'M': -1,
		'F': -2, 'P': -3, 'S': -1, 'T': -1, 'W': -2, 'Y': -2, 'V': -1, 'B': -3,
		'Z': -3, 'X': -2, '*': -4}
	b62['Q'] = map[rune]int{'A': -1, 'R': 1, 'N': 0, 'D': 0, 'C': -3,
		'Q': 5, 'E': 2, 'G': -2, 'H': 0, 'I': -3, 'L': -2, 'K': 1, 'M': 0,
		'F': -3, 'P': -1, 'S': 0, 'T': -1, 'W': -2, 'Y': -1, 'V': -2, 'B': 0,
		'Z': 3, 'X': -1, '*': -4}
	b62['E'] = map[rune]int{'A': -1, 'R': 0, 'N': 0, 'D': 2, 'C': -4,
		'Q': 2, 'E': 5, 'G': -2, 'H': 0, 'I': -3, 'L': -3, 'K': 1, 'M': -2,
		'F': -3, 'P': -1, 'S': 0, 'T': -1, 'W': -3, 'Y': -2, 'V': -2, 'B': 1,
		'Z': 4, 'X': -1, '*': -4}
	b62['G'] = map[rune]int{'A': 0, 'R': -2, 'N': 0, 'D': -1, 'C': -3,
		'Q': -2, 'E': -2, 'G': 6, 'H': -2, 'I': -4, 'L': -4, 'K': -2, 'M': -3,
		'F': -3, 'P': -2, 'S': 0, 'T': -2, 'W': -2, 'Y': -3, 'V': -3, 'B': -1,
		'Z': -2, 'X': -1, '*': -4}
	b62['H'] = map[rune]int{'A': -2, 'R': 0, 'N': 1, 'D': -1, 'C': -3,
		'Q': 0, 'E': 0, 'G': -2, 'H': 8, 'I': -3, 'L': -3, 'K': -1, 'M': -2,
		'F': -1, 'P': -2, 'S': -1, 'T': -2, 'W': -2, 'Y': 2, 'V': -3, 'B': 0,
		'Z': 0, 'X': -1, '*': -4}
	b62['I'] = map[rune]int{'A': -1, 'R': -3, 'N': -3, 'D': -3, 'C': -1,
		'Q': -3, 'E': -3, 'G': -4, 'H': -3, 'I': 4, 'L': 2, 'K': -3, 'M': 1,
		'F': 0, 'P': -3, 'S': -2, 'T': -1, 'W': -3, 'Y': -1, 'V': 3, 'B': -3,
		'Z': -3, 'X': -1, '*': -4}
	b62['L'] = map[rune]int{'A': -1, 'R': -2, 'N': -3, 'D': -4, 'C': -1,
		'Q': -2, 'E': -3, 'G': -4, 'H': -3, 'I': 2, 'L': 4, 'K': -2, 'M': 2,
		'F': 0, 'P': -3, 'S': -2, 'T': -1, 'W': -2, 'Y': -1, 'V': 1, 'B': -4,
		'Z': -3, 'X': -1, '*': -4}
	b62['K'] = map[rune]int{'A': -1, 'R': 2, 'N': 0, 'D': -1, 'C': -3,
		'Q': 1, 'E': 1, 'G': -2, 'H': -1, 'I': -3, 'L': -2, 'K': 5, 'M': -1,
		'F': -3, 'P': -1, 'S': 0, 'T': -1, 'W': -3, 'Y': -2, 'V': -2, 'B': 0,
		'Z': 1, 'X': -1, '*': -4}
	b62['M'] = map[rune]int{'A': -1, 'R': -1, 'N': -2, 'D': -3, 'C': -1,
		'Q': 0, 'E': -2, 'G': -3, 'H': -2, 'I': 1, 'L': 2, 'K': -1, 'M': 5,
		'F': 0, 'P': -2, 'S': -1, 'T': -1, 'W': -1, 'Y': -1, 'V': 1, 'B': -3,
		'Z': -1, 'X': -1, '*': -4}
	b62['F'] = map[rune]int{'A': -2, 'R': -3, 'N': -3, 'D': -3, 'C': -2,
		'Q': -3, 'E': -3, 'G': -3, 'H': -1, 'I': 0, 'L': 0, 'K': -3, 'M': 0,
		'F': 6, 'P': -4, 'S': -2, 'T': -2, 'W': 1, 'Y': 3, 'V': -1, 'B': -3,
		'Z': -3, 'X': -1, '*': -4}
	b62['P'] = map[rune]int{'A': -1, 'R': -2, 'N': -2, 'D': -1, 'C': -3,
		'Q': -1, 'E': -1, 'G': -2, 'H': -2, 'I': -3, 'L': -3, 'K': -1, 'M': -2,
		'F': -4, 'P': 7, 'S': -1, 'T': -1, 'W': -4, 'Y': -3, 'V': -2, 'B': -2,
		'Z': -1, 'X': -2, '*': -4}
	b62['S'] = map[rune]int{'A': 1, 'R': -1, 'N': 1, 'D': 0, 'C': -1,
		'Q': 0, 'E': 0, 'G': 0, 'H': -1, 'I': -2, 'L': -2, 'K': 0, 'M': -1,
		'F': -2, 'P': -1, 'S': 4, 'T': 1, 'W': -3, 'Y': -2, 'V': -2, 'B': 0,
		'Z': 0, 'X': 0, '*': -4}
	b62['T'] = map[rune]int{'A': 0, 'R': -1, 'N': 0, 'D': -1, 'C': -1,
		'Q': -1, 'E': -1, 'G': -2, 'H': -2, 'I': -1, 'L': -1, 'K': -1, 'M': -1,
		'F': -2, 'P': -1, 'S': 1, 'T': 5, 'W': -2, 'Y': -2, 'V': 0, 'B': -1,
		'Z': -1, 'X': 0, '*': -4}
	b62['W'] = map[rune]int{'A': -3, 'R': -3, 'N': -4, 'D': -4, 'C': -2,
		'Q': -2, 'E': -3, 'G': -2, 'H': -2, 'I': -3, 'L': -2, 'K': -3, 'M': -1,
		'F': 1, 'P': -4, 'S': -3, 'T': -2, 'W': 11, 'Y': 2, 'V': -3, 'B': -4,
		'Z': -3, 'X': -2, '*': -4}
	b62['Y'] = map[rune]int{'A': -2, 'R': -2, 'N': -2, 'D': -3, 'C': -2,
		'Q': -1, 'E': -2, 'G': -3, 'H': 2, 'I': -1, 'L': -1, 'K': -2, 'M': -1,
		'F': 3, 'P': -3, 'S': -2, 'T': -2, 'W': 2, 'Y': 7, 'V': -1, 'B': -3,
		'Z': -2, 'X': -1, '*': -4}
	b62['V'] = map[rune]int{'A': 0, 'R': -3, 'N': -3, 'D': -3, 'C': -1,
		'Q': -2, 'E': -2, 'G': -3, 'H': -3, 'I': 3, 'L': 1, 'K': -2, 'M': 1,
		'F': -1, 'P': -2, 'S': -2, 'T': 0, 'W': -3, 'Y': -1, 'V': 4, 'B': -3,
		'Z': -2, 'X': -1, '*': -4}
	b62['B'] = map[rune]int{'A': -2, 'R': -1, 'N': 3, 'D': 4, 'C': -3,
		'Q': 0, 'E': 1, 'G': -1, 'H': 0, 'I': -3, 'L': -4, 'K': 0, 'M': -3,
		'F': -3, 'P': -2, 'S': 0, 'T': -1, 'W': -4, 'Y': -3, 'V': -3, 'B': 4,
		'Z': 1, 'X': -1, '*': -4}
	b62['Z'] = map[rune]int{'A': -1, 'R': 0, 'N': 0, 'D': 1, 'C': -3,
		'Q': 3, 'E': 4, 'G': -2, 'H': 0, 'I': -3, 'L': -3, 'K': 1, 'M': -1,
		'F': -3, 'P': -1, 'S': 0, 'T': -1, 'W': -3, 'Y': -2, 'V': -2, 'B': 1,
		'Z': 4, 'X': -1, '*': -4}
	b62['X'] = map[rune]int{'A': 0, 'R': -1, 'N': -1, 'D': -1, 'C': -2,
		'Q': -1, 'E': -1, 'G': -1, 'H': -1, 'I': -1, 'L': -1, 'K': -1, 'M': -1,
		'F': -1, 'P': -2, 'S': 0, 'T': 0, 'W': -2, 'Y': -1, 'V': -1, 'B': -1,
		'Z': -1, 'X': -1, '*': -4}
	b62['*'] = map[rune]int{'A': -4, 'R': -4, 'N': -4, 'D': -4, 'C': -4,
		'Q': -4, 'E': -4, 'G': -4, 'H': -4, 'I': -4, 'L': -4, 'K': -4, 'M': -4,
		'F': -4, 'P': -4, 'S': -4, 'T': -4, 'W': -4, 'Y': -4, 'V': -4, 'B': -4,
		'Z': -4, 'X': -4, '*': 1}
	return b62
}

func Connect(conf Env) (*sql.DB, error) {
	var db *sql.DB
	var err error
	params := fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable",
		conf.DbUser, conf.DbHost, conf.Db)
	db, err = sql.Open("postgres", params)
	return db, err
}
