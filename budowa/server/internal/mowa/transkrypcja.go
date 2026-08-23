// Odpowiedzialność pliku: zamiana jednego nagrania na tekst — złożenie
// wywołania, odczyt odpowiedzi pomocnika i rozstrzygnięcie, czy wynik wolno
// oddać jako transkrypcję.
//
// Pusta transkrypcja jest wynikiem, gdy przetworzenie się odbyło. Atrapą byłoby
// oddanie pustego tekstu zamiast rozpoznania, którego nie wykonano; nie jest nią
// pusty tekst będący stwierdzeniem, że w nagraniu nie ma mowy. Rozróżnia je pole
// `Stan`, a odmowa zostaje wyłącznie dla przypadku trzeciego: nie przetworzono.
//
// Dźwięk nie wychodzi z maszyny Operatora. To warunek pakietu, nie preferencja.
// Ten plik nie ma i nie będzie miał klienta HTTP: jedyną drogą nagrania jest
// ścieżka na dysku podana pomocnikowi, który pracuje lokalnie na procesorze.
//
// Trwanie to długość nagrania, nie czas przetwarzania. Pomocnik oddaje jedno
// i drugie da się pomylić przy czytaniu; pomyłka byłaby kłamstwem w dzienniku,
// bo na słabym procesorze przetwarzanie trwa dłużej niż samo nagranie.
package mowa

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
)

// granicaTranskrypcji to granica czasu jednego rozpoznania.
//
// Hojna rozmyślnie: pierwsze uruchomienie pobiera wagi modelu (rząd setek
// megabajtów), a rozpoznanie na procesorze bez akceleracji bywa wolniejsze od
// odtwarzania nagrania. Granica ma odciąć proces zawieszony, a nie pracę, która
// po prostu trwa.
const granicaTranskrypcji = 30 * time.Minute

// Przełączniki pomocnika w trybie rozpoznawania.
const (
	przelacznikPliku         = "--plik"
	przelacznikModelu        = "--model"
	przelacznikJezyka        = "--jezyk"
	przelacznikKataloguModel = "--katalog-modeli"
)

// Zlecenie opisuje jedno rozpoznanie.
//
// Okno, zasady i obszar jadą razem ze zleceniem, bo pomocnik przechodzi przez
// tę samą bramę izolacji, co każdy inny proces okna (uruchomienie.go). Model
// i język puste znaczą wartość z ustawień okna, nie brak wskazania — brak
// danych ma dawać poprawny wynik, nie odmowę.
type Zlecenie struct {
	// Odnosnik — wskazanie nagrania; ścieżka pliku na dysku Operatora.
	Odnosnik string
	// Jezyk — kod języka rozpoznawania; pusty bierze wartość z ustawień.
	Jezyk string
	// Model — rozmiar modelu; pusty bierze wartość z ustawień.
	Model string
	// Okno — okno komunikacji, na którego rzecz idzie rozpoznanie.
	Okno session.Okno
	// Zasady — punkty izolacji obowiązujące przy uruchomieniu.
	Zasady session.Zasady
	// Obszar — wydzielone dla okna miejsce na dysku.
	Obszar session.Obszar
}

// Transkrypcja jest wynikiem rozpoznania jednego nagrania.
//
// Pola przepisane co do znaku z odpowiedzi pomocnika
// (pomocniki/transkrypcja/transkrypcja.py, tryb `--plik`).
type Transkrypcja struct {
	// Tekst — rozpoznana wypowiedź.
	Tekst string `json:"tekst"`
	// Znakow — długość tekstu; liczona przez pomocnika, nie zgadywana tutaj.
	Znakow int `json:"znakow"`
	// TrwanieMs — długość nagrania w milisekundach, nie czas przetwarzania.
	TrwanieMs int64 `json:"trwanie_ms"`
	// Model — model, którym rozpoznano.
	Model string `json:"model"`
	// Jezyk — język, w którym rozpoznano.
	Jezyk string `json:"jezyk"`
	// Stan niesie stanRozpoznano albo StanBezMowy — patrz niżej.
	Stan string `json:"stan"`
	// Blad — odmowa pomocnika; wypełniona wyłącznie, gdy rozpoznania nie ma.
	Blad string `json:"blad"`
}

