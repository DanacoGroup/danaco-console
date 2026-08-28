// Plik utrwala okno komunikacji w chwili jego założenia w bazie, będącej
// źródłem prawdy o oknie, ponieważ rejestr nadzorcy jest jedynie jej widokiem
// na czas jednego uruchomienia rdzenia.
package core

import (
	"context"
	"errors"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
)

// utrwalZalozone zapisuje wiersz okna zaraz po jego założeniu w rejestrze;
// niepowodzenie zapisu nie przerywa zakładania okna, które wtedy nie przetrwa
// restartu rdzenia.
func (a *adapterOkien) utrwalZalozone(ctx context.Context, okno session.Okno) {
	if a.trwalosc == nil || a.trwalosc.okna == nil || a.trwalosc.sesje == nil {
		return
	}
	if _, err := a.trwalosc.okna.PoIdentyfikatorze(ctx, okno.Id); err == nil {
		return // wiersz już jest — czynność jest idempotentna
	} else if !errors.Is(err, dane.ErrBrakWiersza) {
		a.trwalosc.odnotuj("utrwalenie okna %s: %v", okno.Id, err)
		return
	}
	wiersz, err := a.wierszZalozonegoOkna(ctx, okno)
	if err != nil {
		a.trwalosc.odnotuj("utrwalenie okna %s: %v", okno.Id, err)
		return
	}
	if _, err := a.trwalosc.okna.Utworz(ctx, wiersz); err != nil {
		a.trwalosc.odnotuj("zapis okna %s: %v", okno.Id, err)
	}
}

// wierszZalozonegoOkna składa wiersz okna z bytu rejestru, rozstrzygając więzy
// obce sesji, modułu i kanału modelu, nie zakładając przy tym wiersza samej
// sesji.
func (a *adapterOkien) wierszZalozonegoOkna(ctx context.Context,
	okno session.Okno) (dane.Okno, error) {

	sesja, err := a.trwalosc.sesje.PoIdentyfikatorze(ctx, okno.IdSesji)
	if err != nil {
		return dane.Okno{}, err
	}
	modulID, err := a.wierszModuluOkna(ctx, okno.Ustawienia.Modul)
	if err != nil {
		return dane.Okno{}, err
	}
	kanalID, err := a.wierszKanaluOkna(ctx, okno.Ustawienia.KanalModelu)
	if err != nil {
		return dane.Okno{}, err
	}
	wiersz := dane.NoweOkno(sesja.ID, modulID, kanalID)
	identyfikator := okno.Id
	wiersz.IdentyfikatorZewnetrzny = &identyfikator
	wiersz.Stan = okno.Stan
	u := okno.Ustawienia
	if u.Tytul != "" {
		tytul := u.Tytul
		wiersz.Tytul = &tytul
	}
	if len(u.KatalogiRobocze) > 0 {
		wiersz.KatalogiRobocze = append([]string(nil), u.KatalogiRobocze...)
	}
	if u.SrodowiskoWykonania != "" {
		wiersz.SrodowiskoWykonania = u.SrodowiskoWykonania
	}
	if u.TrybUprawnien != "" {
		wiersz.TrybUprawnien = u.TrybUprawnien
	}
	if u.RolaOkna != "" {
		wiersz.RolaOkna = u.RolaOkna
	}
	if u.Agent != "" {
		agent := u.Agent
		wiersz.AgentKod = &agent
	}
	// Więzi koordynatora nie ustawia się tutaj: ustanawia ją przekazanie okna.
	return wiersz, nil
}

// wierszModuluOkna przekłada kod modułu na wiersz katalogu modułów zapisany
// w bazie, potrzebny do ustalenia więzu obcego wiersza okna.
func (a *adapterOkien) wierszModuluOkna(ctx context.Context, kod string) (int64, error) {
	moduly, err := a.modulyKatalogu(ctx)
	if err != nil {
		return 0, err
	}
	for _, modul := range moduly {
		if modul.Kod == kod {
			return modul.ID, nil
		}
	}
	return moduly[0].ID, nil
}

// modulyKatalogu oddaje katalog modułów platformy zapisanych w bazie albo błąd,
// gdy w bazie nie ma z czego wybierać.
func (a *adapterOkien) modulyKatalogu(ctx context.Context) ([]dane.Modul, error) {
	if a.moduly == nil {
		return nil, errors.New("core: katalog modułów niewpięty")
	}
	moduly, err := a.moduly.Lista(ctx)
	if err != nil {
		return nil, err
	}
	if len(moduly) == 0 {
		return nil, errors.New("core: katalog modułów pusty")
	}
	return moduly, nil
}

// wierszKanaluOkna przekłada wskazanie kanału na wiersz rejestru. Rdzeń podaje
// zwykle kod, ale przyjmowany jest także identyfikator wiersza w postaci
// tekstowej — tą samą drogą, co przy utrwalaniu rozmowy.
func (a *adapterOkien) wierszKanaluOkna(ctx context.Context, wskazanie string) (int64, error) {
	if a.kanaly == nil {
		return 0, errors.New("core: rejestr kanałów niewpięty")
	}
	kanaly, err := a.kanaly.Lista(ctx, false)
	if err != nil {
		return 0, err
	}
	if len(kanaly) == 0 {
		return 0, errors.New("core: rejestr kanałów pusty")
	}
	if id, err := strconv.ParseInt(wskazanie, 10, 64); err == nil {
		for _, kanal := range kanaly {
			if kanal.ID == id {
				return kanal.ID, nil
			}
		}
	}
	for _, kanal := range kanaly {
		if kanal.Kod == wskazanie {
			return kanal.ID, nil
		}
	}
	return kanaly[0].ID, nil
}
