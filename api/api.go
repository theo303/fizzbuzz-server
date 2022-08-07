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

	"github.com/rs/zerolog/log"
)

type formattedError struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

// In case of internal error, do not send the explicit error to the client
var internalError formattedError = formattedError{Code: http.StatusInternalServerError, Desc: "internal error"}

// Start starts the server
func Start(conf config.Conf) error {
	http.HandleFunc("/fizzbuzz", handlerFizzbuzz)

	log.Info().Int("port", conf.Port).Msg("starting server")
	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}

func handlerFizzbuzz(w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("address", r.RemoteAddr).Msg("received request")

	if r.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(getErrorBody(formattedError{Code: http.StatusMethodNotAllowed, Desc: "method not allowed"}))
		return
	}

	// read body
	reqBody, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Warn().Err(errRead).Msg("error while reading body")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getErrorBody(internalError))
		return
	}

	// retrieve and check params
	params, errParams := getParamsFizzbuzz(reqBody)
	if errParams != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getErrorBody(formattedError{Code: http.StatusBadRequest, Desc: errParams.Error()}))
		return
	}

	// execute fizzbuzz process
	output, errExec := fizzbuzz.ExecFizzbuzz(params)
	if errExec != nil {
		log.Error().Err(errExec).Msg("error while executing fizzbuzz process")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getErrorBody(internalError))
		return
	}

	// write response
	w.WriteHeader(http.StatusOK)
	body, errJson := json.Marshal(output)
	if errJson != nil {
		log.Error().Err(errJson).Str("output", strings.Join(output, ",")).Msg("error while marshalling output")

	}
	_, errWrite := w.Write(body)
	if errWrite != nil {
		log.Error().Err(errWrite).Str("body", string(body)).Msg("error while writing body")
	}
}

func getErrorBody(fErr formattedError) []byte {
	body, errJson := json.Marshal(fErr)
	if errJson != nil {
		log.Error().Err(errJson).Msg("error while creating error body")
		return []byte{}
	}
	return body
}

func getParamsFizzbuzz(body []byte) (fizzbuzz.Params, error) {
	params := fizzbuzz.Params{}
	errJson := json.Unmarshal(body, &params)
	if errJson != nil {
		log.Warn().Err(errJson).Str("body", string(body)).Msg("error while unmarshalling json")
		return fizzbuzz.Params{}, fmt.Errorf("invalid params")
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
		return params, errors.New(errStr)
	}

	return params, nil
}
