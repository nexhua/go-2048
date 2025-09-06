package game

import (
	"errors"
	"fmt"
)

// TODO add tests for merges
// Merge a horizontal slice
// Merge direction is 0 or SIZE-1, any other value will be rejected
func MergeSlice(slice []Cell, to int) (int, error) {
	mergeScore := 0

	if !(to == 0 || to == CELL_COUNT-1) {
		return 0, errors.New("invalid to argument")
	}

	if len(slice) < 2 {
		return 0, errors.New("slice is too small to merge")
	}

	// merge to start of slice
	if to == 0 {
		for i := 0; i < len(slice)-1; i++ {
			lc, rc := &slice[i], &slice[i+1]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				mergeScore += lc.val * 2
				lc.val += rc.val
				rc.isRendered = false
				rc.val = 0
			}
		}
	} else {
		for i := len(slice) - 1; i > 0; i-- {
			lc, rc := &slice[i-1], &slice[i]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				mergeScore += lc.val * 2
				rc.val += lc.val
				lc.isRendered = false
				lc.val = 0
			}
		}
	}

	return mergeScore, nil
}

// Merge a vertical (accepts a ref slice)
// Merge direction is 0 or SIZE-1, any other value will be rejected
func MergeSliceRef(slice []*Cell, to int) (int, error) {
	mergeScore := 0

	if !(to == 0 || to == CELL_COUNT-1) {
		fmt.Println("invalid to argument")
		return 0, errors.New("invalid to argument")
	}

	if len(slice) < 2 {
		fmt.Println("slice is too small to merge")
		return 0, errors.New("slice is too small to merge")
	}

	// merge to start of slice
	if to == 0 {
		for i := 0; i < len(slice)-1; i++ {
			lc, rc := slice[i], slice[i+1]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				mergeScore += lc.val * 2
				lc.val += rc.val
				rc.isRendered = false
				rc.val = 0
			}
		}
	} else {
		for i := len(slice) - 1; i > 0; i-- {
			lc, rc := slice[i-1], slice[i]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				mergeScore += lc.val * 2
				rc.val += lc.val
				lc.isRendered = false
				lc.val = 0
			}
		}
	}

	return mergeScore, nil
}

// assumes no empty cell exists to check possible merges
func HasPossibleMerge(cells [][]Cell) bool {
	for _, row := range cells {
		for i := 0; i < len(row)-1; i++ {
			lc, rc := &row[i], &row[i+1]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				return true
			}
		}
	}

	for i := range len(cells) {
		col, err := TakeVerticalSlice(cells, i)

		if err == nil {
			for i := 0; i < len(col)-1; i++ {
				lc, rc := col[i], col[i+1]
				if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
					return true
				}
			}
		}

	}

	return false
}
