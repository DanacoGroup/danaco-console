// Odpowiedzialność pliku: odmowy modułu Roundtable dobudowanego migracjami
// 190–199 — każda nazwana, każda mówiąca trzy rzeczy: co odmówiło, dlaczego
// i czym to zmienić.
//
// Odmowa bez trzeciej części jest komunikatem o porażce, a nie wskazówką.
// Operator, który czyta „nie udało się", wie tyle samo co przed naciśnięciem
// przycisku.
package core

import (
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// bladNieznanegoZespolu odróżnia zespół nieistniejący od usterki odczytu.
func bladNieznanegoZespolu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: zespół debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}

// bladNieznanejWypowiedzi odróżnia wypowiedź nieistniejącą od usterki odczytu.
func bladNieznanejWypowiedzi(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: wypowiedź debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}

// bladNieznanegoWezla odróżnia węzeł grafu nieistniejący od usterki odczytu.
func bladNieznanegoWezla(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: węzeł grafu argumentów nie istnieje: "+kod+
				". Graf powstaje analizą — uruchom wydobycie argumentów "+
				"(roundtable.analysis.run, rodzaj argumentMining)."))
	}
	return bladDebaty(err)
}

// bladNieznanegoGlosowania odróżnia głosowanie nieistniejące od usterki.
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

// bladNieznanejRubryki odróżnia rubrykę nieistniejącą od usterki odczytu.
func bladNieznanejRubryki(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: rubryka oceny nie istnieje: "+kod+
				". Załóż ją komendą roundtable.rubric.set."))
	}
	return bladDebaty(err)
}

// bladNieznanejMacierzy odróżnia macierz nieistniejącą od usterki odczytu.
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

// bladNieznanegoStanowiska odróżnia brak stanowiska od usterki odczytu.
func bladNieznanegoStanowiska(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: stanowisko końcowe nie istnieje ("+kod+"). "+
				"Powstaje po zamknięciu pierwszej tury albo po zapisaniu treści "+
				"komendą roundtable.consensus.set."))
	}
	return bladDebaty(err)
}

// bladNieznanegoBledu odmawia zawężenia katalogu do kodu spoza katalogu.
func bladNieznanegoBledu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: katalog błędów logicznych nie zna kodu "+kod+". "+
			"Wykaz kodów wydaje roundtable.fallacy.catalog.get."))
}

// bladNieznanegoWariantu odmawia głosu na wariant spoza głosowania.
func bladNieznanegoWariantu(kod, glosowanie string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: głosowanie "+glosowanie+" nie ma wariantu "+kod+". "+
			"Warianty wydaje roundtable.vote.get."))
}

// bladNieznanegoKryterium odmawia oceny w kryterium, którego macierz nie ma.
func bladNieznanegoKryterium(kryterium, wariant string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: wariant "+wariant+" niesie ocenę w kryterium "+kryterium+
			", którego macierz nie ma. Wynik ważony byłby wtedy sumą ocen dzieloną "+
			"przez niepełną sumę wag — dopisz kryterium albo zdejmij ocenę."))
}

// odmowaTuryDebatyWBiegu leży w `adapter_modul_roundtable_przeklad.go`; tutaj
// stoją odmowy stanu dobudowane wraz z resztą modułu.

// odmowaZamknietegoGlosowania odmawia głosu po rozstrzygnięciu.
func odmowaZamknietegoGlosowania(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Roundtable: głosowanie "+kod+" jest zamknięte i nie przyjmuje głosów. "+
			"Głos oddany po ogłoszeniu wyniku przestawiłby rozstrzygnięcie, "+
			"które już zapadło — otwórz nowe głosowanie (roundtable.vote.start)."))
}

