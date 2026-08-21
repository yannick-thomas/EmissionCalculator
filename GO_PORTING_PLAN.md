# Go-Portierung: EmissionCalculator

## Ziel
Die bestehende Python-Desktopanwendung in Go zu portieren, ohne die fachlichen Funktionen zu verlieren:
- Berechnung von Briketts
- Berechnung von Heizöl
- Anzeige von Emissionen, CO2-Kosten und Energiegehalt
- Druck/PDF-Ausgabe

## Projektstatus der aktuellen Version
Die aktuelle App ist ein kleines Desktop-Tool mit:
- Tkinter/CustomTkinter GUI
- Berechnungslogik in `calculator/calc`
- Renderer in `calculator/renderer`
- PDF-Druck in `calculator/misc/print_handler.py`

## Ausgangslage
Die Berechnungslogik ist klein und gut isoliert. Der größte Aufwand liegt bei der GUI-Umsetzung und beim PDF-Export.

## Empfohlene Technologie
- Sprache: Go
- UI: Fyne
- PDF/Druck: gofpdf oder vergleichbares Paket
- Build: Go modules
- Tests: `go test` mit Unit-Tests für die Berechnungslogik

## Grundsatz
Die Portierung soll sauber und modular erfolgen. Es darf keine Python-Logik in Go als „Listen-Resultate“ weitergeführt werden. Stattdessen sollen echte Datenstrukturen und klar getrennte Pakete verwendet werden.

---

# Architektur

## Paketstruktur

```text
app/
  main.go
  app.go
  navigation.go

internal/
  calculation/
    calc.go
    briquettes.go
    oil.go
    formatting.go
  models/
    result.go
  validation/
    input.go
  ui/
    root.go
    sidebar.go
    briquettes_view.go
    oil_view.go
    shared_result.go
  pdf/
    export.go
    template.go
```

## Funktionsaufteilung

### app/
- startet die Go-App
- initialisiert Fyne
- verwaltet Fenster, Theme und Navigation

### internal/calculation/
- berechnet alle Formeln
- enthält gemeinsame Formatierung
- kapselt die Logik für Briketts und Heizöl

### internal/models/
- definiert Datenmodelle
- enthält `CalculationResult`

### internal/validation/
- parst und validiert Inputs
- akzeptiert Komma und Punkt
- liefert saubere Fehler

### internal/ui/
- rendert Sidebar und Ansichten
- zeigt Eingaben und Ergebnisse an
- orchestriert Berechnung und Druck

### internal/pdf/
- erzeugt PDF-Ausgabe
- stellt den Druck-Label-Inhalt zusammen

---

# Datenmodelle

## `CalculationResult`

```go
type CalculationResult struct {
    Valid            bool
    Emissions       string
    EmissionCost    string
    EnergyContent    string
    CO2PerKWh       string
    ErrorMessage    string
}
```

## Warum dieses Modell?
Das verhindert Python-ähnliche Listen mit Positionslogik wie:
- result[0]
- result[1]
- result[2]

In Go sollten Werte immer semantisch benannt sein.

---

# Berechnungslogik

## Gemeinsame Funktionen

### `ParseQuantity(input string) (float64, error)`
- akzeptiert `1,5` und `1.5`
- validiert leere/ungültige Eingaben
- liefert Fehler bei Parsing-Fehlern

### `FormatNumber(value float64, energyMode bool) string`
- formatiert Zahlen mit Komma/Dezimaltrennung
- liefert bei Energiegehalt passende Ausgabe

---

## Briketts-Formel
Aus der aktuellen Python-Version:

```python
emissions = 19 * 0.0992 * 1000 * quantity
emission_component_result = emissions * 45 * 1.19 / 1000
energy_content = 19 * 277.778 * quantity
```

Diese Formeln werden in Go als Funktionen umgesetzt:

```go
func CalculateBriquettes(quantity float64) CalculationResult
```

Erwartete Werte:
- Emissions: `kg CO2`
- Emission cost: `€`
- Energy content: `kWh`
- CO2 per kWh: `0,3571`

---

## Heizöl-Formel
Aus der aktuellen Python-Version:

```python
emissions = 42.8 * 0.074 * 0.845 * quantity
emission_component_result = round(emissions, 2) / 1000 * 45 * 1.19
energy_content = 42.8 * 0.845 / 1000 * quantity * 277.778
```

Diese Formeln werden in Go als Funktionen umgesetzt:

```go
func CalculateOil(quantity float64) CalculationResult
```

Erwartete Werte:
- Emissions: `kg CO2`
- Emission cost: `€`
- Energy content: `kWh`
- CO2 per kWh: `0,2664`

---

# UI-Architektur

## App-Start
- Fyne App initialisieren
- Fenster mit Titel „Emissionsrechner“
- feste Größe oder responsive Layout

## Layout
- linke Sidebar
- Hauptbereich rechts
- zwei Seiten: Briketts und Heizöl

## Seitenkomponenten
### Sidebar
- Button: Briketts
- Button: Heizöl

### Briketts-View
- Label: Liefermenge
- Input-Feld
- Einheit: `t`
- Button: Berechnen
- Ergebnisbereich mit:
  - Brennstoffemissionen
  - Preisbestandteil CO2
  - Energiegehalt
