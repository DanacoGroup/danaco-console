package core

import (
	"strconv"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przekład wierszy słowników platformy i wierszy utrwalonych na byty kontraktu.
//
// przeklad.go przenosi byty pakietu sesji — te żyją w pamięci rdzenia i znają
// czas jako time.Time. Tutaj przechodzą wiersze bazy: identyfikator jest liczbą,
// czas napisem ISO 8601, a kolejność wyświetlania kolumną albo pozycją w wykazie
// uporządkowanym zapytaniem. Drugiego miejsca tego przekładu nie ma.

// formatZnacznikaBazy odpowiada wyrażeniu strftime schematu: czas UTC w ISO 8601
// z częścią ułamkową sekundy.
const formatZnacznikaBazy = "2006-01-02T15:04:05.999Z"

// chwilaBazy przekłada znacznik czasu bazy na milisekundy epoki kontraktu.
// Znacznik nieczytelny daje zero — wykaz ma się wyświetlić, a nie zniknąć
// z powodu jednej kolumny.
func chwilaBazy(znacznik string) int64 {
	chwila, err := time.Parse(formatZnacznikaBazy, znacznik)
	if err != nil {
		return 0
	}
	return chwila.UnixMilli()
}

// identyfikatorWiersza zwraca identyfikator, którym byt wychodzi kontraktem:
// nadany przez rdzeń, a gdy wiersz powstał poza rdzeniem — klucz główny wiersza.
func identyfikatorWiersza(zewnetrzny *string, id int64) string {
	if zewnetrzny != nil && *zewnetrzny != "" {
		return *zewnetrzny
	}
	return strconv.FormatInt(id, 10)
}

// srodowiskoKontraktu przekłada wiersz słownika środowisk na kartę wejścia.
// Kody modułów wchodzą do wyniku wyłącznie na żądanie, ale rodzaj nawigacji
// rozstrzyga macierz widoczności zawsze — bez niej środowisko z panelem
// orkiestracji byłoby nieodróżnialne od środowiska z pustą listą modułów.
func srodowiskoKontraktu(s dane.Srodowisko, pozycja int, kodyModulow []string, zModulami bool) shared.Environment {
	srodowisko := shared.Environment{
		Id:             strconv.FormatInt(s.ID, 10),
		Code:           s.Kod,
		Name:           s.Nazwa,
		Description:    s.Opis,
		Motto:          s.Motto,
		Order:          kolejnoscWyswietlania(s.Kolejnosc, pozycja),
		NavigationKind: rodzajNawigacji(kodyModulow),
	}
	if zModulami {
		srodowisko.ModuleCodes = kodyModulow
	}
	return srodowisko
}

// rodzajNawigacji rozstrzyga rodzaj bocznej nawigacji z macierzy `srodowisko_modul`:
// środowisko bez ani jednego modułu widocznego wystawia panel orkiestracji,
// pozostałe — listę modułów. Rozstrzyga wiersz macierzy, nie kod środowiska
// wpisany w warunek.
func rodzajNawigacji(kodyModulow []string) shared.NavigationKind {
	if len(kodyModulow) == 0 {
		return shared.NavigationKindOrchestration
	}
	return shared.NavigationKindModules
}

// kolejnoscWyswietlania bierze kolejność z kolumny, a gdy kolumna jej nie niesie
// — z pozycji w wykazie uporządkowanym zapytaniem. Kontrakt liczy pozycje od 1.
func kolejnoscWyswietlania(zKolumny, pozycja int) int {
	if zKolumny > 0 {
		return zKolumny
	}
	return pozycja
}

// modulKontraktu przekłada wiersz słownika modułów wraz z macierzą widoczności
// i katalogiem okien operacyjnych.
//
// Pozycja pochodzi z porządku zapytania, bo kolejność modułu jest własnością
// macierzy `srodowisko_modul`, nie samego modułu. Poza wykazem jednego
// środowiska kolejność nie ma treści — wtedy zostaje zerem.
func modulKontraktu(m dane.Modul, pozycja int, kodySrodowisk, kodyOkien []string) shared.Module {
	return shared.Module{
		Id:                     strconv.FormatInt(m.ID, 10),
		Code:                   m.Kod,
		Name:                   m.Nazwa,
		Description:            m.Opis,
		Order:                  pozycja,
		EnvironmentCodes:       kodySrodowisk,
		OperationalWindowCodes: kodyOkien,
		// Rodzaj wychodzi napisem przepuszczonym bez sprawdzania wykazu.
		// Kontrakt opisuje cztery wartości `kind` wyłącznie w treści opisu —
		// pole jest zwykłym `string`, nie nazwanym wyliczeniem, więc nie ma
		// tu typu, który mógłby czegokolwiek pilnować; kolumna też nie ma
		// warunku CHECK. Wartość spoza wykazu przechodzi zatem taka, jaka
		// jest: przekład nie jest miejscem na osąd danych, a podmiana
		// nieznanego rodzaju na zastępczy zatarłaby ślad błędu w słowniku.
		// NULL daje pusty napis — kontrakt wymaga wartości, a pusty napis
		// czytelnie znaczy „rodzaju nie rozstrzygnięto".
		Kind: wartoscTekstu(m.Rodzaj),
		// Ikona idzie wskaźnikiem wprost: pole kontraktu jest opcjonalne,
		// więc brak w kolumnie zostaje brakiem w kopercie, bez zamiany na
		// pusty napis.
		Icon: m.Ikona,
		// Flaga jest w bazie NOT NULL DEFAULT 0, więc rdzeń zawsze ma
		// wartość — stąd pole kontraktu wymagane i przypisanie bez osłony.
		ConfiguredOnHome: m.KonfigurowanyNaStronieGlownej,
	}
}

// sesjaWierszaKontraktu przekłada utrwalony wiersz sesji na sesję kontraktu.
// Okien nie dokłada: wykaz kart nie potrzebuje ich identyfikatorów, a każdy
// wymagałby osobnego zapytania.
func sesjaWierszaKontraktu(s dane.Sesja) shared.Session {
	sesja := shared.Session{
		Id:        identyfikatorWiersza(s.IdentyfikatorZewnetrzny, s.ID),
		ProjectId: s.Projekt,
		Status:    s.Stan,
		CreatedAt: chwilaBazy(s.Utworzono),
		UpdatedAt: chwilaBazy(s.Zaktualizowano),
	}
	if s.Tytul != "" {
		tytul := s.Tytul
		sesja.Title = &tytul
	}
	return sesja
}

// oknoWierszaKontraktu przekłada utrwalony wiersz okna komunikacji. Moduł
// i kanał modelu wychodzą kodem, nie kluczem obcym — kontrakt zna kody, a klucz
// jest wewnętrzną sprawą schematu.
func oknoWierszaKontraktu(o dane.Okno, kodModulu, kodKanalu string) shared.Window {
	okno := shared.Window{
		Id:             identyfikatorWiersza(o.IdentyfikatorZewnetrzny, o.ID),
		ModuleId:       kodModulu,
		ModelChannelId: kodKanalu,
		WorkingDirs:    listaKatalogow(o.KatalogiRobocze),
		ExecutionEnv:   o.SrodowiskoWykonania,
		PermissionMode: o.TrybUprawnien,
		WindowRole:     o.RolaOkna,
		Title:          o.Tytul,
		Status:         o.Stan,
		CreatedAt:      chwilaBazy(o.Utworzono),
		UpdatedAt:      chwilaBazy(o.Zaktualizowano),
	}
	if okno.WorkingDirs == nil {
		okno.WorkingDirs = []string{}
	}
	return okno
}

// identyfikatoryOkien wylicza identyfikatory okien — relacja 1:N sesji.
func identyfikatoryOkien(okna []shared.Window) []string {
	wykaz := make([]string, 0, len(okna))
	for _, okno := range okna {
		wykaz = append(wykaz, okno.Id)
	}
	return wykaz
}
