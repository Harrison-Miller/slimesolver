package game

import "testing"

func TestSpikes(t *testing.T) {
	tt := []testCase{
		{
			name:   "spike flip flop",
			state:  `-^`,
			inputs: []Direction{Right},
			want:   `^-`,
		},
		{
			name:   "spike flip flop flip flop",
			state:  `-^`,
			inputs: []Direction{Right, Right},
			want:   `-^`,
		},
		{
			name:   "spikes don't kill box",
			state:  `.@B-`,
			inputs: []Direction{Right},
			want:   `..@B`,
		},
	}

	testCases(t, tt)
}
