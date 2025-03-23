package sorter

import (
	"math"
)

type Gender int

const (
	Male Gender = iota
	Female
	Other
)

type Pupil struct {
	Id        string
	Gender    Gender
	HasGirfme bool
}

type Class struct {
	name       string
	Pupils     map[string]*Pupil // pupil id -> pupil
	targetSize int
}

func (c *Class) addPupil(p *Pupil) error {
	c.Pupils[p.Id] = p
	return nil
}

func (c *Class) GirfmeCount() int {
	total := 0
	for _, p := range c.Pupils {
		if p.HasGirfme {
			total++
		}
	}
	return total
}

func (c *Class) ScoreWithRules(rules map[string]string) int {
	total := 0
	for p1, p2 := range rules {
		_, p1InClass := c.Pupils[p1]
		_, p2InClass := c.Pupils[p2]
		// score if both pupils are in class
		if p1InClass && p2InClass {
			total += withScore
		}
	}
	return total
}

func (c *Class) ScoreAvoidRules(rules map[string]string) int {
	total := 0
	for p1, p2 := range rules {
		_, p1InClass := c.Pupils[p1]
		_, p2InClass := c.Pupils[p2]
		// score if p1 is in  class and p2 is not, 0 otherwise
		if p1InClass && !p2InClass {
			total += avoidScore
		}
	}
	return total
}

func (c *Class) ScoreClassSize() int {
	delta := len(c.Pupils) - c.targetSize
	absDelta := math.Abs(float64(delta))
	finalScore := classSizeScore - int(absDelta) // remove a point for every pupil we are apart from the target of 30
	return finalScore
}

func (c *Class) ScoreGenderBalance() int {
	males := 0
	females := 0
	for _, p := range c.Pupils {
		if p.Gender == Male {
			males++
		} else if p.Gender == Female {
			females++
		}
	}

	total := males + females
	if total > 0 {
		malePct := males / total
		femalePct := females / total
		diff := malePct - femalePct
		if math.Abs(float64(diff)) < 0.1 {
			return genderBalanceScore
		}
	}
	return 0
}
