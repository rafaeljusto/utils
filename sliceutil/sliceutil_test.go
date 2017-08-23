package sliceutil

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/teamwork/test/diff"
)

func TestIntsToString(t *testing.T) {
	cases := []struct {
		in       []int64
		expected string
	}{
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			"1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8",
		},
		{
			[]int64{-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8},
			"-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8",
		},
		{
			[]int64{},
			"",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := JoinInt(tc.in)
			if got != tc.expected {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniqInt64(t *testing.T) {
	cases := []struct {
		in       []int64
		expected []int64
	}{
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			[]int64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			[]int64{1, 3, 8, 3, 8},
			[]int64{1, 3, 8},
		},
		{
			[]int64{1, 2, 3},
			[]int64{1, 2, 3},
		},
		{
			[]int64{},
			nil,
		},
		{
			nil,
			nil,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := UniqInt64(tc.in)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniqueMergeSlices(t *testing.T) {
	var tests = []struct {
		in       [][]int64
		expected []int64
	}{
		{
			generate2dintslice([]int64{1, 2, 3}),
			[]int64{1, 2, 3},
		},
		{
			generate2dintslice([]int64{0, 1, 2, 3, -1, -10}),
			[]int64{0, 1, 2, 3, -1, -10},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := UniqueMergeSlices(tc.in)
			if !int64slicesequal(got, tc.expected) {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniqString(t *testing.T) {
	var tests = []struct {
		in       []string
		expected []string
	}{
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			[]string{"a", "b", "c", "a", "b", "n", "a", "aaa", "n", "x"},
			[]string{"a", "b", "c", "n", "aaa", "x"},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := UniqString(tc.in)
			if !stringslicesequal(got, tc.expected) {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func int64slicesequal(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for _, ia := range a {
		var found bool
		for _, ib := range b {
			if ib == ia {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func stringslicesequal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, ia := range a {
		var found bool
		for _, ib := range b {
			if ib == ia {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func generate2dintslice(in []int64) [][]int64 {
	var (
		result [][]int64
		loops  = int(rand.Int63n(int64(len(in) * 2)))
	)

	for i := 0; i < loops; i++ {
		var s []int64
		for i := 0; i < loops; i++ {
			s = append(s, in[rand.Intn(len(in))])
		}
		result = append(result, s)
	}

	return result
}

func TestCSVtoInt64Slice(t *testing.T) {
	tests := []struct {
		in          string
		expected    []int64
		expectedErr error
	}{
		{
			"1,2,3",
			[]int64{1, 2, 3},
			nil,
		},
		{
			"",
			[]int64(nil),
			nil,
		},
		{
			"1,				2, \n3",
			[]int64{1, 2, 3},
			nil,
		},
		{
			"1,				2,nope",
			[]int64(nil),
			errors.New("invalid syntax"),
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got, err := CSVtoInt64Slice(tc.in)

			if err != nil {
				if numErrorer, ok := err.(*strconv.NumError); ok {
					err = numErrorer.Err
				}
			}

			if err != tc.expectedErr && err.Error() != tc.expectedErr.Error() {
				t.Errorf(diff.Cmp(tc.expectedErr.Error(), err.Error()))
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestInStringSlice(t *testing.T) {
	tests := []struct {
		list     []string
		find     string
		expected bool
	}{
		{[]string{"hello"}, "hello", true},
		{[]string{"hello"}, "hell", false},
		{[]string{"hello", "world", "test"}, "world", true},
		{[]string{"hello", "world", "test"}, "", false},
		{[]string{}, "", false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := InStringSlice(tc.list, tc.find)
			if got != tc.expected {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestInIntSlice(t *testing.T) {
	tests := []struct {
		list     []int
		find     int
		expected bool
	}{
		{[]int{42}, 42, true},
		{[]int{42}, 4, false},
		{[]int{42, 666, 14159}, 666, true},
		{[]int{42, 666, 14159}, 0, false},
		{[]int{}, 0, false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := InIntSlice(tc.list, tc.find)
			if got != tc.expected {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestInInt64Slice(t *testing.T) {
	tests := []struct {
		list     []int64
		find     int64
		expected bool
	}{
		{[]int64{42}, 42, true},
		{[]int64{42}, 4, false},
		{[]int64{42, 666, 14159}, 666, true},
		{[]int64{42, 666, 14159}, 0, false},
		{[]int64{}, 0, false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := InInt64Slice(tc.list, tc.find)
			if got != tc.expected {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestDifference(t *testing.T) {
	cases := []struct {
		inSet    []int64
		inOthers [][]int64
		expected []int64
	}{
		{[]int64{}, [][]int64{}, []int64{}},
		{nil, [][]int64{}, []int64{}},
		{[]int64{}, nil, []int64{}},
		{nil, nil, []int64{}},
		{[]int64{1}, [][]int64{{1}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2, 2, 3}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2}, {3}}, []int64{}},
		{[]int64{1, 2}, [][]int64{{1}}, []int64{2}},
		{[]int64{1, 2, 3}, [][]int64{{1}}, []int64{2, 3}},
		{[]int64{1, 2, 3}, [][]int64{{}, {1}}, []int64{2, 3}},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Difference(tc.inSet, tc.inOthers...)
			if !reflect.DeepEqual(tc.expected, out) {
				t.Errorf("\nout:      %#v\nexpected: %#v\n", out, tc.expected)
			}
		})
	}
}