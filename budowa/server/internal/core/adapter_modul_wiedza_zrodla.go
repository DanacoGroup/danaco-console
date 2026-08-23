// Odpowiedzialność pliku: zebranie treści, która ma wejść do wskaźnika
// znaczenia — po jednym czytelniku na każdy z trzech zakresów kontraktu.
//
// Treść czyta się przez istniejące repozytoria i istniejące czytniki, nie
// własnym SQL-em: druga droga do wierszy pliku biblioteki byłaby drugą prawdą
// o tym, co Operator w niej ma, i rozjechałaby się przy pierwszej zmianie
// schematu. Dlatego biblioteka idzie przez `dane.RepozytoriumBiblioteki`,
// historia przez `dane.RepozytoriumHistorii`, a treść pliku spod odwołania
// czyta ten sam `trescPodgladuBiblioteki`, którym czyta ją podgląd modułu
// Library.
//
// Zakresy mają różne wymagania wobec `windowId`. `library` jest zbiorem całej
// maszyny — repozytorium wiedzy Operatora nie należy do żadnego okna, więc
// wskazanie okna nic tu nie znaczy i jest pomijane. `history` i `workspace` są
// przeciwnie: historia jest historią konkretnego okna
// (`dane.RepozytoriumHistorii.Pozycje` innego kształtu nie ma), a przestrzeń
// robocza jest katalogiem ustalanym po zasięgu, którego najwęższym poziomem
// jest okno. Dla tych dwóch brak `windowId` jest odmową nazywającą brak, a nie
// cichym zaindeksowaniem zera pozycji: „zbudowano wskaźnik, 0 pozycji" wygląda
// jak pusta historia i nie da się tego odróżnić od pomyłki wołającego.
//
// Każdy czytelnik ma sufit liczby dokumentów. Nie jest to ostrożność na zapas:
// osadzenie liczy się na procesorze i katalog z pięćdziesięcioma tysiącami
// plików zająłby maszynę Operatora na godziny bez żadnego śladu postępu.
// Przekroczenie sufitu nie jest odmową — wskaźnik obejmuje tyle, ile obejmuje,
// a powtórzony przebieg po sprzątnięciu katalogu obejmie resztę.
package core

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

const (
	// granicaDokumentowBiblioteki — sufit liczby plików repozytorium wiedzy
	// branych do jednego przebiegu.
	granicaDokumentowBiblioteki = 2000
	// granicaPozycjiHistorii — sufit liczby wypowiedzi jednego okna.
	granicaPozycjiHistorii = 5000
	// granicaPlikowPrzestrzeni — sufit liczby plików katalogu roboczego. Ta sama
	// liczba, którą stosuje przegląd katalogu projektu w module Workspace
	// (`granicaPrzegladuKatalogu`) — dwie różne granice dla tego samego katalogu
	// byłyby dwiema prawdami o tym, co Operator w nim widzi.
	granicaPlikowPrzestrzeni = granicaPrzegladuKatalogu
	// granicaTresciDokumentu — ile znaków jednego dokumentu wchodzi do
	// wskaźnika. Sufit istnieje, bo plik dziennika o rozmiarze gigabajta dałby
	// milion fragmentów i zjadłby wskaźnik reszty biblioteki.
	granicaTresciDokumentu = 256 << 10
)

// dokumentWiedzy to jedna jednostka treści przed podziałem na fragmenty.
type dokumentWiedzy struct {
	// Zakres — kontraktowa nazwa zakresu, z którego dokument pochodzi.
	Zakres string
	// Zrodlo — nazwa czytelna dla człowieka; wchodzi wprost do pola `source`.
	Zrodlo string
	// ZrodloKod — identyfikator, którym da się po dokument sięgnąć.
	ZrodloKod string
	// Tresc — cały tekst dokumentu.
	Tresc string
}

