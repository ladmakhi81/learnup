package utils

import (
	"fmt"
	"strconv"
)

func ToString(val any) string {
	return fmt.Sprintf("%v", val)
}

func ToUint(val string) (uint, error) {
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
