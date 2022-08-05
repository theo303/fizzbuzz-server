package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fizzbuzz-server/config"
	"fizzbuzz-server/fizzbuzz"
)

// Start starts the server
func Start(conf config.Conf) error {
	http.HandleFunc("/fizzbuzz", handlerFizzbuzz)

	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}

func handlerFizzbuzz(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("405 - Method %s not allowed for the route %s", r.Method, r.URL.Path)))
		return
	}

	params, errParams := getParamsFizzbuzz(r.URL.Query())
	if errParams != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Bad request, %s", errParams.Error())))
	}
	fmt.Printf("%+v\n", params)
}

// getParamsFizzbuzz verify that values contains all the mandatory parameters,
// that their types are correct and then returns them
func getParamsFizzbuzz(values url.Values) (fizzbuzz.Params, error) {
	errStrBuilder := strings.Builder{}
	params := fizzbuzz.Params{}
	var errConv error

	if !values.Has("int1") {
		errStrBuilder.WriteString("int1 missing, ")
	} else if params.Int1, errConv = strconv.Atoi(values.Get("int1")); errConv != nil {
		errStrBuilder.WriteString(fmt.Sprintf("int1: could not convert to int: %s, ", errConv.Error()))
	}
	if !values.Has("int2") {
		errStrBuilder.WriteString("int2 missing, ")
	} else if params.Int2, errConv = strconv.Atoi(values.Get("int2")); errConv != nil {
		errStrBuilder.WriteString(fmt.Sprintf("int2: could not convert to int: %s, ", errConv.Error()))
	}
	if !values.Has("limit") {
		errStrBuilder.WriteString("limit missing, ")
	} else if params.Limit, errConv = strconv.Atoi(values.Get("limit")); errConv != nil {
		errStrBuilder.WriteString(fmt.Sprintf("limit: could not convert to int: %s, ", errConv.Error()))
	}

	if !values.Has("str1") {
		errStrBuilder.WriteString("str1 missing, ")
	}
	params.Str1 = values.Get("str1")
	if !values.Has("str2") {
		errStrBuilder.WriteString("str2 missing, ")
	}
	params.Str2 = values.Get("str2")

	if errStr := errStrBuilder.String(); errStr != "" {
		// remove trailing comma and space
		errStr = errStr[:len(errStr)-2]
		return params, fmt.Errorf("invalid parameters: %s", errStr)
	}

	return params, nil
}
