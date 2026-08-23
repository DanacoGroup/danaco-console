// Odpowiedzialność pliku: dobór doradcy pod sufitem siły — kto może zostać
// zapytany, gdy o radę prosi model w trakcie tury, a nie Operator.
//
// Model proszący o radę z własnej inicjatywy nie sięga po model silniejszy od
// siebie. Doradca zostaje, konsultacja zostaje, znika wyłącznie ruch modelu
// w górę. Ruch w górę robi Operator albo nikt.
//
// Dwie reguły doboru i ich kolejność:
//
//  1. Prośba modelu podlega sufitowi. Kanał, o który prosi sam wołający,
//     przechodzi tylko wtedy, gdy jest dopuszczony do radzenia (czynny wiersz
//     z parametrem `doradca`, albo kanał samego pytającego) i gdy da się
//     wykazać, że nie jest silniejszy od kanału pytającego. Prośba, której nie
//     da się wykazać, kończy się odmową opisującą brak, a nie cichym
//     podstawieniem słabszego.
//  2. Bez prośby rdzeń bierze najsilniejszego kandydata o sile nie wyższej niż
//     pytający, a gdy takiego nie ma — kanał samego pytającego, czyli dokładnie
//     ten sam model. Sufit zwężony do równości jest zawsze wykonalny, więc
//     reguła 2 nie potrzebuje odmowy.
//
// Reguły „wskazanie Operatora znosi sufit" tu nie ma: produkt nie ma ani pola
// „doradca okna", ani komendy, którą Operator by doradcę wskazał. `modelChannelId`
// okna mówi, którym modelem pracuje okno, a nie kto jest doradcą, więc sufitu nie
// znosi.
//
// Skąd porządek modeli — z danych, nie z kodu. Jedyną miarą siły, jaką produkt
// ma, jest parametr `sila` wiersza rejestru kanałów (`kanal_modelu.parametry`).
// Nie ma w produkcie ani katalogu modeli z rangami, ani pola rangi
// w `shared/contract.json`, ani zestawu początkowego, który by `sila` wypełniał.
// Siła niewpisana jest więc nieznana, a nie zerowa.
//
// Siła nieznana znaczy: nie ma czym zmierzyć, więc w górę się nie idzie
// i kandydatem taki kanał nie jest. Zero wpuszczałoby kanał bez `sila` (a także
// literówkę czy wartość ułamkową w parametrze) wszędzie, otwierając sufit
// w obie strony. Stan produkcyjny (nikt nie wpisał `sila`) daje wtedy kanał
// pytającego, czyli ten sam model — dokładnie tyle, ile porządek z danych
// pozwala orzec.
package podagenci

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"danacoconsole/server/internal/models"
)

// Parametry wiersza rejestru kanałów, którymi Operator opisuje doradcę.
// Wartości są danymi (kolumna `kanal_modelu.parametry`), nie nazwami typów:
// `doradca` czytany jest jako prawda logiczna, `sila` jako liczba całkowita
// (brak albo zapis nieczytelny znaczy: siła nieznana — patrz nagłówek).
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

// Odmowy doboru. Każda opisuje brak — czego zabrakło, żeby konsultacja mogła
// się odbyć — i każda jest osobna, bo wołający przekłada je na różne kody
// kontraktu: brak wiersza to `not_found`, brak podstawy do sięgnięcia w górę to
// odmowa trwała, której ponawiać nie ma po co (`adapter_doradcy.go`).
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
// miejscem, w którym stoi sufit siły. Dwie reguły z nagłówka pliku widać tu
// jako dwa kolejne rozgałęzienia; kolejność jest zamierzona, nie przypadkowa.
//
// `zadanyPrzezModel` to prośba samego wołającego — wolno jej zawęzić wybór,
// nie wolno jej podnieść sufitu.
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

// podSufit przykłada sufit do jednego kandydata i nazywa powód odmowy.
//
// Kanał samego pytającego przechodzi bez mierzenia i jest to jedyne odstępstwo:
// równość z samym sobą zachodzi z definicji, więc nie ma czego wykazywać nawet
// wtedy, gdy `sila` nie jest wpisana. Dzięki temu model, który prosi wprost
// o swój własny kanał, dostaje to samo, co dostałby bez prośby.
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

// podSufitem wybiera najsilniejszego kandydata o sile nie wyższej od progu.
//
// Kanał pytającego jest tu zapasem, a nie wykluczeniem: model o takich samych
// parametrach jest dozwolonym rozmówcą, a przy braku innych kandydatów jedynym
// możliwym. Konsultacja u modelu równego nadal ma sens, bo doradca dostaje ramę
// konsultacji i czyste pytanie zamiast całej historii tury.
//
// Próg nieznany nie wpuszcza nikogo obcego: wtedy o żadnym kandydacie nie da
// się orzec, że nie jest silniejszy, więc zostaje kanał pytającego. Tak wygląda
// stan produkcyjny, dopóki Operator nie wpisze `sila` — i tak ma wyglądać.
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

// kandydatPytajacego opisuje kanał pytającego jako doradcę samego dla siebie.
// Wiersza spoza rejestru nie zmyślamy bogaciej, niż wiemy: zostaje kod, który
// podał wołający, a nazwa i model puste — ślad powie wtedy prawdę
// o tym, że kanał nie stał w wykazie.
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
//
// Sprawdzenie parametru `doradca` pilnuje, żeby prośba modelu nie omijała
// wykazu kandydatów i nie sięgała po dowolny czynny kanał rejestru — także
// taki, którego Operator do radzenia nie dopuścił.
//
// Wiersz wyłączony jest tu tym samym co nieistniejący: kanał, który Operator
// zgasił, nie odpowie. Kanał samego pytającego przechodzi bez parametru
// `doradca` — pyta wtedy sam siebie, a na to nie potrzeba dopuszczenia, którego
// reguła 2 też nie wymaga.
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

// sila czyta parametr siły wiersza. Drugi wynik mówi, czy wiemy: parametr
// pusty, nieliczbowy, ułamkowy albo przepełniający zakres `int` daje „nie
// wiem", nigdy zero. Zero jest wartością siły najsłabszej i orzekanie go
// z zapisu, którego nie umiemy przeczytać, wpuszczałoby taki kanał wszędzie.
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
