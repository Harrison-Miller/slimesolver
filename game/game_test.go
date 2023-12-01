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

			g := Game{}
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

		g := Game{}
		err := g.Parse(state)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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

func TestBasicMovement(t *testing.T) {
	tt := []testCase{
		{
			name: "no move",
			state: `...
					.@.
					...`,
			inputs: []Direction{},
			want: `...
 				   .@.
			       ...`,
		},
		{
			name: "move up",
			state: `.
					@`,
			inputs: []Direction{Up},
			want: `@
            	   .`,
		},
		{
			name: "move down",
			state: `@
					.`,
			inputs: []Direction{Down},
			want: `.
                   @`,
		},
		{
			name:   "move left",
			state:  `.@`,
			inputs: []Direction{Left},
			want:   `@.`,
		},
		{
			name:   "move right",
			state:  `@.`,
			inputs: []Direction{Right},
			want:   `.@`,
		},
		{
			name:   "move into wall",
			state:  `@.#.`,
			inputs: []Direction{Right, Right},
			want:   `.@#.`,
		},
		{
			name: "move around wall",
			state: `@#.
					...`,
			inputs: []Direction{Down, Right, Right, Up},
			want: `.#@
                   ...`,
		},
		{
			name: "move multiple slimes",
			state: `@#@
					.#.`,
			inputs: []Direction{Down},
			want: `.#.
			       @#@`,
		},
	}

	testCases(t, tt)
}

