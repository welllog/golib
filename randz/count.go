package randz

import (
	"sort"

	"github.com/welllog/golib/hashz"
)

// CountGenerator is a generator for count
type CountGenerator struct {
	rules []rule
}

type rule struct {
	period           int // time period
	periodEndMaxIncr int // time period end max increase
	interval         int // current time interval
	intervalMaxIncr  int // max increase value per interval
	quickNum         int
	// like: last hour, max increase(1~2) per 10 seconds, reach 1 hour increase(1~100)
	// rule{period: 3600, periodEndMaxIncr: 100, interval: 10, , intervalMaxIncr: 2}
}

// AddRule add count increase rule
func (r *CountGenerator) AddRule(period, periodEndMaxIncr, interval, intervalMaxIncr int) {
	r.rules = append(r.rules, rule{
		period:           period,
		periodEndMaxIncr: periodEndMaxIncr,
		interval:         interval,
		intervalMaxIncr:  intervalMaxIncr,
	})
	sort.Slice(r.rules, func(i, j int) bool {
		return r.rules[i].period < r.rules[j].period
	})
}

// Generate generate count
func (r *CountGenerator) Generate(id string, diff int) int {
	if diff <= 0 {
		return 0
	}

	hn := hashz.BKDRHash(id)
	var count, lastPeriod int
	for _, v := range r.rules {
		multi := r.getRand(hn, v.intervalMaxIncr)
		if diff < v.period {
			return (diff-lastPeriod)/v.interval*multi + count
		}
		count += (v.period-lastPeriod)/v.interval*multi + r.getRand(hn, v.periodEndMaxIncr)
		lastPeriod = v.period
	}
	return count
}

// Max return max count
func (r *CountGenerator) Max(diff int) int {
	if diff <= 0 {
		return 0
	}

	var count, lastGradient int
	for _, v := range r.rules {
		if diff < v.period {
			return (diff-lastGradient)/v.interval*v.intervalMaxIncr + count
		}
		count += (v.period-lastGradient)/v.interval*v.intervalMaxIncr + v.periodEndMaxIncr
		lastGradient = v.period
	}
	return count
}

// Min return min count
func (r *CountGenerator) Min(diff int) int {
	if diff <= 0 {
		return 0
	}

	var count, lastGradient int
	for _, v := range r.rules {
		if diff < v.period {
			return (diff-lastGradient)/v.interval + count
		}
		count += (v.period-lastGradient)/v.interval + 1
		lastGradient = v.period
	}
	return count
}

func (r *CountGenerator) getRand(n uint32, max int) int {
	if max == 0 {
		return 0
	}
	return int(n%uint32(max) + 1)
}
