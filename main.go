package main

import (
	"flag"
	"math/rand"
	"os"
	"time"
)

var (
	buildVersion string
	buildType    string

	printVersion bool
	seed int64
	popSize int
)

func init() {
	flag.BoolVar(&printVersion, "v", false, "print version information and exit")
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "allows specifying a custom seed for all random number generation operations")
	flag.IntVar(&popSize, "pops", 64, "allows specifying the population size to perform the genetic search upon")
}

func main() {
	flag.Parse()

	EnableLogging(os.Stdout)

	l.Printf("Seeding random number generator with %v\n", seed)
	rand.Seed(seed)

	if printVersion {
		l.Printf("%s-%s", buildType, buildVersion)
		return
	}

	population := newPopulation(popSize)
	if candidate := population.genticSearch(); candidate != nil {
		l.Printf("Solution found:")
		candidate.print(os.Stdout)
	}
}
