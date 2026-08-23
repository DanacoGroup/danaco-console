// Odpowiedzialność pliku: rodzina `snippet.*` — słownik skrótów tekstowych
// rozwijanych w dłuższą treść.
//
// ── Dlaczego słownik należy do rdzenia ──────────────────────────────────────
// Skrót rozwija się we WSZYSTKICH polach tekstowych platformy, a nie w jednym
// oknie. Gdyby mieszkał w kliencie, ta sama fraza rozwijałaby się inaczej na
// dwóch maszynach tego samego Operatora, a przeniesienie pracy na inną maszynę
// oznaczałoby przepisywanie słownika od nowa.
//
// Samo ROZWINIĘCIE robi okno, u siebie, w chwili pisania — rdzeń nie widzi
// pola tekstowego i widzieć go nie musi. Rdzeń trzyma słownik i pilnuje, żeby
// jeden profil nie miał dwóch rozwinięć tego samego skrótu.
//
// ── Pola szablonu ───────────────────────────────────────────────────────────
// Skrót bywa szablonem z polami do wypełnienia (`variables`). Rdzeń zna ich
// NAZWY, ale ich nie wypełnia: wypełnia je Operator w chwili rozwinięcia, a
// wartości bywają różne przy każdym użyciu. Podstawienie czegokolwiek po
// stronie rdzenia dałoby szablon rozwinięty raz na zawsze.
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

// adapterSkrotowTekstowych wypełnia port `SkrotyTekstowe`.
type adapterSkrotowTekstowych struct {
	repozytorium dane.RepozytoriumSkrotow
}

// nowyAdapterSkrotowTekstowych wiąże port ze słownikiem skrótów.
func nowyAdapterSkrotowTekstowych(repozytorium dane.RepozytoriumSkrotow) *adapterSkrotowTekstowych {
	return &adapterSkrotowTekstowych{repozytorium: repozytorium}
}

// WykazSkrotow obsługuje `snippet.list`.
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

// ZapiszSkrot obsługuje `snippet.set`.
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
		// Skrót ze spacją nie rozwinąłby się nigdy: okno rozpoznaje skrót po
		// jednym słowie, a fraza wieloczłonowa jest zwykłym zdaniem.
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

// UsunSkrot obsługuje `snippet.delete`.
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

// zapisPolSzablonu składa nazwy pól w zapis strukturalny kolumny.
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

// odczytPolSzablonu rozkłada zapis strukturalny kolumny na nazwy pól.
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

// skrotKontraktu przekłada wiersz na pozycję słownika kontraktu.
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

// bladZapleczaSkrotow nazywa brak słownika po stronie rdzenia.
func bladZapleczaSkrotow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"skróty tekstowe: rdzeń nie ma wpiętego słownika — naprawa: podpiąć "+
			"repozytorium skrótów przy składaniu rdzenia"))
}

// bladWskazaniaSkrotu nazywa niepoprawne żądanie.
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
