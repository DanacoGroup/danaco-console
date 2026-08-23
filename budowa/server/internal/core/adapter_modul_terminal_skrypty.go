// Komendy `terminal.script.save`, `terminal.script.list` i
// `terminal.script.remove` — biblioteka skryptów i snippetów okna Script Library.
//
// ── Czego brakowało ─────────────────────────────────────────────────────────
// Biblioteka żyła jedno posiedzenie, a jedyną drogą jej zachowania był wywóz do
// pliku. Skrypt uruchamiany na maszynach Operatora jest treścią, do której trzeba
// móc wrócić, więc każdy zapis zakłada KOLEJNĄ WERSJĘ, a nie nadpisuje
// poprzedniej (migracja 247).
//
// ── Numer wersji nadaje rdzeń, nie żądanie ──────────────────────────────────
// Kontrakt mówi to wprost, a powód jest współbieżnościowy: numer odczytany przed
// zapisem rozjechałby się przy dwóch zapisach naraz. Numer nadaje więc baza
// w jednej transakcji z wpisem wersji (`dane/terminal_wyposazenie_zapis.go`).
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekSkryptu znakuje identyfikator pozycji biblioteki.
const przedrostekSkryptu = "tscr-"

// ZapiszSkrypt obsługuje `terminal.script.save`.
func (a *adapterTerminala) ZapiszSkrypt(ctx context.Context,
	z shared.TerminalScriptSaveRequest) (shared.TerminalScriptSaveResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalScriptSaveResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Script.Name)
	if nazwa == "" {
		return shared.TerminalScriptSaveResponse{}, bladZadaniaTerminala(
			"pozycja biblioteki wymaga nazwy")
	}
	if z.Script.Shell != "" && !CzyPowlokaZnana(z.Script.Shell) {
		return shared.TerminalScriptSaveResponse{}, bladZadaniaTerminala(
			"powłoka " + string(z.Script.Shell) + " nie należy do wykazu wykonawczego rdzenia")
	}
	rodzaj := z.Script.Kind
	if rodzaj == "" {
		rodzaj = shared.TerminalScriptKindScript
	}

	kod := strings.TrimSpace(z.Script.Id)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekSkryptu)
	}
	wersja, powstala, err := dziennik.ZapiszSkrypt(ctx, dane.SkryptTerminala{
		Kod:       kod,
		Nazwa:     nazwa,
		Rodzaj:    rodzaj,
		Powloka:   z.Script.Shell,
		Tresc:     z.Script.Content,
		Znaczniki: strings.TrimSpace(wartoscTekstu(z.Script.Tags)),
		Alias:     wskaznikTekstu(strings.TrimSpace(wartoscTekstu(z.Script.Alias))),
	})
	if err != nil {
		return shared.TerminalScriptSaveResponse{}, err
	}
	zapisana, err := dziennik.Skrypt(ctx, kod)
	if err != nil {
		return shared.TerminalScriptSaveResponse{}, err
	}
	zapisana.Wersja = wersja
	return shared.TerminalScriptSaveResponse{
		Script:  skryptKontraktu(zapisana),
		Created: powstala,
	}, nil
}

// WykazSkryptow obsługuje `terminal.script.list`. Oddaje wersję najnowszą każdej
// pozycji — historia wersji jest dziennikiem, nie treścią wykazu.
func (a *adapterTerminala) WykazSkryptow(ctx context.Context,
	z shared.TerminalScriptListRequest) (shared.TerminalScriptListResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalScriptListResponse{}, err
	}
	filtr := dane.FiltrSkryptowTerminala{
		Znacznik: strings.TrimSpace(wartoscTekstu(z.Tag)),
		Fraza:    strings.TrimSpace(wartoscTekstu(z.Query)),
	}
	if z.Kind != nil {
		filtr.Rodzaj = *z.Kind
	}
	wiersze, err := dziennik.Skrypty(ctx, filtr)
	if err != nil {
		return shared.TerminalScriptListResponse{}, err
	}
	wykaz := make([]shared.TerminalScript, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, skryptKontraktu(wiersz))
	}
	return shared.TerminalScriptListResponse{Scripts: wykaz, Total: len(wykaz)}, nil
}

// UsunSkrypt obsługuje `terminal.script.remove` — usuwa pozycję wraz ze
// wszystkimi jej wersjami (kasowanie kaskadowe, migracja 247).
func (a *adapterTerminala) UsunSkrypt(ctx context.Context,
	z shared.TerminalScriptRemoveRequest) (shared.TerminalScriptRemoveResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalScriptRemoveResponse{}, err
	}
	kod := strings.TrimSpace(z.ScriptId)
	if kod == "" {
		return shared.TerminalScriptRemoveResponse{}, bladZadaniaTerminala(
			"usunięcie pozycji wymaga jej wskazania")
	}
	usunieta, err := dziennik.UsunSkrypt(ctx, kod)
	if err != nil {
		return shared.TerminalScriptRemoveResponse{}, err
	}
	return shared.TerminalScriptRemoveResponse{Removed: usunieta}, nil
}

// skryptKontraktu przekłada wiersz biblioteki na byt kontraktu.
func skryptKontraktu(w dane.SkryptTerminala) shared.TerminalScript {
	pozycja := shared.TerminalScript{
		Id:        w.Kod,
		Name:      w.Nazwa,
		Kind:      w.Rodzaj,
		Shell:     w.Powloka,
		Content:   w.Tresc,
		Tags:      wskaznikTekstu(w.Znaczniki),
		Alias:     w.Alias,
		Version:   int(w.Wersja),
		CreatedAt: chwilaBazy(w.Utworzono),
		UpdatedAt: chwilaBazy(w.Zaktualizowano),
	}
	if w.Uruchomiono != nil {
		chwila := chwilaBazy(*w.Uruchomiono)
		pozycja.LastRunAt = &chwila
	}
	return pozycja
}

// skryptZDziennika odnajduje pozycję biblioteki albo odmawia jej brakiem.
func (a *adapterTerminala) skryptZDziennika(ctx context.Context,
	kod string) (dane.SkryptTerminala, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return dane.SkryptTerminala{}, err
	}
	pozycja, err := dziennik.Skrypt(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.SkryptTerminala{}, bladBrakuZasobuTerminala("pozycja biblioteki " + kod)
	}
	return pozycja, err
}
