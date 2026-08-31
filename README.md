# EmissionCalculator

Eine lokale Go-Desktopanwendung zur Berechnung direkter CO₂-Emissionen, des Energiegehalts und des CO₂-Kostenanteils von Brennstofflieferungen.

## Unterstützte Brennstoffe

| Brennstoff | Eingabeeinheit | Systemgrenze |
|---|---:|---|
| Heizöl | Liter | direkte Verbrennung |
| Braunkohlebriketts | Tonnen | direkte Verbrennung |
| Erdgas | kWh | direkte Verbrennung |
| Flüssiggas | Kilogramm | direkte Verbrennung |

Die physikalischen Faktoren liegen in einem validierten und versionierten Datenpaket unter `internal/calculation/data/`. Jede Berechnung speichert den verwendeten Faktor-Snapshot, das Faktorenpaket, Preis- und Quellenstand sowie einen vollständigen Rechenweg. Holz und Pellets sind bewusst nicht enthalten: Für biogene Brennstoffe muss zuerst eine fachlich eindeutige Systemgrenze zwischen direkten, biogenen und vorgelagerten Emissionen festgelegt werden.

## Funktionen

- nachvollziehbare Ansicht „So wurde gerechnet“ mit Einzelschritten und Quellenlinks
- konfigurierbarer Vergleich von CO₂-Preisen und Berechnungsjahren
- lokale Projekte und automatisch gespeicherte Berechnungshistorie
- PDF-Nachweis mit eingebetteter Unicode-Schrift
- CSV-/JSON-Verarbeitung über die internen Batch-Pakete
- deutsche Zahleneingabe mit Prüfung mehrdeutiger und nicht endlicher Werte

Die Historie wird ausschließlich lokal als private JSON-Datei im Konfigurationsverzeichnis des Betriebssystems gespeichert (Dateimodus `0600`, Verzeichnis `0700`). Es werden keine Berechnungsdaten an einen Server übertragen.

## Lokale Entwicklung

Voraussetzung ist die in `go.mod` angegebene Go-Version. Fyne benötigt je nach Betriebssystem zusätzliche Grafikbibliotheken.

```sh
go run ./cmd/emissioncalculator
go test ./...
go test -race ./...
go vet ./...
```

## Faktoren aktualisieren

Ein veröffentlichtes Faktorenpaket wird nicht nachträglich verändert. Für einen neuen Quellenstand wird eine neue Paketversion angelegt, gegen das Schema validiert und durch Referenztests abgesichert. Neue Brennstoffe werden nur ergänzt, wenn Einheit, Umrechnung, Systemgrenze, Quelle, Quellenjahr und Gültigkeitsbeginn eindeutig dokumentiert sind.

## Fachlicher Hinweis

Die Anwendung ist ein Rechen- und Dokumentationswerkzeug, keine amtliche Zertifizierung. Der CO₂-Kostenanteil enthält 19 % Umsatzsteuer; Preisannahmen und offizielle Preiskorridore werden im Ergebnis ausdrücklich unterschieden.
