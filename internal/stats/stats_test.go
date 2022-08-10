package stats

import (
	"fizzbuzz-server/internal/fizzbuzz"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FizzbuzzCounter_Inc(t *testing.T) {
	tests := map[string]struct {
		fbc    FizzbuzzCounter
		params fizzbuzz.Params
		want   int
	}{
		"new params": {
			fbc: FizzbuzzCounter{
				{
					Int1:  3,
					Int2:  5,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 1,
			},
			params: fizzbuzz.Params{
				Int1:  2,
				Int2:  5,
				Limit: 16,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: 1,
		},
		"already present in map": {
			fbc: FizzbuzzCounter{
				{
					Int1:  3,
					Int2:  5,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 1,
			},
			params: fizzbuzz.Params{
				Int1:  3,
				Int2:  5,
				Limit: 16,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: 2,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)

			tt.fbc.Inc(tt.params)
			got, ok := tt.fbc[tt.params]
			assertions.True(ok, "key not found")
			assertions.Equal(tt.want, got, "wrong value")
		})
	}
}

func Test_FizzbuzzCounter_Get(t *testing.T) {
	tests := map[string]struct {
		fbc    FizzbuzzCounter
		params fizzbuzz.Params
		want   int
	}{
		"not present in map": {
			fbc: map[fizzbuzz.Params]int{
				{
					Int1:  3,
					Int2:  5,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 1,
			},
			params: fizzbuzz.Params{
				Int1:  2,
				Int2:  5,
				Limit: 16,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: 0,
		},
		"present in map": {
			fbc: map[fizzbuzz.Params]int{
				{
					Int1:  3,
					Int2:  5,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 1,
			},
			params: fizzbuzz.Params{
				Int1:  3,
				Int2:  5,
				Limit: 16,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: 1,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)

			got := tt.fbc.Get(tt.params)
			assertions.Equal(tt.want, got, "wrong value")
		})
	}
}

func Test_FizzbuzzCounter_MostFrequentRequest(t *testing.T) {
	tests := map[string]struct {
		fbc  FizzbuzzCounter
		want MostFrequentReq
	}{
		"one params": {
			fbc: map[fizzbuzz.Params]int{
				{
					Int1:  3,
					Int2:  5,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 2,
				{
					Int1:  3,
					Int2:  6,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 1,
			},
			want: MostFrequentReq{
				Count: 2,
				Params: []fizzbuzz.Params{
					{
						Int1:  3,
						Int2:  5,
						Limit: 16,
						Str1:  "fizz",
						Str2:  "buzz",
					},
				},
			},
		},
		"two params": {
			fbc: map[fizzbuzz.Params]int{
				{
					Int1:  3,
					Int2:  5,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 2,
				{
					Int1:  3,
					Int2:  6,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 2,
				{
					Int1:  3,
					Int2:  7,
					Limit: 16,
					Str1:  "fizz",
					Str2:  "buzz",
				}: 1,
			},
			want: MostFrequentReq{
				Count: 2,
				Params: []fizzbuzz.Params{
					{
						Int1:  3,
						Int2:  5,
						Limit: 16,
						Str1:  "fizz",
						Str2:  "buzz",
					},
					{
						Int1:  3,
						Int2:  6,
						Limit: 16,
						Str1:  "fizz",
						Str2:  "buzz",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)

			got := tt.fbc.MostFrequentReq()
			assertions.Equal(tt.want.Count, got.Count, "count different")
			assertions.ElementsMatch(tt.want.Params, got.Params, "params different")
		})
	}
}
