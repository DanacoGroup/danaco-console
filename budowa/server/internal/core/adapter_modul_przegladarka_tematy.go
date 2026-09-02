// Porządkowanie zebranego materiału: zestawy tematyczne źródeł i wątki notatek;
// skład zestawu i wątku liczy się z kolumny przynależności, nie z osobnej listy.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

func (a *adapterPrzegladarki) UsunZrodlo(ctx context.Context,
	z shared.BrowserSourceRemoveRequest) (shared.BrowserSourceRemoveResponse, error) {

	if strings.TrimSpace(z.SourceId) == "" {
		return shared.BrowserSourceRemoveResponse{}, bladWskazaniaPrzegladarki(
			"komenda source.remove bez źródła")
	}
	usuniete, err := a.repozytorium.UsunZrodlo(ctx, z.SourceId)
	if err != nil {
		return shared.BrowserSourceRemoveResponse{}, bladPrzegladarki(err)
	}
	if !usuniete {
		return shared.BrowserSourceRemoveResponse{}, bladNieznanegoBytu("źródła przeglądania", z.SourceId)
	}
	return shared.BrowserSourceRemoveResponse{Removed: true}, nil
}

func (a *adapterPrzegladarki) ZmienNotatke(ctx context.Context,
	z shared.BrowserNoteUpdateRequest) (shared.BrowserNoteUpdateResponse, error) {

	if strings.TrimSpace(z.NoteId) == "" {
		return shared.BrowserNoteUpdateResponse{}, bladWskazaniaPrzegladarki("komenda note.update bez notatki")
	}
	zmiana := dane.ZmianaNotatki{
		Kod:       z.NoteId,
		Tresc:     z.Content,
		Cytat:     z.Quote,
		ZrodloID:  z.SourceId,
		Watek:     z.ThreadId,
		Przypieta: z.Pinned,
	}
	if z.Classification != nil {
		klasyfikacja := string(*z.Classification)
		zmiana.Klasyfikacja = &klasyfikacja
	}
	zapisana, err := a.repozytorium.AktualizujNotatke(ctx, zmiana)
	if err != nil {
		return shared.BrowserNoteUpdateResponse{}, bladWierszaPrzegladania("notatki przeglądania", z.NoteId, err)
	}
	return shared.BrowserNoteUpdateResponse{Note: notatkaKontraktu(zapisana)}, nil
}

func (a *adapterPrzegladarki) UstawZestawZrodel(ctx context.Context,
	z shared.BrowserSourceGroupSetRequest) (shared.BrowserSourceGroupSetResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Name) == "" {
		return shared.BrowserSourceGroupSetResponse{}, bladWskazaniaPrzegladarki(
			"komenda source.group.set bez okna albo nazwy zestawu")
	}
	kod := wartoscTekstuLubPusta(z.GroupId)

	if z.Removed != nil && *z.Removed {
		if kod == "" {
			return shared.BrowserSourceGroupSetResponse{}, bladWskazaniaPrzegladarki(
				"zdjęcie zestawu źródeł bez jego wskazania")
		}
		zestaw, err := a.repozytorium.ZestawZrodel(ctx, kod)
		if err != nil {
			return shared.BrowserSourceGroupSetResponse{}, bladWierszaPrzegladania("zestawu źródeł", kod, err)
		}
		if _, err := a.repozytorium.UsunZestawZrodel(ctx, kod); err != nil {
			return shared.BrowserSourceGroupSetResponse{}, bladPrzegladarki(err)
		}
		return shared.BrowserSourceGroupSetResponse{Group: zestawKontraktu(zestaw, nil)}, nil
	}

	if kod == "" {
		kod = nowyIdentyfikator(przedrostekZestawuZrodel)
	}
	zapisany, err := a.repozytorium.ZapiszZestawZrodel(ctx, dane.ZestawZrodel{
		Kod: kod, Okno: z.WindowId, Nazwa: z.Name,
	})
	if err != nil {
		return shared.BrowserSourceGroupSetResponse{}, bladPrzegladarki(err)
	}
	for _, zrodlo := range z.SourceIds {
		if err := a.repozytorium.PrzypiszZrodloDoZestawu(ctx, z.WindowId, zrodlo, kod); err != nil {
			return shared.BrowserSourceGroupSetResponse{}, bladWierszaPrzegladania("źródła przeglądania", zrodlo, err)
		}
	}
	sklad, err := a.repozytorium.KodyZrodelZestawu(ctx, kod)
	if err != nil {
		return shared.BrowserSourceGroupSetResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserSourceGroupSetResponse{Group: zestawKontraktu(zapisany, sklad)}, nil
}