// odmowaNieuprawnionegoGlosu odmawia głosu spoza wykazu uprawnionych.
func odmowaNieuprawnionegoGlosu(wyborca, glosowanie string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Roundtable: "+wyborca+" nie jest uprawniony do głosu w głosowaniu "+
			glosowanie+". Wykaz uprawnionych ustalono przy jego otwarciu; "+
			"otwórz nowe głosowanie bez zawężenia albo z tym wyborcą w wykazie."))
}

// odmowaPustegoGlosu odmawia głosu, który niczego nie wskazuje.
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

// odmowaAnalizyBezZapisu odmawia analizy debaty, w której nikt nie zabrał głosu.
func odmowaAnalizyBezZapisu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: debata okna "+okno+" nie ma ani jednej wypowiedzi do analizy. "+
			"Uruchom turę (roundtable.debate.start) i poczekaj na odpowiedzi uczestników."))
}

// odmowaAnalizyBezWyniku odmawia zapisania analizy, z której nic nie wyszło.
//
// Odmowa jest tu właściwsza od pustej odpowiedzi udanej: kanał odpowiedział,
// ale z jego odpowiedzi nie da się odczytać ani jednego ustalenia. Odpowiedź
// `ok` z pustym wykazem byłaby nie do odróżnienia od analizy, która stwierdziła
// brak — a to dwie różne rzeczy.
func odmowaAnalizyBezWyniku(kanal string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: z odpowiedzi kanału "+kanal+" nie da się odczytać ani jednego "+
			"ustalenia analizy. Powtórz analizę albo wskaż inny kanał (channelId) — "+
			"rdzeń nie dopisuje ustaleń, których model nie wypowiedział."))
}

// odmowaScaleniaBezGrafu odmawia scalenia powtórzeń przed wydobyciem argumentów.
func odmowaScaleniaBezGrafu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: okno "+okno+" nie ma grafu argumentów, więc nie ma czego scalać. "+
			"Uruchom wydobycie argumentów (roundtable.analysis.run, rodzaj argumentMining)."))
}

// odmowaPustegoKatalogu odmawia wykrywania błędów przy katalogu wyłączonym
// w całości.
func odmowaPustegoKatalogu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: w oknie "+okno+" wyłączono wykrywanie wszystkich błędów logicznych. "+
			"Włącz przynajmniej jeden kod komendą roundtable.fallacy.catalog.set."))
}

// odmowaOcenyBezWypowiedzi odmawia oceny sędziowskiej bez materiału.
func odmowaOcenyBezWypowiedzi(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: w oknie "+okno+" nie ma wypowiedzi do oceny. "+
			"Wskaż wypowiedzi (targetStatementIds) albo uruchom turę debaty."))
}

// odmowaSumyWag odmawia rubryce, której wagi nie sumują się do jedności.
func odmowaSumyWag(suma float64) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: wagi kryteriów rubryki sumują się do "+sformatujUlamek(suma)+
			", a mają do jedności. Wynik werdyktu liczony inną sumą wag nie da się "+
			"zestawić z wynikiem z żadnej innej rubryki."))
}

// odmowaNieczytelnychKryteriow odmawia rubryce i macierzy o nieczytelnym
// ładunku kryteriów.
func odmowaNieczytelnychKryteriow(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeValidationFailed,
		errors.New("moduł Roundtable: pola criteria nie da się odczytać jako wykazu "+
			"kryteriów o kształcie [{\"name\":…,\"weight\":…}]: "+err.Error())))
}

// odmowaNieczytelnychWariantow odmawia macierzy o nieczytelnym ładunku wariantów.
func odmowaNieczytelnychWariantow(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeValidationFailed,
		errors.New("moduł Roundtable: pola options nie da się odczytać jako wykazu "+
			"wariantów o kształcie [{\"label\":…,\"scores\":{…}}]: "+err.Error())))
}

// odmowaNieczytelnychGlosow odmawia odsłuchu o nieczytelnym przypisaniu głosów.
func odmowaNieczytelnychGlosow(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeValidationFailed,
		errors.New("moduł Roundtable: pola voiceByParticipant nie da się odczytać jako "+
			"odwzorowania {\"kod uczestnika\":\"nazwa głosu\"}: "+err.Error())))
}

