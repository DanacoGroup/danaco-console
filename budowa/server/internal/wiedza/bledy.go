// Brak silnika osadzeń kończy się odmową nazywającą brak, a nie cichym zejściem
// na wyszukiwanie po słowach (`library.file.search`). Treść odmowy niesie trzy
// człony: czego nie ma, ile to waży, i co zrobić, żeby było.
package wiedza

import "strings"

// Rodzaje braku, którymi odpowiada pomocnik. Napisy są kopią wartości pola
// `brak` w `pomocnik_osadzen.py` co do znaku — rozjazd zamieniłby nazwany brak
// w brak nierozpoznany, czyli w gorszą odmowę.
const (
	// brakBiblioteki — interpreter Pythona stoi, ale nie ma w nim zainstalowanej
	// biblioteki o nazwie fastembed.
	brakBiblioteki = "biblioteka"
	// brakModelu — biblioteka fastembed stoi gotowa, ale wag wskazanego modelu
	// nie da się w niej przygotować.
	brakModelu = "model"
	// brakWagStojacych — w katalogu wskazanym nastawą model leży, ale nie da
	// się na nim postawić silnika osadzeń.
	brakWagStojacych = "wagi"
	// BrakInterpretera — nie udało się w ogóle uruchomić samego interpretera
	// Pythona na tej maszynie roboczej.
	BrakInterpretera = "interpreter"
)

// Silniki pakietu. Wartość pusta znaczy silnik osadzeń — pierwszy i jedyny
// przez cały czas, gdy pakiet miał jeden model, więc odmowy składane bez tego
// pola mówią dalej to samo, co mówiły.
const (
	// SilnikDlaPrzesiewu — krzyżowy koder drugiego przebiegu wyszukiwania,
	// zdefiniowany w pliku przesiew.go.
	SilnikDlaPrzesiewu = "przesiew"
	// SilnikDlaObrazu — model dwuwieżowy osi obrazu, zdefiniowany w pliku
	// obraz.go tego samego pakietu wiedzy.
	SilnikDlaObrazu = "obraz"
)

// BrakSilnika mówi, że wskaźnika znaczenia nie ma czym zbudować ani przeszukać.
// Osobny typ, tak samo jak `zewnetrzne.BrakNarzedzia`: adapter rozpoznaje go
// i znakuje jako niedostępność zaplecza, a nie jako usterkę wewnętrzną.
type BrakSilnika struct {
	// Silnik — który z trzech silników pakietu odmówił. Puste znaczy silnik
	// osadzeń.
	Silnik string
	// Rodzaj — jedna ze stałych wyżej.
	Rodzaj string
	// Model — nazwa modelu, którego dotyczy brak.
	Model string
	// WagaMb — ile waży do dociągnięcia; 0 znaczy „waga nieznana".
	WagaMb int
	// Powod — to, co powiedział o sobie sam pomocnik albo system.
	Powod string
}

// Error składa trzyczłonową odmowę: czego nie ma, ile to waży, i co dokładnie
// zrobić, żeby to naprawić.
func (b *BrakSilnika) Error() string {
	zdanie := strings.Builder{}
	zdanie.WriteString("wskaźnik znaczenia: ")
	zdanie.WriteString(b.opisBraku())
	if b.Powod != "" {
		zdanie.WriteString("; pomocnik powiedział: " + skroc(b.Powod, 400))
	}
	zdanie.WriteString(b.opisSkutku())
	zdanie.WriteString("; naprawa: " + b.opisNaprawy())
	return zdanie.String()
}

// opisSkutku mówi, czego wobec tego nie będzie — i dlaczego rdzeń nie podstawia
// w to miejsce niczego innego.
func (b *BrakSilnika) opisSkutku() string {
	switch b.Silnik {
	case SilnikDlaPrzesiewu:
		return "; przesiewu nie da się wykonać bez krzyżowego kodera, " +
			"a rdzeń nie odda po cichu kolejności z pierwszego przebiegu — " +
			"wołający prosił o kolejność ułożoną NA NOWO i dostałby tę samą, " +
			"którą miał bez pytania"
	case SilnikDlaObrazu:
		return "; osi obrazu nie da się wykonać bez modelu wiążącego obraz ze zdaniem, " +
			"a rdzeń nie zejdzie po cichu na dopasowanie NAZW plików — " +
			"wołający prosił o to, co na obrazie widać, a nie o to, jak plik nazwano"
	default:
		return "; wyszukiwania po ZNACZENIU nie da się wykonać bez silnika, " +
			"a rdzeń nie zejdzie po cichu na wyszukiwanie po SŁOWACH — " +
			"model dostałby trafienia po literach w miejscu, w którym prosił o trafienia po sensie"
	}
}

