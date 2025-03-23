package main

import (
	"log"

	"github.com/neilo40/classroom_sorter/internal/importer"
	"github.com/neilo40/classroom_sorter/internal/runner"
)

func main() {
	log.Println("Starting")
	i := importer.Import("testinput.csv")
	runner.Zookeeper(2, i)
}
