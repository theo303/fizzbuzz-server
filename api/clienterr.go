package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type clientError struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

// In case of internal error, do not send the explicit error to the client
var internalError clientError = clientError{Code: http.StatusInternalServerError, Desc: "internal error"}

func (fErr clientError) getErrorBody() []byte {
	body, errJson := json.Marshal(fErr)
	if errJson != nil {
		log.Error().Err(errJson).Msg("error while creating error body")
		return []byte{}
	}
	return body
}
