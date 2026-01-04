package validator

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestCompareDecimal(t *testing.T) {
	tests := []struct {
		name     string
		first    decimal.Decimal
		second   decimal.Decimal
		operator string
		expected bool
		hasError bool
	}{
		{"less than true", decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5), "<", true, false},
		{"less than false", decimal.NewFromFloat(2.5), decimal.NewFromFloat(1.5), "<", false, false},
		{"greater than true", decimal.NewFromFloat(2.5), decimal.NewFromFloat(1.5), ">", true, false},
		{"greater than false", decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5), ">", false, false},
		{"less than or equal true", decimal.NewFromFloat(1.5), decimal.NewFromFloat(1.5), "<=", true, false},
		{"less than or equal true 2", decimal.NewFromFloat(1.0), decimal.NewFromFloat(1.5), "<=", true, false},
		{"greater than or equal true", decimal.NewFromFloat(1.5), decimal.NewFromFloat(1.5), ">=", true, false},
		{"equal true", decimal.NewFromFloat(1.5), decimal.NewFromFloat(1.5), "==", true, false},
		{"equal false", decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5), "==", false, false},
		{"invalid operator", decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5), "!=", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compareDecimal(tt.first, tt.second, tt.operator)
			if (err != nil) != tt.hasError {
				t.Errorf("compareDecimal() error = %v, hasError %v", err, tt.hasError)
				return
			}
			if result != tt.expected {
				t.Errorf("compareDecimal() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsDecimalRequired(t *testing.T) {
	if IsDecimalRequired(decimal.Zero) {
		t.Error("zero should be false")
	}
	if !IsDecimalRequired(decimal.NewFromFloat(1.5)) {
		t.Error("non-zero should be true")
	}
	if !IsDecimalRequired(decimal.NewFromFloat(-1.5)) {
		t.Error("negative non-zero should be true")
	}
}

func TestIsDecimalMin(t *testing.T) {
	min := decimal.NewFromFloat(10.0)

	if !IsDecimalMin(decimal.NewFromFloat(15.0), min) {
		t.Error("15 >= 10 should be true")
	}
	if !IsDecimalMin(decimal.NewFromFloat(10.0), min) {
		t.Error("10 >= 10 should be true")
	}
	if IsDecimalMin(decimal.NewFromFloat(5.0), min) {
		t.Error("5 >= 10 should be false")
	}
}

func TestIsDecimalMax(t *testing.T) {
	max := decimal.NewFromFloat(10.0)

	if !IsDecimalMax(decimal.NewFromFloat(5.0), max) {
		t.Error("5 <= 10 should be true")
	}
	if !IsDecimalMax(decimal.NewFromFloat(10.0), max) {
		t.Error("10 <= 10 should be true")
	}
	if IsDecimalMax(decimal.NewFromFloat(15.0), max) {
		t.Error("15 <= 10 should be false")
	}
}

func TestIsDecimalBetween(t *testing.T) {
	min := decimal.NewFromFloat(1.0)
	max := decimal.NewFromFloat(10.0)

	if !IsDecimalBetween(decimal.NewFromFloat(5.0), min, max) {
		t.Error("5 between 1-10 should be true")
	}
	if !IsDecimalBetween(decimal.NewFromFloat(1.0), min, max) {
		t.Error("1 between 1-10 should be true")
	}
	if !IsDecimalBetween(decimal.NewFromFloat(10.0), min, max) {
		t.Error("10 between 1-10 should be true")
	}
	if IsDecimalBetween(decimal.NewFromFloat(0.5), min, max) {
		t.Error("0.5 between 1-10 should be false")
	}
	if IsDecimalBetween(decimal.NewFromFloat(10.5), min, max) {
		t.Error("10.5 between 1-10 should be false")
	}
}

func TestIsDecimalGtGteLtLte(t *testing.T) {
	threshold := decimal.NewFromFloat(5.0)

	// Gt
	if !IsDecimalGt(decimal.NewFromFloat(6.0), threshold) {
		t.Error("6 > 5 should be true")
	}
	if IsDecimalGt(decimal.NewFromFloat(5.0), threshold) {
		t.Error("5 > 5 should be false")
	}

	// Gte
	if !IsDecimalGte(decimal.NewFromFloat(6.0), threshold) {
		t.Error("6 >= 5 should be true")
	}
	if !IsDecimalGte(decimal.NewFromFloat(5.0), threshold) {
		t.Error("5 >= 5 should be true")
	}
	if IsDecimalGte(decimal.NewFromFloat(4.0), threshold) {
		t.Error("4 >= 5 should be false")
	}

	// Lt
	if !IsDecimalLt(decimal.NewFromFloat(4.0), threshold) {
		t.Error("4 < 5 should be true")
	}
	if IsDecimalLt(decimal.NewFromFloat(5.0), threshold) {
		t.Error("5 < 5 should be false")
	}

	// Lte
	if !IsDecimalLte(decimal.NewFromFloat(4.0), threshold) {
		t.Error("4 <= 5 should be true")
	}
	if !IsDecimalLte(decimal.NewFromFloat(5.0), threshold) {
		t.Error("5 <= 5 should be true")
	}
	if IsDecimalLte(decimal.NewFromFloat(6.0), threshold) {
		t.Error("6 <= 5 should be false")
	}
}

func TestIsDecimalPositiveNegative(t *testing.T) {
	if !IsDecimalPositive(decimal.NewFromFloat(1.0)) {
		t.Error("1.0 should be positive")
	}
	if IsDecimalPositive(decimal.NewFromFloat(-1.0)) {
		t.Error("-1.0 should not be positive")
	}
	if IsDecimalPositive(decimal.Zero) {
		t.Error("0 should not be positive")
	}

	if !IsDecimalNegative(decimal.NewFromFloat(-1.0)) {
		t.Error("-1.0 should be negative")
	}
	if IsDecimalNegative(decimal.NewFromFloat(1.0)) {
		t.Error("1.0 should not be negative")
	}
	if IsDecimalNegative(decimal.Zero) {
		t.Error("0 should not be negative")
	}
}

func TestIsDecimalEqual(t *testing.T) {
	if !IsDecimalEqual(decimal.NewFromFloat(1.5), decimal.NewFromFloat(1.5)) {
		t.Error("1.5 == 1.5 should be true")
	}
	if IsDecimalEqual(decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5)) {
		t.Error("1.5 == 2.5 should be false")
	}

	// Test precision
	a := decimal.NewFromFloat(0.1).Add(decimal.NewFromFloat(0.2))
	b := decimal.NewFromFloat(0.3)
	if !IsDecimalEqual(a, b) {
		t.Errorf("0.1 + 0.2 should equal 0.3, got %s vs %s", a.String(), b.String())
	}
}

