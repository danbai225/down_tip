package utils

import "strconv"

func ComplementZero(i uint) string {
	itoa := strconv.Itoa(int(i))
	if i < 10 {
		return "0" + itoa
	}
	return itoa
}
