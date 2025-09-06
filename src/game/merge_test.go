package game

import "testing"

type MergeTest struct {
	name                  string
	expectedMergePossible bool
	inputCells            [][]Cell
}

func TestHasPossibleMerge(t *testing.T) {
	testCases := []MergeTest{
		{
			name:                  "empty cells",
			expectedMergePossible: false,
			inputCells:            [][]Cell{{m(0), m(0), m(0), m(0)}, {m(0), m(0), m(0), m(0)}, {m(0), m(0), m(0), m(0)}, {m(0), m(0), m(0), m(0)}},
		},
		{
			name:                  "full grid but no merge possible",
			expectedMergePossible: false,
			inputCells:            [][]Cell{{m(2), m(4), m(8), m(16)}, {m(32), m(64), m(128), m(256)}, {m(512), m(1024), m(2048), m(4096)}, {m(8192), m(16384), m(32768), m(65536)}},
		},
		{
			name:                  "full grid with horizontal merge possible",
			expectedMergePossible: true,
			inputCells:            [][]Cell{{m(2), m(4), m(8), m(16)}, {m(32), m(64), m(128), m(256)}, {m(512), m(512), m(2048), m(4096)}, {m(8192), m(16384), m(32768), m(65536)}},
		},
		{
			name:                  "full grid with vertical merge possible",
			expectedMergePossible: true,
			inputCells:            [][]Cell{{m(2), m(4), m(8), m(16)}, {m(32), m(1024), m(128), m(256)}, {m(512), m(1024), m(2048), m(4096)}, {m(8192), m(16384), m(32768), m(65536)}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualMergePossible := HasPossibleMerge(tc.inputCells)
			if actualMergePossible != tc.expectedMergePossible {
				t.Errorf("expected %t, found %t", tc.expectedMergePossible, actualMergePossible)
			}
		})
	}
}

func m(val int) Cell {
	if val == 0 {
		return Cell{val: 0, isRendered: false}
	} else {
		return Cell{val: val, isRendered: true}
	}
}
