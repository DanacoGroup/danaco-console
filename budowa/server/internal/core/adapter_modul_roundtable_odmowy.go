// Plik zawiera odmowy modułu Roundtable dobudowanego migracjami 190–199. Każda
// odmowa nazywa trzy rzeczy: co odmówiło, dlaczego i czym to naprawić — bo
// komunikat bez trzeciej części zostawia operatora tam, gdzie był przed
// naciśnięciem przycisku.
package core

import (
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// bladNieznanegoZespolu odróżnia zespół debaty, który nie istnieje w magazynie,
// od usterki samego odczytu, żeby operator dostał komunikat trafny w każdym
// z dwóch przypadków.
func bladNieznanegoZespolu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: zespół debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}

// bladNieznanejWypowiedzi odróżnia wypowiedź debaty, której magazyn nie zna, od
// usterki samego odczytu, żeby operator dostał komunikat trafny w każdym
// z dwóch przypadków.
func bladNieznanejWypowiedzi(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: wypowiedź debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}

// bladNieznanegoWezla odróżnia węzeł grafu argumentów, którego magazyn nie zna,
// od usterki samego odczytu, i wskazuje komendę wydobycia argumentów jako
// drogę założenia grafu.
func bladNieznanegoWezla(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: węzeł grafu argumentów nie istnieje: "+kod+
				". Graf powstaje analizą — uruchom wydobycie argumentów "+
				"(roundtable.analysis.run, rodzaj argumentMining)."))
	}
	return bladDebaty(err)
}

// bladNieznanegoGlosowania odróżnia głosowanie nieistniejące od usterki odczytu
// i od okna bez ani jednego głosowania, dobierając dla każdego z tych stanów
// osobny komunikat.
func bladNieznanegoGlosowania(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		opis := "moduł Roundtable: to okno nie ma ani jednego głosowania — otwórz je " +
			"komendą roundtable.vote.start"
		if kod != "" {
			opis = "moduł Roundtable: głosowanie nie istnieje: " + kod
		}
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound, opis))
	}
	return bladDebaty(err)
}

// bladNieznanejRubryki odróżnia rubrykę oceny, której magazyn nie zna, od
// usterki samego odczytu, i wskazuje komendę zapisu rubryki jako drogę jej
// założenia.
func bladNieznanejRubryki(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: rubryka oceny nie istnieje: "+kod+
				". Załóż ją komendą roundtable.rubric.set."))
	}
	return bladDebaty(err)
}

// bladNieznanejMacierzy odróżnia macierz decyzyjną nieistniejącą od usterki
// odczytu i od okna bez ani jednej macierzy, dobierając dla każdego z tych
// stanów osobny komunikat.
func bladNieznanejMacierzy(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		opis := "moduł Roundtable: to okno nie ma ani jednej macierzy decyzyjnej — " +
			"załóż ją komendą roundtable.decision.matrix.set"
		if kod != "" {
			opis = "moduł Roundtable: macierz decyzyjna nie istnieje: " + kod
		}
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound, opis))
	}
	return bladDebaty(err)
}

// bladNieznanegoStanowiska odróżnia brak stanowiska końcowego od usterki
// odczytu i nazywa dwie drogi, którymi stanowisko powstaje: zamknięcie
// pierwszej tury albo zapis treści wprost.
func bladNieznanegoStanowiska(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: stanowisko końcowe nie istnieje ("+kod+"). "+
				"Powstaje po zamknięciu pierwszej tury albo po zapisaniu treści "+
				"komendą roundtable.consensus.set."))
	}
	return bladDebaty(err)
}

// bladNieznanegoBledu odmawia zawężenia katalogu błędów logicznych do kodu
// spoza katalogu i wskazuje komendę, która wydaje pełny wykaz dozwolonych
// kodów.
func bladNieznanegoBledu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: katalog błędów logicznych nie zna kodu "+kod+". "+
			"Wykaz kodów wydaje roundtable.fallacy.catalog.get."))
}

