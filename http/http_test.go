package http

import (
	"fizzbuzz-server/fizzbuzz"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getParamsFizzbuzz(t *testing.T) {
	tests := map[string]struct {
		values     url.Values
		want       fizzbuzz.Params
		wantErrStr []string
	}{
		"OK": {
			values: url.Values{
				"int1":  []string{"1"},
				"int2":  []string{"2"},
				"limit": []string{"3"},
				"str1":  []string{"str1"},
				"str2":  []string{"str2"},
			},
			want: fizzbuzz.Params{
				Int1:  1,
				Int2:  2,
				Limit: 3,
				Str1:  "str1",
				Str2:  "str2",
			},
		},
		"KO - missing all parameters": {
			values:     url.Values{},
			wantErrStr: []string{"int1 missing", "int2 missing", "limit missing", "str1 missing", "str2 missing"},
		},
		"KO - invalid integers": {
			values: url.Values{
				"int1":  []string{"a"},
				"int2":  []string{"2.7"},
				"limit": []string{"9,8"},
				"str1":  []string{"str1"},
				"str2":  []string{"str2"},
			},
			wantErrStr: []string{"int1: could not convert to int", "int2: could not convert to int", "limit: could not convert to int"},
		},
		"KO - invalid limit": {
			values: url.Values{
				"int1":  []string{"1"},
				"int2":  []string{"2"},
				"limit": []string{"0"},
				"str1":  []string{"str1"},
				"str2":  []string{"str2"},
			},
			wantErrStr: []string{"limit: invalid value (must be superior or equal to 1)"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)

			got, errGot := getParamsFizzbuzz(tt.values)
			if len(tt.wantErrStr) != 0 {
				for i, wantStr := range tt.wantErrStr {
					assertions.Contains(errGot.Error(), wantStr, fmt.Sprintf("invalid error %d", i))
				}
			} else {
				assertions.NoError(errGot, "unexpected error")
				assertions.Equal(tt.want, got, "params not equal")
			}
		})
	}
}
