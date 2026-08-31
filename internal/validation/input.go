package validation

import (
	"fmt"
	"math"
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

	normalized, err := normalizeQuantity(trimmed)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("Bitte gültige Zahl eingeben")
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("Die Liefermenge muss eine endliche Zahl sein")
	}
	if value <= 0 {
		return 0, fmt.Errorf("Die Liefermenge muss größer als 0 sein")
	}
	if value > MaxQuantity {
		return 0, fmt.Errorf("Die Liefermenge darf %.0f nicht überschreiten", MaxQuantity)
	}
	return value, nil
}

func normalizeQuantity(input string) (string, error) {
	sign := ""
	unsigned := input
	if strings.HasPrefix(unsigned, "+") || strings.HasPrefix(unsigned, "-") {
		sign, unsigned = unsigned[:1], unsigned[1:]
	}
	if unsigned == "" || strings.ContainsAny(unsigned, " \t\r\n") {
		return "", fmt.Errorf("Bitte gültige Zahl eingeben")
	}

	commaCount := strings.Count(unsigned, ",")
	dotCount := strings.Count(unsigned, ".")
	if commaCount > 1 {
		return "", fmt.Errorf("Bitte nur ein Dezimaltrennzeichen verwenden")
	}

	if commaCount == 1 {
		parts := strings.SplitN(unsigned, ",", 2)
		if !validIntegerPart(parts[0]) || !allDigits(parts[1]) {
			return "", fmt.Errorf("Bitte gültige Zahl eingeben")
		}
		return sign + strings.ReplaceAll(parts[0], ".", "") + "." + parts[1], nil
	}

	if dotCount == 0 {
		if !allDigits(unsigned) {
			return "", fmt.Errorf("Bitte gültige Zahl eingeben")
		}
		return sign + unsigned, nil
	}

	if dotCount == 1 {
		parts := strings.SplitN(unsigned, ".", 2)
		if !allDigits(parts[0]) || !allDigits(parts[1]) {
			return "", fmt.Errorf("Bitte gültige Zahl eingeben")
		}
		if len(parts[1]) == 3 {
			return "", fmt.Errorf("Die Eingabe ist mehrdeutig: Bitte Dezimalstellen mit Komma oder die Zahl ohne Tausenderpunkt eingeben")
		}
		return sign + unsigned, nil
	}

	if !validIntegerPart(unsigned) {
		return "", fmt.Errorf("Bitte gültige Zahl eingeben")
	}
	return sign + strings.ReplaceAll(unsigned, ".", ""), nil
}

func validIntegerPart(value string) bool {
	groups := strings.Split(value, ".")
	if len(groups) == 1 {
		return allDigits(groups[0])
	}
	if len(groups[0]) < 1 || len(groups[0]) > 3 || !allDigits(groups[0]) {
		return false
	}
	for _, group := range groups[1:] {
		if len(group) != 3 || !allDigits(group) {
			return false
		}
	}
	return true
}

func allDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}