// bladNieznanegoWariantu odmawia głosu na wariant, którego wskazane głosowanie
// nie ma, i wskazuje komendę, która wydaje wykaz wariantów tego głosowania.
func bladNieznanegoWariantu(kod, glosowanie string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: głosowanie "+glosowanie+" nie ma wariantu "+kod+". "+
			"Warianty wydaje roundtable.vote.get."))
}

// bladNieznanegoKryterium odmawia oceny wariantu w kryterium, którego macierz
// decyzyjna nie zawiera, bo wynik ważony liczony niepełną sumą wag nie dałby
// się porównać z innym wynikiem.
func bladNieznanegoKryterium(kryterium, wariant string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: wariant "+wariant+" niesie ocenę w kryterium "+kryterium+
			", którego macierz nie ma. Wynik ważony byłby wtedy sumą ocen dzieloną "+
			"przez niepełną sumę wag — dopisz kryterium albo zdejmij ocenę."))
}

// odmowaTuryDebatyWBiegu leży w adapter_modul_roundtable_przeklad.go, wśród
// odmów tego modułu.

// odmowaZamknietegoGlosowania odmawia głosu w głosowaniu już zamkniętym, bo głos
// oddany po ogłoszeniu wyniku przestawiłby rozstrzygnięcie, które już
// zapadło.
func odmowaZamknietegoGlosowania(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Roundtable: głosowanie "+kod+" jest zamknięte i nie przyjmuje głosów. "+
			"Głos oddany po ogłoszeniu wyniku przestawiłby rozstrzygnięcie, "+
			"które już zapadło — otwórz nowe głosowanie (roundtable.vote.start)."))
}

// odmowaNieuprawnionegoGlosu odmawia głosu wyborcy spoza wykazu uprawnionych
// ustalonego przy otwarciu głosowania i wskazuje otwarcie nowego głosowania
// jako drogę naprawy.
func odmowaNieuprawnionegoGlosu(wyborca, glosowanie string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Roundtable: "+wyborca+" nie jest uprawniony do głosu w głosowaniu "+
			glosowanie+". Wykaz uprawnionych ustalono przy jego otwarciu; "+
			"otwórz nowe głosowanie bez zawężenia albo z tym wyborcą w wykazie."))
}

// odmowaPustegoGlosu odmawia głosu, który nie wskazuje żadnego wariantu,
// i nazywa kształt danych, jakiego oczekuje metoda głosowania wskazana
// w żądaniu.
func odmowaPustegoGlosu(metoda string) error {
	czym := "wykazem zaaprobowanych wariantów (approvals)"
	switch metoda {
	case shared.RoundtableVoteMethodIrv, shared.RoundtableVoteMethodSchulze:
		czym = "kolejnością preferencji (ranking)"
	case shared.RoundtableVoteMethodScore, shared.RoundtableVoteMethodQuadratic:
		czym = "punktami przypisanymi wariantom (scores)"
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: głos nie wskazuje żadnego wariantu. Metoda "+metoda+
			" oczekuje głosu opisanego "+czym+"."))
}

// odmowaAnalizyBezZapisu odmawia analizy debaty, w której nikt jeszcze nie
// zabrał głosu, i wskazuje uruchomienie tury jako warunek, który musi być
// spełniony wcześniej.
func odmowaAnalizyBezZapisu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: debata okna "+okno+" nie ma ani jednej wypowiedzi do analizy. "+
			"Uruchom turę (roundtable.debate.start) i poczekaj na odpowiedzi uczestników."))
}

// odmowaAnalizyBezWyniku odmawia zapisania analizy, z której nie dało się
// odczytać ani jednego ustalenia, bo odpowiedź powodzenia z pustym wykazem
// byłaby nie do odróżnienia od analizy, która stwierdziła brak.
func odmowaAnalizyBezWyniku(kanal string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: z odpowiedzi kanału "+kanal+" nie da się odczytać ani jednego "+
			"ustalenia analizy. Powtórz analizę albo wskaż inny kanał (channelId) — "+
			"rdzeń nie dopisuje ustaleń, których model nie wypowiedział."))
}

