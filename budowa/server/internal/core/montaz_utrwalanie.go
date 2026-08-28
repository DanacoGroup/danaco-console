// Plik wpina utrwalanie rozmowy: przełożenie okna żywego z rejestru nadzorcy na opis warstwy
// danych, będąc jedynym miejscem, w którym warstwa danych i pakiet sesji się spotykają.
package core

import (
	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
)

// utrwalaczRozmow wpina repozytoria sesji, okien i wiadomości w dziennik rozmów, domykając łańcuch
// środowisko, moduł, karta sesji, sesja, okno komunikacji, wiadomość.
func utrwalaczRozmow(repozytoria *dane.Zestaw, nadzorca *session.Nadzorca) *dane.UtrwalaczRozmowy {
	return dane.NowyUtrwalaczRozmowy(repozytoria, func(idOkna string) (dane.OpisOkna, bool) {
		okno, err := nadzorca.Rejestr().Okno(idOkna)
		if err != nil {
			return dane.OpisOkna{}, false
		}
		return opisOkna(okno, nadzorca), true
	})
}

// opisOkna przenosi okno rejestru na opis warstwy danych wraz z tytułem
// i projektem sesji nadrzędnej. Brak sesji w rejestrze nie unieważnia opisu —
// wiersz sesji powstanie z samego identyfikatora.
func opisOkna(okno session.Okno, nadzorca *session.Nadzorca) dane.OpisOkna {
	opis := dane.OpisOkna{
		Id:                  okno.Id,
		IdSesji:             okno.IdSesji,
		Modul:               okno.Modul,
		KanalModelu:         okno.KanalModelu,
		Tytul:               okno.Tytul,
		KatalogiRobocze:     okno.KatalogiRobocze,
		SrodowiskoWykonania: okno.SrodowiskoWykonania,
		TrybUprawnien:       okno.TrybUprawnien,
		RolaOkna:            okno.RolaOkna,
		Agent:               okno.Agent,
	}
	if sesja, err := nadzorca.Rejestr().Sesja(okno.IdSesji); err == nil {
		opis.TytulSesji = sesja.Tytul
		opis.Projekt = sesja.IdProjektu
	}
	return opis
}
