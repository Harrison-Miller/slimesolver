package game

import "testing"

func TestGoThroughDoorDifferentHeight(t *testing.T) {
	tt := []testCase{
		{
			name: "go through door",
			state: `.#@BD..
			  	    xB@....`,
			inputs: []Direction{Left, Right, Right, Right},
			want: `.#.._@B
				   B...@..`,
		},
		{
			name:   "go through door one row",
			state:  `xB@#@BD..`,
			inputs: []Direction{Left, Right, Right, Right},
			want:   `B.@#.._@B`,
		},
	}

	testCases(t, tt)
}
