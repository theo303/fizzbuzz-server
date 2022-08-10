package stats

import "fizzbuzz-server/internal/fizzbuzz"

// FizzbuzzCounter keeps count of the number of request for a set of parameters
type FizzbuzzCounter map[fizzbuzz.Params]int

type MostFrequentReq struct {
	Count  int               `json:"count"`
	Params []fizzbuzz.Params `json:"params"`
}

func NewFizzbuzzCounter() FizzbuzzCounter {
	return make(map[fizzbuzz.Params]int)
}

// Inc increments the counter for these parameters
func (fbc FizzbuzzCounter) Inc(params fizzbuzz.Params) {
	_, found := fbc[params]
	if !found {
		fbc[params] = 0
	}
	fbc[params]++
}

// Get retrieve the numbers of request received for these parameters
func (fbc FizzbuzzCounter) Get(params fizzbuzz.Params) int {
	_, found := fbc[params]
	if !found {
		return 0
	}
	return fbc[params]
}

// MostFrequentReq retrieves the number and the parameters (one or multiple) of the most frequent request
func (fbc FizzbuzzCounter) MostFrequentReq() MostFrequentReq {
	max := 0
	maxParams := []fizzbuzz.Params{}
	for params, count := range fbc {
		if count > max {
			max = count
			maxParams = []fizzbuzz.Params{}
			maxParams = append(maxParams, params)
		} else if count == max {
			maxParams = append(maxParams, params)
		}
	}
	return MostFrequentReq{Count: max, Params: maxParams}
}
