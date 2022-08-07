package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"fizzbuzz-server/config"
	"fizzbuzz-server/internal/fizzbuzz"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type clientError struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

// In case of internal error, do not send the explicit error to the client
var internalError clientError = clientError{Code: http.StatusInternalServerError, Desc: "internal error"}

func getErrorBody(fErr clientError) []byte {
	body, errJson := json.Marshal(fErr)
	if errJson != nil {
		log.Error().Err(errJson).Msg("error while creating error body")
		return []byte{}
	}
	return body
}

// Start starts the server
func Start(conf config.Conf) error {
	http.HandleFunc("/fizzbuzz", handlerFizzbuzz)

	log.Info().Int("port", conf.Port).Msg("starting server")
	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}

func handlerFizzbuzz(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New()
	log.Info().
		Str("address", r.RemoteAddr).
		Str("requestID", reqID.String()).
		Msg("received request")

	code, headersMap, body, errProcess := processFizzbuzz(r)
	if errProcess != nil {
		log.Warn().Err(errProcess).Str("requestID", reqID.String()).Msg("error while processing request")
	}
	for headerKey, headers := range headersMap {
		for _, header := range headers {
			w.Header().Add(headerKey, header)
		}
	}
	w.WriteHeader(code)
	w.Write(body)
}

func processFizzbuzz(r *http.Request) (code int, headers map[string][]string, body []byte, err error) {
	// check method
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed,
			map[string][]string{"Allow": {"GET"}},
			getErrorBody(clientError{Code: http.StatusMethodNotAllowed, Desc: "method not allowed"}),
			errors.New("invalid method")
	}

	// read body
	reqBody, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			getErrorBody(internalError),
			fmt.Errorf("error reading body: %w", errRead)
	}

	// retrieve and check params
	params, clientErr, errParams := getParamsFizzbuzz(reqBody)
	if errParams != nil {
		return http.StatusBadRequest,
			map[string][]string{},
			getErrorBody(clientErr),
			fmt.Errorf("invalid params: %w", errParams)
	}

	// execute fizzbuzz
	output, errExec := fizzbuzz.ExecFizzbuzz(params)
	if errExec != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			getErrorBody(internalError),
			fmt.Errorf("error executing fizzbuzz: %w", errExec)
	}

	// write response
	body, errJson := json.Marshal(output)
	if errJson != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			getErrorBody(internalError),
			fmt.Errorf("error marshalling json: %w", errJson)
	}
	return http.StatusOK,
		map[string][]string{},
		body,
		nil
}

// getParamsFizzbuzz retrieves and checks params from the body
// it returns two versions of the error if needed, one for the client and one more precise for internal use
func getParamsFizzbuzz(body []byte) (fizzbuzz.Params, clientError, error) {
	params := fizzbuzz.Params{}
	errJson := json.Unmarshal(body, &params)
	if errJson != nil {
		return fizzbuzz.Params{},
			clientError{Code: http.StatusBadRequest, Desc: "invalid params"},
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
		strBuilder.WriteString("limit missing (can't inferior to one), ")
	} else if params.Limit < 0 {
		strBuilder.WriteString("limit must be superior to one, ")
	}
	if errStr := strBuilder.String(); errStr != "" {
		// remove trailing comma and space
		errStr = errStr[:len(errStr)-2]
		return params, clientError{Code: http.StatusBadRequest, Desc: errStr}, errors.New(errStr)
	}

	return params, clientError{}, nil
}
