// Odpowiedzialność pliku: doprowadzenie łańcucha więzów do wiersza okna komunikacji, zakładając brakujące ogniwa zamiast odmawiać zapisu.
package dane

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"danacoconsole/shared"
)

const (
	// nazwaKartyDomyslnej to karta, pod którą trafiają sesje założone przez rdzeń
	// bez wskazania karty. Karta jest kontenerem obowiązkowym — sesja bez niej nie
	// przejdzie więzu klucza obcego.
	nazwaKartyDomyslnej = "Karta domyślna"

	// kodModuluZastepczego wskazuje moduł przypisywany oknu, którego moduł nie ma odpowiednika w słowniku modułów.
	kodModuluZastepczego = "workspace"
)

// wierszOkna zwraca identyfikator wiersza okna, zakładając brakujące ogniwa
// łańcucha. Odwzorowanie identyfikatorów jest zapamiętywane, więc kolejne
// wiadomości tego samego okna nie odpytują bazy.
func (u *UtrwalaczRozmowy) wierszOkna(ctx context.Context, idOkna string) (int64, error) {
	if idOkna == "" {
		return 0, fmt.Errorf("dane: wiadomość bez okna komunikacji nie da się utrwalić")
	}
	u.mu.Lock()
	defer u.mu.Unlock()

	if id, jest := u.okna[idOkna]; jest {
		return id, nil
	}
	okno, err := u.zestaw.Okna.PoIdentyfikatorze(ctx, idOkna)
	if err == nil {
		u.okna[idOkna] = okno.ID
		return okno.ID, nil
	}
	if !errors.Is(err, ErrBrakWiersza) {
		return 0, err
	}
	id, err := u.zalozOkno(ctx, idOkna)
	if err != nil {
		return 0, err
	}
	u.okna[idOkna] = id
	return id, nil
}

// zalozOkno zapisuje wiersz okna wraz z całym łańcuchem więzów, na którym ten sam wiersz właśnie stoi.
func (u *UtrwalaczRozmowy) zalozOkno(ctx context.Context, idOkna string) (int64, error) {
	if u.zrodlo == nil {
		return 0, fmt.Errorf("dane: okno %q nie ma opisu w rejestrze serwera", idOkna)
	}
	opis, jest := u.zrodlo(idOkna)
	if !jest {
		return 0, fmt.Errorf("dane: okno %q nie jest znane rejestrowi serwera", idOkna)
	}
	sesjaID, err := u.wierszSesji(ctx, opis)
	if err != nil {
		return 0, err
	}
	modulID, err := u.wierszModulu(ctx, opis.Modul)
	if err != nil {
		return 0, err
	}
	kanalID, err := u.wierszKanalu(ctx, opis.KanalModelu)
	if err != nil {
		return 0, err
	}
	return u.zestaw.Okna.Utworz(ctx, oknoZOpisu(opis, sesjaID, modulID, kanalID))
}

// oknoZOpisu składa wiersz okna z opisu rejestru rdzenia, gotowy do zapisania w tej samej bazie danych.
func oknoZOpisu(opis OpisOkna, sesjaID, modulID, kanalID int64) Okno {
	okno := NoweOkno(sesjaID, modulID, kanalID)
	identyfikator := opis.Id
	okno.IdentyfikatorZewnetrzny = &identyfikator
	if opis.Tytul != "" {
		tytul := opis.Tytul
		okno.Tytul = &tytul
	}
	if len(opis.KatalogiRobocze) > 0 {
		okno.KatalogiRobocze = append([]string(nil), opis.KatalogiRobocze...)
	}
	if opis.SrodowiskoWykonania != "" {
		okno.SrodowiskoWykonania = opis.SrodowiskoWykonania
	}
	if opis.TrybUprawnien != "" {
		okno.TrybUprawnien = opis.TrybUprawnien
	}
	if opis.RolaOkna != "" {
		okno.RolaOkna = opis.RolaOkna
	}
	if opis.Agent != "" {
		agent := opis.Agent
		okno.AgentKod = &agent
	}
	return okno
}

