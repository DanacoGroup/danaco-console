// Odpowiedzialność pliku: wpięcie utrwalania rozmowy — przełożenie okna żywego
// z rejestru nadzorcy na opis warstwy danych.
//
// Warstwa danych nie zna pakietu sesji, a pakiet sesji nie zna bazy.
// Ten plik jest jedynym miejscem, w którym oba się spotykają, i nie robi nic
// poza tym spotkaniem.
package core

import (
	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
)

// utrwalaczRozmow wpina repozytoria sesji, okien i wiadomości w dziennik rozmów.
// Utrwalacz domyka łańcuch
// `srodowisko → modul → karta_sesji → sesja → okno_komunikacji → wiadomosc`.
//
// Warstwa danych nie zna pakietu sesji, więc opis okna podaje tutejsza
// funkcja czytająca rejestr nadzorcy. Okno nieznane rejestrowi daje fałsz —
// wtedy wiadomość zostaje w buforze pamięci, a rozmowa toczy się dalej.
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
