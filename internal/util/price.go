package util

import "fmt"

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
