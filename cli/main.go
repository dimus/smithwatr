package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	. "github.com/dimus/smithwatr"
)

var githash = "n/a"
var buildstamp = "n/a"

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	switch command {
	case "version":
		fmt.Printf(" Git commit hash: %s\n UTC Build Time: %s\n\n",
			githash, buildstamp)
	case "align":
		if len(os.Args) > 3 {
			var genome1, genome2 int
			log.Println("Importing data")
			b62 := InitBlosum62()
			conf := EnvVars()
			db, err := Connect(conf)
			Check(err)
			genome1, err = strconv.Atoi(os.Args[2])
			Check(err)
			genome2, err = strconv.Atoi(os.Args[3])
			Check(err)
			ImportData(db, conf)
			log.Println("Importing jobs")
			ImportJobs(db, genome1)
			log.Println("Aligning genomes")
			Align(db, genome2, -1, b62, conf)
		} else {
			fmt.Printf("Not enough arguments. Example:\n\n%s align 1 2", os.Args[0])
		}
	default:
		fmt.Printf("Usage:\n\n%s align 3 2\n\n", os.Args[0])
	}
}
