package game

import (
	"testing"
)

type TestCellCmp struct {
	name                   string
	expectedNumOfMovements int
	expectedCells          []Cell
	actualCells            []Cell
	expectedCellsRef       []*Cell
	actualCellsRef         []*Cell
}

func TestShiftRight(t *testing.T) {
	testCases := []TestCellCmp{
		{
			name:                   "shift single cell at rightest position to right",
			expectedNumOfMovements: 0,
			expectedCells:          []Cell{Cell{}, Cell{}, Cell{}, Cell{val: 2, isRendered: true}},
			actualCells:            []Cell{Cell{}, Cell{}, Cell{}, Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift cell right",
			expectedNumOfMovements: 1,
			expectedCells:          []Cell{Cell{}, Cell{}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
			actualCells:            []Cell{Cell{val: 2, isRendered: true}, Cell{}, Cell{}, Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift multiple in tandem cells to right",
			expectedNumOfMovements: 2,
			expectedCells:          []Cell{Cell{}, Cell{}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
			actualCells:            []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{}, Cell{}},
		},
		{
			name:                   "shift multiple cells to right",
			expectedNumOfMovements: 2,
			expectedCells:          []Cell{Cell{}, Cell{}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
			actualCells:            []Cell{Cell{val: 2, isRendered: true}, Cell{}, Cell{val: 2, isRendered: true}, Cell{}},
		},
		{
			name:                   "shift when all cells exists to right",
			expectedNumOfMovements: 0,
			expectedCells:          []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
			actualCells:            []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualNumOfMovements := ShiftRight(tc.actualCells)
			compareCells(t, tc, actualNumOfMovements)
		})
	}
}

func TestShiftLeft(t *testing.T) {
	testCases := []TestCellCmp{
		{
			name:                   "shift single cell at leftest position to left",
			expectedNumOfMovements: 0,
			expectedCells:          []Cell{Cell{val: 2, isRendered: true}, Cell{}, Cell{}, Cell{}},
			actualCells:            []Cell{Cell{val: 2, isRendered: true}, Cell{}, Cell{}, Cell{}},
		},
		{
			name:                   "shift cell left",
			expectedNumOfMovements: 1,
			expectedCells:          []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{}, Cell{}},
			actualCells:            []Cell{Cell{val: 2, isRendered: true}, Cell{}, Cell{}, Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift multiple in tandem cells to left",
			expectedNumOfMovements: 2,
			expectedCells:          []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{}, Cell{}},
			actualCells:            []Cell{Cell{}, Cell{}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift multiple cells to left",
			expectedNumOfMovements: 2,
			expectedCells:          []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{}, Cell{}},
			actualCells:            []Cell{Cell{}, Cell{val: 2, isRendered: true}, Cell{}, Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift when all cells exists to left",
			expectedNumOfMovements: 0,
			expectedCells:          []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
			actualCells:            []Cell{Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}, Cell{val: 2, isRendered: true}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualNumOfMovements := ShiftLeft(tc.actualCells)
			compareCells(t, tc, actualNumOfMovements)
		})
	}
}

func TestShiftUp(t *testing.T) {
	testCases := []TestCellCmp{
		{
			name:                   "shift single cell at top position up",
			expectedNumOfMovements: 0,
			expectedCellsRef:       []*Cell{&Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}, &Cell{}},
			actualCellsRef:         []*Cell{&Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}, &Cell{}},
		},
		{
			name:                   "shift cell up",
			expectedNumOfMovements: 1,
			expectedCellsRef:       []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}},
			actualCellsRef:         []*Cell{&Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift multiple in tandem cells up",
			expectedNumOfMovements: 2,
			expectedCellsRef:       []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}},
			actualCellsRef:         []*Cell{&Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift multiple cells up",
			expectedNumOfMovements: 2,
			expectedCellsRef:       []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}},
			actualCellsRef:         []*Cell{&Cell{}, &Cell{val: 2, isRendered: true}, &Cell{}, &Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift when all cells exists up",
			expectedNumOfMovements: 0,
			expectedCellsRef:       []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
			actualCellsRef:         []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualNumOfMovements := ShiftUp(tc.actualCellsRef)
			compareCells(t, tc, actualNumOfMovements)
		})
	}
}

func TestShiftDown(t *testing.T) {
	testCases := []TestCellCmp{
		{
			name:                   "shift single cell at bottom position down",
			expectedNumOfMovements: 0,
			expectedCellsRef:       []*Cell{&Cell{}, &Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}},
			actualCellsRef:         []*Cell{&Cell{}, &Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift cell down",
			expectedNumOfMovements: 1,
			expectedCellsRef:       []*Cell{&Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
			actualCellsRef:         []*Cell{&Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}},
		},
		{
			name:                   "shift multiple in tandem cells down",
			expectedNumOfMovements: 2,
			expectedCellsRef:       []*Cell{&Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
			actualCellsRef:         []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{}, &Cell{}},
		},
		{
			name:                   "shift multiple cells down",
			expectedNumOfMovements: 2,
			expectedCellsRef:       []*Cell{&Cell{}, &Cell{}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
			actualCellsRef:         []*Cell{&Cell{val: 2, isRendered: true}, &Cell{}, &Cell{val: 2, isRendered: true}, &Cell{}},
		},
		{
			name:                   "shift when all cells exists down",
			expectedNumOfMovements: 0,
			expectedCellsRef:       []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
			actualCellsRef:         []*Cell{&Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}, &Cell{val: 2, isRendered: true}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualNumOfMovements := ShiftDown(tc.actualCellsRef)
			compareCells(t, tc, actualNumOfMovements)
		})
	}
}

func compareCells(t *testing.T, testCase TestCellCmp, actualNumOfMovements int) {
	if actualNumOfMovements != testCase.expectedNumOfMovements {
		t.Errorf("Expected %d, found %d", testCase.expectedNumOfMovements, testCase.expectedNumOfMovements)
	}

	if len(testCase.actualCells) != len(testCase.expectedCells) {
		t.Errorf("Expected %d, found %d", len(testCase.expectedCells), len(testCase.actualCells))
	}

	for i := 0; i < len(testCase.actualCells); i++ {
		c := &testCase.actualCells[i]
		e := &testCase.expectedCells[i]

		if e.isRendered != c.isRendered || e.val != c.val {
			t.Errorf("cell movement failed, at pos %d cells are not equal", i)
		}
	}
}
