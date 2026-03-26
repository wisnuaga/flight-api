package util

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// FormatPrice formats an integer amount with thousands separator (for backward compatibility)
func FormatPrice(amount int64, currency string) string {
	str := fmt.Sprintf("%d", amount)
	var result []byte
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return fmt.Sprintf("%s %s", currency, string(result))
}

// FormatPriceDecimal formats a decimal.Decimal amount with thousands separator and decimal places
// Example: 1500000.50 → "IDR 1.500.000,50"
// Handles both integer and fractional parts
func FormatPriceDecimal(price decimal.Decimal, currency string) string {
	// Get the integer and decimal parts
	intPart := price.IntPart()
	fracPart := price.Mul(decimal.NewFromInt(100)).IntPart() % 100

	// Format the integer part with thousands separator (dot)
	intStr := fmt.Sprintf("%d", intPart)
	var intResult []byte
	for i, c := range intStr {
		if i > 0 && (len(intStr)-i)%3 == 0 {
			intResult = append(intResult, '.')
		}
		intResult = append(intResult, byte(c))
	}

	// Format fractional part with 2 decimal places using comma separator
	formattedInt := string(intResult)
	if fracPart > 0 {
		return fmt.Sprintf("%s %s,%02d", currency, formattedInt, fracPart)
	}
	return fmt.Sprintf("%s %s,00", currency, formattedInt)
}

// CurrencySymbols maps currency codes to their symbols
var CurrencySymbols = map[string]string{
	"IDR": "Rp",
	"USD": "$",
	"SGD": "S$",
	"MYR": "RM",
	"THB": "฿",
	"PHP": "₱",
}

// GetCurrencySymbol returns the symbol for a currency code, or the code itself as fallback
func GetCurrencySymbol(currencyCode string) string {
	if symbol, ok := CurrencySymbols[currencyCode]; ok {
		return symbol
	}
	return currencyCode
}

// FormatPriceWithSymbol formats a price with currency symbol instead of code
// Example: 1500000.50 → "Rp 1.500.000,50"
func FormatPriceWithSymbol(price decimal.Decimal, currencyCode string) string {
	formatted := FormatPriceDecimal(price, currencyCode)
	// Replace currency code with symbol
	symbol := GetCurrencySymbol(currencyCode)
	return strings.Replace(formatted, currencyCode, symbol, 1)
}
