package mathutil

import "sort"

// ComplementsInt takes two int64 slices, and returns the complements of the
// two lists.
func ComplementsInt(a, b []int64) (aOnly, bOnly []int64) {
	aCopy := make([]int64, len(a))
	copy(aCopy, a)

	bCopy := make([]int64, len(b))
	copy(bCopy, b)

	sort.Sort(int64Sorted(aCopy))
	sort.Sort(int64Sorted(bCopy))

	var aIndex, bIndex int
	for aIndex < len(aCopy) && bIndex < len(bCopy) {
		if aCopy[aIndex] < bCopy[bIndex] {
			aOnly = append(aOnly, aCopy[aIndex])
			aIndex++

		} else if bCopy[bIndex] < aCopy[aIndex] {
			bOnly = append(bOnly, bCopy[bIndex])
			bIndex++

		} else {
			aIndex++
			bIndex++
		}
	}

	for aIndex < len(aCopy) {
		aOnly = append(aOnly, aCopy[aIndex])
		aIndex++
	}

	for bIndex < len(bCopy) {
		bOnly = append(bOnly, bCopy[bIndex])
		bIndex++
	}

	return aOnly, bOnly
}

// int64Sorted add the sort capability to a slice of int64.
type int64Sorted []int64

// Len is the number of elements in the collection.
func (i64 int64Sorted) Len() int {
	return len(i64)
}

// Less reports whether the element with index i should sort before the element
// with index j.
func (i64 int64Sorted) Less(i, j int) bool {
	return i64[i] < i64[j]
}

// Swap swaps the elements with indexes i and j.
func (i64 int64Sorted) Swap(i, j int) {
	i64[i], i64[j] = i64[j], i64[i]
}
