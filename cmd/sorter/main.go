package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/neilo40/classroom_sorter/internal/importer"
	"github.com/neilo40/classroom_sorter/internal/runner"
)

func main() {
	workersPtr := flag.Int("workers", 6, "number of workers")
	inputPtr := flag.String("input", "classlist_23_3_25.csv", "input filename")
	flag.Parse()

	// performance profiling
	// go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Println("Starting")
	i := importer.Import(*inputPtr)
	r := runner.Runner{}
	r.Zookeeper(*workersPtr, i)
}
