package sortz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestAsc(t *testing.T) {
	testsInt := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 3, 2}, []int{1, 2, 3}},
		{[]int{1, 2, 3, -1, 3, 9, 8, -2, 5}, []int{-2, -1, 1, 2, 3, 3, 5, 8, 9}},
	}

	for _, tt := range testsInt {
		Asc(tt.in)
		testz.Equal(t, tt.want, tt.in, tt.want)
	}

	testsFloat := []struct {
		in   []float64
		want []float64
	}{
		{[]float64{1.1, 3.3, 2.2}, []float64{1.1, 2.2, 3.3}},
		{[]float64{1.1, 2.2, 3.3, -1.1, 3.3, 9.9, 8.8, -2.2, 5.5}, []float64{-2.2, -1.1, 1.1, 2.2, 3.3, 3.3, 5.5, 8.8, 9.9}},
	}

	for _, tt := range testsFloat {
		Asc(tt.in)
		testz.Equal(t, tt.want, tt.in, tt.want)
	}
}

func TestSortDesc(t *testing.T) {
	testsInt := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 3, 2}, []int{3, 2, 1}},
		{[]int{1, 2, 3, -1, 3, 9, 8, -2, 5}, []int{9, 8, 5, 3, 3, 2, 1, -1, -2}},
	}

	for _, tt := range testsInt {
		Desc(tt.in)
		testz.Equal(t, tt.want, tt.in, tt.want)
	}

	testsFloat := []struct {
		in   []float64
		want []float64
	}{
		{[]float64{1.1, 3.3, 2.2}, []float64{3.3, 2.2, 1.1}},
		{[]float64{1.1, 2.2, 3.3, -1.1, 3.3, 9.9, 8.8, -2.2, 5.5}, []float64{9.9, 8.8, 5.5, 3.3, 3.3, 2.2, 1.1, -1.1, -2.2}},
	}

	for _, tt := range testsFloat {
		Desc(tt.in)
		testz.Equal(t, tt.want, tt.in, tt.want)
	}
}
