package fizzbuzzhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"fizzbuzz-server/api/clienterr"
	"fizzbuzz-server/internal/fizzbuzz"
	"fizzbuzz-server/internal/stats"
)

// ProcessFizzbuzz does all the process of a fizzbuzz request
func ProcessFizzbuzz(r *http.Request, counter stats.FizzbuzzCounter) (code int, headers map[string][]string, body []byte, err error) {
	// check method
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed,
			map[string][]string{"Allow": {"GET"}},
			clienterr.ClientError{Code: http.StatusMethodNotAllowed, Desc: "method not allowed"}.GetErrorBody(),
			errors.New("invalid method")
	}

	// read body
	reqBody, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			clienterr.InternalError.GetErrorBody(),
			fmt.Errorf("error reading body: %w", errRead)
	}

	// retrieve and check params
	params, clientErr, errParams := getParamsFizzbuzz(reqBody)
	if errParams != nil {
		return http.StatusBadRequest,
			map[string][]string{},
			clientErr.GetErrorBody(),
			fmt.Errorf("invalid params: %w", errParams)
	}

	// increment counter
	counter.Inc(params)

	// execute fizzbuzz
	output, errExec := fizzbuzz.ExecFizzbuzz(params)
	if errExec != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			clienterr.InternalError.GetErrorBody(),
			fmt.Errorf("error executing fizzbuzz: %w", errExec)
	}

	// create response
	body, errJson := json.Marshal(output)
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

// getParamsFizzbuzz retrieves and checks params from the body
// it returns two versions of the error if needed, one for the client and one more precise for internal use
func getParamsFizzbuzz(body []byte) (fizzbuzz.Params, clienterr.ClientError, error) {
	params := fizzbuzz.Params{}
	errJson := json.Unmarshal(body, &params)
	if errJson != nil {
		return fizzbuzz.Params{},
			clienterr.ClientError{Code: http.StatusBadRequest, Desc: "invalid params"},
			fmt.Errorf("unmarshalling json: %w", errJson)
	}

	strBuilder := strings.Builder{}
	if params.Int1 == 0 {
		strBuilder.WriteString("int1 missing (can't be zero), ")
	}
	if params.Int2 == 0 {
		strBuilder.WriteString("int2 missing (can't be zero), ")
	}
	if params.Limit == 0 {
		strBuilder.WriteString("limit missing (can't be inferior to one), ")
	} else if params.Limit < 0 {
		strBuilder.WriteString("limit must be superior to one, ")
	}
	if errStr := strBuilder.String(); errStr != "" {
		// remove trailing comma and space
		errStr = errStr[:len(errStr)-2]
		return params, clienterr.ClientError{Code: http.StatusBadRequest, Desc: errStr}, errors.New(errStr)
	}

	return params, clienterr.ClientError{}, nil
}
