package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_clientError_getErrorBody(t *testing.T) {
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
			assert.Equal(t, tt.want, tt.err.getErrorBody())
		})
	}
}
