package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Usun wyprowadza wskazane sesje z historii do kosza.
//
// Pozostaje to jedyną drogą utraty danych sesji w całym produkcie:
// znacznik kosza zdejmuje sesję z historii od ręki, a zapis
// kasuje trwale czyszczenie startowe po terminie kosza (trwalosc_kosza.go)
// — druga faza tego samego usuwania, nie nowa droga utraty. Pomyłkę
// naprawia odwracalność, nie bramka: w oknie terminu sesja
// wraca w całości komendą `session.restore` (adapter_sesje_kosz.go).
// Zamknięcie okna i zamknięcie sesji zmieniają wyłącznie stan — wiadomości,
// okna, artefakty i katalog roboczy zostają nietknięte.
//
// Wskazań może być wiele, bo „usuń zaznaczone" i „usuń jedną" to w historii
// sesji ten sam gest. Wskazanie bez odpowiednika nie jest błędem — wraca
// w `missingIds`; usuwanie zbiorcze nie może paść przez jedną pozycję, której
// ktoś usunął wcześniej z drugiego okna.
func (a *adapterSesji) Usun(ctx context.Context, z shared.SessionDeleteRequest) (shared.SessionDeleteResponse, error) {
	if !z.Confirm {
		return shared.SessionDeleteResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeValidationFailed,
			"usunięcie trwałe wymaga potwierdzenia — pole confirm nie zostało ustawione"))
	}
	wynik := shared.SessionDeleteResponse{DeletedIds: []string{}, MissingIds: []string{}}
	for _, identyfikator := range z.SessionIds {
		a.zwolnijZasobySesji(identyfikator)
		if err := a.usunZapisSesji(ctx, identyfikator); err != nil {
			if errors.Is(err, dane.ErrBrakWiersza) {
				wynik.MissingIds = append(wynik.MissingIds, identyfikator)
				continue
			}
			return shared.SessionDeleteResponse{}, err
		}
		wynik.DeletedIds = append(wynik.DeletedIds, identyfikator)
	}
	wynik.DeletedCount = len(wynik.DeletedIds)
	return wynik, nil
}

// zwolnijZasobySesji zatrzymuje procesy okien i zdejmuje sesję z rejestru
// żywego. Robimy to przed skasowaniem zapisu, żeby nie zostawić biegnącego
// procesu bez sesji, do której należy.
//
// Niepowodzenie nie wstrzymuje usuwania: sesja, której nie ma już w rejestrze
// żywym, i tak ma zniknąć z historii.
func (a *adapterSesji) zwolnijZasobySesji(identyfikator string) {
	okna, err := a.nadzorca.Rejestr().UsunSesje(identyfikator)
	if err != nil {
		return
	}
	_ = a.nadzorca.Procesy().ZatrzymajOkna(identyfikatoryOkienSesji(okna))
}

// usunZapisSesji przenosi wiersz sesji do kosza — zapis zostaje w całości,
// znika wyłącznie z wykazów (dane/sesje.go). Kasowanie fizyczne należy do
// czyszczenia startowego po terminie. Brak utrwalacza znaczy rdzeń bez bazy —
// nie ma wtedy czego usuwać i nie jest to błąd.
func (a *adapterSesji) usunZapisSesji(ctx context.Context, identyfikator string) error {
	if a.trwalosc == nil {
		return nil
	}
	return a.trwalosc.PrzeniesSesjeDoKosza(ctx, identyfikator)
}

// ZapewnieniemSesji nazywa czynność utrwalenia sesji zaraz po jej założeniu.
// Port jest wąski celowo — adapter sesji nie potrzebuje całego utrwalacza
// rozmowy, tylko tej jednej gwarancji.
type ZapewnienieSesji interface {
	ZapewnijSesje(kontekst context.Context, idSesji, tytul, projekt string) (int64, error)
}

// ZZapewnieniem wpina utrwalanie sesji przy zakładaniu. Bez tego portu sesja
// materializuje się w bazie dopiero przy pierwszej wiadomości — a więc sesja
// utworzona i nieużyta ginie po restarcie, a czynności historii nie mają czego
// dotknąć (zapis trwa do ręcznego usunięcia).
func (a *adapterSesji) ZZapewnieniem(z ZapewnienieSesji) *adapterSesji {
	a.zapewnienie = z
	return a
}

// utrwalZalozona zapisuje świeżo założoną sesję. Niepowodzenie nie przerywa
// zakładania: sesja żyje w rejestrze i pracuje, tylko nie przetrwa restartu —
// odmowa założenia byłaby dla Operatora gorsza.
func (a *adapterSesji) utrwalZalozona(kontekst context.Context, sesja session.Sesja) {
	if a.zapewnienie == nil {
		return
	}
	if _, err := a.zapewnienie.ZapewnijSesje(kontekst, sesja.Id, sesja.Tytul, sesja.IdProjektu); err != nil {
		_ = err
	}
}
