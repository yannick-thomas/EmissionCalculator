package calculation

import (
	"fmt"
	"strings"
)

func FormatNumber(value float64, energyMode bool) string {
	if energyMode {
		return strings.ReplaceAll(fmt.Sprintf("%.3f", value), ".", ",")
	}
	return strings.ReplaceAll(fmt.Sprintf("%.2f", value), ".", ",")
}

func FormatWithLocale(value float64, energyMode bool) string {
	return FormatNumber(value, energyMode)
}
