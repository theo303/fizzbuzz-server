package fizzbuzz

import (
	"errors"
	"strconv"
)

// Params are the required parameters for the fizzbuzz process
type Params struct {
	Int1  int
	Int2  int
	Limit int
	Str1  string
	Str2  string
}

// ExecFizzbuzz starts the fizzbuzz process
func ExecFizzbuzz(params Params) ([]string, error) {
	var output []string
	if params.Int1 == 0 || params.Int2 == 0 || params.Limit < 1 {
		return []string{}, errors.New("invalid params")
	}
	for i := 1; i <= params.Limit; i++ {
		str := ""
		empty := true
		if i%params.Int1 == 0 {
			str = params.Str1
			empty = false
		}
		if i%params.Int2 == 0 {
			str = str + params.Str2
			empty = false
		}
		if empty {
			str = strconv.Itoa(i)
		}
		output = append(output, str)
	}
	return output, nil
}
