package cloudburst

import "sort"

// ScalingResult represents the calculated demand for instances
type ScalingResult struct {
	Result []*ResultValue // the rounded demand for new instances
}

type ResultValue struct {
	Provider       string
	Weight         float32
	InstanceDemand int
	Instances      []*Instance
}

type ByResultValues func(p1, p2 *ResultValue) bool

func (by ByResultValues) Sort(vars []*ResultValue) {
	ps := &resultSorter{
		vars: vars,
		by:   by,
	}
	sort.Sort(ps)
}

type resultSorter struct {
	vars []*ResultValue
	by   func(p1, p2 *ResultValue) bool
}

func (s *resultSorter) Len() int {
	return len(s.vars)
}

func (s *resultSorter) Swap(i, j int) {
	s.vars[i], s.vars[j] = s.vars[j], s.vars[i]
}

func (s *resultSorter) Less(i, j int) bool {
	return s.by(s.vars[i], s.vars[j])
}
