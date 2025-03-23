package sorter

import "math/rand/v2"

type Option func(*Sorter)

func New(options ...Option) *Sorter {
	s := &Sorter{}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func WithClassNames(names map[string]int) Option {
	return func(s *Sorter) {
		s.ClassNames = names
	}
}

func WithClasses(cs []*Class) Option {
	return func(s *Sorter) {
		s.Classes = cs
	}
}

func WithWithRules(rs map[string]string) Option {
	return func(s *Sorter) {
		s.WithRules = rs
	}
}

func WithAvoidRules(rs map[string]string) Option {
	return func(s *Sorter) {
		s.AvoidRules = rs
	}
}

func WithSiblingRules(rs map[string]string) Option {
	return func(s *Sorter) {
		s.SiblingRules = rs
	}
}

func WithPupils(ps []*Pupil) Option {
	return func(s *Sorter) {
		s.Pupils = ps
	}
}

func WithRando(r *rand.Rand) Option {
	return func(s *Sorter) {
		s.Rando = r
	}
}
