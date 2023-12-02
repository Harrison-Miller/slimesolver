package game

import "testing"

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

func TestSmallSlime(t *testing.T) {
	tt := []testCase{
		{
			name:   "small slime can't push box",
			state:  `oB.`,
			inputs: []Direction{Right},
			want:   `oB.`,
		},
		{
			name:   "small slime can't activate switch",
			state:  `oxD`,
			inputs: []Direction{Right},
			want:   `.oD`,
		},
		{
			name:   "small slime dies in spike",
			state:  `o-`,
			inputs: []Direction{Right},
			want:   `.^`,
		},
		{
			name:   "small slime combines with another small slime",
			state:  `oo`,
			inputs: []Direction{Right},
			want:   `.@`,
		},
		{
			name:   "combine against door",
			state:  `ooD`,
			inputs: []Direction{Right},
			want:   `.@D`,
		},
		{
			name:   "don't combine against door that opens",
			state:  `ooD#@x`,
			inputs: []Direction{Right},
			want:   `.oo#.@`,
		},
		{
			name:   "combine against box",
			state:  `ooB.`,
			inputs: []Direction{Right},
			want:   `.@B.`,
		},
		{
			name:   "combine against big slime",
			state:  `oo@`,
			inputs: []Direction{Right},
			want:   `.@@`,
		},
		{
			name:   "don't combine against big slime that moves",
			state:  `oo@.`,
			inputs: []Direction{Right},
			want:   `.oo@`,
		},
		{
			name:   "can't move into small slime and vice a versa",
			state:  `o@#@o`,
			inputs: []Direction{Right},
			want:   `o@#@o`,
		},
		{
			name:   "big slime splits into small slimes",
			state:  `.@-`,
			inputs: []Direction{Right, Left},
			want:   `oo-`,
		},
		{
			name:   "combining slime activates switch",
			state:  `oox#D`,
			inputs: []Direction{Right, Right},
			want:   `..@#_`,
		},
		{
			name:   "combining slime against box activates switch",
			state:  `ooxBD`,
			inputs: []Direction{Right, Right},
			want:   `..@B_`,
		},
	}

	testCases(t, tt)
}