// odmowaScaleniaBezGrafu odmawia scalenia powtórzeń w oknie bez grafu
// argumentów i wskazuje wydobycie argumentów jako komendę, która graf
// zakłada.
func odmowaScaleniaBezGrafu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: okno "+okno+" nie ma grafu argumentów, więc nie ma czego scalać. "+
			"Uruchom wydobycie argumentów (roundtable.analysis.run, rodzaj argumentMining)."))
}

// odmowaPustegoKatalogu odmawia wykrywania błędów logicznych w oknie, w którym
// wyłączono wszystkie kody katalogu, i wskazuje komendę włączenia przynajmniej
// jednego.
func odmowaPustegoKatalogu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: w oknie "+okno+" wyłączono wykrywanie wszystkich błędów logicznych. "+
			"Włącz przynajmniej jeden kod komendą roundtable.fallacy.catalog.set."))
}

// odmowaOcenyBezWypowiedzi odmawia oceny sędziowskiej w oknie bez wypowiedzi do
// oceny i wskazuje dwie drogi naprawy: wskazanie wypowiedzi wprost albo
// uruchomienie tury debaty.
func odmowaOcenyBezWypowiedzi(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: w oknie "+okno+" nie ma wypowiedzi do oceny. "+
			"Wskaż wypowiedzi (targetStatementIds) albo uruchom turę debaty."))
}

// odmowaSumyWag odmawia rubryce, której wagi kryteriów nie sumują się do
// jedności, bo wynik werdyktu liczony inną sumą wag nie dałby się zestawić
// z wynikiem innej rubryki.
func odmowaSumyWag(suma float64) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: wagi kryteriów rubryki sumują się do "+sformatujUlamek(suma)+
			", a mają do jedności. Wynik werdyktu liczony inną sumą wag nie da się "+
			"zestawić z wynikiem z żadnej innej rubryki."))
}

// odmowaNieczytelnychKryteriow odmawia rubryce i macierzy, których pole
// kryteriów nie da się odczytać jako wykazu par nazwy i wagi.
func odmowaNieczytelnychKryteriow(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeValidationFailed,
		errors.New("moduł Roundtable: pola criteria nie da się odczytać jako wykazu "+
			"kryteriów o kształcie [{\"name\":…,\"weight\":…}]: "+err.Error())))
}

// odmowaNieczytelnychWariantow odmawia macierzy, której pole wariantów nie da
// się odczytać jako wykazu par etykiety i wykazu ocen.
func odmowaNieczytelnychWariantow(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeValidationFailed,
		errors.New("moduł Roundtable: pola options nie da się odczytać jako wykazu "+
			"wariantów o kształcie [{\"label\":…,\"scores\":{…}}]: "+err.Error())))
}

// odmowaNieczytelnychGlosow odmawia odsłuchu, którego pole przypisania głosów
// nie da się odczytać jako odwzorowania kodu uczestnika na nazwę głosu.
func odmowaNieczytelnychGlosow(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeValidationFailed,
		errors.New("moduł Roundtable: pola voiceByParticipant nie da się odczytać jako "+
			"odwzorowania {\"kod uczestnika\":\"nazwa głosu\"}: "+err.Error())))
}

// odmowaPustejRegeneracji odmawia zastąpienia wypowiedzi ciszą, bo kanał, który
// przy powtórnym wywołaniu nie zwrócił ani jednego znaku, nie ma prawa
// skasować tego, co uczestnik już powiedział.
func odmowaPustejRegeneracji(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Roundtable: powtórne wywołanie kanału nie przyniosło ani jednego znaku, "+
			"więc wypowiedź "+kod+" zostaje bez zmiany. Kanał, który zamilkł, nie ma prawa "+
			"skasować tego, co uczestnik powiedział za pierwszym razem."))
}

// odmowaKanaluDebaty znakuje niepowodzenie wywołania kanału modelu w toku
// debaty, niosąc dalej usterkę zgłoszoną przez samą warstwę kanału.
func odmowaKanaluDebaty(kanal string, err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeChannelUnavailable,
		errors.New("moduł Roundtable: kanał "+kanal+" nie odpowiedział: "+err.Error())))
}

