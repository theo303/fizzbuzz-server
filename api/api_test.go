package api

import (
	"bytes"
	"encoding/json"
	"fizzbuzz-server/api/fizzbuzzhandler"
	"fizzbuzz-server/api/mostfreqreqhandler"
	"fizzbuzz-server/config"
	"fizzbuzz-server/internal/fizzbuzz"
	"fizzbuzz-server/internal/stats"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// integration tests
func Test_Fizzbuzz(t *testing.T) {
	assertions := assert.New(t)

	params1 := fizzbuzz.Params{Int1: 3, Int2: 5, Limit: 16, Str1: "fizz", Str2: "buzz"}
	params2 := fizzbuzz.Params{Int1: 2, Int2: 7, Limit: 16, Str1: "fazz", Str2: "bozz"}

	api := Init(config.Conf{})

	// request 1: mostfreqreq without any previous requests
	gotCode, gotMostFreqReq, gotErr := getMostFreqReq(api)
	assertions.NoError(gotErr, "req 1 - error")
	assertions.Equal(http.StatusOK, gotCode, "req 1 - wrong code")
	assertions.Equal(stats.MostFrequentReq{Count: 0, Params: []fizzbuzz.Params{}}, gotMostFreqReq, "req 1 - wrong body")

	// request 2: fizzbuzz with correct params
	gotCode, gotFizzbuzz, gotErr := getFizzbuzz(api, params1)
	assertions.NoError(gotErr, "req 2 - error")
	assertions.Equal(http.StatusOK, gotCode, "req 2 - wrong code")
	assertions.Equal(
		[]string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16"},
		gotFizzbuzz,
		"req 2 - wrong body",
	)

	// request 3: fizzbuzz with other params
	gotCode, gotFizzbuzz, gotErr = getFizzbuzz(api, params2)
	assertions.NoError(gotErr, "req 3 - error")
	assertions.Equal(http.StatusOK, gotCode, "req 3 - wrong code")
	assertions.Equal(
		[]string{"1", "fazz", "3", "fazz", "5", "fazz", "bozz", "fazz", "9", "fazz", "11", "fazz", "13", "fazzbozz", "15", "fazz"},
		gotFizzbuzz,
		"req 3 - wrong body",
	)

	// request 4: mostfreqreq - get two params set
	gotCode, gotMostFreqReq, gotErr = getMostFreqReq(api)
	assertions.NoError(gotErr, "req 4 - error")
	assertions.Equal(http.StatusOK, gotCode, "req 4 - wrong code")
	// MostFreqReq struct contains a array so test each element
	assertions.Equal(1, gotMostFreqReq.Count, "req 4 - wrong count")
	assertions.ElementsMatch([]fizzbuzz.Params{params1, params2}, gotMostFreqReq.Params, "req 4 - wrong params")

	// request 5: fizzbuzz (same as 2)
	gotCode, gotFizzbuzz, gotErr = getFizzbuzz(api, params1)
	assertions.NoError(gotErr, "req 5 - error")
	assertions.Equal(http.StatusOK, gotCode, "req 5 - wrong code")
	assertions.Equal(
		[]string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16"},
		gotFizzbuzz,
		"req 5 - wrong body",
	)

	// request 6: mostfreqreq - get one params set
	gotCode, gotMostFreqReq, gotErr = getMostFreqReq(api)
	assertions.NoError(gotErr, "req 6 - error")
	assertions.Equal(http.StatusOK, gotCode, "req 6 - wrong code")
	assertions.Equal(stats.MostFrequentReq{Count: 2, Params: []fizzbuzz.Params{params1}}, gotMostFreqReq, "req 6 - wrong body")
}

func getMostFreqReq(api *Api) (int, stats.MostFrequentReq, error) {
	rr := httptest.NewRecorder()

	req, errReq := http.NewRequest("GET", "/mostfreqreq", nil)
	if errReq != nil {
		return 0, stats.MostFrequentReq{}, fmt.Errorf("creating request: %w", errReq)
	}

	http.HandlerFunc(api.handlerWithLogs(mostfreqreqhandler.ProcessMostFrequentReq)).ServeHTTP(rr, req)

	var response stats.MostFrequentReq
	body, errRead := ioutil.ReadAll(rr.Result().Body)
	if errRead != nil {
		return 0, stats.MostFrequentReq{}, fmt.Errorf("reading body: %w", errRead)
	}

	errJson := json.Unmarshal(body, &response)
	if errJson != nil {
		return 0, stats.MostFrequentReq{}, fmt.Errorf("decoding json: %w", errJson)
	}
	return rr.Code, response, nil
}

func getFizzbuzz(api *Api, params fizzbuzz.Params) (int, []string, error) {
	rr := httptest.NewRecorder()

	body, errJson := json.Marshal(params)
	if errJson != nil {
		return 0, []string{}, fmt.Errorf("encoding json: %w", errJson)
	}

	req, errReq := http.NewRequest("GET", "/fizzbuzz", ioutil.NopCloser(bytes.NewReader(body)))
	if errReq != nil {
		return 0, []string{}, fmt.Errorf("creating request: %w", errReq)
	}

	http.HandlerFunc(api.handlerWithLogs(fizzbuzzhandler.ProcessFizzbuzz)).ServeHTTP(rr, req)

	var response []string
	body, errRead := ioutil.ReadAll(rr.Result().Body)
	if errRead != nil {
		return 0, []string{}, fmt.Errorf("reading body: %w", errRead)
	}

	errJson = json.Unmarshal(body, &response)
	if errJson != nil {
		return 0, []string{}, fmt.Errorf("decoding json: %w", errJson)
	}
	return rr.Code, response, nil
}
