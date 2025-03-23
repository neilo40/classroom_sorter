package graph

import "fmt"

type Pupil struct {
    With []*Pupil
    Avoid []*Pupil
}

func Dosomething(s string) error {
	fmt.Printf("%s\n", s)
	return nil
}
