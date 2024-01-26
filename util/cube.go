package util

import (
	"fmt"
)

func FindRange(min int, max int, step int, pointer int) (rangeMin int, err error) {
	if !(min <= pointer && pointer <= max) {
		return 0, fmt.Errorf("pointer is out of range")
	}
	for i := min; i <= max; i += step {
		if i <= pointer && pointer <= i+step-1 {
			rangeMin = i
			return rangeMin, nil
		}
	}
	return 0, fmt.Errorf("pointer is out of range")
}
