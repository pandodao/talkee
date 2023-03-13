package util

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var OneSatoshi decimal.Decimal

func init() {
	OneSatoshi, _ = decimal.NewFromString("0.00000001")
}

func ValidateMixinAmount(amount decimal.Decimal) bool {
	return !amount.LessThan(OneSatoshi)
}

func RoundUp(amount decimal.Decimal) decimal.Decimal {
	return amount.Round(8)
}

func RoundDown(amount decimal.Decimal) decimal.Decimal {
	return amount.RoundDown(8)
}

func MixinAmount(amount decimal.Decimal) decimal.Decimal {
	if !ValidateMixinAmount(amount) {
		RoundUp(OneSatoshi)
	}
	return RoundUp(amount)
}

func TrimTextForDisplay(input string, max int) string {
	// xxxxx...xxxx
	if max < 12 {
		max = 12
	}
	asRunes := []rune(input)
	if len(asRunes) <= max {
		return input
	}
	prefix := asRunes[0:5]
	suffix := asRunes[len(asRunes)-4 : len(asRunes)]
	return fmt.Sprintf("%s...%s", string(prefix), string(suffix))
}

func TrimTailZero(input string) string {
	return strings.TrimRight(strings.TrimRight(input, "0"), ".")
}
