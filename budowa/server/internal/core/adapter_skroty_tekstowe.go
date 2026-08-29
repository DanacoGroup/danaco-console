// Plik wypełnia rodzinę komend snippet.* przechowującą słownik skrótów tekstowych rozwijanych w
// treść dłuższą, wspólny dla wszystkich pól tekstowych platformy i niezależny od maszyny.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const przedrostekSkrotuTekstowego = "skrot-"

// adapterSkrotowTekstowych wypełnia port SkrotyTekstowe, wiążąc go z repozytorium skrótów
// przechowywanym w bazie danych rdzenia.
type adapterSkrotowTekstowych struct {
	repozytorium dane.RepozytoriumSkrotow
}

// nowyAdapterSkrotowTekstowych wiąże port ze słownikiem skrótów, tworząc adapter gotowy do obsługi
// komend rodziny snippet.
func nowyAdapterSkrotowTekstowych(repozytorium dane.RepozytoriumSkrotow) *adapterSkrotowTekstowych {
	return &adapterSkrotowTekstowych{repozytorium: repozytorium}
}

// WykazSkrotow obsługuje snippet.list — zwraca skróty pasujące do zapytania wraz z łączną liczbą
// pozycji w słowniku.
func (a *adapterSkrotowTekstowych) WykazSkrotow(ctx context.Context,
	z shared.SnippetListRequest) (shared.SnippetListResponse, error) {

	if a.repozytorium == nil {
		return shared.SnippetListResponse{}, bladZapleczaSkrotow()
	}
	skroty, wszystkich, err := a.repozytorium.SkrotyTekstowe(ctx,
		wartoscTekstu(z.Query), wartoscTekstu(z.ProfileId), wartoscLiczby(z.Limit))
	if err != nil {
		return shared.SnippetListResponse{}, bladMagazynuSkrotow(err)
	}
	wykaz := make([]shared.TextSnippet, 0, len(skroty))
	for _, skrot := range skroty {
		wykaz = append(wykaz, skrotKontraktu(skrot))
	}
	return shared.SnippetListResponse{Snippets: wykaz, Total: wszystkich}, nil
}

// ZapiszSkrot obsługuje snippet.set — zapisuje nowy skrót albo aktualizuje istniejący, zwracając
// zapisaną pozycję słownika.
func (a *adapterSkrotowTekstowych) ZapiszSkrot(ctx context.Context,
	z shared.SnippetSetRequest) (shared.SnippetSetResponse, error) {

	if a.repozytorium == nil {
		return shared.SnippetSetResponse{}, bladZapleczaSkrotow()
	}
	skrot := strings.TrimSpace(z.Shortcut)
	if skrot == "" {
		return shared.SnippetSetResponse{}, bladWskazaniaSkrotu(
			"skrót bez frazy wyzwalającej — nie ma czego wpisać, żeby się rozwinął")
	}
	if strings.ContainsAny(skrot, " \t\n") {
		// Skrót ze spacją nie rozwinąłby się: rozpoznawanie działa po jednym słowie.
		return shared.SnippetSetResponse{}, bladWskazaniaSkrotu(
			"skrót „" + skrot + "” zawiera odstęp — rozwijanie rozpoznaje skrót po " +
				"jednym słowie, więc fraza z odstępem nigdy by się nie rozwinęła")
	}
	if strings.TrimSpace(z.Content) == "" {
		return shared.SnippetSetResponse{}, bladWskazaniaSkrotu(
			"skrót bez treści rozwinięcia — rozwinąłby się w nic")
	}

	teraz := time.Now().UnixMilli()
	kod := strings.TrimSpace(wartoscTekstu(z.SnippetId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekSkrotuTekstowego)
	}
	wiersz := dane.SkrotTekstowy{
		Kod:            kod,
		ProfilKod:      strings.TrimSpace(wartoscTekstu(z.ProfileId)),
		Skrot:          skrot,
		Tresc:          z.Content,
		Opis:           z.Description,
		PolaJSON:       zapisPolSzablonu(z.Variables),
		Czynny:         z.Enabled == nil || *z.Enabled,
		Utworzono:      teraz,
		Zaktualizowano: teraz,
	}
	zapisany, err := a.repozytorium.ZapiszSkrotTekstowy(ctx, wiersz)
	if err != nil {
		return shared.SnippetSetResponse{}, bladMagazynuSkrotow(err)
	}
	return shared.SnippetSetResponse{Snippet: skrotKontraktu(zapisany)}, nil
}