// Trzy stany transkrypcji, z których dwa są wynikiem, a trzeci odmową.
//
// Nagranie ciszy albo szumu daje pusty tekst, który jest faktem, a nie atrapą:
// odmowa znaczyłaby „rdzeń nie umie rozpoznać”, podczas gdy rdzeń umie i
// stwierdził, że nie ma czego rozpoznać. Rozróżnienie przebiega więc nie po
// długości tekstu, lecz po tym, czy przetworzenie się odbyło:
//
//	stanRozpoznano — przetworzono, mowa jest, tekst niepusty
//	StanBezMowy    — przetworzono, mowy brak, tekst pusty i to jest wynik
//	(odmowa)       — nie przetworzono; jedyny przypadek błędu, bez Transkrypcji
const (
	// stanRozpoznano — w nagraniu rozpoznano wypowiedź.
	stanRozpoznano = "rozpoznano"
	// StanBezMowy — nagranie przetworzono i nie zawiera mowy.
	StanBezMowy = "bez_mowy"
)

// Rozpoznano odpowiada, czy transkrypcja niesie wypowiedź.
//
// Wołający, który chce wstawić tekst do pola polecenia, pyta o to, a nie
// o długość tekstu — bo pusty wynik jest tu poprawną odpowiedzią, a nie
// brakiem odpowiedzi.
func (t Transkrypcja) Rozpoznano() bool {
	return t.Stan == stanRozpoznano
}

// Transkrybuj zamienia nagranie na tekst.
//
// Kolejność sprawdzeń nie jest dowolna: najpierw nagranie, potem pomocnik.
// Ścieżka, której nie ma, jest pomyłką Operatora i ma wrócić natychmiast,
// zanim rdzeń zacznie szukać interpretera i uruchamiać proces — komunikat
// o brakującym Pythonie w odpowiedzi na literówkę w ścieżce wskazywałby
// naprawę zupełnie nie tam, gdzie leży usterka.
//
// Dziennik zapisuje także odmowy. Odmowa bez śladu jest nie do zdiagnozowania,
// a to właśnie odmowy Operator zgłasza.
func (s *Silnik) Transkrybuj(ctx context.Context, z Zlecenie) (Transkrypcja, error) {
	model := pierwszeNiepuste(z.Model, s.ustawienia.Model)
	jezyk := pierwszeNiepuste(z.Jezyk, s.ustawienia.Jezyk)

	nagranie, err := OtworzNagranie(z.Odnosnik)
	if err != nil {
		return Transkrypcja{}, s.odnotujOdmowe(ctx, z, model, jezyk, err)
	}

	pomocnik, err := OdnajdzPomocnika(s.ustawienia.Program)
	if err != nil {
		return Transkrypcja{}, s.odnotujOdmowe(ctx, z, model, jezyk, err)
	}

	wynik, err := Uruchom(ctx, s.uruchamiacz, z.Okno, z.Zasady, z.Obszar, pomocnik,
		pomocnik.Argumenty(argumentyRozpoznania(nagranie, model, jezyk, s.ustawienia.KatalogModeli)...),
		s.katalogPracy(z.Obszar), granicaTranskrypcji)
	if err != nil {
		return Transkrypcja{}, s.odnotujOdmowe(ctx, z, model, jezyk,
			odmowaUruchomienia(s.ustawienia.Program, err, wynik))
	}

	transkrypcja, err := odczytajRozpoznanie(wynik)
	if err != nil {
		return Transkrypcja{}, s.odnotujOdmowe(ctx, z, model, jezyk, err)
	}

	s.odnotujGotowa(ctx, z, transkrypcja)
	return transkrypcja, nil
}

