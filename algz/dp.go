package algz

import "math"

type DpSolvers[T any] map[int][]T

// Best returns the best solution for the given maximum value.
func (s DpSolvers[T]) Best(maxValue int) []T {
	best, ok := s[maxValue]
	if ok {
		return best
	}

	minDiff := math.MaxInt
	for value, solver := range s {
		if diff := maxValue - value; diff >= 0 && diff < minDiff {
			best = solver
			minDiff = diff
		}
	}
	return best
}

// BestAllowMinOverflow returns the best solution for the given maximum value, but allows a minimum overflow.
func (s DpSolvers[T]) BestAllowMinOverflow(maxValue int) []T {
	best, ok := s[maxValue]
	if ok {
		return best
	}

	minDiff := math.MaxInt
	for value, solver := range s {
		diff := maxValue - value
		if diff < 0 {
			if minDiff > 0 || diff > minDiff {
				best = solver
				minDiff = diff
			}

			continue
		}

		if diff < minDiff {
			best = solver
			minDiff = diff
		}
	}

	return best
}

// FindDpSolvers finds all possible solutions to the 0-1 knapsack problem.
func FindDpSolvers[T any](maxValue int, items []T, valueFunc func(T) int, allowOverOnce bool,
	tieBreaker ...func(old []T, new []T) (replace bool)) DpSolvers[T] {
	var breaker func([]T, []T) bool
	if len(tieBreaker) > 0 {
		breaker = tieBreaker[0]
	}

	dp := map[int][]T{
		0: {},
	}

	var overflow int
	dpTmp := make(map[int][]T, 6)
	var tmpPool slicesPool[T]
	for _, item := range items {
		value := valueFunc(item)
		for currentValue, solver := range dp {
			newValue := currentValue + value
			if newValue > maxValue {
				if !allowOverOnce || (overflow > 0 && newValue > overflow) {
					continue
				}
				overflow = newValue
			}

			oldSolver, ok := dp[newValue]
			if ok && breaker == nil {
				continue
			}

			newSolver := tmpPool.Get(len(solver) + 1)
			newSolver = append(newSolver, solver...)
			newSolver = append(newSolver, item)
			if ok && !breaker(oldSolver, newSolver) {
				tmpPool.Put(newSolver)
				continue
			}

			dpTmp[newValue] = newSolver
		}

		for v, solver := range dpTmp {
			if old, ok := dp[v]; ok {
				tmpPool.Put(old)
			}

			dp[v] = solver
			delete(dpTmp, v)
		}
	}

	return dp
}

type knapsack[T any] struct {
	score int
	items []T
}

// Knapsack solves the 0-1 knapsack problem.
// maxWeight is the maximum weight the knapsack can carry.
// items is a slice of items to choose from.
// weightFunc is a function that returns the weight of an item.
// valueFunc is a function that returns the value of an item.
// tieBreaker is an optional function that is used to break ties when multiple solutions have the same value.
func Knapsack[T any](maxWeight int, items []T, weightFunc, valueFunc func(T) int,
	tieBreaker ...func(old []T, new []T) (replace bool)) []T {
	var breaker func([]T, []T) bool
	if len(tieBreaker) > 0 {
		breaker = tieBreaker[0]
	}

	dp := make([]knapsack[T], maxWeight+1)
	var tmp []T
	for _, item := range items {
		w := weightFunc(item)
		value := valueFunc(item)
		for i := maxWeight; i >= w; i-- {
			newScore := dp[i-w].score + value
			if newScore > dp[i].score {
				tmp = append(tmp[:0], dp[i-w].items...)
				tmp = append(tmp, item)
				dp[i].items = append(dp[i].items[:0], tmp...)
				dp[i].score = newScore
			} else if newScore == dp[i].score && breaker != nil {
				tmp = append(tmp[:0], dp[i-w].items...)
				tmp = append(tmp, item)
				if breaker(dp[i].items, tmp) {
					dp[i].items = append(dp[i].items[:0], tmp...)
					dp[i].score = newScore
				}
			}
		}
	}

	return dp[maxWeight].items
}

type slicesPool[T any] struct {
	entries [][]T
}

func (p *slicesPool[T]) Get(initCap int) []T {
	if len(p.entries) == 0 {
		return make([]T, 0, initCap)
	}

	last := len(p.entries) - 1
	slice := p.entries[last]
	p.entries = p.entries[:last]
	return slice[:0]
}

func (p *slicesPool[T]) Put(slice []T) {
	p.entries = append(p.entries, slice)
}
