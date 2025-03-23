package sorter_test

import (
	"log"
	"testing"

	"github.com/neilo40/classroom_sorter/internal/sorter"
)

func TestInit(t *testing.T) {
	s := sorter.New(
		sorter.WithClassNames(map[string]int{"1a1": 30, "1a2": 30, "1a3": 30, "1b1": 30, "1b2": 30, "1c1": 30, "1c2": 30}),
		sorter.WithAvoidRules(map[string]string{"1": "3", "2": "4"}),
		sorter.WithWithRules(map[string]string{"1": "2", "3": "22"}),
		sorter.WithSiblingRules(map[string]string{"8": "1b1", "43": "1c2"}),
		sorter.WithPupils([]*sorter.Pupil{
			{Id: "1", HasGirfme: false, Gender: sorter.Male},
			{Id: "2", HasGirfme: false, Gender: sorter.Male},
			{Id: "3", HasGirfme: true, Gender: sorter.Female},
			{Id: "4", HasGirfme: false, Gender: sorter.Male},
			{Id: "8", HasGirfme: false, Gender: sorter.Female},
			{Id: "22", HasGirfme: true, Gender: sorter.Male},
			{Id: "43", HasGirfme: false, Gender: sorter.Other},
		}),
	)

	err := s.InitClasses()
	if err != nil {
		t.Errorf("Error generating classes: %v", err)
	}

	score := s.ScoreAllRules()
	log.Println(score)
}
