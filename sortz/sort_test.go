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

func TestAscByKey(t *testing.T) {
	type p struct {
		Age  int
		Name string
	}
	tests := struct {
		in   []p
		want []int
	}{
		in: []p{
			{
				Age:  10,
				Name: "bob",
			},
			{
				Age:  13,
				Name: "jack",
			},
			{
				Age:  9,
				Name: "monika",
			},
			{
				Age:  10,
				Name: "julia",
			},
			{
				Age:  10,
				Name: "tim",
			},
			{
				Age:  8,
				Name: "steven",
			},
		},
		want: []int{8, 9, 10, 10, 10, 13},
	}

	AscByKey(tests.in, func(p p) int {
		return p.Age
	})

	for i, age := range tests.in {
		testz.Equal(t, tests.want[i], age.Age, tests.want[i])
	}
}

func TestDescByKey(t *testing.T) {
	type p struct {
		Age  int
		Name string
	}
	tests := struct {
		in   []p
		want []int
	}{
		in: []p{
			{
				Age:  10,
				Name: "bob",
			},
			{
				Age:  13,
				Name: "jack",
			},
			{
				Age:  9,
				Name: "monika",
			},
			{
				Age:  10,
				Name: "julia",
			},
			{
				Age:  10,
				Name: "tim",
			},
			{
				Age:  8,
				Name: "steven",
			},
		},
		want: []int{13, 10, 10, 10, 9, 8},
	}

	DescByKey(tests.in, func(p p) int {
		return p.Age
	})

	for i, age := range tests.in {
		testz.Equal(t, tests.want[i], age.Age, tests.want[i])
	}
}

func TestAscStableByKey(t *testing.T) {
	type p struct {
		Age  int
		Name string
	}
	tests := struct {
		in   []p
		want []p
	}{
		in: []p{
			{
				Age:  10,
				Name: "bob",
			},
			{
				Age:  13,
				Name: "jack",
			},
			{
				Age:  9,
				Name: "monika",
			},
			{
				Age:  10,
				Name: "julia",
			},
			{
				Age:  10,
				Name: "tim",
			},
			{
				Age:  8,
				Name: "steven",
			},
		},
		want: []p{
			{
				Age:  8,
				Name: "steven",
			},
			{
				Age:  9,
				Name: "monika",
			},
			{
				Age:  10,
				Name: "bob",
			},
			{
				Age:  10,
				Name: "julia",
			},
			{
				Age:  10,
				Name: "tim",
			},
			{
				Age:  13,
				Name: "jack",
			},
		},
	}

	AscStableByKey(tests.in, func(p p) int {
		return p.Age
	})

	for i, age := range tests.in {
		testz.Equal(t, tests.want[i].Age, age.Age, tests.want[i].Age)
		testz.Equal(t, tests.want[i].Name, age.Name, tests.want[i].Name)
	}
}

func TestDescStableByKey(t *testing.T) {
	type p struct {
		Age  int
		Name string
	}
	tests := struct {
		in   []p
		want []p
	}{
		in: []p{
			{
				Age:  10,
				Name: "bob",
			},
			{
				Age:  13,
				Name: "jack",
			},
			{
				Age:  9,
				Name: "monika",
			},
			{
				Age:  10,
				Name: "julia",
			},
			{
				Age:  10,
				Name: "tim",
			},
			{
				Age:  8,
				Name: "steven",
			},
		},
		want: []p{
			{
				Age:  13,
				Name: "jack",
			},
			{
				Age:  10,
				Name: "bob",
			},
			{
				Age:  10,
				Name: "julia",
			},
			{
				Age:  10,
				Name: "tim",
			},
			{
				Age:  9,
				Name: "monika",
			},
			{
				Age:  8,
				Name: "steven",
			},
		},
	}

	DescStableByKey(tests.in, func(p p) int {
		return p.Age
	})

	for i, age := range tests.in {
		testz.Equal(t, tests.want[i].Age, age.Age, tests.want[i].Age)
		testz.Equal(t, tests.want[i].Name, age.Name, tests.want[i].Name)
	}
}

func TestAscByKey1(t *testing.T) {
	testsInt := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 3, 2}, []int{1, 2, 3}},
		{[]int{1, 2, 3, -1, 3, 9, 8, -2, 5}, []int{-2, -1, 1, 2, 3, 3, 5, 8, 9}},
	}

	for _, tt := range testsInt {
		AscByKey(tt.in, func(v int) int {
			return v + 1
		})
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
		AscByKey(tt.in, func(v float64) float64 {
			return v + 1
		})
		testz.Equal(t, tt.want, tt.in, tt.want)
	}
}

func TestDescByKey1(t *testing.T) {
	testsInt := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 3, 2}, []int{3, 2, 1}},
		{[]int{1, 2, 3, -1, 3, 9, 8, -2, 5}, []int{9, 8, 5, 3, 3, 2, 1, -1, -2}},
	}

	for _, tt := range testsInt {
		DescByKey(tt.in, func(v int) int {
			return 2 * v
		})
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
		DescByKey(tt.in, func(v float64) float64 {
			return 2 * v
		})
		testz.Equal(t, tt.want, tt.in, tt.want)
	}
}
