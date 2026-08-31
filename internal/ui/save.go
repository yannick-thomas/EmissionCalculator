package ui

import (
	"emissioncalculator/internal/models"
	"emissioncalculator/internal/pdf"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

var saveLabelAs = defaultSaveLabelAs

// defaultSaveLabelAs opens a native "save as" dialog and writes the rendered PDF to the chosen
// location. onResult is called with the saved path, or an empty path if the user cancelled.
func defaultSaveLabelAs(window fyne.Window, record models.CalculationRecord, onResult func(path string, err error)) {
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			onResult("", err)
			return
		}
		if writer == nil {
			onResult("", nil)
			return
		}
		defer writer.Close()
		data, renderErr := pdf.RenderLabel(record)
		if renderErr != nil {
			onResult("", renderErr)
			return
		}
		if _, writeErr := writer.Write(data); writeErr != nil {
			onResult("", writeErr)
			return
		}
		onResult(writer.URI().Path(), nil)
	}, window)
	saveDialog.SetFileName(fmt.Sprintf("emission_label_%s.pdf", record.FuelType))
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	saveDialog.Show()
}
