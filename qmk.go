package qmk

import (
	"github.com/leep-frog/command"
)

func CLI() *QMK {
	return &QMK{}
}

const (
	QMKEnvArg = "LEEP_QMK"
)

var (
	basicKeyboardBindings = []string{
		`bind '"\C-h":backward-delete-char'`,
	}
	qmkKeyboardBindings = []string{
		`bind '"\C-h":backward-kill-word'`,
	}
)

type QMK struct {
	Basic   bool
	changed bool
}

func (q *QMK) Setup() []string { return nil }
func (q *QMK) Name() string    { return "q" }
func (q *QMK) Changed() bool   { return q.changed }
func (q *QMK) Node() command.Node {
	return command.SerialNodes(
		command.Description("CLI for qmk/basic keyboard bindings in bash"),
		&command.BranchNode{
			Branches: map[string]command.Node{
				"toggle t": command.SerialNodes(
					command.Description("Toggles QMK mode"),
					&command.ExecutorProcessor{F: q.Toggle},
					command.ExecutableProcessor(q.loadBashBindings),
				),
				"load-bindings": command.SerialNodes(
					command.Description("Load bash bindings based on keyboard mode"),
					command.ExecutableProcessor(q.loadBashBindings),
				),
			},
		},
	)
}

func (q *QMK) loadBashBindings(o command.Output, d *command.Data) ([]string, error) {
	if q.Basic {
		o.Stdoutln(`Loading basic keyboard bash bindings...`)
		return basicKeyboardBindings, nil
	}
	o.Stdoutln(`Loading QMK keyboard bash bindings...`)
	return qmkKeyboardBindings, nil
}

func (q *QMK) Toggle(output command.Output, data *command.Data) error {
	q.Basic = !q.Basic
	q.changed = true
	return nil
}
