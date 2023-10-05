package utils

import (
	"math/big"
	"regexp"
)

func CheckIfValidAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	return re.MatchString(addr)
}

func HexToEth(hex string) *big.Float {
	n := new(big.Int)
	n.SetString(hex, 0)
	f := new(big.Float).SetInt(n)
	gwei := big.NewFloat(1000000000000000000)
	return f.Quo(f, gwei)
}
