package mostfreqreqhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"fizzbuzz-server/api/clienterr"
	"fizzbuzz-server/internal/stats"
)

// ProcessMostFrequentReq does all the process of a mostfreqreq request
func ProcessMostFrequentReq(r *http.Request, counter stats.FizzbuzzCounter) (int, map[string][]string, []byte, error) {
	// check method
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed,
			map[string][]string{"Allow": {"GET"}},
			clienterr.ClientError{Code: http.StatusMethodNotAllowed, Desc: "method not allowed"}.GetErrorBody(),
			errors.New("invalid method")
	}

	// retrieve most frequent request
	mostFreReq := counter.MostFrequentReq()

	// create response
	body, errJson := json.Marshal(mostFreReq)
	if errJson != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			clienterr.InternalError.GetErrorBody(),
			fmt.Errorf("error marshalling json: %w", errJson)
	}
	return http.StatusOK,
		map[string][]string{},
		body,
		nil
}