// dokumenty rozdziela pytanie o treść na trzech czytelników.
func (a *adapterWiedzy) dokumenty(ctx context.Context, zakres string,
	idOkna *string) ([]dokumentWiedzy, error) {

	switch zakres {
	case shared.KnowledgeScopeLibrary:
		return a.dokumentyBiblioteki(ctx)
	case shared.KnowledgeScopeHistory:
		return a.dokumentyHistorii(ctx, idOkna)
	case shared.KnowledgeScopeWorkspace:
		return a.dokumentyPrzestrzeni(idOkna)
	default:
		return nil, bladZadaniaWiedzy("zakres `" + zakres + "` nie ma czytelnika treści")
	}
}

// dokumentyBiblioteki czyta pliki repozytorium wiedzy przez repozytorium modułu
// Library, a ich treść — spod odwołania, tak jak czyta ją podgląd.
//
// Plik bez odwołania jest pomijany, nie odmawiany: wiersz bez bajtów w magazynie
// treści niesie same metadane, a jeden taki plik nie ma prawa
// odebrać Operatorowi wskaźnika dla pozostałych. Tak samo pomijany
// jest plik, którego treści nie da się odczytać jako tekstu — obraz osadzony
// jako ciąg bajtów dałby wektor, który do niczego nie pasuje i nigdy nie
// wygrywa, czyli koszt bez pożytku.
func (a *adapterWiedzy) dokumentyBiblioteki(ctx context.Context) ([]dokumentWiedzy, error) {
	if a.biblioteka == nil {
		return nil, bladWiedzyBezZrodla("biblioteka")
	}
	pliki, _, err := a.biblioteka.Pliki(ctx, dane.FiltrPlikow{Limit: granicaDokumentowBiblioteki})
	if err != nil {
		return nil, bladWiedzy(err)
	}

	dokumenty := make([]dokumentWiedzy, 0, len(pliki))
	granica := granicaTresciDokumentu
	for _, plik := range pliki {
		if plik.TrescOdwolanie == nil || strings.TrimSpace(*plik.TrescOdwolanie) == "" {
			continue
		}
		tresc, _, err := trescPodgladuBiblioteki(*plik.TrescOdwolanie, &granica)
		if err != nil || strings.TrimSpace(tresc) == "" {
			continue
		}
		dokumenty = append(dokumenty, dokumentWiedzy{
			Zakres:    shared.KnowledgeScopeLibrary,
			Zrodlo:    plik.Nazwa,
			ZrodloKod: plik.Kod,
			Tresc:     tresc,
		})
	}
	return dokumenty, nil
}

// dokumentyHistorii czyta wypowiedzi jednego okna rozmowy.
//
// Jedna wypowiedź to jeden dokument, a nie cała rozmowa sklejona w jeden tekst.
// Powód jest ten sam, dla którego dokument dzieli się na fragmenty: Operator
// pyta o rzecz, którą kiedyś powiedział, i ma dostać tę wypowiedź wraz z jej
// identyfikatorem, a nie całą rozmowę, w której gdzieś ona jest. Identyfikator
// wypowiedzi wchodzi do `sourceId`, więc trafienie da się sprawdzić.
func (a *adapterWiedzy) dokumentyHistorii(ctx context.Context,
	idOkna *string) ([]dokumentWiedzy, error) {

	if a.historia == nil {
		return nil, bladWiedzyBezZrodla("historia rozmów")
	}
	okno := wskazanieOkna(idOkna)
	if okno == "" {
		return nil, bladZadaniaWiedzy(
			"zakres `history` bez wskazania okna — historia jest historią KONKRETNEJ " +
				"rozmowy i rdzeń nie ma której przeczytać; naprawa: podać `windowId`")
	}
	pozycje, err := a.historia.Pozycje(ctx, okno, 0, granicaPozycjiHistorii)
	if err != nil {
		return nil, bladWiedzy(err)
	}

	dokumenty := make([]dokumentWiedzy, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if strings.TrimSpace(pozycja.Tresc) == "" {
			continue
		}
		dokumenty = append(dokumenty, dokumentWiedzy{
			Zakres: shared.KnowledgeScopeHistory,
			// Nazwa czytelna niesie rolę mówiącego i okno. Sam identyfikator
			// wypowiedzi nic nie mówi człowiekowi, a model czytający trafienie
			// ma od razu wiedzieć, czy cytuje Operatora, czy samego siebie.
			Zrodlo:    "rozmowa " + okno + ", " + string(pozycja.Rola),
			ZrodloKod: pozycja.Identyfikator,
			Tresc:     pozycja.Tresc,
		})
	}
	return dokumenty, nil
}