func TestValidateDecimalFunctions(t *testing.T) {
	// Test ValidateDecimalBetween
	if !ValidateDecimalBetween(decimal.NewFromFloat(5.0), decimal.NewFromFloat(1.0), decimal.NewFromFloat(10.0)) {
		t.Error("ValidateDecimalBetween failed")
	}

	// Test ValidateDecimalMin
	if !ValidateDecimalMin(decimal.NewFromFloat(10.0), decimal.NewFromFloat(5.0)) {
		t.Error("ValidateDecimalMin failed")
	}

	// Test ValidateDecimalMax
	if !ValidateDecimalMax(decimal.NewFromFloat(5.0), decimal.NewFromFloat(10.0)) {
		t.Error("ValidateDecimalMax failed")
	}

	// Test ValidateDecimalZero
	if !ValidateDecimalZero(decimal.Zero) {
		t.Error("ValidateDecimalZero failed for zero")
	}
	if ValidateDecimalZero(decimal.NewFromFloat(1.0)) {
		t.Error("ValidateDecimalZero should fail for non-zero")
	}
}

// Struct validation tests
func TestValidateStructDecimal(t *testing.T) {
	type Product struct {
		Name  string          `valid:"required"`
		Price decimal.Decimal `valid:"required,min=0.01,max=9999.99"`
	}

	tests := []struct {
		name      string
		product   Product
		wantError bool
	}{
		{
			name:      "valid product",
			product:   Product{Name: "Widget", Price: decimal.NewFromFloat(99.99)},
			wantError: false,
		},
		{
			name:      "zero price fails required",
			product:   Product{Name: "Widget", Price: decimal.Zero},
			wantError: true,
		},
		{
			name:      "price below min",
			product:   Product{Name: "Widget", Price: decimal.NewFromFloat(0.001)},
			wantError: true,
		},
		{
			name:      "price above max",
			product:   Product{Name: "Widget", Price: decimal.NewFromFloat(10000.00)},
			wantError: true,
		},
		{
			name:      "price at min boundary",
			product:   Product{Name: "Widget", Price: decimal.NewFromFloat(0.01)},
			wantError: false,
		},
		{
			name:      "price at max boundary",
			product:   Product{Name: "Widget", Price: decimal.NewFromFloat(9999.99)},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.product)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStructDecimalBetween(t *testing.T) {
	type Order struct {
		Quantity decimal.Decimal `valid:"between=1|100"`
	}

	tests := []struct {
		name      string
		order     Order
		wantError bool
	}{
		{
			name:      "quantity in range",
			order:     Order{Quantity: decimal.NewFromFloat(50)},
			wantError: false,
		},
		{
			name:      "quantity below range",
			order:     Order{Quantity: decimal.NewFromFloat(0.5)},
			wantError: true,
		},
		{
			name:      "quantity above range",
			order:     Order{Quantity: decimal.NewFromFloat(101)},
			wantError: true,
		},
		{
			name:      "quantity at min boundary",
			order:     Order{Quantity: decimal.NewFromFloat(1)},
			wantError: false,
		},
		{
			name:      "quantity at max boundary",
			order:     Order{Quantity: decimal.NewFromFloat(100)},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.order)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStructDecimalComparison(t *testing.T) {
	type Transaction struct {
		Amount decimal.Decimal `valid:"gt=0"`
	}

	tests := []struct {
		name      string
		tx        Transaction
		wantError bool
	}{
		{
			name:      "amount greater than zero",
			tx:        Transaction{Amount: decimal.NewFromFloat(10.50)},
			wantError: false,
		},
		{
			name:      "amount zero fails gt",
			tx:        Transaction{Amount: decimal.Zero},
			wantError: true,
		},
		{
			name:      "negative amount fails gt",
			tx:        Transaction{Amount: decimal.NewFromFloat(-5)},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.tx)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStructDecimalFieldComparison(t *testing.T) {
	type PriceRange struct {
		MinPrice decimal.Decimal `valid:"required"`
		MaxPrice decimal.Decimal `valid:"required,gt=MinPrice"`
	}

	tests := []struct {
		name      string
		pr        PriceRange
		wantError bool
	}{
		{
			name: "max greater than min",
			pr: PriceRange{
				MinPrice: decimal.NewFromFloat(10),
				MaxPrice: decimal.NewFromFloat(100),
			},
			wantError: false,
		},
		{
			name: "max equals min fails gt",
			pr: PriceRange{
				MinPrice: decimal.NewFromFloat(50),
				MaxPrice: decimal.NewFromFloat(50),
			},
			wantError: true,
		},
		{
			name: "max less than min fails gt",
			pr: PriceRange{
				MinPrice: decimal.NewFromFloat(100),
				MaxPrice: decimal.NewFromFloat(50),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.pr)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStructDecimalSame(t *testing.T) {
	type ConfirmAmount struct {
		Amount        decimal.Decimal `valid:"required"`
		ConfirmAmount decimal.Decimal `valid:"required,same=Amount"`
	}

	tests := []struct {
		name      string
		ca        ConfirmAmount
		wantError bool
	}{
		{
			name: "amounts match",
			ca: ConfirmAmount{
				Amount:        decimal.NewFromFloat(100.50),
				ConfirmAmount: decimal.NewFromFloat(100.50),
			},
			wantError: false,
		},
		{
			name: "amounts differ",
			ca: ConfirmAmount{
				Amount:        decimal.NewFromFloat(100.50),
				ConfirmAmount: decimal.NewFromFloat(100.51),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.ca)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// Benchmarks
func BenchmarkIsDecimalMin(b *testing.B) {
	value := decimal.NewFromFloat(100.0)
	min := decimal.NewFromFloat(10.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsDecimalMin(value, min)
	}
}

func BenchmarkIsDecimalBetween(b *testing.B) {
	value := decimal.NewFromFloat(50.0)
	min := decimal.NewFromFloat(1.0)
	max := decimal.NewFromFloat(100.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsDecimalBetween(value, min, max)
	}
}

func BenchmarkCompareDecimal(b *testing.B) {
	first := decimal.NewFromFloat(100.0)
	second := decimal.NewFromFloat(50.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compareDecimal(first, second, ">")
	}
}
