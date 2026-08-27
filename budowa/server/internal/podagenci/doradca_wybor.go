// Dobór doradcy pod sufitem siły: kto może zostać zapytany, gdy o radę prosi
// model w trakcie tury, a nie Operator; model z własnej inicjatywy nie sięga
// po model silniejszy od siebie.
package podagenci

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"danacoconsole/server/internal/models"
)

// Parametry wiersza rejestru kanałów, którymi Operator opisuje doradcę;
// doradca jest prawdą logiczną, sila liczbą całkowitą.
const (
	parametrDoradcy = "doradca"
	parametrSily    = "sila"
)

// Kandydaci wybiera z wykazu kanały dopuszczone do roli doradcy: czynne
// i z parametrem `doradca`. Porządkowanie należy do wyboru, nie odczytu.
func Kandydaci(wykaz []models.Definicja) []Kandydat {
	kandydaci := make([]Kandydat, 0, len(wykaz))
	for _, wiersz := range wykaz {
		if !wiersz.Aktywny || !prawda(wiersz.Parametr(parametrDoradcy)) {
			continue
		}
		kandydaci = append(kandydaci, kandydatZWiersza(wiersz))
	}
	return kandydaci
}

// Odmowy doboru. Każda opisuje brak, czego zabrakło, żeby konsultacja mogła
// się odbyć, i każda jest osobna dla wołającego.
var (
	// ErrEskalacjaBezWskazania: żądany kanał jest silniejszy od pytającego.
	// Nazwany brak, nie zastępczy wybór — dlatego wołający dostaje odmowę,
	// a nie po cichu słabszego doradcę.
	ErrEskalacjaBezWskazania = errors.New(
		"podagenci: konsultacja u modelu silniejszego wymaga wcześniejszego wskazania Operatora; " +
			"z własnej inicjatywy model konsultuje model równy sobie albo słabszy")

	// ErrSilaNieopisana: sufitu nie da się przyłożyć, bo któraś z dwóch sił nie
	// jest opisana liczbą. Odmowa idzie w stronę zamkniętą: „nie wiem, czy
	// równy" nie jest tym samym co „równy", a sufit dopuszcza parametry takie
	// same albo niższe — nie nieznane.
	ErrSilaNieopisana = errors.New(
		"podagenci: siła kanału nie jest opisana liczbą w parametrze `sila`, " +
			"więc nie ma czym zmierzyć sufitu — Operator wpisuje `sila` obu kanałom albo model pyta bez prośby")

	// ErrDoradcaNieznany: wskazanego kanału nie ma w rejestrze albo jest
	// wygaszony. Milczące zejście na innego doradcę byłoby podmianą modelu za
	// plecami proszącego.
	ErrDoradcaNieznany = errors.New(
		"podagenci: żądany kanał doradcy nie jest czynnym wierszem rejestru kanałów")

	// ErrDoradcaNiedopuszczony: wiersz stoi i jest czynny, ale Operator nie
	// dopuścił go do radzenia parametrem `doradca`. Bez tego sprawdzenia
	// deklaracja „kto może być doradcą" (`doradca.go`) nie obowiązuje na
	// jedynej drodze, którą chodzi model.
	ErrDoradcaNiedopuszczony = errors.New(
		"podagenci: żądany kanał nie jest dopuszczony do radzenia — brak parametru `doradca` w wierszu rejestru")

	// ErrBrakPytajacego zostaje odmową doboru bez żadnych wskazań: bez kanału
	// pytającego nie ma ani progu sufitu, ani kanału zapasowego równego jemu.
	ErrBrakPytajacego = errors.New(
		"podagenci: pytanie do doradcy bez kanału pytającego — nie ma czym zmierzyć sufitu siły")
)

// WybierzDoradce rozstrzyga, kogo wolno zapytać o radę, i jest jedynym
// miejscem, w którym stoi sufit siły, dwiema regułami w ustalonej kolejności.
func WybierzDoradce(wykaz []models.Definicja, kanalPytajacego, zadanyPrzezModel string) (Kandydat, error) {
	pytajacy := strings.TrimSpace(kanalPytajacego)
	if pytajacy == "" {
		return Kandydat{}, ErrBrakPytajacego
	}
	prog, progZnany := silaKanalu(wykaz, pytajacy)

	// Reguła 1 — prośba modelu przechodzi wyłącznie pod sufitem.
	if zadany := strings.TrimSpace(zadanyPrzezModel); zadany != "" {
		kandydat, err := dopuszczonyZWykazu(wykaz, zadany, pytajacy)
		if err != nil {
			return Kandydat{}, err
		}
		if err := podSufit(kandydat, pytajacy, prog, progZnany); err != nil {
			return Kandydat{}, err
		}
		return kandydat, nil
	}

	// Reguła 2 — najsilniejszy kandydat pod sufitem, a w jego braku pytający sam.
	return podSufitem(wykaz, pytajacy, prog, progZnany), nil
}