func (a *adapterPrzegladarki) WykazZestawowZrodel(ctx context.Context,
	z shared.BrowserSourceGroupListRequest) (shared.BrowserSourceGroupListResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.BrowserSourceGroupListResponse{}, bladWskazaniaPrzegladarki(
			"komenda source.group.list bez okna")
	}
	wiersze, err := a.repozytorium.ZestawyZrodel(ctx, z.WindowId, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.BrowserSourceGroupListResponse{}, bladPrzegladarki(err)
	}
	zestawy := make([]shared.BrowserSourceGroup, 0, len(wiersze))
	for _, wiersz := range wiersze {
		sklad, err := a.repozytorium.KodyZrodelZestawu(ctx, wiersz.Kod)
		if err != nil {
			return shared.BrowserSourceGroupListResponse{}, bladPrzegladarki(err)
		}
		zestawy = append(zestawy, zestawKontraktu(wiersz, sklad))
	}
	return shared.BrowserSourceGroupListResponse{Groups: zestawy}, nil
}

func (a *adapterPrzegladarki) UstawWatekNotatek(ctx context.Context,
	z shared.BrowserNoteThreadSetRequest) (shared.BrowserNoteThreadSetResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Name) == "" {
		return shared.BrowserNoteThreadSetResponse{}, bladWskazaniaPrzegladarki(
			"komenda note.thread.set bez okna albo nazwy wątku")
	}
	kod := wartoscTekstuLubPusta(z.ThreadId)

	if z.Removed != nil && *z.Removed {
		if kod == "" {
			return shared.BrowserNoteThreadSetResponse{}, bladWskazaniaPrzegladarki(
				"zdjęcie wątku notatek bez jego wskazania")
		}
		watek, err := a.repozytorium.WatekNotatek(ctx, kod)
		if err != nil {
			return shared.BrowserNoteThreadSetResponse{}, bladWierszaPrzegladania("wątku notatek", kod, err)
		}
		if _, err := a.repozytorium.UsunWatekNotatek(ctx, kod); err != nil {
			return shared.BrowserNoteThreadSetResponse{}, bladPrzegladarki(err)
		}
		return shared.BrowserNoteThreadSetResponse{Thread: watekKontraktu(watek, nil)}, nil
	}

	if kod == "" {
		kod = nowyIdentyfikator(przedrostekWatkuNotatek)
	}
	zapisany, err := a.repozytorium.ZapiszWatekNotatek(ctx, dane.WatekNotatek{
		Kod: kod, Okno: z.WindowId, Nazwa: z.Name,
	})
	if err != nil {
		return shared.BrowserNoteThreadSetResponse{}, bladPrzegladarki(err)
	}
	for _, notatka := range z.NoteIds {
		if err := a.repozytorium.PrzypiszNotatkeDoWatku(ctx, z.WindowId, notatka, kod); err != nil {
			return shared.BrowserNoteThreadSetResponse{}, bladPrzegladarki(err)
		}
	}
	sklad, err := a.repozytorium.KodyNotatekWatku(ctx, kod)
	if err != nil {
		return shared.BrowserNoteThreadSetResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserNoteThreadSetResponse{Thread: watekKontraktu(zapisany, sklad)}, nil
}

func (a *adapterPrzegladarki) WykazWatkowNotatek(ctx context.Context,
	z shared.BrowserNoteThreadListRequest) (shared.BrowserNoteThreadListResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.BrowserNoteThreadListResponse{}, bladWskazaniaPrzegladarki(
			"komenda note.thread.list bez okna")
	}
	wiersze, err := a.repozytorium.WatkiNotatek(ctx, z.WindowId, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.BrowserNoteThreadListResponse{}, bladPrzegladarki(err)
	}
	watki := make([]shared.BrowserNoteThread, 0, len(wiersze))
	for _, wiersz := range wiersze {
		sklad, err := a.repozytorium.KodyNotatekWatku(ctx, wiersz.Kod)
		if err != nil {
			return shared.BrowserNoteThreadListResponse{}, bladPrzegladarki(err)
		}
		watki = append(watki, watekKontraktu(wiersz, sklad))
	}
	return shared.BrowserNoteThreadListResponse{Threads: watki}, nil
}

func zestawKontraktu(w dane.ZestawZrodel, sklad []string) shared.BrowserSourceGroup {
	zestaw := shared.BrowserSourceGroup{
		Id:        w.Kod,
		WindowId:  w.Okno,
		Name:      w.Nazwa,
		CreatedAt: chwilaBazy(w.Utworzono),
		UpdatedAt: chwilaBazy(w.Zaktualizowano),
	}
	if len(sklad) > 0 {
		zestaw.SourceIds = sklad
	}
	return zestaw
}

func watekKontraktu(w dane.WatekNotatek, sklad []string) shared.BrowserNoteThread {
	watek := shared.BrowserNoteThread{
		Id:        w.Kod,
		WindowId:  w.Okno,
		Name:      w.Nazwa,
		CreatedAt: chwilaBazy(w.Utworzono),
		UpdatedAt: chwilaBazy(w.Zaktualizowano),
	}
	if len(sklad) > 0 {
		watek.NoteIds = sklad
	}
	return watek
}