// wierszSesji odnajduje sesję po identyfikatorze rdzenia albo zakłada ją wraz
// z kartą sesji w pierwszym czynnym środowisku.
func (u *UtrwalaczRozmowy) wierszSesji(ctx context.Context, opis OpisOkna) (int64, error) {
	sesja, err := u.zestaw.Sesje.PoIdentyfikatorze(ctx, opis.IdSesji)
	if err == nil {
		return sesja.ID, nil
	}
	if !errors.Is(err, ErrBrakWiersza) {
		return 0, err
	}
	srodowisko, err := u.zestaw.Srodowiska.Pierwsze(ctx)
	if err != nil {
		return 0, err
	}
	kartaID, err := u.zestaw.KartySesji.Zapewnij(ctx, srodowisko.ID, nazwaKartyDomyslnej)
	if err != nil {
		return 0, err
	}
	return u.zestaw.Sesje.Utworz(ctx, sesjaZOpisu(opis, kartaID))
}

// sesjaZOpisu składa wiersz sesji z opisu okna. Tytuł jest kolumną obowiązkową,
// więc sesja bez tytułu dostaje własny identyfikator — nigdy pustą wartość.
func sesjaZOpisu(opis OpisOkna, kartaID int64) Sesja {
	identyfikator := opis.IdSesji
	sesja := Sesja{
		KartaSesjiID:            kartaID,
		Tytul:                   opis.TytulSesji,
		Stan:                    shared.SessionStatusActive,
		IdentyfikatorZewnetrzny: &identyfikator,
	}
	if sesja.Tytul == "" {
		sesja.Tytul = opis.IdSesji
	}
	if opis.Projekt != "" {
		projekt := opis.Projekt
		sesja.Projekt = &projekt
	}
	return sesja
}

// wierszModulu przekłada kod modułu okna na wiersz słownika. Kod nieznany —
// na przykład wskazujący środowisko zamiast modułu — schodzi na moduł zastępczy,
// bo utrata całej rozmowy byłaby ceną nieproporcjonalną.
func (u *UtrwalaczRozmowy) wierszModulu(ctx context.Context, kod string) (int64, error) {
	if kod != "" {
		modul, err := u.zestaw.Moduly.PoKodzie(ctx, kod)
		if err == nil {
			return modul.ID, nil
		}
		if !errors.Is(err, ErrBrakWiersza) {
			return 0, err
		}
	}
	zastepczy, err := u.zestaw.Moduly.PoKodzie(ctx, kodModuluZastepczego)
	if err != nil {
		return 0, fmt.Errorf("dane: słownik modułów pusty, okna nie da się zapisać: %w", err)
	}
	return zastepczy.ID, nil
}

// wierszKanalu przekłada wskazanie kanału okna na wiersz rejestru kanałów.
// Rdzeń podaje identyfikator wiersza w postaci tekstowej, ale przyjmowany jest
// także kod kanału. Wskazanie nierozpoznane schodzi na pierwszy kanał rejestru.
func (u *UtrwalaczRozmowy) wierszKanalu(ctx context.Context, wskazanie string) (int64, error) {
	kanaly, err := u.zestaw.Kanaly.Lista(ctx, false)
	if err != nil {
		return 0, err
	}
	if len(kanaly) == 0 {
		return 0, fmt.Errorf("dane: rejestr kanałów pusty, okna nie da się zapisać")
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

// ZapewnijSesje utrwala sesję natychmiast po jej założeniu, żeby czynności na historii miały czego dotknąć bez pierwszej wiadomości.
func (u *UtrwalaczRozmowy) ZapewnijSesje(ctx context.Context, idSesji, tytul, projekt string) (int64, error) {
	sesja, err := u.zestaw.Sesje.PoIdentyfikatorze(ctx, idSesji)
	if err == nil {
		return sesja.ID, nil
	}
	if !errors.Is(err, ErrBrakWiersza) {
		return 0, err
	}
	srodowisko, err := u.zestaw.Srodowiska.Pierwsze(ctx)
	if err != nil {
		return 0, err
	}
	kartaID, err := u.zestaw.KartySesji.Zapewnij(ctx, srodowisko.ID, nazwaKartyDomyslnej)
	if err != nil {
		return 0, err
	}
	return u.zestaw.Sesje.Utworz(ctx, sesjaZOpisu(OpisOkna{
		IdSesji: idSesji, TytulSesji: tytul, Projekt: projekt,
	}, kartaID))
}
