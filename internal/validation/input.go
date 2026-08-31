package validation

import (
	"fmt"
	"strconv"
	"strings"
)

// MaxQuantity is the largest accepted delivery quantity (applies to both litres and tonnes).
const MaxQuantity = 1_000_000.0

func ParseQuantity(input string) (float64, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return 0, fmt.Errorf("Bitte gültige Liefermenge eingeben")
	}

	var normalized string
	switch {
	case strings.Contains(trimmed, ","):
		// German format: periods are thousands separators, comma is decimal
		normalized = strings.ReplaceAll(trimmed, ".", "")
		normalized = strings.ReplaceAll(normalized, ",", ".")
	case strings.Contains(trimmed, "."):
		// Period with exactly 3 trailing digits → German thousands separator
		lastDot := strings.LastIndex(trimmed, ".")
		if len(trimmed)-lastDot-1 == 3 {
			normalized = strings.ReplaceAll(trimmed, ".", "")
		} else {
			normalized = trimmed
		}
	default:
		normalized = trimmed
	}

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("Bitte gültige Zahl eingeben")
	}
	if value <= 0 {
		return 0, fmt.Errorf("Die Liefermenge muss größer als 0 sein")
	}
	if value > MaxQuantity {
		return 0, fmt.Errorf("Die Liefermenge darf %.0f nicht überschreiten", MaxQuantity)
	}
	return value, nil
}
