package sorter

import (
	"fmt"
	"log"
)

func (s *Sorter) PrintSummary() {
	fmt.Println()
	log.Printf("Summary of best class placement found (Total score %d)...\n", s.Score)
	log.Printf("    GIRFME balance score %d\n", s.ScoreGirfmeBalance())
	log.Printf("    Class size score %d\n", s.ScoreClassSizes())
	fmt.Println()

	for _, c := range s.Classes {
		log.Printf("Class %s (%d pupils, %d GIRFMEs) \n", c.name, len(c.Pupils), c.GirfmeCount())
		gs, malepct, femalepct := c.ScoreGenderBalance()
		log.Printf("\tGender balance: %.1f%% Male : %.1f%% Female (score %d)\n", malepct*100, femalepct*100, gs)
		met, unmet, reasons := c.MetUnmetWithRules(s.WithRules, s.SiblingRules)
		log.Printf("\tWith Rules: %d met, %d unmet, (score %d)\n", met, unmet, c.ScoreWithRules(s.WithRules))
		for _, r := range reasons {
			log.Printf("\t\t%s\n", r)
		}
		met, unmet = c.MetUnmetAvoidRules(s.AvoidRules)
		log.Printf("\tAvoid Rules: %d met, %d unmet, (score %d)\n", met, unmet, c.ScoreAvoidRules(s.AvoidRules))
	}
}
