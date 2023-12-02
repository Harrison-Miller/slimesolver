package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestBasicBoard(t *testing.T) {
	tt := []struct {
		name  string
		state string
		err   bool
	}{
		{
			name: "mismatched width",
			state: `###
					##`,
			err: true,
		},
		{
			name: "just walls",
			state: `##
					##`,
		},
		{
			name: "just pits",
			state: `OO
					OO`,
		},
		{
			name: "just empty",
			state: `..
					..`,
		},
		{
			name: "standard room",
			state: `#####
					#..##
					#...#	
					#...#
					#####`,
		},
	}

	// test that the board is parsed correctly
	// and String returns the same board
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// we add whitespace for readability
			state := cleanState(tc.state)

			g := NewGame(true)
			err := g.Parse(state)
			if tc.err && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.err {
				return
			}

			result := strings.TrimSpace(g.String())
			if result != state {
				t.Fatalf("expected:\n%s\ngot\n%s", state, result)
			}
		})

	}
}

type testCase struct {
	name   string
	state  string
	inputs []Direction
	want   string
}

func testCases(t *testing.T, cases []testCase) {
	for _, tc := range cases {
		testGame(t, tc)
	}
}

func testGame(t *testing.T, tc testCase) {
	t.Run(tc.name, func(t *testing.T) {
		state := cleanState(tc.state)

		g := NewGame(true)
		err := g.Parse(state)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		fmt.Println(g.String())

		for _, dir := range tc.inputs {
			g.Move(dir)
			fmt.Println(g.String())
		}

		want := cleanState(tc.want)
		result := strings.TrimSpace(g.String())
		if result != want {
			t.Fatalf("expected:\n%s\ngot\n%s", want, result)
		}
	})
}

func TestPit(t *testing.T) {
	testGame(t, testCase{
		name: "pit",
		state: `.@.
				..O`,
		inputs: []Direction{Down, Right},
		want: `...
			   ..O`,
	})
}