// opisBraku nazywa brak wraz z wagą tego, czego brakuje, do dociągnięcia
// z sieci albo z dysku lokalnego.
func (b *BrakSilnika) opisBraku() string {
	if b.Silnik != "" {
		return b.opisBrakuDolozonego()
	}
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

// opisBrakuDolozonego nazywa brak dwóch silników dołożonych do osadzarki.
// Jedna droga dla obu, bo różnią się wyłącznie nazwą zdolności i wykazem
// bibliotek — a te podaje `opisZdolnosci` i `opisPakietow`.
func (b *BrakSilnika) opisBrakuDolozonego() string {
	switch b.Rodzaj {
	case brakBiblioteki:
		return "na tej maszynie nie ma bibliotek, którymi liczy się " + b.opisZdolnosci() +
			" (" + b.opisPakietow() + " to około 900 MB pobrania, a wagi modelu " +
			b.Model + " " + b.opisWagi() + ")"
	case brakModelu:
		return "biblioteki stoją, ale wag modelu " + b.Model + " potrzebnych do " +
			b.opisZdolnosci() + " nie ma na dysku i nie dało się ich przygotować (" +
			b.opisWagi() + " do pobrania)"
	case brakWagStojacych:
		return "w katalogu wskazanym ustawieniem `" + b.kluczKatalogu() + "` leżą wagi, " +
			"ale nie da się na nich postawić modelu " + b.Model
	case BrakInterpretera:
		return "nie ma czym uruchomić pomocnika liczącego " + b.opisZdolnosci() +
			" — interpreter Pythona nie wystartował"
	default:
		return "silnik liczący " + b.opisZdolnosci() + " nie odpowiedział zrozumiale"
	}
}

// opisZdolnosci nazywa rzecz, której Operator nie dostanie — dopełniaczem, bo
// wchodzi w środek zdania.
func (b *BrakSilnika) opisZdolnosci() string {
	if b.Silnik == SilnikDlaObrazu {
		return "oś obrazu"
	}
	return "przesiew wyników"
}

// opisPakietow wymienia biblioteki, których brakuje — nazwami, którymi się je
// instaluje, a nie opisowo.
func (b *BrakSilnika) opisPakietow() string {
	if b.Silnik == SilnikDlaObrazu {
		return "`torch`, `transformers` i `pillow`"
	}
	return "`torch` i `transformers`"
}

// kluczKatalogu oddaje nazwę ustawienia wskazującego katalog wag tego
// silnika, w którym leżą pliki modeli.
func (b *BrakSilnika) kluczKatalogu() string {
	switch b.Silnik {
	case SilnikDlaPrzesiewu:
		return KluczKatalogPrzesiewu
	case SilnikDlaObrazu:
		return KluczKatalogObrazu
	default:
		return KluczKatalogModeli
	}
}

// kluczModelu oddaje nazwę ustawienia wskazującego model tego silnika,
// którego dotyczy zgłoszony brak wag.
func (b *BrakSilnika) kluczModelu() string {
	switch b.Silnik {
	case SilnikDlaPrzesiewu:
		return KluczModelPrzesiewu
	case SilnikDlaObrazu:
		return KluczModelObrazu
	default:
		return KluczModel
	}
}

// opisNaprawy mówi, co dokładnie zrobić — z nazwą pakietu do zainstalowania
// i nazwą ustawienia do poprawienia.
func (b *BrakSilnika) opisNaprawy() string {
	if b.Silnik != "" {
		return b.opisNaprawyDolozonej()
	}
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

// opisNaprawyDolozonej mówi, co zrobić, żeby dołożony silnik ruszył. Trzy
// drogi, bo trzy różne przyczyny: brak wag, wagi nieczytelne, brak bibliotek.
func (b *BrakSilnika) opisNaprawyDolozonej() string {
	switch b.Rodzaj {
	case brakModelu:
		return "dać maszynie dostęp do sieci przy pierwszym użyciu albo przenieść " +
			"pobrane wagi do katalogu ustawienia `" + b.kluczKatalogu() + "`"
	case brakWagStojacych:
		return "sprawdzić, czy w katalogu ustawienia `" + b.kluczKatalogu() +
			"` leży komplet wydania modelu `" + b.kluczModelu() + "` — plik wag " +
			"`model.safetensors` wraz z ustrojem `config.json` i opisem podziału na " +
			"tokeny; albo wskazać tym ustawieniem katalog pusty i pozwolić bibliotece " +
			"pobrać wydanie własne"
	default:
		return "zainstalować biblioteki poleceniem `python3 -m pip install " +
			b.pakietyDoInstalacji() + "` w interpreterze wskazanym ustawieniem `" +
			KluczProgram + "` (pusta wartość znaczy `python3` ze ścieżki wyszukiwania " +
			"systemu)"
	}
}

// pakietyDoInstalacji wymienia pakiety w postaci, w której idą do polecenia
// instalacji — bez znaków wyróżnienia, bo wchodzą do polecenia, nie do zdania.
func (b *BrakSilnika) pakietyDoInstalacji() string {
	if b.Silnik == SilnikDlaObrazu {
		return "torch transformers pillow"
	}
	return "torch transformers"
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
