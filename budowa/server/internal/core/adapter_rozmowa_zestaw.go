// Plik to miejsce, w którym zestaw narzędzi tury przestaje być stały: tu spotykają się podstawa eksperta i doraźne dołożenia sesji, wchodząc do wpisu danaco konfiguracji MCP okna.
package core

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/narzedzia"
	"danacoconsole/server/internal/session"
)

// DolozeniaNarzedziSesji podaje doraźne dołożenia — narzędzia dorzucone do tury przez Operatora w obrębie jednej sesji, poza definicją eksperta; sesja bez dołożeń oddaje listę pustą.
type DolozeniaNarzedziSesji interface {
	DolozeniaNarzedzi(ctx context.Context, idSesji string) (nazwy []string, err error)
}

// ZDolozeniamiSesji wpina źródło doraźnych dołożeń; adapter bez tego portu prowadzi turę samą podstawą eksperta i melduje powód.
func (a *adapterRozmowy) ZDolozeniamiSesji(d DolozeniaNarzedziSesji) *adapterRozmowy {
	a.dolozeniaSesji = d
	return a
}

// zestawTury składa opis zestawu narzędzi tej jednej tury z obu źródeł: podstawa bierze się z pola Agent okna, a dołożenia z sesji okna, bo dołożenie żyje sesją, nie oknem.
func (a *adapterRozmowy) zestawTury(ctx context.Context, okno session.Okno) injection.ZestawTury {
	if a == nil {
		return injection.ZestawTury{}
	}
	return injection.ZlozZestawTury(okno.Agent, a.dolozeniaTury(ctx, okno))
}

// dolozeniaTury pyta port o doraźne dołożenia sesji; meldunek o niewpiętym porcie wychodzi wyłącznie przy oknie z ekspertem, bo tylko wtedy brak dołożeń ma skutek dający się nazwać.
func (a *adapterRozmowy) dolozeniaTury(ctx context.Context, okno session.Okno) []string {
	if okno.IdSesji == "" {
		return nil
	}
	if a.dolozeniaSesji == nil {
		if okno.Agent != "" {
			a.zglosZestaw("brak-portu-dolozen-sesji",
				"doraźne dołożenia narzędzi poza turą: okno %q pracuje na zestawie eksperta %q,"+
					" ale port DolozeniaNarzedziSesji nie jest wpięty — narzędzia dołożone"+
					" w sesji %q do tury nie wejdą; tura jedzie samą podstawą eksperta",
				okno.Id, okno.Agent, okno.IdSesji)
		}
		return nil
	}
	nazwy, err := a.dolozeniaSesji.DolozeniaNarzedzi(ctx, okno.IdSesji)
	if err != nil {
		a.zglosZestaw("blad-dolozen-sesji",
			"doraźne dołożenia narzędzi nieodczytane (okno %q, sesja %q): %v —"+
				" tura idzie samą podstawą eksperta",
			okno.Id, okno.IdSesji, err)
		return nil
	}
	return nazwy
}

// argumentyZestawu przekłada opis zestawu na argumenty wpisu danaco; dołożenia bez podstawy nie jadą i mówią, dlaczego, bo okno bez eksperta nie ma gdzie ich dołożyć.
func (a *adapterRozmowy) argumentyZestawu(zestaw injection.ZestawTury, okno session.Okno) []string {
	argumenty := narzedzia.ArgumentyEksperta(zestaw.Ekspert)
	if len(zestaw.Dolozenia) == 0 {
		return argumenty
	}
	if zestaw.Ekspert == "" {
		a.zglosZestaw("dolozenia-bez-podstawy",
			"doraźne dołożenia sesji %q nie weszły do tury okna %q: okno nie ma nałożonego"+
				" eksperta, a zawężony wykaz — jedyne miejsce, w którym dołożenie ma co"+
				" dołożyć — powstaje tylko dla okna z ekspertem; tura idzie pełnym wykazem",
			okno.IdSesji, okno.Id)
		return argumenty
	}
	return append(argumenty, zestaw.ArgumentyDolozen()...)
}

// zZestawemTury dopisuje do wpisu danaco argumenty niosące zestaw tury; zestaw pusty nie dokłada nic, a tura jedzie pełnym wykazem w zasięgu roli okna.
func zZestawemTury(tekst string, argumenty []string) string {
	if len(argumenty) == 0 {
		return tekst
	}
	var konfiguracja KonfiguracjaMostu
	if err := json.Unmarshal([]byte(tekst), &konfiguracja); err != nil {
		return tekst
	}
	wpis, jest := konfiguracja.McpServers[narzedzia.KluczWpisu]
	if !jest {
		return tekst
	}
	wpis.Args = append(append([]string{}, wpis.Args...), argumenty...)
	konfiguracja.McpServers[narzedzia.KluczWpisu] = wpis
	tresc, err := json.MarshalIndent(konfiguracja, "", "  ")
	if err != nil {
		return tekst
	}
	return string(tresc)
}

// zKonfiguracjaZestawu jest całą drogą złożoną w jedno wywołanie, żeby miejsce
// wpięcia w `adapter_rozmowa_srodowisko.go` pozostało jedną linią.
func (a *adapterRozmowy) zKonfiguracjaZestawu(ctx context.Context, okno session.Okno, tekst string) string {
	zestaw := a.zestawTury(ctx, okno)
	if !zestaw.Skladany() {
		return tekst
	}
	return zZestawemTury(tekst, a.argumentyZestawu(zestaw, okno))
}

// zglosZestaw wpisuje do dziennika rdzenia powód, dla którego zestaw narzędzi tury jest taki, jaki jest; kluczem jest sam powód, raz na powód, nie raz na turę.
func (a *adapterRozmowy) zglosZestaw(powod, wzorzec string, wartosci ...any) {
	if a == nil || a.mosty == nil || a.mosty.dziennik == nil || powod == "" {
		return
	}
	if _, juzBylo := a.zgloszoneZestawy.LoadOrStore(powod, true); juzBylo {
		return
	}
	a.mosty.dziennik.Printf(wzorzec, wartosci...)
}
