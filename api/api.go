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

// Start starts the server
func Start(conf config.Conf) error {
	http.HandleFunc("/fizzbuzz", handlerFizzbuzz)

	log.Info().Int("port", conf.Port).Msg("starting server")
	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}

func handlerFizzbuzz(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New()
	log.Debug().
		Str("address", r.RemoteAddr).
		Str("requestID", reqID.String()).
		Msg("received request")

	if r.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(getErrorBody(clientError{Code: http.StatusMethodNotAllowed, Desc: "method not allowed"}))
		return
	}

	// read body
	reqBody, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Warn().Err(errRead).
			Str("requestID", reqID.String()).
			Msg("error while reading body")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getErrorBody(internalError))
		return
	}

	// retrieve and check params
	params, clientErr, errParams := getParamsFizzbuzz(reqBody)
	if errParams != nil {
		log.Warn().Err(errParams).
			Str("requestID", reqID.String()).
			Str("body", string(reqBody)).
			Msg("invalid params")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getErrorBody(clientErr))
		return
	}

	// execute fizzbuzz process
	output, errExec := fizzbuzz.ExecFizzbuzz(params)
	if errExec != nil {
		log.Error().
			Err(errExec).
			Str("requestID", reqID.String()).
			Msg("error while executing fizzbuzz process")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getErrorBody(internalError))
		return
	}

	// write response
	w.WriteHeader(http.StatusOK)
	body, errJson := json.Marshal(output)
	if errJson != nil {
		log.Error().Err(errJson).
			Str("requestID", reqID.String()).
			Str("output", strings.Join(output, ",")).
			Msg("error while marshalling output")

	}
	_, errWrite := w.Write(body)
	if errWrite != nil {
		log.Error().Err(errWrite).
			Str("requestID", reqID.String()).
			Str("body", string(body)).
			Msg("error while writing body")
	}
}

func getErrorBody(fErr clientError) []byte {
	body, errJson := json.Marshal(fErr)
	if errJson != nil {
		log.Error().Err(errJson).Msg("error while creating error body")
		return []byte{}
	}
	return body
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
