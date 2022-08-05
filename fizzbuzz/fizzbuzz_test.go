package fizzbuzz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecFizzbuzz(t *testing.T) {
	tests := map[string]struct {
		params  Params
		want    []string
		wantErr bool
	}{
		"KO - invalid limit": {
			params: Params{
				Int1:  3,
				Int2:  5,
				Limit: 0,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want:    []string{},
			wantErr: true,
		},
		"OK": {
			params: Params{
				Int1:  3,
				Int2:  5,
				Limit: 16,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: []string{
				"1", "2", "fizz", "4", "buzz", "fizz", "7", "8",
				"fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)

			got, gotErr := ExecFizzbuzz(tt.params)
			if tt.wantErr {
				assertions.Error(gotErr)
			} else {
				assertions.NoError(gotErr)
			}
			assertions.Equal(got, tt.want)
		})
	}
}