// podSufit przykłada sufit do jednego kandydata i nazywa powód odmowy; kanał
// samego pytającego przechodzi bez mierzenia.
func podSufit(kandydat Kandydat, pytajacy string, prog int, progZnany bool) error {
	if kandydat.Kanal == pytajacy {
		return nil
	}
	if !progZnany || !kandydat.SilaZnana {
		return fmt.Errorf("%w (kanał żądany %q, kanał pytającego %q)",
			ErrSilaNieopisana, kandydat.Kanal, pytajacy)
	}
	if kandydat.Sila > prog {
		return fmt.Errorf("%w (żądany kanał %q ma siłę %d, kanał pytającego %q siłę %d)",
			ErrEskalacjaBezWskazania, kandydat.Kanal, kandydat.Sila, pytajacy, prog)
	}
	return nil
}

// podSufitem wybiera najsilniejszego kandydata o sile nie wyższej od progu,
// a przy braku innych kandydatów kanał pytającego.
func podSufitem(wykaz []models.Definicja, pytajacy string, prog int, progZnany bool) Kandydat {
	var wybrany Kandydat
	if progZnany {
		for _, kandydat := range Kandydaci(wykaz) {
			if !kandydat.SilaZnana || kandydat.Sila > prog {
				continue
			}
			if wybrany.Kanal == "" || kandydat.Sila > wybrany.Sila {
				wybrany = kandydat
			}
		}
	}
	if wybrany.Kanal != "" {
		return wybrany
	}
	return kandydatPytajacego(wykaz, pytajacy)
}

// kandydatPytajacego opisuje kanał pytającego jako doradcę samego dla siebie,
// bez zmyślania danych wiersza spoza rejestru.
func kandydatPytajacego(wykaz []models.Definicja, pytajacy string) Kandydat {
	for _, wiersz := range wykaz {
		if wiersz.Kod == pytajacy {
			return kandydatZWiersza(wiersz)
		}
	}
	return Kandydat{Kanal: pytajacy}
}

// dopuszczonyZWykazu odnajduje czynny wiersz o podanym kodzie i sprawdza, czy
// Operator dopuścił go do radzenia.
func dopuszczonyZWykazu(wykaz []models.Definicja, kod, pytajacy string) (Kandydat, error) {
	for _, wiersz := range wykaz {
		if wiersz.Kod != kod || !wiersz.Aktywny {
			continue
		}
		if !prawda(wiersz.Parametr(parametrDoradcy)) && kod != pytajacy {
			return Kandydat{}, fmt.Errorf("%w (kanał %q)", ErrDoradcaNiedopuszczony, kod)
		}
		return kandydatZWiersza(wiersz), nil
	}
	return Kandydat{}, fmt.Errorf("%w (kanał %q)", ErrDoradcaNieznany, kod)
}

func kandydatZWiersza(d models.Definicja) Kandydat {
	moc, znana := sila(d)
	return Kandydat{Kanal: d.Kod, Nazwa: d.Nazwa, Model: d.Model, Sila: moc, SilaZnana: znana}
}

// silaKanalu odczytuje siłę wiersza o podanym kodzie; wiersza nieznanego nie
// zgadujemy — jego siła jest nieznana tak samo jak siła wiersza nieopisanego.
func silaKanalu(wykaz []models.Definicja, kod string) (int, bool) {
	for _, wiersz := range wykaz {
		if wiersz.Kod == kod {
			return sila(wiersz)
		}
	}
	return 0, false
}

// sila czyta parametr siły wiersza; drugi wynik mówi, czy wiemy, nigdy nie
// orzeka zero z zapisu nieczytelnego.
func sila(d models.Definicja) (int, bool) {
	zapis := strings.TrimSpace(d.Parametr(parametrSily))
	if zapis == "" {
		return 0, false
	}
	liczba, err := strconv.Atoi(zapis)
	if err != nil {
		return 0, false
	}
	return liczba, true
}

// prawda czyta parametr logiczny; poza zapisem Go uznaje polskie „tak", bo
// parametry wpisuje Operator ręcznie.
func prawda(wartosc string) bool {
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "true", "1", "tak":
		return true
	default:
		return false
	}
}
