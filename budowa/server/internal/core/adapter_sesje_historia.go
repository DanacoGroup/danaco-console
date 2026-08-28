// Plik zawiera odwracalne czynności na wykazie sesji: zmianę nazwy,
// archiwizację i przywrócenie, z zapisem zachowanym w całości w bazie.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ZmienNazwe zmienia nazwę sesji w historii bieżącej konta, utrwalając ją
// w bazie, gdy baza jest dostępna.
func (a *adapterSesji) ZmienNazwe(ctx context.Context, z shared.SessionRenameRequest) (shared.SessionRenameResponse, error) {
	if err := a.zapiszNazwe(ctx, z.SessionId, z.Title); err != nil {
		return shared.SessionRenameResponse{}, err
	}
	sesja, err := a.nadzorca.Rejestr().ZmienTytulSesji(z.SessionId, z.Title)
	if err != nil {
		return shared.SessionRenameResponse{}, bladSesji(err)
	}
	return shared.SessionRenameResponse{Session: sesjaKontraktu(sesja)}, nil
}

// Archiwizuj przenosi sesje do archiwum.
//
// Sesja archiwalna znika z historii bieżącej, ale nic nie traci: wiadomości,
// okna i katalog roboczy zostają nietknięte, a wgląd daje session.archive.list
// z okna ustawień.
func (a *adapterSesji) Archiwizuj(ctx context.Context, z shared.SessionArchiveRequest) (shared.SessionArchiveResponse, error) {
	przeniesione := a.przestawStan(ctx, z.SessionIds, shared.SessionStatusArchived)
	return shared.SessionArchiveResponse{ArchivedIds: przeniesione}, nil
}

// Przywroc wraca sesje do historii bieżącej z archiwum i z kosza, jedną
// komendą kontraktu i jednym znaczeniem przywrócenia niezależnie od stanu
// wyjściowego.
func (a *adapterSesji) Przywroc(ctx context.Context, z shared.SessionRestoreRequest) (shared.SessionRestoreResponse, error) {
	zKosza := a.wyjmijZKosza(ctx, z.SessionIds)
	przywrocone := a.przestawStan(ctx, z.SessionIds, shared.SessionStatusActive)
	a.wniesPrzywroconeDoRejestru(ctx, zKosza)
	return shared.SessionRestoreResponse{RestoredIds: przywrocone}, nil
}

// WykazArchiwum zwraca sesje archiwum konta, jako wgląd udostępniany z okna
// ustawień platformy operatorowi.
func (a *adapterSesji) WykazArchiwum(ctx context.Context, z shared.SessionArchiveListRequest) (shared.SessionArchiveListResponse, error) {
	if a.trwalosc == nil || a.trwalosc.sesje == nil {
		return shared.SessionArchiveListResponse{Sessions: []shared.Session{}}, nil
	}
	wiersze, err := a.trwalosc.sesje.Lista(ctx, 0)
	if err != nil {
		return shared.SessionArchiveListResponse{}, err
	}
	archiwalne := make([]shared.Session, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if wiersz.Stan != shared.SessionStatusArchived {
			continue
		}
		archiwalne = append(archiwalne, sesjaWierszaKontraktu(wiersz))
	}
	return shared.SessionArchiveListResponse{
		Sessions: wycinek(archiwalne, z.Offset, z.Limit),
		Total:    len(archiwalne),
	}, nil
}

// przestawStan zmienia stan wskazanych sesji i zwraca te, które faktycznie
// zmieniono. Wskazanie bez odpowiednika jest pomijane, nie jest błędem —
// czynność zbiorcza nie może paść przez jedną pozycję.
func (a *adapterSesji) przestawStan(ctx context.Context, wskazania []string,
	stan shared.SessionStatus) []string {

	zmienione := make([]string, 0, len(wskazania))
	for _, identyfikator := range wskazania {
		if a.trwalosc == nil || a.trwalosc.sesje == nil {
			continue
		}
		wiersz, err := a.trwalosc.sesje.PoIdentyfikatorze(ctx, identyfikator)
		if err != nil {
			continue
		}
		if err := a.trwalosc.sesje.ZmienStan(ctx, wiersz.ID, stan); err != nil {
			continue
		}
		zmienione = append(zmienione, identyfikator)
	}
	return zmienione
}

// zapiszNazwe utrwala nazwę sesji. Brak utrwalacza znaczy rdzeń bez bazy —
// nazwa żyje wtedy wyłącznie w rejestrze i nie jest to błąd.
func (a *adapterSesji) zapiszNazwe(ctx context.Context, identyfikator, nazwa string) error {
	if a.trwalosc == nil || a.trwalosc.sesje == nil {
		return nil
	}
	wiersz, err := a.trwalosc.sesje.PoIdentyfikatorze(ctx, identyfikator)
	if err != nil {
		// Odczyt sesji owija brak wiersza, więc porównanie wprost nie rozpozna sesji nieznanej bazie.
		if errors.Is(err, dane.ErrBrakWiersza) {
			return nil
		}
		return err
	}
	return a.trwalosc.sesje.ZmienTytul(ctx, wiersz.ID, nazwa)
}
