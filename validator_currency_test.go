package validator

import "testing"

func TestIsCurrencyFiat(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"CNY", true},
		{"JPY", true},
		{"TWD", true},
		{"KRW", true},
		{"HKD", true},
		{"SGD", true},
		{"AUD", true},
		{"CAD", true},
		{"CHF", true},
		{"usd", true}, // lowercase should work
		{"eur", true},
		{"XXX", true}, // No currency (valid ISO 4217)
		{"XAU", true}, // Gold
		{"XAG", true}, // Silver
		{"ABC", false},
		{"", false},
		{"US", false},
		{"USDD", false},
		{"BTC", false}, // crypto not fiat
		{"ETH", false},
	}
	for _, test := range tests {
		result := IsCurrencyFiat(test.input)
		if result != test.expected {
			t.Errorf("IsCurrencyFiat(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsCurrencyCrypto(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Major cryptocurrencies
		{"BTC", true},
		{"ETH", true},
		{"USDT", true},
		{"USDC", true},
		{"BNB", true},
		{"XRP", true},
		{"ADA", true},
		{"DOGE", true},
		{"SOL", true},
		{"DOT", true},
		// Case insensitive
		{"btc", true},
		{"eth", true},
		{"Btc", true},
		// DeFi tokens
		{"UNI", true},
		{"AAVE", true},
		{"LINK", true},
		// Meme coins
		{"SHIB", true},
		{"PEPE", true},
		{"BONK", true},
		{"WIF", true},
		// AI tokens
		{"FET", true},
		{"RNDR", true},
		// Invalid
		{"USD", false},  // fiat currency
		{"EUR", false},  // fiat currency
		{"ABC", false},  // not a currency
		{"FAKE", false}, // not a known crypto
		{"", false},
	}
	for _, test := range tests {
		result := IsCurrencyCrypto(test.input)
		if result != test.expected {
			t.Errorf("IsCurrencyCrypto(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsCurrencyAll(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Fiat currencies
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"JPY", true},
		{"CNY", true},
		// Cryptocurrencies
		{"BTC", true},
		{"ETH", true},
		{"USDT", true},
		{"SOL", true},
		{"DOGE", true},
		// Case insensitive
		{"usd", true},
		{"btc", true},
		// Invalid
		{"ABC", false},
		{"FAKE", false},
		{"", false},
	}
	for _, test := range tests {
		result := IsCurrencyAll(test.input)
		if result != test.expected {
			t.Errorf("IsCurrencyAll(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// Benchmark tests
func BenchmarkIsCurrencyFiat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsCurrencyFiat("USD")
	}
}

func BenchmarkIsCurrencyCrypto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsCurrencyCrypto("BTC")
	}
}

func BenchmarkIsCurrencyAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsCurrencyAll("BTC")
	}
}
