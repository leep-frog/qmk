package qmk

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/leep-frog/command"
)

func TestQMKAutocomplete(t *testing.T) {
	for _, test := range []struct {
		name string
		ctc  *command.CompleteTestCase
	}{
		{
			name: "completes subcommands",
			ctc: &command.CompleteTestCase{
				Want: []string{
					"load-bindings",
					"toggle",
				},
			},
		},
		/* Useful for commenting out tests. */
	} {
		t.Run(test.name, func(t *testing.T) {
			test.ctc.Node = CLI().Node()
			command.CompleteTest(t, test.ctc)
		})
	}
}

func TestQMKExecution(t *testing.T) {
	for _, test := range []struct {
		name  string
		basic bool
		etc   *command.ExecuteTestCase
		want  *QMK
	}{
		{
			name: "Requires argument",
			etc: &command.ExecuteTestCase{
				WantStderr: "Branching argument must be one of [load-bindings toggle t]\n",
				WantErr:    fmt.Errorf("Branching argument must be one of [load-bindings toggle t]"),
			},
		},
		{
			name: "Toggles from qmk to basic",
			etc: &command.ExecuteTestCase{
				Args:       []string{"toggle"},
				WantStdout: "Loading QMK keyboard bash bindings...\n",
				WantExecuteData: &command.ExecuteData{
					Executable: qmkKeyboardBindings,
				},
			},
			want: &QMK{true, true},
		},
		{
			name:  "Toggles from basic to qmk",
			basic: true,
			etc: &command.ExecuteTestCase{
				Args:       []string{"toggle"},
				WantStdout: "Loading basic keyboard bash bindings...\n",
				WantExecuteData: &command.ExecuteData{
					Executable: basicKeyboardBindings,
				},
			},
			want: &QMK{false, true},
		},
		/* Useful for commenting out tests. */
	} {
		t.Run(test.name, func(t *testing.T) {
			q := &QMK{test.basic, false}
			test.etc.Node = q.Node()
			command.ExecuteTest(t, test.etc)
			command.ChangeTest(t, test.want, q, cmpopts.IgnoreFields(QMK{}, "changed"))
		})
	}
}

func TestQMKMetadata(t *testing.T) {
	q := CLI()
	command.UsageTest(t, &command.UsageTestCase{
		Node: q.Node(),
		WantString: []string{
			"CLI for qmk/basic keyboard bindings in bash",
			"<",
			"",
			"  Load bash bindings based on keyboard mode",
			"  load-bindings",
			"",
			"  Toggles QMK mode",
			"  [toggle|t]",
			"",
			"Symbols:",
			command.BranchDesc,
		},
	})

	if setup := q.Setup(); setup != nil {
		t.Errorf("Setup() returned %v; expected nil", setup)
	}

	if name := q.Name(); name != "q" {
		t.Errorf("Setup() returned %q; expected %q", name, "q")
	}
}
