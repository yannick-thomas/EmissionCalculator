package validation

import (
	"math"
	"strings"
	"testing"
)

func TestParseQuantityAcceptsUnambiguousFormats(t *testing.T) {
	tests := map[string]float64{
		"1250":      1250,
		"1250,5":    1250.5,
		"1250.5":    1250.5,
		"12.500,50": 12500.5,
		"1.000.000": 1000000,
		"+42":       42,
	}
	for input, expected := range tests {
		actual, err := ParseQuantity(input)
		if err != nil {
			t.Fatalf("ParseQuantity(%q) returned %v", input, err)
		}
		if math.Abs(actual-expected) > 0.000001 {
			t.Fatalf("ParseQuantity(%q) = %v, expected %v", input, actual, expected)
		}
	}
}

func TestParseQuantityRejectsAmbiguousAndNonFiniteValues(t *testing.T) {
	for _, input := range []string{"1.234", "12.500", "NaN", "+Inf", "-Inf", "1e3"} {
		_, err := ParseQuantity(input)
		if err == nil {
			t.Fatalf("ParseQuantity(%q) unexpectedly succeeded", input)
		}
		if strings.Contains(input, ".") && !strings.Contains(err.Error(), "mehrdeutig") {
			t.Fatalf("ParseQuantity(%q) returned an unclear error: %v", input, err)
		}
	}
}

func TestParseQuantityRejectsMalformedGrouping(t *testing.T) {
	for _, input := range []string{"1.23.456", "12.34,5", "1,2,3", "1 000", ".5", "5."} {
		if _, err := ParseQuantity(input); err == nil {
			t.Fatalf("ParseQuantity(%q) unexpectedly succeeded", input)
		}
	}
}