// UsunSkrot obsługuje snippet.delete — usuwa wskazany skrót ze słownika i zwraca informację, czy
// pozycja istniała.
func (a *adapterSkrotowTekstowych) UsunSkrot(ctx context.Context,
	z shared.SnippetDeleteRequest) (shared.SnippetDeleteResponse, error) {

	if a.repozytorium == nil {
		return shared.SnippetDeleteResponse{}, bladZapleczaSkrotow()
	}
	usuniety, err := a.repozytorium.UsunSkrotTekstowy(ctx, strings.TrimSpace(z.SnippetId))
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.SnippetDeleteResponse{Deleted: false}, nil
		}
		return shared.SnippetDeleteResponse{}, bladMagazynuSkrotow(err)
	}
	return shared.SnippetDeleteResponse{Deleted: usuniety}, nil
}

// zapisPolSzablonu składa nazwy pól szablonu w zapis strukturalny kolumny bazy, zwracając pustą
// tablicę dla braku pól.
func zapisPolSzablonu(pola []string) string {
	if len(pola) == 0 {
		return "[]"
	}
	bajty, err := json.Marshal(pola)
	if err != nil {
		return "[]"
	}
	return string(bajty)
}

// odczytPolSzablonu rozkłada zapis strukturalny kolumny bazy na nazwy pól szablonu, zwracając
// pustą wartość dla zapisu bez treści.
func odczytPolSzablonu(zapis string) []string {
	if strings.TrimSpace(zapis) == "" {
		return nil
	}
	var pola []string
	if err := json.Unmarshal([]byte(zapis), &pola); err != nil {
		return nil
	}
	return pola
}

// skrotKontraktu przekłada wiersz repozytorium na pozycję słownika kontraktu, dołączając kod
// profilu tylko dla skrótu do niego przypisanego.
func skrotKontraktu(s dane.SkrotTekstowy) shared.TextSnippet {
	pozycja := shared.TextSnippet{
		Id: s.Kod, Shortcut: s.Skrot, Content: s.Tresc, Description: s.Opis,
		Variables: odczytPolSzablonu(s.PolaJSON), Enabled: s.Czynny,
		CreatedAt: s.Utworzono, UpdatedAt: s.Zaktualizowano,
	}
	if strings.TrimSpace(s.ProfilKod) != "" {
		profil := s.ProfilKod
		pozycja.ProfileId = &profil
	}
	return pozycja
}

// bladZapleczaSkrotow nazywa brak wpiętego słownika skrótów po stronie rdzenia i wskazuje sposób
// naprawy.
func bladZapleczaSkrotow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"skróty tekstowe: serwer nie ma wpiętego słownika — naprawa: podpiąć "+
			"repozytorium skrótów przy składaniu serwera"))
}

// bladWskazaniaSkrotu nazywa niepoprawne żądanie dotyczące skrótu, na przykład skrót bez frazy
// wyzwalającej.
func bladWskazaniaSkrotu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"skróty tekstowe: "+powod))
}

// bladMagazynuSkrotow nazywa niepowodzenie zapisu albo odczytu.
//
// Kolizja pary (profil, skrót) jest tu odmową walidacji, a nie awarią: to
// Operator próbuje przypisać drugie rozwinięcie temu samemu skrótowi, a nie
// rdzeń zawodzi.
func bladMagazynuSkrotow(err error) error {
	if strings.Contains(strings.ToLower(err.Error()), "unique") {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"skróty tekstowe: ten skrót jest już w słowniku tego profilu — "+
				"jeden skrót ma jedno rozwinięcie, inaczej rozstrzygałaby o nim kolejność odczytu"))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"skróty tekstowe: "+err.Error()))
}
