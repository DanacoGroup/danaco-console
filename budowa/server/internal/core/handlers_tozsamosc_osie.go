// Plik odczytuje osie rozstrzygania okna rozmowy: modelu i konta, dla których liczona jest
// nakładka tego okna, na podstawie kanału modelu wskazanego przez okno.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
)

// wskazanieOsiOkna odczytuje osie rozstrzygania jednego okna rozmowy tej platformy dla konta użytkownika.
type wskazanieOsiOkna struct {
	okna   dane.RepozytoriumOkien
	kanaly dane.RepozytoriumKanalow
}

// nowaWskazanieOsiOkna wiąże odczyt osi rozstrzygania z repozytoriami okien komunikacji i kanałów modelu.
func nowaWskazanieOsiOkna(okna dane.RepozytoriumOkien,
	kanaly dane.RepozytoriumKanalow) *wskazanieOsiOkna {

	return &wskazanieOsiOkna{okna: okna, kanaly: kanaly}
}

// dlaOkna zwraca osie rozstrzygania okna komunikacji wskazanego identyfikatorem rdzenia tej platformy.
func (w *wskazanieOsiOkna) dlaOkna(ctx context.Context, idOkna string) (ZapytanieTozsamosci, error) {
	if w == nil || w.okna == nil {
		return ZapytanieTozsamosci{}, nil
	}
	okno, err := w.okna.PoIdentyfikatorze(ctx, idOkna)
	if brakWiersza(err) {
		return ZapytanieTozsamosci{}, nil
	}
	if err != nil {
		return ZapytanieTozsamosci{}, err
	}
	return w.osieKanalu(ctx, okno.KanalModeluID)
}

// osieKanalu odczytuje model i konto kanału. Rejestr kanałów jest wykazem
// kilkunastu wierszy czytanym w całości — osobnego odczytu po numerze wiersza
// repozytorium nie ma i nie potrzebuje.
func (w *wskazanieOsiOkna) osieKanalu(ctx context.Context, kanalID int64) (ZapytanieTozsamosci, error) {
	if w.kanaly == nil || kanalID == 0 {
		return ZapytanieTozsamosci{}, nil
	}
	kanaly, err := w.kanaly.Lista(ctx, false)
	if err != nil {
		return ZapytanieTozsamosci{}, err
	}
	for _, kanal := range kanaly {
		if kanal.ID != kanalID {
			continue
		}
		zapytanie := ZapytanieTozsamosci{Model: kanal.IdentyfikatorModelu}
		if kanal.KontoID != nil {
			zapytanie.Konto = strconv.FormatInt(*kanal.KontoID, 10)
		}
		return zapytanie, nil
	}
	return ZapytanieTozsamosci{}, nil
}
