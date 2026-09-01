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

## Oberflächenkonzept

Die Oberfläche ist als präziser Arbeitsbereich gestaltet: Liefermenge eingeben, Berechnung auslösen, Ergebnis prüfen und bei Bedarf dokumentieren. Der Rechner vermeidet dekorative Flächen und hält den fachlichen Ablauf im Vordergrund.

### Navigation

- Das Menü zeigt anfangs nur Heizöl und Briketts.
- Über „Brennstoff hinzufügen“ können weitere bereits unterstützte Brennstoffe in das Menü aufgenommen werden.
- Einstellungen und Verlauf bleiben als sekundäre, per Tooltip beschriebene Icon-Aktionen erreichbar.

### Eingabe und Berechnung

- Die Liefermenge wird in einem klar abgegrenzten Eingabefeld mit festem Einheitssuffix erfasst.
- Der Hinweis „Beispiel: 1.250 oder 1250,5“ erläutert die akzeptierte deutsche Zahleneingabe.
- Eine Berechnung wird ausschließlich durch „Jetzt berechnen“ oder Enter ausgelöst.
- Änderungen an einer bereits berechneten Eingabe markieren das sichtbare Ergebnis als veraltet; Export, Drucken, Szenarien und Rechenweg werden bis zur erneuten Berechnung deaktiviert.
- Fehler erscheinen direkt am Eingabefeld und enthalten eine konkrete Handlungsempfehlung.

### Ergebnis und Nachweis

- Die Gesamtemissionen in `kg CO₂` sind der visuell wichtigste Wert.
- CO₂-Kosten in `€ brutto`, Energiegehalt in `kWh` und Emissionsintensität in `kg CO₂/kWh` werden kompakt und gleichwertig dargestellt.
- Die Berechnungsgrundlage zeigt Emissionsfaktor, CO₂-Preis, Umsatzsteuersatz und Berechnungsjahr in lesbarer Form.
- Nach einer erfolgreichen Berechnung stehen „Rechenweg“, „Szenarien“, „PDF speichern“ und „Drucken“ als sekundäre Aktionen bereit.

### Gestaltung und Barrierefreiheit

- Lime kennzeichnet Ergebnis und aktive Auswahl; Blau steht für primäre Interaktion und Fokus. Rot wird nur für Fehler verwendet.
- Buttons und Eingabefelder nutzen zurückhaltende Radien und klare Kontraste statt großflächiger Dekoration.
- Alle Aktionen sind per Tastatur erreichbar. Fokus, Fehler, Erfolg und deaktivierte Zustände werden nicht allein über Farbe vermittelt.
- Das Layout wechselt bei schmalen Fenstern von zwei Spalten zu einem vertikalen Ablauf, ohne Inhalte horizontal abzuschneiden.

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
