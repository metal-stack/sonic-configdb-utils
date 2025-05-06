package platform

import (
	"fmt"
	"strconv"
	"strings"
)

func stringToIntSlice(input string) ([]int, error) {
	ints := make([]int, 0)

	for n := range strings.SplitSeq(input, ",") {
		number, err := strconv.Atoi(n)
		if err != nil {
			return []int{}, err
		}
		ints = append(ints, number)
	}

	return ints, nil
}

// cutBetween finds the first occurence of <leftDelim>x<rightDelim> and returns x if found.
// If no such substring is found cutBetween returns s, false.
func cutBetween(s, leftDelim, rightDelim string) (between string, found bool) {
	_, between, found = strings.Cut(s, leftDelim)
	if !found {
		return s, found
	}

	between, _, found = strings.Cut(between, rightDelim)
	if !found {
		return s, found
	}

	return between, found
}

// stringToSpeed takes a string of the form xG,where x must be a positive integer, and returns x*1000.
// If the provided string does not comply with the expected format stringToSpeed returns 0 and an error.
func stringToSpeed(s string) (int, error) {
	parseError := fmt.Errorf("speed must be of format xG, where x is a positive integer")

	speedString, found := strings.CutSuffix(s, "G")
	if !found {
		return 0, parseError
	}

	speed, err := strconv.Atoi(speedString)
	if err != nil {
		return 0, parseError
	}

	if speed <= 0 {
		return 0, parseError
	}

	return speed * 1000, nil
}
