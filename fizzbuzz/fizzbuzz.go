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
	if params.Limit < 1 {
		return []string{}, errors.New("limit must be superior or equal to 1")
	}
	for i := 1; i <= params.Limit; i++ {
		str := ""
		if i%params.Int1 == 0 {
			str = "fizz"
		}
		if i%params.Int2 == 0 {
			str = str + "buzz"
		}
		if str == "" {
			str = strconv.Itoa(i)
		}
		output = append(output, str)
	}
	return output, nil
}
