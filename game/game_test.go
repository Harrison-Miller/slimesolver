package game

import (
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

func TestBasicMovement(t *testing.T) {
	tt := []struct {
		name   string
		state  string
		inputs []Direction
		want   string
	}{
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

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			state := cleanState(tc.state)

			g := Game{}
			err := g.Parse(state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, dir := range tc.inputs {
				g.Move(dir)
			}

			want := cleanState(tc.want)
			result := strings.TrimSpace(g.String())
			if result != want {
				t.Fatalf("expected:\n%s\ngot\n%s", want, result)
			}
		})
	}
}

func TestMovingTowardsEachother(t *testing.T) {
	tt := []struct {
		name   string
		state  string
		inputs []Direction
		want   string
	}{
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

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			state := cleanState(tc.state)

			g := Game{}
			err := g.Parse(state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, dir := range tc.inputs {
				g.Move(dir)
			}

			want := cleanState(tc.want)
			result := strings.TrimSpace(g.String())
			if result != want {
				t.Fatalf("expected:\n%s\ngot\n%s", want, result)
			}
		})
	}
}

func TestPit(t *testing.T) {
	game := Game{}
	err := game.Parse(`.@.
					   		 ..O`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inputs := []Direction{Down, Right}
	for _, dir := range inputs {
		game.Move(dir)
	}

	want := cleanState(`...
							  ..O`)
	result := strings.TrimSpace(game.String())
	if result != want {
		t.Fatalf("expected:\n%s\ngot\n%s", want, result)
	}
}

func TestBoxes(t *testing.T) {
	tt := []struct {
		name   string
		state  string
		inputs []Direction
		want   string
	}{
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

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			state := cleanState(tc.state)

			g := Game{}
			err := g.Parse(state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, dir := range tc.inputs {
				g.Move(dir)
			}

			want := cleanState(tc.want)
			result := strings.TrimSpace(g.String())
			if result != want {
				t.Fatalf("expected:\n%s\ngot\n%s", want, result)
			}
		})
	}
}

func TestBoxWithPit(t *testing.T) {
	game := Game{}
	err := game.Parse(`.@.B.O.#`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inputs := []Direction{Right, Right, Right, Right, Right}
	for _, dir := range inputs {
		game.Move(dir)
	}

	// box fills pit
	want := cleanState(`......@#`)
	result := strings.TrimSpace(game.String())
	if result != want {
		t.Fatalf("expected:\n%s\ngot\n%s", want, result)
	}
}
