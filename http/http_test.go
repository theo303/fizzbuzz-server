package http

import (
	"fizzbuzz-server/fizzbuzz"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getErrorBody(t *testing.T) {
	tests := map[string]struct {
		err  formattedError
		want []byte
	}{
		"OK": {
			err: formattedError{
				Code: http.StatusForbidden,
				Desc: "test error",
			},
			want: []byte(`{"code":403,"desc":"test error"}`),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, getErrorBody(tt.err))
		})
	}
}

func Test_getParamsFizzbuzz(t *testing.T) {
	tests := map[string]struct {
		body    []byte
		want    fizzbuzz.Params
		wantErr []string
	}{
		"OK": {
			body: []byte(`{
				"int1":3,
				"int2":5,
				"limit":16,
				"str1":"fizz"
			}`),
			want: fizzbuzz.Params{
				Int1:  3,
				Int2:  5,
				Limit: 16,
				Str1:  "fizz",
				Str2:  "",
			},
		},
		"KO - missing integers": {
			body:    []byte(`{}`),
			wantErr: []string{"int1 missing", "int2 missing", "limit missing"},
		},
		"KO - limit negative": {
			body: []byte(`{
				"int1":3,
				"limit":-16
			}`),
			wantErr: []string{"int2 missing", "limit must be superior to one"},
		},
		"KO - invalid JSON": {
			body:    []byte(`aaa`),
			wantErr: []string{"invalid params"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)
			got, errGot := getParamsFizzbuzz(tt.body)
			if len(tt.wantErr) != 0 {
				for _, errStr := range tt.wantErr {
					assertions.Contains(errGot.Error(), errStr, "error not found")
				}
			} else {
				assertions.NoError(errGot, "unexpected error")
				assertions.Equal(tt.want, got)
			}
		})
	}
}