// dokumentyPrzestrzeni czyta pliki katalogu roboczego okna.
//
// Katalog ustala ten sam `KatalogRoboczy`, którym jadą Terminal i Developer —
// po poziomach zasięgu, z oknem jako poziomem najwęższym. Wpisanie tu własnej
// ścieżki byłoby drugą prawdą o tym, gdzie pracuje okno.
func (a *adapterWiedzy) dokumentyPrzestrzeni(idOkna *string) ([]dokumentWiedzy, error) {
	okno := wskazanieOkna(idOkna)
	if okno == "" {
		return nil, bladZadaniaWiedzy(
			"zakres `workspace` bez wskazania okna — przestrzeń robocza jest " +
				"katalogiem KONKRETNEGO okna; naprawa: podać `windowId`")
	}
	if a.katalog == nil {
		return nil, bladWiedzyBezZrodla("ustalacz katalogu roboczego")
	}
	korzen := strings.TrimSpace(a.katalog.Ustal(konfig.Kontekst{Okno: okno}, "").Sciezka)
	if korzen == "" {
		return nil, bladZadaniaWiedzy(
			"okno " + okno + " nie ma ustalonego katalogu roboczego — nie ma czego " +
				"przeszukać; naprawa: ustawić podstawę katalogu roboczego w konfiguracji")
	}

	dokumenty := []dokumentWiedzy{}
	granica := granicaTresciDokumentu
	przejrzane := 0
	_ = filepath.WalkDir(korzen, func(sciezka string, wpis fs.DirEntry, err error) error {
		if err != nil {
			// Wpis nieczytelny pomijamy: jeden plik bez uprawnień nie ma prawa
			// odebrać Operatorowi wskaźnika pozostałych.
			return nil
		}
		if przejrzane >= granicaPlikowPrzestrzeni {
			return filepath.SkipAll
		}
		przejrzane++
		if wpis.IsDir() {
			return pomijanyKatalog(wpis.Name())
		}
		if strings.HasPrefix(wpis.Name(), ".") {
			return nil
		}
		wzgledna, blad := filepath.Rel(korzen, sciezka)
		if blad != nil {
			return nil
		}
		tresc, _, blad := trescPodgladuBiblioteki(sciezka, &granica)
		if blad != nil || strings.TrimSpace(tresc) == "" {
			return nil
		}
		dokumenty = append(dokumenty, dokumentWiedzy{
			Zakres: shared.KnowledgeScopeWorkspace,
			Zrodlo: wzgledna,
			// Kodem źródła jest ścieżka względna, a nie bezwzględna: bezwzględna
			// wynosiłaby modelowi układ dysku Operatora, a względna wystarcza,
			// żeby po plik sięgnąć modułem Developer.
			ZrodloKod: wzgledna,
			Tresc:     tresc,
		})
		return nil
	})
	return dokumenty, nil
}

// wskazanieOkna oczyszcza wskaźnik okna z żądania.
func wskazanieOkna(idOkna *string) string {
	if idOkna == nil {
		return ""
	}
	return strings.TrimSpace(*idOkna)
}

// bladWiedzyBezZrodla nazywa brak wpięcia — rdzeń złożony bez tego ogniwa nie ma
// czego czytać, a milczące oddanie zera dokumentów wyglądałoby jak pusta wiedza.
func bladWiedzyBezZrodla(nazwa string) error {
	return odmowaWiedzy(shared.ErrorCodeInternalError,
		"rdzeń nie ma wpiętego źródła `"+nazwa+"` — wskaźnik nie ma z czego powstać; "+
			"naprawa: podpiąć je przy składaniu rdzenia")
}
