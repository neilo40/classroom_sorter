package sorter

import (
	"math"
	"math/rand/v2"
)

const (
	// score per pupil / rule
	withScore    = 1 // We could score up to 90 times 30 pupils, max 3 "with" buddies
	avoidScore   = 5 // could also be up to 90 but we prefer satisfying this over the "with"s
	siblingScore = 2
	// score per class
	classSizeScore     = 10 // -2 for each pupil the class is under or over 30
	genderBalanceScore = 10 // score if male : female ratio is within 10%
	// score overall
	girfmeScore      = 100 // -10 for each % away from optimum balance:withScore
	maxClassOversize = 2   // can go to max this amount over the average class size
)

type Sorter struct {
	ClassNames   map[string]int // class name -> size
	Classes      []*Class
	WithRules    map[string]string // pupil -> pupil
	AvoidRules   map[string]string // pupil -> pupil
	SiblingRules map[string]string // pupil -> class name
	Score        int
	Pupils       []*Pupil
	Rando        *rand.Rand
}

func (s *Sorter) InitClasses() error {
	pupilCount := len(s.Pupils)
	numClasses := len(s.ClassNames)
	maxClassSize := (pupilCount / numClasses) + maxClassOversize

	classes := s.Classes
	if classes == nil {
		classes = make([]*Class, 0, numClasses)
		for name, size := range s.ClassNames {
			classes = append(classes, &Class{name: name, Pupils: make(map[string]*Pupil), targetSize: size})
		}
		s.Classes = classes
	}

	err := s.placeSiblings()
	if err != nil {
		return err
	}

	for i := range pupilCount {
		pupil := s.Pupils[i]
		// skip if pupil has sibling - already placed
		_, ok := s.SiblingRules[pupil.Id]
		if ok {
			continue
		}

		// skip a class if it is full (has reached pupils/classes + maxClassOversize)
		for {
			n := s.Rando.Int32N(int32(numClasses))
			if len(classes[n].Pupils) < maxClassSize { // TODO: test if pre-computing the class sizes speeds things up
				classes[n].addPupil(pupil)
				break
			}
		}
	}

	return nil
}

func (s *Sorter) placeSiblings() error {
	// TODO: move this map out if needed elsewhere
	classMap := make(map[string]*Class)
	for _, c := range s.Classes {
		classMap[c.name] = c
	}
	for _, p := range s.Pupils {
		cName, ok := s.SiblingRules[p.Id]
		if ok {
			classMap[cName].addPupil(p)
		}
	}
	return nil
}

func (s *Sorter) ScoreAllRules() int {
	total := 0
	total += s.ScoreWithRules()
	total += s.ScoreAvoidRules()
	total += s.ScoreClassSizes()
	total += s.ScoreGenderBalance()
	total += s.ScoreGirfmeBalance()
	s.Score = total
	return total
}

func (s *Sorter) ScoreWithRules() int {
	total := 0
	for _, c := range s.Classes {
		total += c.ScoreWithRules(s.WithRules)
	}
	return total
}

func (s *Sorter) ScoreAvoidRules() int {
	total := 0
	for _, c := range s.Classes {
		total += c.ScoreAvoidRules(s.AvoidRules)
	}
	return total
}

func (s *Sorter) ScoreClassSizes() int {
	total := 0
	for _, c := range s.Classes {
		total += c.ScoreClassSize()
	}
	return total

}

func (s *Sorter) ScoreGenderBalance() int {
	total := 0
	for _, c := range s.Classes {
		total += c.ScoreGenderBalance()
	}
	return total

}

func (s *Sorter) ScoreGirfmeBalance() int {
	girfmes := make([]int, 0, len(s.Classes))
	totalGirfmes := 0
	for _, c := range s.Classes {
		gc := c.GirfmeCount()
		girfmes = append(girfmes, gc)
		totalGirfmes += gc
	}
	score := girfmeScore
	optimumCount := totalGirfmes / len(s.Classes)
	for _, gc := range girfmes {
		delta := optimumCount - gc
		score -= int(math.Abs(float64(delta)))
	}
	return score
}
