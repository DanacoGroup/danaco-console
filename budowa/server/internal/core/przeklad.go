package core

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Plik przekłada byty pakietu sesji na kształt kontraktu w jednym miejscu, bez rozproszenia.

// sesjaKontraktu przekłada sesję pakietu sesji na sesję kontraktu, przenosząc identyfikator, stan, wykaz okien oraz znaczniki czasu utworzenia i aktualizacji.
func sesjaKontraktu(s session.Sesja) shared.Session {
	sesja := shared.Session{
		Id:        s.Id,
		Status:    s.Stan,
		WindowIds: s.IdOkien,
		CreatedAt: s.Utworzono.UnixMilli(),
		UpdatedAt: s.Zaktualizowano.UnixMilli(),
	}
	/*
		Tytuł równy identyfikatorowi znaczy sesję bez nazwy.

		Kolumna tytułu w magazynie nie przyjmuje pustej wartości, więc sesja bez
		nazwy dostaje tam własny identyfikator (`dane.sesjaZOpisu`). To ograniczenie
		magazynu, nie nazwa nadana przez Operatora — podane oknu czytałoby się jako
		nazwa i tak właśnie stawało w szynie sesji. Okno ma własny stan pusty i to
		on ma się pokazać.
	*/
	if s.Tytul != "" && s.Tytul != s.Id {
		sesja.Title = &s.Tytul
	}
	if s.IdProjektu != "" {
		sesja.ProjectId = &s.IdProjektu
	}
	return sesja
}

// oknoKontraktu przekłada okno komunikacji na okno kontraktu, zachowując powiązania z sesją, modułem i kanałem modelu oraz pozostałe pola opisujące jego miejsce w interfejsie.
func oknoKontraktu(o session.Okno) shared.Window {
	okno := shared.Window{
		Id:             o.Id,
		SessionId:      o.IdSesji,
		ModuleId:       o.Modul,
		ModelChannelId: o.KanalModelu,
		WorkingDirs:    listaKatalogow(o.KatalogiRobocze),
		ExecutionEnv:   o.SrodowiskoWykonania,
		PermissionMode: o.TrybUprawnien,
		WindowRole:     o.RolaOkna,
		Status:         o.Stan,
		CreatedAt:      o.Utworzono.UnixMilli(),
		UpdatedAt:      o.Zaktualizowano.UnixMilli(),
	}
	if o.OknoKoordynatora != "" {
		okno.CoordinatorWindowId = &o.OknoKoordynatora
	}
	if o.Tytul != "" {
		okno.Title = &o.Tytul
	}
	// Ekspert wraca do klienta tylko, gdy jest ustawiony — pole puste oznacza model surowy.
	if o.Agent != "" {
		okno.AgentId = &o.Agent
	}
	return okno
}

// oknaKontraktu przekłada wykaz okien pakietu sesji na wykaz okien kontraktu, wywołując przekład pojedynczego okna dla każdej pozycji z osobna.
func oknaKontraktu(okna []session.Okno) []shared.Window {
	wykaz := make([]shared.Window, 0, len(okna))
	for _, okno := range okna {
		wykaz = append(wykaz, oknoKontraktu(okno))
	}
	return wykaz
}

// ustawieniaOkna przekłada żądanie założenia okna na ustawienia pakietu sesji.
// Wartości puste zostają puste — uzupełni je wartościami domyślnymi właściciel
// pojęcia, a nie warstwa przekładu.
func ustawieniaOkna(z shared.WindowCreateRequest) session.Ustawienia {
	u := session.Ustawienia{
		Modul:               z.ModuleId,
		KanalModelu:         z.ModelChannelId,
		KatalogiRobocze:     z.WorkingDirs,
		SrodowiskoWykonania: z.ExecutionEnv,
		TrybUprawnien:       z.PermissionMode,
		RolaOkna:            z.WindowRole,
	}
	if z.CoordinatorWindowId != nil {
		u.OknoKoordynatora = *z.CoordinatorWindowId
	}
	if z.Title != nil {
		u.Tytul = *z.Title
	}
	return u
}

// zmianaOkna przekłada żądanie zmiany okna na zmianę wybiórczą pakietu sesji.
// Pole niewypełnione w żądaniu zostaje polem niewypełnionym w zmianie, więc
// ustawienie nietknięte przez Operatora nie jest nadpisywane.
func zmianaOkna(z shared.WindowUpdateRequest) session.Zmiana {
	return session.Zmiana{
		Modul:               z.ModuleId,
		KanalModelu:         z.ModelChannelId,
		KatalogiRobocze:     z.WorkingDirs,
		SrodowiskoWykonania: z.ExecutionEnv,
		TrybUprawnien:       z.PermissionMode,
		RolaOkna:            z.WindowRole,
		OknoKoordynatora:    z.CoordinatorWindowId,
		Tytul:               z.Title,
		// Ekspert okna zapisuje się jak każde inne ustawienie okna, a nie ginie po drodze.
		Agent: z.AgentId,
	}
}

// bladSesji nadaje błędowi pakietu sesji kod kontraktu. Kod rozstrzyga właściciel pojęcia, rdzeń go wyłącznie przenosi.
func bladSesji(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.NowyBlad(session.Kod(err), err.Error()))
}

// listaKatalogow pilnuje, żeby katalogi robocze wyszły jako tablica, nigdy jako brak wartości, ponieważ kontrakt zapowiada to pole bezwarunkowo, a wycinek pusty w Go koduje się do wartości null.
func listaKatalogow(katalogi []string) []string {
	if katalogi == nil {
		return []string{}
	}
	return katalogi
}
