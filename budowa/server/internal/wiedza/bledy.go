// Brak silnika osadzeń kończy się odmową nazywającą brak, a nie cichym zejściem
// na wyszukiwanie po słowach (`library.file.search`). Odpowiedź trafień po
// literach ma ten sam kształt co odpowiedź trafień po znaczeniu, więc pytający
// nie rozpoznałby podmiany i uznałby, że szukanej treści w wiedzy nie ma.
//
// Treść odmowy niesie trzy człony — czego nie ma, ile waży to, czego nie ma,
// i co zrobić, żeby było — tak samo jak odmowy w `zewnetrzne/wolanie.go`.
package wiedza

import "strings"

// Rodzaje braku, którymi odpowiada pomocnik. Napisy są kopią wartości pola
// `brak` w `pomocnik_osadzen.py` co do znaku — rozjazd zamieniłby nazwany brak
// w brak nierozpoznany, czyli w gorszą odmowę.
const (
	// brakBiblioteki — interpreter stoi, ale nie ma w nim `fastembed`.
	brakBiblioteki = "biblioteka"
	// brakModelu — biblioteka stoi, ale wag modelu nie da się przygotować.
	brakModelu = "model"
	// brakWagStojacych — w katalogu wskazanym nastawą model LEŻY, ale nie da się
	// na nim postawić silnika. Osobny rodzaj od `brakModelu`, bo naprawa jest
	// odwrotna: tam trzeba wagi ściągnąć, tu są już na dysku i odesłanie po nie
	// kierowałoby Operatora po to, co ma.
	brakWagStojacych = "wagi"
	// BrakInterpretera — nie udało się uruchomić samego Pythona.
	BrakInterpretera = "interpreter"
)

// BrakSilnika mówi, że wskaźnika znaczenia nie ma czym zbudować ani przeszukać.
//
// Osobny typ, tak samo jak `zewnetrzne.BrakNarzedzia`: brak usuwa się jedną
// instalacją, więc adapter rozpoznaje ten typ i znakuje go jako niedostępność
// zaplecza (`adapter_modul_wiedza.go`), a nie jako usterkę wewnętrzną. Różnica
// jest cała w treści, którą czytelnik dostaje: usterka wewnętrzna nie mówi ani
// czego brak, ani ile to waży, ani co zainstalować.
type BrakSilnika struct {
	// Rodzaj — jedna ze stałych wyżej.
	Rodzaj string
	// Model — nazwa modelu, którego dotyczy brak.
	Model string
	// WagaMb — ile waży do dociągnięcia; 0 znaczy „waga nieznana" i wtedy człon
	// o wadze po prostu nie wchodzi do zdania.
	WagaMb int
	// Powod — to, co powiedział o sobie sam pomocnik albo system.
	Powod string
}

// Error składa trzyczłonową odmowę: czego nie ma, ile to waży, co zrobić.
func (b *BrakSilnika) Error() string {
	zdanie := strings.Builder{}
	zdanie.WriteString("wskaźnik znaczenia: ")
	zdanie.WriteString(b.opisBraku())
	if b.Powod != "" {
		zdanie.WriteString("; pomocnik powiedział: " + skroc(b.Powod, 400))
	}
	zdanie.WriteString("; wyszukiwania po ZNACZENIU nie da się wykonać bez silnika, " +
		"a rdzeń nie zejdzie po cichu na wyszukiwanie po SŁOWACH — " +
		"model dostałby trafienia po literach w miejscu, w którym prosił o trafienia po sensie")
	zdanie.WriteString("; naprawa: " + b.opisNaprawy())
	return zdanie.String()
}

