package utils

import "regexp"

func CheckIfValidAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	return re.MatchString(addr)
}
