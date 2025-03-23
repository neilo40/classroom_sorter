package runner

import (
	"log"
	"time"

	"math/rand/v2"

	"github.com/neilo40/classroom_sorter/internal/importer"
	"github.com/neilo40/classroom_sorter/internal/sorter"
)

func run(ch chan<- *sorter.Sorter, i *importer.Importer) {
	log.Println("Runner started")
	topScore := 0
	for {
		// TODO: is creating this struct each loop too slow?
		// TODO: pre-generate the base struct and change sorter to not mutate initial conditions?
		s := sorter.New(
			sorter.WithClassNames(map[string]int{"1a1": 30, "1a2": 30, "1d1": 30, "1d2": 30, "1s1": 30, "1s2": 30, "1s3": 20}),
			sorter.WithAvoidRules(i.GetAvoidRules()),
			sorter.WithWithRules(i.GetWithRules()),
			sorter.WithSiblingRules(i.GetSiblingRules()),
			sorter.WithPupils(i.GetPupils()),
			sorter.WithRando(rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))), // TODO: why 2 seeds?
		)
		err := s.InitClasses()
		if err != nil {
			log.Println(err)
			continue
		}
		score := s.ScoreAllRules()
		if score > topScore {
			ch <- s
			topScore = score
		} else {
			ch <- nil // receiver will count how many attempts were made
		}
	}
}

// TODO: handle interruption / shutdown
func receiver(ch <-chan *sorter.Sorter) {
	log.Println("Receiver started")
	var s *sorter.Sorter
	totalAttempts := 0
	bestScore := 0
	var bestSort *sorter.Sorter

	for {
		s = <-ch
		totalAttempts += 1
		if s != nil {
			if s.Score > bestScore {
				bestScore = s.Score
				bestSort = s
				log.Printf("New best score: %d, attempts: %d\n", bestSort.Score, totalAttempts)
			}
		}
		// TODO: output this after X time, not attempts
		if (totalAttempts % 10000) == 0 {
			log.Printf("Attempts: %d\n", totalAttempts)
		}
	}
}

func Zookeeper(runCount int, i *importer.Importer) {
	ch := make(chan *sorter.Sorter, runCount)
	go receiver(ch)
	for range runCount {
		go run(ch, i)
	}
	// TODO: listen for interrupt and gracefully quit, printing best sort
	for {
	}
}