// argumentyRozpoznania składa wiersz wywołania pomocnika w trybie rozpoznawania.
//
// Katalog modeli dokładany jest tylko wtedy, gdy Operator go wskazał: wartość
// pusta znaczy katalog domyślny biblioteki, a przekazanie pustego przełącznika
// kazałoby pomocnikowi szukać modelu w katalogu o nazwie pustej.
func argumentyRozpoznania(nagranie Nagranie, model, jezyk, katalogModeli string) []string {
	argumenty := []string{
		przelacznikPliku, nagranie.Sciezka,
		przelacznikModelu, model,
		przelacznikJezyka, jezyk,
	}
	if strings.TrimSpace(katalogModeli) != "" {
		argumenty = append(argumenty, przelacznikKataloguModel, katalogModeli)
	}
	return argumenty
}

// odczytajRozpoznanie rozbiera odpowiedź pomocnika i pilnuje, żeby nie oddać
// pustego wyniku jako udanego rozpoznania.
//
// Nazwa mówi „rozpoznanie", a nie „transkrypcję", bo `odczytajTranskrypcje`
// należy już do dziennika i czyta wiersz bazy. Dwie funkcje o jednej nazwie
// w jednym pakiecie nie tylko się nie skompilują — myliłyby dwa różne odczyty.
func odczytajRozpoznanie(wynik Wynik) (Transkrypcja, error) {
	tresc := strings.TrimSpace(string(wynik.Wyjscie))
	if tresc == "" {
		return Transkrypcja{}, &BrakSilnika{Powod: "pomocnik transkrypcji nie oddał wyniku" +
			ogonDiagnostyki(wynik.Diagnostyka) +
			"; naprawa: sprawdzić gotowość silnika mowy przed nagrywaniem"}
	}

	var odpowiedz Transkrypcja
	if err := json.Unmarshal([]byte(tresc), &odpowiedz); err != nil {
		return Transkrypcja{}, &BrakSilnika{Powod: "odpowiedź pomocnika transkrypcji" +
			" jest nieczytelna (" + err.Error() + ")" + ogonDiagnostyki(wynik.Diagnostyka) +
			"; naprawa: zgłosić usterkę pomocnika — rdzeń oczekuje odpowiedzi w JSON"}
	}
	if odmowa := strings.TrimSpace(odpowiedz.Blad); odmowa != "" {
		return Transkrypcja{}, &BrakSilnika{Powod: odmowa}
	}

	// Pusty tekst nie jest odmową — patrz komentarz przy stałych stanów.
	// Stan uzupełniamy, gdy pomocnik go nie podał: starsze wydanie pomocnika
	// pola nie miało, a odmowa z powodu brakującego pola unieważniłaby wynik,
	// który jest kompletny. Rozstrzyga wtedy sama treść.
	if strings.TrimSpace(odpowiedz.Stan) == "" {
		odpowiedz.Stan = stanRozpoznano
		if strings.TrimSpace(odpowiedz.Tekst) == "" {
			odpowiedz.Stan = StanBezMowy
		}
	}
	return odpowiedz, nil
}

// OpisDziennika składa zdanie, którym transkrypcja melduje się w dzienniku.
//
// Kształt wprost z rozpoznania procesu głównego: „dyktowanie: rozpoznano N zn.
// (model …)”. Liczba znaków jest miarą, którą widać bez zaglądania w treść —
// dziennik nie powtarza wypowiedzi Operatora.
func OpisDziennika(t Transkrypcja) string {
	return "dyktowanie: rozpoznano " + strconv.Itoa(t.Znakow) + " zn. (model " + t.Model + ")"
}

// pierwszeNiepuste oddaje wskazanie zlecenia albo — gdy go nie ma — ustawienie.
func pierwszeNiepuste(wskazane, domyslne string) string {
	if strings.TrimSpace(wskazane) != "" {
		return wskazane
	}
	return domyslne
}