// odmowaPustejRegeneracji odmawia zastąpienia wypowiedzi ciszą.
func odmowaPustejRegeneracji(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Roundtable: powtórne wywołanie kanału nie przyniosło ani jednego znaku, "+
			"więc wypowiedź "+kod+" zostaje bez zmiany. Kanał, który zamilkł, nie ma prawa "+
			"skasować tego, co uczestnik powiedział za pierwszym razem."))
}

// odmowaKanaluDebaty znakuje niepowodzenie wywołania kanału.
func odmowaKanaluDebaty(kanal string, err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeChannelUnavailable,
		errors.New("moduł Roundtable: kanał "+kanal+" nie odpowiedział: "+err.Error())))
}

// odmowaTranskryptuBezZapisu odmawia wydania transkryptu debaty, której nie było.
func odmowaTranskryptuBezZapisu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: debata okna "+okno+" nie ma ani jednej tury do wydania. "+
			"Uruchom turę komendą roundtable.debate.start."))
}

// odmowaWydaniaBezGrafu odmawia wydania grafu, którego nie ma.
func odmowaWydaniaBezGrafu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: okno "+okno+" nie ma grafu argumentów do wydania. "+
			"Graf powstaje analizą — uruchom roundtable.analysis.run w rodzaju argumentMining."))
}

// odmowaOdsluchuBezZapisu odmawia odsłuchu, gdy nie ma czego odczytać.
func odmowaOdsluchuBezZapisu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: w oknie "+okno+" nie ma wypowiedzi objętych odsłuchem. "+
			"Zdejmij zawężenie do uczestnika albo do tury."))
}

// odmowaSzablonuBezTury odmawia zapisania szablonu z debaty, która nie ruszyła.
func odmowaSzablonuBezTury(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: debata okna "+okno+" nie ma ani jednej tury, więc nie ma formatu "+
			"ani granic do zapisania jako szablon. Uruchom turę i zapisz szablon po niej."))
}

// odmowaPustegoArtefaktu odmawia zapisania artefaktu o zerowej długości.
func odmowaPustegoArtefaktu(rodzaj string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Roundtable: wydanie rodzaju "+rodzaj+" dało zero bajtów. Artefakt zerowej "+
			"długości przeszedłby każdy sprawdzian istnienia i nie miałby czego pokazać, "+
			"więc nie zostaje zapisany."))
}

// odmowaBrakuMagazynuDebaty odmawia wydania bez wskazanego magazynu.
func odmowaBrakuMagazynuDebaty() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Roundtable: rdzeń nie ma katalogu magazynu artefaktów debaty; "+
			"naprawa: podpiąć katalog danych przy składaniu rdzenia (ZKatalogiemDanych)."))
}

// odmowaBrakuUruchamiacza odmawia czynności wymagającej programu serwerowego.
func odmowaBrakuUruchamiacza() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Roundtable: rdzeń nie ma uruchamiacza procesów, więc nie ruszy programu "+
			"serwerowego; naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia."))
}

// odmowaZamianyFormatu znakuje niepowodzenie zamiany formatu dokumentu.
func odmowaZamianyFormatu(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		errors.New("moduł Roundtable: zamiana transkryptu na dokument biurowy nie powiodła się: "+
			err.Error())))
}

// odmowaSyntezyMowy znakuje niepowodzenie silnika mowy.
func odmowaSyntezyMowy(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		errors.New("moduł Roundtable: synteza mowy nie powiodła się: "+err.Error())))
}

// odmowaSkladaniaDokumentu znakuje niepowodzenie biblioteki dokumentu.
func odmowaSkladaniaDokumentu(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		errors.New("moduł Roundtable: nie można złożyć dokumentu transkryptu: "+err.Error())))
}
