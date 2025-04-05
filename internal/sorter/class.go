package sorter

import (
	"fmt"
	"math"
)

type Pupil struct {
	Id        string
	Gender    string
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

// determine whether pupil has withrules going back to someone placed due to siblingRules
func (c *Class) hasSiblingRuleAncestor(siblingRules map[string]string, withRules map[string]string, pupilId string) bool {
	_, isInClass := c.Pupils[pupilId]
	if !isInClass {
		//log.Printf("Pupil %s is not in class %s\n", pupilId, c.name)
		return false
	}

	_, isSibling := siblingRules[pupilId]
	if isSibling {
		//log.Printf("Pupil %s has a sibling\n", pupilId)
		return true
	}

	for parent, child := range withRules {
		if pupilId == child {
			//log.Printf("Checking whether parent %s of %s has sibling\n", parent, child)
			if c.hasSiblingRuleAncestor(siblingRules, withRules, parent) {
				return true
			}
		}
	}

	//log.Printf("Pupil %s does not have a sibling ancestor\n", pupilId)
	return false
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

func (c *Class) MetUnmetWithRules(rules map[string]string, siblingRules map[string]string) (int, int, []string) {
	met := 0
	unmet := 0
	unmetReasons := make([]string, 0)
	for p1, p2 := range rules {
		_, p1InClass := c.Pupils[p1]
		_, p2InClass := c.Pupils[p2]
		// met if both pupils are in class
		if p1InClass && p2InClass {
			met++
			continue
		}
		// TODO: unmet should exclude impossible rules (fixed in separate classes)
		// TODO: return count of impossible rules
		// unmet, if p1 is in class and p2 is not
		if p1InClass && !p2InClass {
			unmet++
			p1forced := ""
			p2forced := ""
			_, p1f := siblingRules[p1]
			if p1f {
				p1forced = fmt.Sprintf(" (forced)")
			}
			p2c, p2f := siblingRules[p2]
			if p2f {
				p2forced = fmt.Sprintf("(forced to %s)", p2c)
			}
			reason := fmt.Sprintf("pupil %s is in class%s, but %s is not %s", p1, p1forced, p2, p2forced)
			unmetReasons = append(unmetReasons, reason)
		}
	}
	return met, unmet, unmetReasons
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

func (c *Class) MetUnmetAvoidRules(rules map[string]string) (int, int) {
	met := 0
	unmet := 0
	for p1, p2 := range rules {
		_, p1InClass := c.Pupils[p1]
		_, p2InClass := c.Pupils[p2]
		// met if p1 is in  class and p2 is not, 0 otherwise
		if p1InClass && !p2InClass {
			met++
			continue
		}
		// unmet if both are in class
		if p1InClass && p2InClass {
			unmet++
		}
	}
	return met, unmet
}

func (c *Class) ScoreClassSize() int {
	delta := len(c.Pupils) - c.targetSize
	absDelta := math.Abs(float64(delta))
	finalScore := classSizeScore - int(absDelta) // remove a point for every pupil we are apart from the target of 30
	return finalScore
}

func (c *Class) ScoreGenderBalance() (int, float64, float64) {
	males := 0
	females := 0
	for _, p := range c.Pupils {
		if p.Gender == "male" {
			males++
		} else if p.Gender == "female" {
			females++
		}
	}

	total := males + females
	if total > 0 {
		malePct := float64(males) / float64(total)
		femalePct := float64(females) / float64(total)
		diff := malePct - femalePct
		score := genderBalanceScore - 2*(math.Abs(diff)*genderBalanceScore)
		if score < 0 {
			score = 0
		}
		return int(score), malePct, femalePct
	}
	return 0, 0, 0
}
