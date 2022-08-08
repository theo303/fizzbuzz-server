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
	"fizzbuzz-server/internal/stats"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Api represents the API of the fizzbuzz server
type Api struct {
	counter stats.FizzbuzzCounter
}

// Init initialize API with a new counter
func Init() *Api {
	return &Api{counter: stats.NewFizzbuzzCounter()}
}

// Run starts the server
func (a *Api) Run(conf config.Conf) error {
	http.HandleFunc("/fizzbuzz", a.handlerFizzbuzz)
	http.HandleFunc("/mostfreqreq", a.handlerMostFrequentReq)

	log.Info().Int("port", conf.Port).Msg("starting server")
	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}

func (a *Api) handlerFizzbuzz(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New()
	log.Info().
		Str("address", r.RemoteAddr).
		Str("requestID", reqID.String()).
		Msg("received request")

	code, headersMap, body, errProcess := a.processFizzbuzz(r)
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
	log.Info().
		Str("body", string(body)).
		Str("requestID", reqID.String()).
		Msg("sending response")
}

func (a *Api) processFizzbuzz(r *http.Request) (code int, headers map[string][]string, body []byte, err error) {
	// check method
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed,
			map[string][]string{"Allow": {"GET"}},
			clientError{Code: http.StatusMethodNotAllowed, Desc: "method not allowed"}.getErrorBody(),
			errors.New("invalid method")
	}

	// read body
	reqBody, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			internalError.getErrorBody(),
			fmt.Errorf("error reading body: %w", errRead)
	}

	// retrieve and check params
	params, clientErr, errParams := getParamsFizzbuzz(reqBody)
	if errParams != nil {
		return http.StatusBadRequest,
			map[string][]string{},
			clientErr.getErrorBody(),
			fmt.Errorf("invalid params: %w", errParams)
	}

	// increment counter
	a.counter.Inc(params)

	// execute fizzbuzz
	output, errExec := fizzbuzz.ExecFizzbuzz(params)
	if errExec != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			internalError.getErrorBody(),
			fmt.Errorf("error executing fizzbuzz: %w", errExec)
	}

	// create response
	body, errJson := json.Marshal(output)
	if errJson != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			internalError.getErrorBody(),
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

func (a *Api) handlerMostFrequentReq(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New()
	log.Info().
		Str("address", r.RemoteAddr).
		Str("requestID", reqID.String()).
		Msg("received request")

	code, headersMap, body, errProcess := a.processMostFrequentReq(r)
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
	log.Info().
		Str("body", string(body)).
		Str("requestID", reqID.String()).
		Msg("sending response")
}

func (a *Api) processMostFrequentReq(r *http.Request) (code int, headers map[string][]string, body []byte, err error) {
	// check method
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed,
			map[string][]string{"Allow": {"GET"}},
			clientError{Code: http.StatusMethodNotAllowed, Desc: "method not allowed"}.getErrorBody(),
			errors.New("invalid method")
	}

	// retrieve most frequent request
	mostFreReq := a.counter.MostFrequentReq()

	// create response
	body, errJson := json.Marshal(mostFreReq)
	if errJson != nil {
		return http.StatusInternalServerError,
			map[string][]string{},
			internalError.getErrorBody(),
			fmt.Errorf("error marshalling json: %w", errJson)
	}
	return http.StatusOK,
		map[string][]string{},
		body,
		nil
}
