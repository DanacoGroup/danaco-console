// Odpowiedzialność pliku: okno komunikacji i jego obszar widziane przez moduł
// Terminal — tryb uprawnień, sesja, katalog własny — oraz odłożenie karty do
// dziennika i znakowanie odmów kodami kontraktu.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// oknoWykonania zwraca okno komunikacji modułu Terminal, w którym właśnie
// pracuje karta wskazana wywołaniem.
func (a *adapterTerminala) oknoWykonania(oknoKod string) (session.Okno, error) {
	if a.okna == nil {
		return session.Okno{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Terminal: serwer nie ma rejestru okien, więc nie zna trybu uprawnień okna"))
	}
	okno, err := a.okna.Okno(oknoKod)
	if err != nil {
		return session.Okno{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Terminal: okno "+oknoKod+" nie jest otwarte"))
	}
	return okno, nil
}

// sesjaOkna zwraca sesję okna albo pustkę, gdy okna już nie ma. Karta
// odtworzona po restarcie nie musi mieć żywego okna — sesja służy wyłącznie
// zaadresowaniu zdarzenia.
func (a *adapterTerminala) sesjaOkna(oknoKod string) string {
	if a.okna == nil {
		return ""
	}
	okno, err := a.okna.Okno(oknoKod)
	if err != nil {
		return ""
	}
	return okno.IdSesji
}

// obszarOkna składa obszar własny okna terminala z ustalonego katalogu
// roboczego jego własnej powłoki.
func (a *adapterTerminala) obszarOkna(okno session.Okno) session.Obszar {
	if a.katalog == nil {
		return session.Obszar{IdOkna: okno.Id}
	}
	return ObszarOkna(a.katalog.Ustal(konfig.Kontekst{Okno: okno.Id}, okno.IdSesji), okno.Id)
}

// zapiszKarte odkłada kartę do dziennika. Nieudany zapis nie unieważnia karty
// otwartej w pamięci; jedynym skutkiem jest brak jej odtworzenia po restarcie.
func (a *adapterTerminala) zapiszKarte(ctx context.Context, karta *kartaTerminala) {
	if a.repozytorium == nil {
		return
	}
	wiersz := dane.KartaTerminala{
		Kod:            karta.kod,
		OknoKod:        karta.oknoKod,
		Powloka:        karta.powloka,
		Tytul:          wskaznikTekstu(karta.tytul),
		KatalogRoboczy: wskaznikTekstu(karta.katalog),
		Stan:           karta.stan,
		CelZdalny:      karta.celZdalny,
		HostKod:        wskaznikTekstu(karta.hostKod),
	}
	if karta.portZdalny > 0 {
		port := int64(karta.portZdalny)
		wiersz.PortZdalny = &port
	}
	_ = a.repozytorium.ZapiszKarte(ctx, wiersz)
}

// bladBrakuZasobuTerminala odmawia czynności na bycie, którego rdzeń nie zna.
// Nazywa byt, a nie samo „nie znaleziono”: Operator ma wiedzieć, czego brakuje.
func bladBrakuZasobuTerminala(co string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Terminal: "+co+" nie występuje w serwerze"))
}

// protocolBladTerminala składa odmowę modułu wskazanym kodem kontraktu. Jedno
// miejsce, w którym powstaje przedrostek „moduł Terminal” — inaczej odmowy tej
// samej rodziny brzmiałyby różnie w zależności od pliku.
func protocolBladTerminala(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "moduł Terminal: "+powod))
}

// bladZadaniaTerminala znakuje wadę żądania kodem kontraktu, oddzielając ją od
// usterki wewnętrznej rdzenia.
func bladZadaniaTerminala(powod string) error {
	return protocolBladTerminala(shared.ErrorCodeValidationFailed, powod)
}

// bladBrakuProcesu odmawia czynności na procesie, którego rdzeń już nie zna,
// kodem kontraktu `not_found`.
func bladBrakuProcesu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Terminal: proces "+kod+" nie występuje ani w rejestrze serwera, ani w dzienniku"))
}