// odmowaTranskryptuBezZapisu odmawia wydania transkryptu debaty, w której nie
// odbyła się ani jedna tura, i wskazuje uruchomienie tury jako warunek
// wcześniejszy.
func odmowaTranskryptuBezZapisu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: debata okna "+okno+" nie ma ani jednej tury do wydania. "+
			"Uruchom turę komendą roundtable.debate.start."))
}

// odmowaWydaniaBezGrafu odmawia wydania grafu argumentów w oknie, które go nie
// ma, i wskazuje wydobycie argumentów jako komendę, która graf zakłada.
func odmowaWydaniaBezGrafu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: okno "+okno+" nie ma grafu argumentów do wydania. "+
			"Graf powstaje analizą — uruchom roundtable.analysis.run w rodzaju argumentMining."))
}

// odmowaOdsluchuBezZapisu odmawia odsłuchu w oknie bez wypowiedzi objętych
// zawężeniem i wskazuje zdjęcie zawężenia do uczestnika albo do tury jako
// drogę naprawy.
func odmowaOdsluchuBezZapisu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: w oknie "+okno+" nie ma wypowiedzi objętych odsłuchem. "+
			"Zdejmij zawężenie do uczestnika albo do tury."))
}

// odmowaSzablonuBezTury odmawia zapisania szablonu z debaty, w której nie
// odbyła się ani jedna tura, bo format i granice szablonu biorą się właśnie
// z przebiegu tury.
func odmowaSzablonuBezTury(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: debata okna "+okno+" nie ma ani jednej tury, więc nie ma formatu "+
			"ani granic do zapisania jako szablon. Uruchom turę i zapisz szablon po niej."))
}

// odmowaPustegoArtefaktu odmawia zapisania artefaktu, którego wydanie dało zero
// bajtów, bo artefakt zerowej długości przeszedłby każdy sprawdzian istnienia,
// nie mając czego pokazać.
func odmowaPustegoArtefaktu(rodzaj string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Roundtable: wydanie rodzaju "+rodzaj+" dało zero bajtów. Artefakt zerowej "+
			"długości przeszedłby każdy sprawdzian istnienia i nie miałby czego pokazać, "+
			"więc nie zostaje zapisany."))
}

// odmowaBrakuMagazynuDebaty odmawia wydania artefaktu debaty, gdy rdzeń nie ma
// podpiętego katalogu magazynu, i wskazuje miejsce składania rdzenia jako
// miejsce naprawy.
func odmowaBrakuMagazynuDebaty() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Roundtable: rdzeń nie ma katalogu magazynu artefaktów debaty; "+
			"naprawa: podpiąć katalog danych przy składaniu rdzenia (ZKatalogiemDanych)."))
}

// odmowaBrakuUruchamiacza odmawia czynności wymagającej programu serwerowego,
// gdy rdzeń nie ma uruchamiacza procesów, i wskazuje warstwę kanału jako
// brakującą zależność.
func odmowaBrakuUruchamiacza() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Roundtable: rdzeń nie ma uruchamiacza procesów, więc nie ruszy programu "+
			"serwerowego; naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia."))
}

// odmowaZamianyFormatu znakuje niepowodzenie zamiany transkryptu debaty na
// dokument biurowy, niosąc dalej usterkę zgłoszoną przez bibliotekę składania
// dokumentu.
func odmowaZamianyFormatu(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		errors.New("moduł Roundtable: zamiana transkryptu na dokument biurowy nie powiodła się: "+
			err.Error())))
}

// odmowaSyntezyMowy znakuje niepowodzenie silnika syntezy mowy przy odczycie
// transkryptu debaty na głos, niosąc dalej usterkę zgłoszoną przez ten
// silnik.
func odmowaSyntezyMowy(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		errors.New("moduł Roundtable: synteza mowy nie powiodła się: "+err.Error())))
}

// odmowaSkladaniaDokumentu znakuje niepowodzenie biblioteki składania
// dokumentu przy budowie transkryptu debaty jako pliku biurowego.
func odmowaSkladaniaDokumentu(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		errors.New("moduł Roundtable: nie można złożyć dokumentu transkryptu: "+err.Error())))
}
