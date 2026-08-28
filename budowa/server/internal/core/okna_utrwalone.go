package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// oknoUtrwalone odnajduje w warstwie utrwalonej okno po identyfikatorze nadanym przez rdzeń i uzupełnia je kodem modułu oraz kanału modelu; brak wiersza nie jest błędem.
func oknoUtrwalone(ctx context.Context, zestaw *dane.Zestaw, idOkna string) (shared.Window, bool, error) {
	wiersz, err := zestaw.Okna.PoIdentyfikatorze(ctx, idOkna)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.Window{}, false, nil
	}
	if err != nil {
		return shared.Window{}, false, err
	}
	moduly, kanaly, err := slownikiOkna(ctx, zestaw)
	if err != nil {
		return shared.Window{}, false, err
	}
	okno := oknoWierszaKontraktu(wiersz, moduly[wiersz.ModulID], kanaly[wiersz.KanalModeluID])
	sesja, err := zestaw.Sesje.Pobierz(ctx, wiersz.SesjaID)
	if err == nil {
		okno.SessionId = identyfikatorWiersza(sesja.IdentyfikatorZewnetrzny, sesja.ID)
	}
	return okno, true, nil
}

// oknaSesjiUtrwalone zwraca okna jednej sesji w kolejności ich zakładania.
// Identyfikator sesji podaje wywołujący, bo zna go z warstwy wyższej — drugi
// odczyt wiersza sesji byłby zapytaniem po to samo.
func oknaSesjiUtrwalone(ctx context.Context, zestaw *dane.Zestaw, sesjaID int64, idSesji string) ([]shared.Window, error) {
	wiersze, err := zestaw.Okna.ListaSesji(ctx, sesjaID)
	if err != nil {
		return nil, err
	}
	moduly, kanaly, err := slownikiOkna(ctx, zestaw)
	if err != nil {
		return nil, err
	}
	okna := make([]shared.Window, 0, len(wiersze))
	for _, wiersz := range wiersze {
		okno := oknoWierszaKontraktu(wiersz, moduly[wiersz.ModulID], kanaly[wiersz.KanalModeluID])
		okno.SessionId = idSesji
		okna = append(okna, okno)
	}
	return okna, nil
}

// slownikiOkna buduje odwzorowania klucza obcego na kod dla modułu i kanału modelu, potrzebne do przekładu wiersza okna na kontrakt.
func slownikiOkna(ctx context.Context, zestaw *dane.Zestaw) (map[int64]string, map[int64]string, error) {
	wiersze, err := zestaw.Moduly.Lista(ctx)
	if err != nil {
		return nil, nil, err
	}
	moduly := make(map[int64]string, len(wiersze))
	for _, wiersz := range wiersze {
		moduly[wiersz.ID] = wiersz.Kod
	}
	kanalyWierszy, err := zestaw.Kanaly.Lista(ctx, false)
	if err != nil {
		return nil, nil, err
	}
	kanaly := make(map[int64]string, len(kanalyWierszy))
	for _, wiersz := range kanalyWierszy {
		kanaly[wiersz.ID] = wiersz.Kod
	}
	return moduly, kanaly, nil
}

// wybraneOkna zawęża okna do wskazanych przez klienta. Wykaz pusty obejmuje
// wszystkie okna sesji, zgodnie z kontraktem `session.bind`.
func wybraneOkna(okna []shared.Window, wskazane []string) []shared.Window {
	if len(wskazane) == 0 {
		return okna
	}
	zbior := make(map[string]struct{}, len(wskazane))
	for _, id := range wskazane {
		zbior[id] = struct{}{}
	}
	wybrane := make([]shared.Window, 0, len(wskazane))
	for _, okno := range okna {
		if _, jest := zbior[okno.Id]; jest {
			wybrane = append(wybrane, okno)
		}
	}
	return wybrane
}
