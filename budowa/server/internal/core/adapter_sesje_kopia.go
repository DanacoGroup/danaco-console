package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// dopisekKopii odróżnia kopię od źródła w wykazie, gdy Operator nie podał
// własnej nazwy. Bez tego historia pokazywałaby dwie pozycje nie do rozróżnienia.
const dopisekKopii = " (kopia)"

// ZZestawem wpina komplet repozytoriów. Potrzebny wyłącznie kopiowaniu, które
// dotyka trzech tabel naraz — sesji, okien i wiadomości.
func (a *adapterSesji) ZZestawem(z *dane.Zestaw) *adapterSesji {
	a.zestaw = z
	return a
}

// Kopiuj powiela sesję wraz z jej zapisem jako osobny byt historii z własnymi
// identyfikatorami i cyklem życia, niezależny od usunięcia źródła.
func (a *adapterSesji) Kopiuj(ctx context.Context,
	z shared.SessionCopyRequest) (shared.SessionCopyResponse, error) {

	if a.zestaw == nil || a.trwalosc == nil {
		return shared.SessionCopyResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "serwer bez bazy nie ma czego kopiować"))
	}
	zrodlo, err := a.trwalosc.sesje.PoIdentyfikatorze(ctx, z.SessionId)
	if err != nil {
		return shared.SessionCopyResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "sesja "+z.SessionId+" nie istnieje"))
	}

	kopia := a.nadzorca.ZalozSesje(nazwaKopii(zrodlo.Tytul, z.Title), tekstLubPusty(zrodlo.Projekt),
		kontoRejestru(ctx, a.zestaw))
	a.utrwalZalozona(ctx, kopia)

	// Okna bierzemy z rejestru żywego, nie z wierszy, bo obejmuje to także sesje sprzed restartu rdzenia.
	oknaZywe, err := a.nadzorca.Rejestr().OknaSesji(z.SessionId)
	if err != nil {
		return shared.SessionCopyResponse{}, bladSesji(err)
	}
	oknaZrodla, wgOkna, err := a.zalozOknaKopii(ctx, kopia.Id, oknaZywe)
	if err != nil {
		return shared.SessionCopyResponse{}, err
	}
	skopiowanych, err := a.zestaw.KopiujZapisSesji(ctx, oknaZrodla, wgOkna)
	if err != nil {
		return shared.SessionCopyResponse{}, err
	}
	return shared.SessionCopyResponse{
		Session:        sesjaKontraktu(kopia),
		CopiedMessages: skopiowanych,
	}, nil
}

// zalozOknaKopii zakłada w kopii okno za każde okno źródła i zwraca
// odwzorowanie wierszy okna źródłowego na okno docelowe, przenosząc komplet
// parametrów wykonania.
func (a *adapterSesji) zalozOknaKopii(ctx context.Context, idKopii string,
	oknaZywe []session.Okno) ([]dane.Okno, map[int64]int64, error) {

	zrodlowe := make([]dane.Okno, 0, len(oknaZywe))
	wgOkna := make(map[int64]int64, len(oknaZywe))
	for _, zywe := range oknaZywe {
		wiersz, err := a.zestaw.Okna.PoIdentyfikatorze(ctx, zywe.Id)
		if err != nil {
			continue // okno bez wiersza nie ma zapisu do skopiowania
		}
		nowe, err := a.nadzorca.OtworzOkno(idKopii, zywe.Ustawienia)
		if err != nil {
			return nil, nil, bladSesji(err)
		}
		docelowy, err := a.wierszOknaKopii(ctx, nowe, wiersz)
		if err != nil {
			return nil, nil, err
		}
		zrodlowe = append(zrodlowe, wiersz)
		wgOkna[wiersz.ID] = docelowy
	}
	return zrodlowe, wgOkna, nil
}

// wierszOknaKopii zakłada wiersz okna kopii na wzór wiersza źródłowego.
// Identyfikator zewnętrzny bierze się z nowego okna rejestru — inaczej dwa okna
// wskazywałyby ten sam byt rdzenia.
func (a *adapterSesji) wierszOknaKopii(ctx context.Context, nowe session.Okno,
	zrodlowy dane.Okno) (int64, error) {

	wiersz := zrodlowy
	wiersz.ID = 0
	wiersz.IdentyfikatorZewnetrzny = &nowe.Id
	wiersz.OknoKoordynatoraID = nil
	wiersz.Stan = shared.WindowStatusOpen
	if a.trwalosc != nil && a.trwalosc.sesje != nil {
		sesja, err := a.trwalosc.sesje.PoIdentyfikatorze(ctx, nowe.IdSesji)
		if err != nil {
			return 0, err
		}
		wiersz.SesjaID = sesja.ID
	}
	return a.zestaw.Okna.Utworz(ctx, wiersz)
}

// nazwaKopii rozstrzyga nazwę kopii: wskazanie Operatora wygrywa, brak wskazania
// bierze nazwę źródła z dopiskiem.
func nazwaKopii(zrodlo string, wskazana *string) string {
	if wskazana != nil && *wskazana != "" {
		return *wskazana
	}
	return zrodlo + dopisekKopii
}

// tekstLubPusty rozpakowuje wskaźnik kolumny tekstowej, oddając napis pusty
// zamiast wskaźnika pustego.
func tekstLubPusty(wartosc *string) string {
	if wartosc == nil {
		return ""
	}
	return *wartosc
}