func TestMovingTowardsEachother(t *testing.T) {
	tt := []testCase{
		{
			name:   "left 2 slimes",
			state:  `.@@`,
			inputs: []Direction{Left},
			want:   `@@.`,
		},
		{
			name:   "left 2 slimes w/ wall",
			state:  `#.@@`,
			inputs: []Direction{Left, Left},
			want:   `#@@.`,
		},
		{
			name:   "right 2 slimes",
			state:  `@@.`,
			inputs: []Direction{Right},
			want:   `.@@`,
		},
		{
			name:   "right 2 slimes w/ wall",
			state:  `@@.#`,
			inputs: []Direction{Right, Right},
			want:   `.@@#`,
		},
		{
			name: "up 2 slimes",
			state: `.
					@
					@`,
			inputs: []Direction{Up},
			want: `@
			       @
			       .`,
		},
		{
			name: "up 2 slimes w/ wall",
			state: `#
			        .
					@
					@`,
			inputs: []Direction{Up, Up},
			want: `#
			       @
			       @
			       .`,
		},
		{
			name: "down 2 slimes",
			state: `@
					@
					.`,
			inputs: []Direction{Down},
			want: `.
			       @	
			       @`,
		},
		{
			name: "down 2 slimes w/ wall",
			state: `@
					@
					.
					#`,
			inputs: []Direction{Down, Down},
			want: `.
			       @	
				   @
				   #`,
		},
	}

	testCases(t, tt)
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

func TestBoxes(t *testing.T) {
	tt := []testCase{
		{
			name:   "boxes don't move on their own",
			state:  `.#.`,
			inputs: []Direction{Up, Down, Left, Right},
			want:   `.#.`,
		},
		{
			name:   "push right",
			state:  `@B.`,
			inputs: []Direction{Right},
			want:   `.@B`,
		},
		{
			name:   "push right into wall",
			state:  `@B#.`,
			inputs: []Direction{Right},
			want:   `@B#.`,
		},
		{
			name:   "push left",
			state:  `.B@`,
			inputs: []Direction{Left},
			want:   `B@.`,
		},
		{
			name:   "push left into wall",
			state:  `.#B@`,
			inputs: []Direction{Left},
			want:   `.#B@`,
		},
		{
			name: "push up",
			state: `.
					 B
					 @`,
			inputs: []Direction{Up},
			want: `B
			       @	
				   .`,
		},
		{
			name: "push up into wall",
			state: `.
					 #
					 B
					 @`,
			inputs: []Direction{Up},
			want: `.
			       #
				   B
				   @`,
		},
		{
			name: "push down",
			state: `@
					 B
					 .`,
			inputs: []Direction{Down},
			want: `.
			       @
				   B`,
		},
		{
			name: "push down into wall",
			state: `@
					 B
					 #
					 .`,
			inputs: []Direction{Down},
			want: `@
			       B	
				   #
				   .`,
		},
		{
			name:   "can't push multiple boxes",
			state:  `@BB.`,
			inputs: []Direction{Right},
			want:   `@BB.`,
		},
		{
			name:   "push with gap",
			state:  `@B.B.`,
			inputs: []Direction{Right},
			want:   `.@BB.`,
		},
		{
			name:   "chain push",
			state:  `@B@B..`,
			inputs: []Direction{Right, Right},
			want:   `..@B@B`,
		},
	}

	testCases(t, tt)
}

func TestBoxWithPit(t *testing.T) {
	testGame(t, testCase{
		name:   "box falls into pit",
		state:  `.@.B.O.#`,
		inputs: []Direction{Right, Right, Right, Right, Right},
		want:   `......@#`,
	})
}

func TestSwitchAndDoor(t *testing.T) {
	tt := []testCase{
		{
			name:   "door blocks movement",
			state:  `@D.`,
			inputs: []Direction{Right, Right},
			want:   `@D.`,
		},
		{
			name:   "switch doesn't block movement",
			state:  `@x.`,
			inputs: []Direction{Right, Right},
			want:   `.x@`,
		},
		{
			name:   "slime on switch opens door",
			state:  `@xD.`,
			inputs: []Direction{Right},
			want:   `.@_.`,
		},
		{
			name:   "door stays open when slime doesn't move off",
			state:  `@x#D`,
			inputs: []Direction{Right, Right},
			want:   `.@#_`,
		},
		{
			name:   "box on switch opens door",
			state:  `@BxD`,
			inputs: []Direction{Right},
			want:   `.@B_`,
		},
		{
			name:   "box can't go through closed door",
			state:  `@BD.`,
			inputs: []Direction{Right, Right},
			want:   `@BD.`,
		},
		{
			name: "slime can go through open door",
			state: `@x#.
					@.D.`,
			inputs: []Direction{Right, Right, Right},
			want: `.@#.
			       .._@`,
		},
		{
			name: "box can go through open door",
			state: `@x#..
					@BD..`,
			inputs: []Direction{Right, Right, Right},
			want: `.@#..
			       .._@B`,
		},
		{
			name:   "door closes when slime moves off switch",
			state:  `@x.D`,
			inputs: []Direction{Right, Right},
			want:   `.x@D`,
		},
		{
			name:   "door stays open when box moves off but slimes moves on",
			state:  `@Bx..D`,
			inputs: []Direction{Right, Right},
			want:   `..@B._`,
		},
		{
			name: "can't move through a closing door",
			state: `@x.#.
					@.D..`,
			inputs: []Direction{Right, Right, Right},
			want: `.x@#.
			       .@D..`,
		},
		{
			name:   "the problem (tm)",
			state:  `@@D.@Bx..`,
			inputs: []Direction{Right, Right, Right},
			want:   `..D@@.x@B`,
		},
		{
			name:   "slime push box on switch stay still",
			state:  `@Bx#D`,
			inputs: []Direction{Right, Right},
			want:   `.@B#_`,
		},
	}

	testCases(t, tt)
}

func TestBoxThroughDoor(t *testing.T) {
	tt := []testCase{
		{
			name:   "box through door",
			state:  `@x#@BD..`,
			inputs: []Direction{Right, Right, Right},
			want:   `.@#.._@B`,
		},
		{
			name:   "door stays open with box on switch",
			state:  `xB@#D`,
			inputs: []Direction{Left, Right},
			want:   `B.@#_`,
		},
		{
			name:   "move slime through door with box on switch",
			state:  `xB@D.`,
			inputs: []Direction{Left, Right, Right, Right},
			want:   `B.._@`,
		},
		{
			name:   "move box through door with box on switch",
			state:  `xB@BD..`,
			inputs: []Direction{Left, Right, Right, Right, Right},
			want:   `B..._@B`,
		},
		{
			name:   "move box through door with box on switch",
			state:  `xB@BD..`,
			inputs: []Direction{Left, Right, Right, Right},
			want:   `B...@B.`,
		},
		{
			name:   "move box through 2 doors with box on switch",
			state:  `xB@BDD..`,
			inputs: []Direction{Left, Right, Right, Right, Right, Right},
			want:   `B...__@B`,
		},
	}

	testCases(t, tt)
}