- Button: Drucken

### Heizöl-View
- Label: Liefermenge
- Input-Feld
- Einheit: `l`
- Button: Berechnen
- Ergebnisbereich mit:
  - Brennstoffemissionen
  - Preisbestandteil CO2
  - Energiegehalt
- Button: Drucken

## UI-Best Practice
Keine UI-Logik direkt in Berechnungspaketen. UI soll nur:
- Eingaben lesen
- Daten an Rechenfunktionen übergeben
- Ergebnisse rendern

---

# PDF-/Druck-Architektur

## Anforderungen
- PDF-Datei erzeugen
- Ergebnisdaten als Text einbinden
- großes Label/Druckformat wie in der aktuellen App

## Umsetzung
### `internal/pdf/export.go`
- erzeugt Datei im System-Temp-Verzeichnis
- nutzt `gofpdf` oder ein vergleichbares Paket
- schreibt das fertige Dokument

### `internal/pdf/template.go`
- sammelt die Texte zusammen
- formuliert den Inhalt des Labels

## Beispiel-Text
- `CO2-Abgabe ... EUR, Brutto`
- `CO2 kg der Lieferung ...`
- `CO2 kg pro kWh ...`
- `kWh der Lieferung ...`

---

# Fehlerbehandlung

## Ziele
- ungültige Eingaben sauber verarbeiten
- keine unsicheren `except:`-Blöcke mehr
- klare Meldungen an den Nutzer

## Regeln
- Leerstring = Fehler
- Nicht-numerische Werte = Fehler
- Negative Werte = Fehler
- Bei Fehler: Meldung im UI anzeigen

## Beispiel
```go
if quantity <= 0 {
    return CalculationResult{Valid: false, ErrorMessage: "Bitte gültige Liefermenge eingeben"}
}
```

---

# Teststrategie

## Muss-Tests
1. Briketts-Berechnung mit gültiger Eingabe
2. Heizöl-Berechnung mit gültiger Eingabe
3. negative/ungültige Eingabe
4. Formatierung mit Komma
5. Formatierung mit Energiegehalt

## Beispiel-Testcases
```go
func TestCalculateBriquettes(t *testing.T)
func TestCalculateOil(t *testing.T)
func TestParseQuantity(t *testing.T)
func TestFormatNumber(t *testing.T)
```

---

# Migrationsplan

## Phase 1: Projektgrundlage
- Go-Modul starten
- Abhängigkeiten einrichten
- Fyne initialisieren
- grundlegendes Fenster testen

## Phase 2: Berechnungslogik
- `internal/calculation` anlegen
- Formeln überführen
- Unit-Tests schreiben

## Phase 3: UI-Struktur
- Sidebar bauen
- Briketts-View erstellen
- Heizöl-View erstellen
- Ergebnisdarstellung übernehmen

## Phase 4: PDF-Export
- Druckfunktion umsetzen
- Temp-Datei erzeugen
- Ausgabe prüfen

## Phase 5: Verfeinerung
- Layout optimieren
- Fehlerzustände sauber machen
- App-Name/Icons/Build einstellen

---

# Risiken

## 1. UI-Layout stimmt nicht exakt mit Python-Version überein
Das ist wahrscheinlich. Die Go-Version sollte bewusst ein modernes, sauberes Layout statt pixelgenaues Replikat nutzen.

## 2. PDF-Format unterscheidet sich leicht
Das ist akzeptabel. Hauptsache die Inhalte und Funktionalität bleiben erhalten.

## 3. Fehlerbehandlung bleibt ungenau
Das soll durch saubere Validierungslogik und Tests vermieden werden.

---

# Empfehlung
Die Portierung ist sinnvoll, gut machbar und nicht riesig. Der echte Vorteil von Go liegt nicht darin, die Python-Logik blind zu kopieren, sondern die App sauberer, robuster und wartbarer neu zu strukturieren.

## Endziel
Eine moderne Go-Desktop-App mit:
- sauberer Paketstruktur
- guten Unit-Tests
- klarer UI-Aufteilung
- funktionalem PDF-Druck
- besserer Fehlerbehandlung als in Python

---

# TODO-Liste

## Muss erledigt werden
- [ ] Go-Modul initialisieren
- [ ] Fyne als UI-Abhängigkeit hinzufügen
- [ ] Berechnungslogik in Go portieren
- [ ] Ergebnisse als `CalculationResult` modellieren
- [ ] UI-Ansichten für Briketts und Heizöl bauen
- [ ] Validierung für Eingaben hinzufügen
- [ ] PDF-Export implementieren
- [ ] Tests für Berechnungen schreiben
- [ ] App-Build für Zielplattformen vorbereiten

## Optional / später
- [ ] App-Icon
- [ ] Installer
- [ ] Dark/Light Theme
- [ ] zusätzliche Fehler- und Zustandsanzeigen

---

# Abschluss
Diese Datei dient als offizieller Anker für die Übersetzung in Go. Sie sollte bei der Portierung als Referenz und Arbeitsplan verwendet werden.
