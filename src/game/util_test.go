package game

import "testing"

func TestTakeVerticalSlice(t *testing.T) {
	size := 4
	grid := makeGrid(size, false)
	expectedSlice := make([]*Cell, size)
	expectedSlice[0] = &Cell{val: 1}
	expectedSlice[1] = &Cell{val: 5}
	expectedSlice[2] = &Cell{val: 9}
	expectedSlice[3] = &Cell{val: 13}

	actualSlice, _ := TakeVerticalSlice(grid, 1)
	for i := range actualSlice {
		if actualSlice[i].val != expectedSlice[i].val {
			t.Errorf("Slices are not equal. Expected %d, received %d at index %d", expectedSlice[i].val, actualSlice[i].val, i)
		}
	}
}

func TestTakeVerticalSliceErr(t *testing.T) {
	size := 4
	grid := makeGrid(size, false)
	selectedSize := 7

	_, err := TakeVerticalSlice(grid, selectedSize)
	if err == nil {
		t.Errorf("Expected error. %d is not within bounds of grid with size %d", selectedSize, size)
	}

	selectedSize = -1
	_, err = TakeVerticalSlice(grid, selectedSize)
	if err == nil {
		t.Errorf("Expected error. %d is not within bounds of grid with size %d", selectedSize, size)
	}
}

func TestGetRandomCell(t *testing.T) {
	grid := makeGrid(4, false)

	cell, err := GetRandomCell(grid)
	if err != nil || cell.val != 0 {
		t.Errorf("Get random cell failed")
	}
}

func TestGetRandomCellErr(t *testing.T) {
	grid := makeGrid(4, false)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			grid[i][j].isRendered = true
		}
	}

	_, err := GetRandomCell(grid)
	if err == nil {
		t.Errorf("When all cells are filled, get random cell should fail with an error.")
	}

	grid = makeGrid(0, false)
	_, err = GetRandomCell(grid)
	if err == nil {
		t.Errorf("When grid does not have a valid size, get random cell should fail with an error.")
	}
}

func TestHasEmptyCell(t *testing.T) {
	// fully occupied grid
	grid := makeGrid(4, true)
	expectedHasEmptyCell := false
	actualHasEmptyCell := HasEmptyCell(grid)

	if expectedHasEmptyCell != actualHasEmptyCell {
		t.Errorf("expected %t, found %t", expectedHasEmptyCell, actualHasEmptyCell)
	}

	// grid with one cell emptied
	c := &grid[2][2]
	c.val = 0
	c.isRendered = false

	expectedHasEmptyCell = true
	actualHasEmptyCell = HasEmptyCell(grid)

	if expectedHasEmptyCell != actualHasEmptyCell {
		t.Errorf("expected %t, found %t", expectedHasEmptyCell, actualHasEmptyCell)
	}
}

func makeGrid(size int, isRendered bool) [][]Cell {
	grid := make([][]Cell, size)

	for i := range grid {
		grid[i] = make([]Cell, size)
	}

	for i, row := range grid {
		for j := range row {
			grid[i][j] = Cell{val: i*size + j, isRendered: isRendered}
		}
	}

	return grid
}
