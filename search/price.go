package search

import (
	"regexp"
	"strconv"
	"strings"
)

var priceDigits = regexp.MustCompile(`\d[\d,.]*`)

func PriceConversion(price string) (int, error) {
	//strip any contexts of outlying values in the string that are not needed for values sake.
	price = priceDigits.FindString(strings.TrimSpace(price))
	price = strings.ReplaceAll(price, ",", "")
	// Mexican store listings commonly use a decimal point after the peso amount.
	// Drop the cents instead of treating them as part of the price.
	price = strings.Split(price, ".")[0]
	//convert string to int
	cost, err := strconv.Atoi(price)
	if err != nil {
		return 0, err
	}
	return cost, nil
}
