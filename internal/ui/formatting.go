package ui

import (
	"emissioncalculator/internal/calculation"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func formatFloat(value float64, decimals int) string {
	factor := math.Pow10(decimals)
	// add a sub-ULP nudge so that x.xx5 rounds away from zero (banker's rounding mitigation)
	rounded := math.Round((value+math.Copysign(1e-9, value))*factor) / factor
	formatted := strconv.FormatFloat(rounded, 'f', decimals, 64)
	parts := strings.SplitN(formatted, ".", 2)
	integer := parts[0]
	sign := ""
	if strings.HasPrefix(integer, "-") {
		sign = "-"
		integer = strings.TrimPrefix(integer, "-")
	}
	for index := len(integer) - 3; index > 0; index -= 3 {
		integer = integer[:index] + "." + integer[index:]
	}
	if len(parts) == 1 || decimals == 0 {
		return sign + integer
	}
	return sign + integer + "," + parts[1]
}

func formatGermanNumberString(value string, decimals int) string {
	normalized := strings.ReplaceAll(strings.TrimSpace(value), ".", "")
	normalized = strings.ReplaceAll(normalized, ",", ".")
	number, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return value
	}
	factor := math.Pow10(decimals)
	number = math.Round((number+math.Copysign(1e-9, number))*factor) / factor
	formatted := strconv.FormatFloat(number, 'f', decimals, 64)
	parts := strings.SplitN(formatted, ".", 2)
	integer := parts[0]
	sign := ""
	if strings.HasPrefix(integer, "-") {
		sign = "-"
		integer = strings.TrimPrefix(integer, "-")
	}
	for index := len(integer) - 3; index > 0; index -= 3 {
		integer = integer[:index] + "." + integer[index:]
	}
	if len(parts) == 1 || decimals == 0 {
		return sign + integer
	}
	return sign + integer + "," + parts[1]
}

func formatQuantityDisplay(value float64) string {
	whole := int64(math.Floor(value))
	digits := strconv.FormatInt(whole, 10)
	for index := len(digits) - 3; index > 0; index -= 3 {
		digits = digits[:index] + "." + digits[index:]
	}
	fraction := value - float64(whole)
	if math.Abs(fraction) < 0.005 {
		return digits
	}
	decimal := fmt.Sprintf("%.2f", fraction)
	decimal = strings.TrimPrefix(decimal, "0")
	decimal = strings.TrimRight(decimal, "0")
	return digits + strings.ReplaceAll(decimal, ".", ",")
}

func formatMeasurement(value string, decimals int) string {
	parts := strings.SplitN(strings.TrimSpace(value), " ", 2)
	formatted := formatGermanNumberString(parts[0], decimals)
	if len(parts) == 1 {
		return formatted
	}
	return formatted + " " + parts[1]
}

func titleForMode(mode string) string {
	if descriptor, ok := calculation.FuelByType(calculation.FuelType(mode)); ok {
		return descriptor.Label
	}
	return mode
}

func unitForMode(mode string) string {
	if descriptor, ok := calculation.FuelByType(calculation.FuelType(mode)); ok {
		switch descriptor.Unit {
		case "L":
			return "Liter"
		case "t":
			return "Tonnen"
		case "kg":
			return "Kilogramm"
		case "kWh":
			return "kWh"
		default:
			return string(descriptor.Unit)
		}
	}
	return ""
}