// opisBraku nazywa brak wraz z wagą tego, czego nie ma.
func (b *BrakSilnika) opisBraku() string {
	switch b.Rodzaj {
	case brakBiblioteki:
		return "na tej maszynie nie ma biblioteki osadzeń `fastembed` " +
			"(sama biblioteka wraz ze środowiskiem wykonawczym ONNX to około 120 MB pobrania, " +
			"a wagi modelu " + b.Model + " " + b.opisWagi() + ")"
	case brakModelu:
		return "biblioteka osadzeń stoi, ale wag modelu " + b.Model +
			" nie ma na dysku i nie dało się ich przygotować (" + b.opisWagi() + " do pobrania)"
	case brakWagStojacych:
		return "w katalogu wskazanym ustawieniem `" + KluczKatalogModeli + "` leżą wagi, " +
			"ale nie da się na nich postawić silnika modelu " + b.Model
	case BrakInterpretera:
		return "nie ma czym uruchomić pomocnika osadzeń — interpreter Pythona " +
			"nie wystartował"
	default:
		return "silnik osadzeń nie odpowiedział zrozumiale"
	}
}

// opisWagi oddaje wagę słowami. Wartość niedodatnia znaczy „waga nieznana"
// i wtedy żadna liczba nie wchodzi do zdania.
func (b *BrakSilnika) opisWagi() string {
	if b.WagaMb <= 0 {
		return "waga nieznana"
	}
	if b.WagaMb >= 1000 {
		return "około " + liczbaZPrzecinkiem(b.WagaMb) + " GB"
	}
	return "około " + liczba(b.WagaMb) + " MB"
}

// opisNaprawy mówi, co dokładnie zrobić — z nazwą pakietu i nazwą ustawienia.
func (b *BrakSilnika) opisNaprawy() string {
	switch b.Rodzaj {
	case brakModelu:
		return "dać maszynie dostęp do sieci przy pierwszym budowaniu wskaźnika " +
			"albo przenieść pobrane wagi do katalogu ustawienia `" + KluczKatalogModeli + "`"
	case brakWagStojacych:
		return "uzupełnić przy wagach to, czego pomocnik nie znalazł (powód wyżej nazywa " +
			"plik) — wykaz warstw modelu `modules.json` wraz z opisem warstwy łączącej " +
			"tokeny w jeden wektor niosą sposób łączenia, normalizację i wymiar, których " +
			"pomocnik nie zgaduje; albo wskazać ustawieniem `" + KluczKatalogModeli +
			"` katalog pusty i pozwolić bibliotece pobrać własne wydanie modelu `" +
			KluczModel + "`"
	default:
		return "zainstalować bibliotekę poleceniem `python3 -m pip install fastembed` " +
			"w interpreterze wskazanym ustawieniem `" + KluczProgram + "` " +
			"(pusta wartość znaczy `python3` ze ścieżki wyszukiwania systemu)"
	}
}

// liczba wypisuje liczbę całkowitą bez sięgania po strconv — pakiet i tak nie
// formatuje niczego innego.
func liczba(wartosc int) string {
	if wartosc == 0 {
		return "0"
	}
	cyfry := []byte{}
	for wartosc > 0 {
		cyfry = append([]byte{byte('0' + wartosc%10)}, cyfry...)
		wartosc /= 10
	}
	return string(cyfry)
}

// liczbaZPrzecinkiem przelicza megabajty na gigabajty z jedną cyfrą po
// przecinku — „około 1,0 GB" czyta się lepiej niż „około 1000 MB".
func liczbaZPrzecinkiem(megabajty int) string {
	return liczba(megabajty/1000) + "," + liczba((megabajty%1000)/100)
}

// skroc przycina cudzy komunikat, żeby ślad stosu Pythona nie przesłonił
// odmowy. Cięcie cofa się do początku znaku UTF-8, bo granica liczona jest
// w bajtach i wypadłaby w środku litery.
func skroc(tekst string, granica int) string {
	tekst = strings.TrimSpace(tekst)
	if len(tekst) <= granica {
		return tekst
	}
	for granica > 0 && tekst[granica]&0xC0 == 0x80 {
		granica--
	}
	return tekst[:granica] + "…"
}
