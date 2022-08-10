package fizzbuzzhandler

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"fizzbuzz-server/internal/fizzbuzz"
	"fizzbuzz-server/internal/stats"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessFizzbuzz(t *testing.T) {
	tests := map[string]struct {
		req         *http.Request
		counter     stats.FizzbuzzCounter
		wantCode    int
		wantHeaders map[string][]string
		wantBody    []byte
		wantErrStr  string
	}{
		"OK": {
			req: &http.Request{
				Method: "GET",
				Body: ioutil.NopCloser(strings.NewReader(`{
					"int1":3,
					"int2":4,
					"limit":12,
					"str1":"fizz",
					"str2":"buzz"
				}`)),
			},
			counter:     stats.NewFizzbuzzCounter(),
			wantCode:    http.StatusOK,
			wantHeaders: map[string][]string{},
			wantBody:    []byte(`["1","2","fizz","buzz","5","fizz","7","buzz","fizz","10","11","fizzbuzz"]`),
		},
		"KO - method not allowed": {
			req: &http.Request{
				Method: "POST",
				Body: ioutil.NopCloser(strings.NewReader(`{
					"int1":3,
					"int2":4,
					"limit":12,
					"str1":"fizz",
					"str2":"buzz"
				}`)),
			},
			counter:     stats.NewFizzbuzzCounter(),
			wantCode:    http.StatusMethodNotAllowed,
			wantHeaders: map[string][]string{"Allow": {"GET"}},
			wantBody:    []byte(`{"code":405,"desc":"method not allowed"}`),
			wantErrStr:  "invalid method",
		},
		"KO - invalid params": {
			req: &http.Request{
				Method: "GET",
				Body: ioutil.NopCloser(strings.NewReader(`{
					"int1":3,
					"limit":12,
					"str1":"fizz",
					"str2":"buzz"
				}`)),
			},
			counter:     stats.NewFizzbuzzCounter(),
			wantCode:    http.StatusBadRequest,
			wantHeaders: map[string][]string{},
			wantBody:    []byte(`{"code":400,"desc":"int2 missing (can't be zero)"}`),
			wantErrStr:  "invalid params",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)

			gotCode, gotHeaders, gotBody, gotErr := ProcessFizzbuzz(tt.req, tt.counter)

			if tt.wantErrStr != "" {
				assertions.Contains(gotErr.Error(), tt.wantErrStr)
			} else {
				assertions.NoError(gotErr)
			}
			assertions.Equal(tt.wantCode, gotCode)
			assertions.Equal(tt.wantHeaders, gotHeaders)
			assertions.Equal(tt.wantBody, gotBody)
		})
	}
}

func Test_getParamsFizzbuzz(t *testing.T) {
	tests := map[string]struct {
		body          []byte
		want          fizzbuzz.Params
		wantClientErr []string
		wantErr       []string
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
			body:          []byte(`{}`),
			wantClientErr: []string{"int1 missing", "int2 missing", "limit missing"},
			wantErr:       []string{"int1 missing", "int2 missing", "limit missing"},
		},
		"KO - limit negative": {
			body: []byte(`{
				"int1":3,
				"limit":-16
			}`),
			wantClientErr: []string{"int2 missing", "limit must be superior to one"},
			wantErr:       []string{"int2 missing", "limit must be superior to one"},
		},
		"KO - invalid JSON": {
			body:          []byte(`aaa`),
			wantClientErr: []string{"invalid params"},
			wantErr:       []string{"invalid character"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)
			got, errClientGot, errGot := getParamsFizzbuzz(tt.body)
			if len(tt.wantErr) != 0 {
				// test internal error
				for _, errStr := range tt.wantErr {
					assertions.Contains(errGot.Error(), errStr, "error not found")
				}
				// test client error
				for _, errClientStr := range tt.wantClientErr {
					assertions.Contains(errClientGot.Desc, errClientStr, "error not found")
				}
			} else {
				assertions.NoError(errGot, "unexpected error")
				assertions.Equal(tt.want, got)
			}
		})
	}
}
