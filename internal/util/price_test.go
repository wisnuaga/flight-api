package util

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestFormatPrice_Integer(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		want     string
	}{
		{
			name:     "Simple number",
			amount:   1500000,
			currency: "IDR",
			want:     "IDR 1.500.000",
		},
		{
			name:     "Small number",
			amount:   100,
			currency: "IDR",
			want:     "IDR 100",
		},
		{
			name:     "Zero",
			amount:   0,
			currency: "IDR",
			want:     "IDR 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPrice(tt.amount, tt.currency)
			if got != tt.want {
				t.Errorf("FormatPrice(%d, %s) = %s, want %s", tt.amount, tt.currency, got, tt.want)
			}
		})
	}
}

func TestFormatPriceDecimal_Integer(t *testing.T) {
	tests := []struct {
		name     string
		price    decimal.Decimal
		currency string
		want     string
	}{
		{
			name:     "Integer with thousands",
			price:    decimal.NewFromInt(1500000),
			currency: "IDR",
			want:     "IDR 1.500.000,00",
		},
		{
			name:     "Small integer",
			price:    decimal.NewFromInt(100),
			currency: "IDR",
			want:     "IDR 100,00",
		},
		{
			name:     "Zero",
			price:    decimal.NewFromInt(0),
			currency: "IDR",
			want:     "IDR 0,00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPriceDecimal(tt.price, tt.currency)
			if got != tt.want {
				t.Errorf("FormatPriceDecimal(%s, %s) = %s, want %s", tt.price, tt.currency, got, tt.want)
			}
		})
	}
}

func TestFormatPriceDecimal_Fractional(t *testing.T) {
	tests := []struct {
		name     string
		price    string
		currency string
		want     string
	}{
		{
			name:     "With decimal cents",
			price:    "1500000.50",
			currency: "IDR",
			want:     "IDR 1.500.000,50",
		},
		{
			name:     "With decimal cents (01)",
			price:    "100.01",
			currency: "IDR",
			want:     "IDR 100,01",
		},
		{
			name:     "99 cents",
			price:    "50.99",
			currency: "IDR",
			want:     "IDR 50,99",
		},
		{
			name:     "1 cent",
			price:    "1000.01",
			currency: "IDR",
			want:     "IDR 1.000,01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := decimal.NewFromString(tt.price)
			got := FormatPriceDecimal(p, tt.currency)
			if got != tt.want {
				t.Errorf("FormatPriceDecimal(%s, %s) = %s, want %s", tt.price, tt.currency, got, tt.want)
			}
		})
	}
}

func TestGetCurrencySymbol(t *testing.T) {
	tests := []struct {
		name         string
		currencyCode string
		want         string
	}{
		{
			name:         "IDR to Rp",
			currencyCode: "IDR",
			want:         "Rp",
		},
		{
			name:         "USD to $",
			currencyCode: "USD",
			want:         "$",
		},
		{
			name:         "SGD to S$",
			currencyCode: "SGD",
			want:         "S$",
		},
		{
			name:         "Unknown returns code",
			currencyCode: "XXX",
			want:         "XXX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCurrencySymbol(tt.currencyCode)
			if got != tt.want {
				t.Errorf("GetCurrencySymbol(%s) = %s, want %s", tt.currencyCode, got, tt.want)
			}
		})
	}
}

func TestFormatPriceWithSymbol(t *testing.T) {
	tests := []struct {
		name         string
		price        string
		currencyCode string
		want         string
	}{
		{
			name:         "IDR with decimal",
			price:        "1500000.50",
			currencyCode: "IDR",
			want:         "Rp 1.500.000,50",
		},
		{
			name:         "USD with decimal",
			price:        "99.99",
			currencyCode: "USD",
			want:         "$ 99,99",
		},
		{
			name:         "Unknown currency",
			price:        "100.00",
			currencyCode: "XXX",
			want:         "XXX 100,00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := decimal.NewFromString(tt.price)
			got := FormatPriceWithSymbol(p, tt.currencyCode)
			if got != tt.want {
				t.Errorf("FormatPriceWithSymbol(%s, %s) = %s, want %s", tt.price, tt.currencyCode, got, tt.want)
			}
		})
	}
}
