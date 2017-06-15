package query

import (
	"container/list"
)

type Plan struct {
	DryRun bool

	ready *list.List
	want  map[*Edge]struct{}
}

func NewPlan() *Plan {
	return &Plan{
		ready: list.New(),
		want:  make(map[*Edge]struct{}),
	}
}

func (p *Plan) AddTarget(e *Edge) {
	if _, ok := p.want[e]; ok {
		return
	}

	p.want[e] = struct{}{}
	if inputs := e.Input.Inputs(); len(inputs) == 0 {
		p.ready.PushBack(e.Input)
		return
	} else {
		for _, input := range inputs {
			p.AddTarget(input)
		}
	}
}

func (p *Plan) FindWork() Node {
	front := p.ready.Front()
	if front != nil {
		return p.ready.Remove(front).(Node)
	}
	return nil
}

func (p *Plan) ScheduleWork(nodes ...Node) {
	for _, n := range nodes {
		if AllInputsReady(n) {
			p.ready.PushBack(n)
			continue
		}

		// Add each input edge as a target.
		for _, input := range n.Inputs() {
			p.AddTarget(input)
		}
	}
}

// EdgeFinished runs when notified that an Edge has finished running so the
// Edge's output Nodes can be checked to see if their output edges are now
// ready to be run.
func (p *Plan) NodeFinished(n Node) {
	for _, e := range n.Outputs() {
		// The nodes are now considered ready. Check if their output edge is
		// now ready to be executed (if they have one).
		if e.Output != nil && AllInputsReady(e.Output) {
			p.ready.PushBack(e.Output)
		}
	}
}
