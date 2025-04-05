package sorter

import (
	"log"
	"math"
	"math/rand/v2"
	"os"

	"github.com/gocarina/gocsv"
)

const (
	// score per pupil / rule
	withScore    = 10 // We could score up to 90 times 30 pupils, max 3 "with" buddies
	avoidScore   = 0  //10 // could also be up to 90 but we prefer satisfying this over the "with"s
	siblingScore = 0  //2
	// score per class
	classSizeScore     = 0 //20  // -2 for each pupil the class is under or over 30
	genderBalanceScore = 0 //200 // subtract 10 for each 1% away from 50:50 balance
	// score overall
	girfmeScore      = 0 //100 // -10 for each % away from optimum balance:withScore
	maxClassOversize = 2 // can go to max this amount over the average class size
	progressFilename = "latest.csv"
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

func (s *Sorter) InitClasses(classList []Class) error {
	numClasses := len(s.ClassNames)
	var classes []*Class
	placed := make(map[string]bool)

	if classList != nil {
		// if classList was provided, just use it and skip the initialization
		classes = make([]*Class, 0, len(classList))
		for _, c := range classList {
			// need to do a deep copy to retain the original classlist for next iteration
			class := &Class{name: c.name, targetSize: c.targetSize, Pupils: make(map[string]*Pupil)}
			for k, v := range c.Pupils {
				class.Pupils[k] = &Pupil{Id: v.Id, HasGirfme: v.HasGirfme, Gender: v.Gender}
				placed[v.Id] = true
			}
			classes = append(classes, class)
		}
		s.Classes = classes
	} else {
		classes = make([]*Class, 0, numClasses)
		for name, size := range s.ClassNames {
			classes = append(classes, &Class{name: name, Pupils: make(map[string]*Pupil), targetSize: size})
		}

		s.Classes = classes
		err := s.placeSiblings()
		if err != nil {
			return err
		}
	}

	pupilCount := len(s.Pupils)
	maxClassSize := (pupilCount / numClasses) + maxClassOversize

	for i := range pupilCount {
		pupil := s.Pupils[i]
		// skip if pupil has sibling - already placed
		_, hasSibling := s.SiblingRules[pupil.Id]
		_, isPlaced := placed[pupil.Id]
		if hasSibling || isPlaced {
			continue
		}

		// skip a class if it is full (has reached pupils/classes + maxClassOversize)
		for {
			n := s.Rando.Int32N(int32(numClasses))
			if len(classes[n].Pupils) < maxClassSize {
				classes[n].addPupil(pupil)
				break
			}
		}
	}

	return nil
}

type ProgressRow struct {
	PupilId     string `csv:"pid"`
	PupilGender string `csv:"gender"`
	PupilGirfme bool   `csv:"girfme"`
	ClassName   string `csv:"cname"`
	ClassTarget int    `csv:"targetsize"`
}

func (s *Sorter) LoadClasses() ([]Class, error) {
	// populate the classes from disk, if file exists
	// only load the pupils that can be tied back to a force placement (due to SiblingRules)
	progressFile, err := os.OpenFile(progressFilename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, nil // if file can't be opened / doesn't exist, just skip
	}

	rows := []*ProgressRow{}
	err = gocsv.UnmarshalFile(progressFile, &rows)
	if err != nil {
		return nil, err
	}

	classMap := make(map[string]Class)
	for _, r := range rows {
		_, ok := classMap[r.ClassName]
		if !ok {
			classMap[r.ClassName] = Class{name: r.ClassName, Pupils: make(map[string]*Pupil), targetSize: r.ClassTarget}
		}
		classMap[r.ClassName].Pupils[r.PupilId] = &Pupil{Id: r.PupilId, Gender: r.PupilGender, HasGirfme: r.PupilGirfme}
	}

	classes := make([]Class, 0, len(classMap))
	for _, c := range classMap {
		for _, p := range c.Pupils {
			if !c.hasSiblingRuleAncestor(s.SiblingRules, s.WithRules, p.Id) {
				log.Printf("Removing %s from preload - no sibling ancestor", p.Id)
				delete(c.Pupils, p.Id)
			}
		}
		classes = append(classes, c)
	}

	return classes, nil
}

func (s *Sorter) SaveClasses() error {
	// save the classes to disk
	// csv, row per pupil, columns: pupil id, hasgirfme, gender, class name, class targetsize
	rows := make([]ProgressRow, 0)
	for _, c := range s.Classes {
		for _, p := range c.Pupils {
			rows = append(rows, ProgressRow{PupilId: p.Id, PupilGender: p.Gender, PupilGirfme: p.HasGirfme, ClassName: c.name, ClassTarget: c.targetSize})
		}
	}
	progressFile, err := os.OpenFile(progressFilename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	err = gocsv.MarshalFile(&rows, progressFile)
	return err
}

func (s *Sorter) placeSiblings() error {
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
	if total < 0 {
		total = 0
	}
	return total

}

func (s *Sorter) ScoreGenderBalance() int {
	total := 0
	for _, c := range s.Classes {
		t, _, _ := c.ScoreGenderBalance()
		total += t
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
	if score < 0 {
		score = 0
	}
	return score
}
