// Odczyt stanu bramki na zewnątrz: rozpoznanie sesji z powitania,
// odpowiedź „czy bramka w ogóle jest”, przekłady wierszy na kształt
// kontraktu i wspólna postać odmowy rodziny `auth.*`.
package core

import (
	"encoding/json"

	"context"
	"errors"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ── rozpoznanie sesji i stan bramki ──────────────────────────────────────────

// BramkaZalozona mówi, czy sekret bramki istnieje. Brak kotwicy nie jest
// błędem — jest odpowiedzią „nie" i tak wychodzi do powitania.
func (a *adapterUwierzytelnienia) BramkaZalozona(ctx context.Context) (bool, error) {
	if a == nil || a.repozytorium == nil {
		return false, bladBramki(shared.ErrorCodeInternalError, "Magazyn kont jest niedostępny.")
	}
	if _, err := a.kotwica(ctx); err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/*
KontoSesjiBramki oddaje konto, któremu wydano sesję o podanym skrócie. Zero
znaczy sesję nieznaną albo wiersz sprzed rozdzielenia kont — wołający czyta
wtedy pracę konta najstarszego, tak jak stała przed migracją 407.
*/
func (a *adapterUwierzytelnienia) KontoSesjiBramki(ctx context.Context,
	skrotTokenu string) (int64, error) {

	sesja, err := a.repozytorium.SesjaBramkiPoSkrocie(ctx, skrotTokenu)
	if err != nil {
		return 0, err
	}
	return sesja.KontoId, nil
}

// RozpoznajSesjeBramki sprawdza token z powitania i oddaje skrót sesji.
// Token nieznany, unieważniony i wygasły znaczą to samo.
func (a *adapterUwierzytelnienia) RozpoznajSesjeBramki(ctx context.Context,
	token string) (string, bool, error) {

	if a == nil || a.repozytorium == nil {
		return "", false, bladBramki(shared.ErrorCodeInternalError, "Magazyn kont jest niedostępny.")
	}
	if token == "" {
		return "", false, nil
	}
	skrot := skrotTokenu(token)
	sesja, err := a.repozytorium.SesjaBramkiPoSkrocie(ctx, skrot)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if sesja.Uniewazniono != nil || sesja.Wygasa <= time.Now().UnixMilli() {
		return "", false, nil
	}
	return skrot, true, nil
}

// bladBramki składa odmowę rodziny `auth.*` z jednym przedrostkiem, żeby
// Errors Panel pokazywał, czyja to odmowa, bez zgadywania po treści.
func bladBramki(kod protocol.KodBledu, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, powod))
}

// Powody odmowy rozróżniane maszynowo. Sam kod kontraktu bywa za szeroki:
// „conflict" na rejestracji znaczy albo że konto Operatora już stoi, albo że
// podany login lub adres należy do konta istniejącego. Okno musi te dwa
// przypadki nazwać inaczej, a nie ma jak ich rozpoznać po treści komunikatu,
// bo treść od serwera do okna nie dociera.
const (
	// PowodKontoIstnieje oznacza, że konto Operatora zostało już założone.
	PowodKontoIstnieje = "konto-istnieje"
	// PowodKolizjaDanych oznacza, że podany login albo adres jest już zajęty.
	PowodKolizjaDanych = "kolizja-danych"
	/* PowodAdresNiepotwierdzony oznacza konto założone, którego adresu nikt nie
	   potwierdził. Okno ma po nim rozpoznać, że hasło jest dobre, a brakuje
	   wyłącznie kodu z listu — i zaprowadzić Operatora na odsłonę kodu zamiast
	   mówić mu, że danych nie rozpoznano. */
	PowodAdresNiepotwierdzony = "adres-niepotwierdzony"
)

// bladBramkiZPowodem dokłada do odmowy maszynowy powód w polu `details`.
// Pole niesie wartość rozpoznawaną przez okno, nie zdanie do wyświetlenia.
func bladBramkiZPowodem(kod protocol.KodBledu, rozpoznanie string, powod string) error {
	blad := protocol.NowyBlad(kod, powod)
	/* Kontrakt trzyma szczegóły jako surowy JSON, więc powód idzie obiektem
	   o jednym polu — okno czyta wartość, nie rozbiera zdania. */
	blad.Details = json.RawMessage(`{"powod":"` + rozpoznanie + `"}`)
	return protocol.JakoError(blad)
}

// metodaBramkiKontraktu przekłada wiersz na `AuthMethod`. Odwołania do sejfu
// w wyniku nie ma: kontrakt mówi wprost, że odpowiedź nie niesie ani skrótu
// hasła, ani skrótu PIN-u, ani materiału klucza.
func metodaBramkiKontraktu(m dane.MetodaUwierzytelnienia) shared.AuthMethod {
	return shared.AuthMethod{
		Id:         m.Kod,
		Kind:       shared.AuthMethodKind(m.Rodzaj),
		Label:      m.Etykieta,
		DeviceId:   m.UrzadzenieKod,
		DeviceName: m.NazwaUrzadzenia,
		Anchor:     m.Kotwica,
		CreatedAt:  m.Utworzono,
		LastUsedAt: m.OstatnioUzyto,
	}
}

// sesjaBramkiKontraktu przekłada wiersz sesji na `AuthSession`. Token przychodzi
// osobno, bo wiersz trzyma wyłącznie jego skrót.
func sesjaBramkiKontraktu(token string, s dane.SesjaBramki) shared.AuthSession {
	sesja := shared.AuthSession{
		Token:     token,
		ExpiresAt: s.Wygasa,
		DeviceId:  s.UrzadzenieKod,
	}
	if s.RodzajMetody != nil {
		rodzaj := shared.AuthMethodKind(*s.RodzajMetody)
		sesja.Method = &rodzaj
	}
	return sesja
}

// niepustyTekst przepuszcza wskaźnik dalej tylko wtedy, gdy niesie treść.
// Napis pusty i brak wskazania znaczą dla bramki to samo, więc dwie postaci
// braku byłyby rozróżnieniem bez treści.
func niepustyTekst(wartosc *string) *string {
	if wartosc == nil || *wartosc == "" {
		return nil
	}
	return wartosc
}
