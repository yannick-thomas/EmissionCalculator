package validation

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseQuantity(input string) (float64, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return 0, fmt.Errorf("Bitte gültige Liefermenge eingeben")
	}

	normalized := strings.ReplaceAll(trimmed, ",", ".")
	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("Bitte gültige Zahl eingeben")
	}
	if value <= 0 {
		return 0, fmt.Errorf("Die Liefermenge muss größer als 0 sein")
	}
	return value, nil
}
