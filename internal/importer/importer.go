package importer

import (
	"bytes"
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/neilo40/classroom_sorter/internal/sorter"
)

// Input will be CSV:
// Pupil id, Sibling class, gender, girfme, With1,With2,With3,Avoid1,Avoid2,Avoid3,...,AvoidN

type Importer struct {
	Rows []InputRow
}

type InputRow struct {
	Id           string
	SiblingClass string
	Gender       string
	Girfme       string
	With         []string
	Avoid        []string
}

func Import(filename string) *Importer {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	rdr := bytes.NewReader(f)
	r := csv.NewReader(rdr)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	inputRows := make([]InputRow, 0, len(rows))
	for _, r := range rows {
		with := make([]string, 0, 3)
		for _, w := range r[4:7] {
			if w != "" {
				with = append(with, w)
			}
		}
		avoid := make([]string, 0, 6)
		for _, a := range r[7:] {
			if a != "" {
				avoid = append(avoid, a)
			}
		}
		ir := InputRow{Id: r[0], SiblingClass: r[1], Gender: r[2], Girfme: r[3], With: with, Avoid: avoid}
		inputRows = append(inputRows, ir)
	}
	return &Importer{
		Rows: inputRows,
	}
}

func (i *Importer) GetWithRules() map[string]string {
	withRules := make(map[string]string)
	for _, r := range i.Rows {
		for _, wr := range r.With {
			withRules[r.Id] = wr
		}
	}
	return withRules
}

func (i *Importer) GetAvoidRules() map[string]string {
	avoidRules := make(map[string]string)
	for _, r := range i.Rows {
		for _, ar := range r.Avoid {
			avoidRules[r.Id] = ar
		}
	}
	return avoidRules
}

func (i *Importer) GetSiblingRules() map[string]string {
	siblingRules := make(map[string]string)
	for _, r := range i.Rows {
		if r.SiblingClass != "" {
			siblingRules[r.Id] = r.SiblingClass
		}
	}
	return siblingRules
}

func (i *Importer) GetPupils() []*sorter.Pupil {
	pupils := make([]*sorter.Pupil, 0, len(i.Rows))
	for _, r := range i.Rows {
		hasGirfme := false
		if strings.ToLower(r.Girfme) == "yes" {
			hasGirfme = true
		}

		p := sorter.Pupil{Id: r.Id, Gender: strings.TrimSpace(strings.ToLower(r.Gender)), HasGirfme: hasGirfme}
		pupils = append(pupils, &p)
	}
	return pupils
}
