package game

import "testing"

func TestBoxes(t *testing.T) {
	tt := []testCase{
		{
			name:   "boxes don't move on their own",
			state:  `.B.`,
			inputs: []Direction{Up, Down, Left, Right},
			want:   `.B.`,
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
		state:  `@BO.`,
		inputs: []Direction{Right, Right, Right},
		want:   `...@`,
	})
}
