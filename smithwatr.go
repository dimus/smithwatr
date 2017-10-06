package smithwatr

import (
	"fmt"
	"os"
	"strings"
)

// Env is a collection of environment variables.
type Env struct {
	DbHost  string
	DbUser  string
	Db      string
	DataDir string
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
	envVars := [5]string{"POSTGRES_HOST", "POSTGRES_USER", "POSTGRES_DB",
		"DATA_DIR"}
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
	return Env{DbHost: envVars[0], DbUser: envVars[1], Db: envVars[2],
		DataDir: envVars[3]}
}
