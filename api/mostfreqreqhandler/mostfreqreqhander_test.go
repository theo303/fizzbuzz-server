package mostfreqreqhandler

import (
	"fizzbuzz-server/internal/fizzbuzz"
	"fizzbuzz-server/internal/stats"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessMostFrequentReq(t *testing.T) {
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
			},
			counter: stats.FizzbuzzCounter{
				fizzbuzz.Params{Int1: 3, Int2: 4, Limit: 12, Str1: "fizz", Str2: "buzz"}: 1,
				fizzbuzz.Params{Int1: 1, Int2: 2, Limit: 12, Str1: "fizz", Str2: "buzz"}: 2,
			},
			wantCode:    http.StatusOK,
			wantHeaders: map[string][]string{},
			wantBody:    []byte(`{"count":2,"params":[{"int1":1,"int2":2,"limit":12,"str1":"fizz","str2":"buzz"}]}`),
		},
		"KO - method not allowed": {
			req: &http.Request{
				Method: "POST",
			},
			counter: stats.FizzbuzzCounter{
				fizzbuzz.Params{Int1: 3, Int2: 4, Limit: 12, Str1: "fizz", Str2: "buzz"}: 1,
				fizzbuzz.Params{Int1: 1, Int2: 2, Limit: 12, Str1: "fizz", Str2: "buzz"}: 2,
			},
			wantCode:    http.StatusMethodNotAllowed,
			wantHeaders: map[string][]string{"Allow": {"GET"}},
			wantBody:    []byte(`{"code":405,"desc":"method not allowed"}`),
			wantErrStr:  "invalid method",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := assert.New(t)

			gotCode, gotHeaders, gotBody, gotErr := ProcessMostFrequentReq(tt.req, tt.counter)

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
