package rendering

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

const (
	EventTypeHorizontalSplit EventType = iota
	EventTypeVerticalSplit
	EventTypeHorizontalEnd
	EventTypeVerticalSplitEnd
	EventTypeVerticalDocumentEnd
	EventTypeVerticalVerticalSplit

	StateTypeRegular StateType = iota
	StateTypeNested
)

var (
	EventMapping = map[EventType][]byte{
		EventTypeHorizontalSplit: []byte(`
</section>
<section>`),
		EventTypeVerticalSplit: []byte(`
</section>
<section>
<section>
`),
		EventTypeHorizontalEnd: []byte(`
</section>
`),
		EventTypeVerticalSplitEnd: []byte(`
</section>
</section>
<section>`),
		EventTypeVerticalDocumentEnd: []byte(`
</section>
</section>
`),
		EventTypeVerticalVerticalSplit: []byte(`
</section>
</section>
<section>
<section>`),
	}
)

func NewStateMachine(verticalSplit, horizontalSplit string) *StateMachine {
	return &StateMachine{
		CurrentState: StateTypeRegular,
		States: map[StateType]State{
			StateTypeRegular: {
				Transitions: map[string]TransitionResult{
					horizontalSplit: {EventTypeHorizontalSplit, StateTypeRegular},
					fmt.Sprintf("%s%s", horizontalSplit, horizontalSplit): {EventTypeHorizontalSplit, StateTypeRegular},
					fmt.Sprintf("%s%s", horizontalSplit, verticalSplit):   {EventTypeVerticalSplit, StateTypeNested},
					fmt.Sprintf("%s%s", verticalSplit, horizontalSplit):   {EventTypeVerticalSplit, StateTypeNested},
					"": {EventTypeHorizontalEnd, StateTypeRegular},
				},
			},
			StateTypeNested: {
				Transitions: map[string]TransitionResult{
					fmt.Sprintf("%s%s", verticalSplit, verticalSplit): {EventTypeHorizontalSplit, StateTypeNested},
					verticalSplit: {EventTypeHorizontalSplit, StateTypeNested},
					fmt.Sprintf("%s%s", verticalSplit, horizontalSplit): {EventTypeVerticalSplitEnd, StateTypeNested},
					horizontalSplit: {EventTypeVerticalSplitEnd, StateTypeRegular},
					fmt.Sprintf("%s%s", horizontalSplit, verticalSplit):   {EventTypeVerticalVerticalSplit, StateTypeNested},
					fmt.Sprintf("%s%s", horizontalSplit, horizontalSplit): {EventTypeVerticalSplitEnd, StateTypeRegular},
					"": {EventTypeVerticalDocumentEnd, StateTypeRegular},
				},
			},
		},
	}
}

type (
	EventType        uint
	StateType        uint
	TransitionResult struct {
		EventType EventType
		StateType StateType
	}
	State struct {
		Transitions map[string]TransitionResult
	}
	StateMachine struct {
		CurrentState StateType
		States       map[StateType]State
	}
)

func (m *StateMachine) Accept(input string) []byte {
	if result, ok := m.States[m.CurrentState].Transitions[input]; ok {
		m.CurrentState = result.StateType
		return EventMapping[result.EventType]
	}
	return nil
}

func peekNextRuler(node ast.Node) *ast.HorizontalRule {
	if node.AsContainer() == nil {
		node = node.GetParent()
	}

	nodes := node.GetChildren()
	if nodes == nil {
		return nil
	}

	var selfIdx int
	for idx := range nodes {
		if nodes[idx] == node {
			selfIdx = idx
			break
		}
	}

	for idx := selfIdx + 1; idx < len(nodes); idx++ {
		if hr, ok := nodes[idx].(*ast.HorizontalRule); ok {
			return hr
		}
	}
	return nil
}
