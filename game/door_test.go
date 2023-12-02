package game

import "testing"

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
		{
			name:   "door closing on slime kills slime",
			state:  `@x.#@D#`,
			inputs: []Direction{Right, Right},
			want:   `.x@#.D#`,
		},
		{
			name:   "door closing on box kills box",
			state:  `@x.#@BD#`,
			inputs: []Direction{Right, Right},
			want:   `.x@#.@D#`,
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
