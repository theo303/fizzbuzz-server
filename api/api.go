package api

import (
	"fmt"
	"net/http"

	"fizzbuzz-server/api/fizzbuzzhandler"
	"fizzbuzz-server/api/mostfreqreqhandler"
	"fizzbuzz-server/config"
	"fizzbuzz-server/internal/stats"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Api represents the API of the fizzbuzz server
type Api struct {
	*http.Server
	counter stats.FizzbuzzCounter
}

// ProcessFunc is a template func that can be wrapped with 'handlerWithLogs'
type ProcessFunc func(*http.Request, stats.FizzbuzzCounter) (
	statusCode int,
	headers map[string][]string,
	body []byte,
	err error,
)

// Init initialize API server with a new counter
func Init(conf config.Conf) *Api {
	return &Api{
		Server:  &http.Server{Addr: fmt.Sprintf(":%d", conf.Port)},
		counter: stats.NewFizzbuzzCounter(),
	}
}

// Run starts the server
func (a *Api) Run() error {
	http.HandleFunc("/fizzbuzz", a.handlerWithLogs(fizzbuzzhandler.ProcessFizzbuzz))
	http.HandleFunc("/mostfreqreq", a.handlerWithLogs(mostfreqreqhandler.ProcessMostFrequentReq))

	log.Info().Str("addr", a.Addr).Msg("starting server")
	return a.ListenAndServe()
}

func (a *Api) handlerWithLogs(f ProcessFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New()
		log.Info().
			Str("address", r.RemoteAddr).
			Str("requestID", reqID.String()).
			Msg("received request")

		code, headersMap, body, errProcess := f(r, a.counter)
		if errProcess != nil {
			log.Warn().
				Err(errProcess).
				Str("requestID", reqID.String()).
				Msg("error while processing request")
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
}
