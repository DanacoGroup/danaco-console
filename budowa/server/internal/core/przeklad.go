package core

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przekład między bytami pakietu sesji a kształtem kontraktu.
//
// Pakiet session jest właścicielem pojęcia okna komunikacji i jego cyklu życia,
// więc trzyma własne struktury. Kontrakt trzyma własne. Ten plik jest jedynym
// miejscem, w którym jedno przechodzi w drugie — nie ma drugiego przekładu
// rozsianego po obsługiwaczach.

// sesjaKontraktu przekłada sesję pakietu sesji na sesję kontraktu.
func sesjaKontraktu(s session.Sesja) shared.Session {
	sesja := shared.Session{
		Id:        s.Id,
		Status:    s.Stan,
		WindowIds: s.IdOkien,
		CreatedAt: s.Utworzono.UnixMilli(),
		UpdatedAt: s.Zaktualizowano.UnixMilli(),
	}
	if s.Tytul != "" {
		sesja.Title = &s.Tytul
	}
	if s.IdProjektu != "" {
		sesja.ProjectId = &s.IdProjektu
	}
	return sesja
}

// oknoKontraktu przekłada okno komunikacji na okno kontraktu.
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
	// Ekspert wraca do klienta wyłącznie wtedy, gdy jest — pole puste znaczy
	// model surowy, a wskaźnik na pusty napis mówiłby to samo drugim sposobem.
	if o.Agent != "" {
		okno.AgentId = &o.Agent
	}
	return okno
}

// oknaKontraktu przekłada wykaz okien.
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
		// Ekspert okna. Pole przechodzi do `session.Zmiana`, więc wybór eksperta
		// zapisuje się tak samo jak każde inne ustawienie okna, a nie ginie po
		// drodze.
		Agent: z.AgentId,
	}
}

// bladSesji nadaje błędowi pakietu sesji kod kontraktu. Kod rozstrzyga właściciel
// pojęcia — rdzeń go wyłącznie przenosi (patrz wynik.go).
func bladSesji(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.NowyBlad(session.Kod(err), err.Error()))
}

// listaKatalogow pilnuje, żeby katalogi robocze wyszły jako tablica, nigdy jako
// brak wartości.
//
// Kontrakt zapowiada `workingDirs` bezwarunkowo, a wycinek pusty w Go koduje się
// do `null`. Okno bez katalogów jest stanem poprawnym, więc kontrakt nie każe
// odbiorcy przygotowywać się na brak — klient czytający to pole bez osłony
// wywróciłby na `null` wczytanie modułu. Wysyłamy pustą tablicę i różnica znika.
func listaKatalogow(katalogi []string) []string {
	if katalogi == nil {
		return []string{}
	}
	return katalogi
}
