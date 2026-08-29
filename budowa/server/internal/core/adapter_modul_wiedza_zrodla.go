// Odpowiedzialność pliku: zebranie treści, która ma wejść do wskaźnika
// znaczenia — po jednym czytelniku na każdy z trzech zakresów kontraktu:
// `library`, `history` i `workspace`. Każdy czytelnik ma sufit liczby dokumentów.
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
	// granicaDokumentowBiblioteki — sufit plików repozytorium wiedzy branych
	// do jednego przebiegu czytania biblioteki.
	granicaDokumentowBiblioteki = 2000
	// granicaPozycjiHistorii — sufit liczby wypowiedzi jednego okna branych
	// do jednego przebiegu czytania historii rozmowy.
	granicaPozycjiHistorii = 5000
	// granicaPlikowPrzestrzeni — sufit plików katalogu roboczego brany do
	// jednego przebiegu czytania przestrzeni roboczej okna.
	granicaPlikowPrzestrzeni = granicaPrzegladuKatalogu
	// granicaTresciDokumentu — ile znaków jednego dokumentu wchodzi do
	// wskaźnika, żeby plik dziennika o rozmiarze gigabajta nie zjadł wskaźnika.
	granicaTresciDokumentu = 256 << 10
)

// dokumentWiedzy to jedna jednostka treści przed podziałem na fragmenty
// wskaźnika znaczenia, niosąca zakres, źródło i pełną treść.
type dokumentWiedzy struct {
	// Zakres — kontraktowa nazwa zakresu, z którego dokument pochodzi.
	Zakres string
	// Zrodlo — nazwa czytelna dla człowieka; wchodzi wprost do pola `source`.
	Zrodlo string
	// ZrodloKod — identyfikator, którym da się po dokument sięgnąć.
	ZrodloKod string
	// Tresc — cały tekst dokumentu, dzielony na fragmenty dopiero przy
	// osadzeniu wskaźnika znaczenia.
	Tresc string
}

// dokumenty rozdziela pytanie o treść na trzech czytelników, po jednym na
// każdy zakres kontraktu: bibliotekę, historię i przestrzeń roboczą.
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
// Library, a ich treść spod odwołania, tak jak czyta ją podgląd. Plik bez
// odwołania albo nieczytelny jako tekst jest pomijany, nie odmawiany.
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

// dokumentyHistorii czyta wypowiedzi jednego okna rozmowy. Jedna wypowiedź
// to jeden dokument, a nie cała rozmowa sklejona w jeden tekst.
func (a *adapterWiedzy) dokumentyHistorii(ctx context.Context,
	idOkna *string) ([]dokumentWiedzy, error) {

	if a.historia == nil {
		return nil, bladWiedzyBezZrodla("historia rozmów")
	}
	okno := wskazanieOkna(idOkna)
	if okno == "" {
		return nil, bladZadaniaWiedzy(
			"zakres `history` bez wskazania okna — historia jest historią KONKRETNEJ " +
				"rozmowy i serwer nie ma której przeczytać; naprawa: podać `windowId`")
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
			// Nazwa czytelna niesie rolę mówiącego i okno, identyfikator nic nie mówi.
			Zrodlo:    "rozmowa " + okno + ", " + string(pozycja.Rola),
			ZrodloKod: pozycja.Identyfikator,
			Tresc:     pozycja.Tresc,
		})
	}
	return dokumenty, nil
}

// dokumentyPrzestrzeni czyta pliki katalogu roboczego okna. Katalog ustala
// ten sam `KatalogRoboczy`, którym jadą Terminal i Developer.
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
			// Wpis nieczytelny zostaje pominięty, nie odbiera wskaźnika pozostałych.
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
			// Kodem źródła jest ścieżka względna, nie bezwzględna: nie wynosi
			// modelowi układu dysku Operatora.
			ZrodloKod: wzgledna,
			Tresc:     tresc,
		})
		return nil
	})
	return dokumenty, nil
}

// wskazanieOkna oczyszcza wskaźnik okna z żądania, oddając pusty ciąg,
// gdy pole nie zostało podane w żądaniu klienta.
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
		"serwer nie ma wpiętego źródła `"+nazwa+"` — wskaźnik nie ma z czego powstać; "+
			"naprawa: podpiąć je przy składaniu serwera")
}
