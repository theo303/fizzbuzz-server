package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fizzbuzz-server/internal/fizzbuzz"

	"github.com/stretchr/testify/assert"
)

func Test_getErrorBody(t *testing.T) {
	tests := map[string]struct {
		err  clientError
		want []byte
	}{
		"OK": {
			err: clientError{
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

// functional test
func Test_FizzbuzzEndpoint(t *testing.T) {
	tests := map[string]struct {
		method   string
		body     string
		wantBody []byte
		wantCode int
	}{
		"OK": {
			method: "GET",
			body: `{
				"int1":3,
				"int2":5,
				"limit":16,
				"str1":"fazz",
				"str2":"buzz"
			}`,
			wantBody: []byte(`["1","2","fazz","4","buzz","fazz","7","8","fazz","buzz","11","fazz","13","14","fazzbuzz","16"]`),
			wantCode: http.StatusOK,
		},
		"KO - method not allowed": {
			method: "POST",
			body: `{
				"int1":3,
				"int2":5,
				"limit":16,
				"str1":"fazz",
				"str2":"buzz"
			}`,
			wantBody: []byte(`{"code":405,"desc":"method not allowed"}`),
			wantCode: http.StatusMethodNotAllowed,
		},
		"KO - invalid params": {
			method: "GET",
			body: `{
				"int1":0,
				"int2":5,
				"limit":16,
				"str1":"fazz",
				"str2":"buzz"
			}`,
			wantBody: []byte(`{"code":400,"desc":"int1 missing (can't be zero)"}`),
			wantCode: http.StatusBadRequest,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "localhost:1234/fizzbuzz", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			handlerFizzbuzz(w, req)

			resp := w.Result()
			assertions := assert.New(t)

			gotBody, errRead := ioutil.ReadAll(resp.Body)
			if errRead != nil {
				panic(errRead)
			}
			assertions.Equal(tt.wantCode, resp.StatusCode)
			assertions.Equal([]byte(tt.wantBody), gotBody)
		})
	}
}
