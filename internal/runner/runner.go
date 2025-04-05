package runner

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"math/rand/v2"

	"github.com/neilo40/classroom_sorter/internal/importer"
	"github.com/neilo40/classroom_sorter/internal/sorter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Runner struct {
	bestSort sorter.Sorter
}

func (r *Runner) run(ch chan<- *sorter.Sorter, i *importer.Importer, ctx context.Context) {
	log.Println("Runner started")
	topScore := 0
	s := sorter.New(
		sorter.WithClassNames(map[string]int{"1A1": 30, "1A2": 30, "1D1": 30, "1D2": 30, "1S1": 30, "1S2": 30, "1S3": 20}),
		sorter.WithAvoidRules(i.GetAvoidRules()),
		sorter.WithWithRules(i.GetWithRules()),
		sorter.WithSiblingRules(i.GetSiblingRules()),
		sorter.WithPupils(i.GetPupils()),
		sorter.WithRando(rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))),
	)
	baseClassLists, err := s.LoadClasses()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := s.InitClasses(baseClassLists)
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
}

func (r *Runner) receiver(ch <-chan *sorter.Sorter) {
	log.Println("Receiver started")
	var s *sorter.Sorter
	totalAttempts := 0
	bestScore := 0
	start := time.Now()
	last := time.Now()
	lastBest := time.Now()
	pp := message.NewPrinter(language.BritishEnglish)

	for {
		s = <-ch
		totalAttempts += 1
		if s != nil {
			if s.Score > bestScore {
				bestScore = s.Score
				r.bestSort = sorter.Sorter{Score: s.Score, Pupils: s.Pupils, Classes: s.Classes, WithRules: s.WithRules, AvoidRules: s.AvoidRules, SiblingRules: s.SiblingRules}
				pp.Printf("New best score: %d; attempts: %d\n", r.bestSort.Score, totalAttempts)
				lastBest = time.Now()
			}
		}
		// TODO: output this after X time, not attempts
		if (totalAttempts % 60000) == 0 {
			pp.Printf("Attempts: %d; Elapsed: %s; Attempts/sec: %.0f (last best %d; %s ago)\n",
				totalAttempts, time.Since(start).Round(time.Second), 60000/time.Since(last).Seconds(), r.bestSort.Score, time.Since(lastBest).Round(time.Second))
			last = time.Now()
		}
	}
}

func (r *Runner) dumpAndPrintBestClass() {
	r.bestSort.PrintSummary()
	// write the class memberships to disk
	err := r.bestSort.SaveClasses()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Runner) Zookeeper(runCount int, i *importer.Importer) {
	ch := make(chan *sorter.Sorter, runCount)
	go r.receiver(ch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for range runCount {
		go r.run(ch, i, ctx)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for range sigChan {
		cancel()
		r.dumpAndPrintBestClass()
		break
	}
}
