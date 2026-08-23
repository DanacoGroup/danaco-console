package session

// Polecenie opisuje uruchomienie procesu okna. Pakiet session nie wie, jak
// zbudować wiersz poleceń kanału modelu — to należy do warstwy kanału
// (internal/injection dla kanału głównego). Rozdzielenie idzie przez interfejs,
// nie przez wywołanie w głąb cudzego pakietu.
//
// Polecenie jest wspólnym kształtem, nie punktem uruchomienia: wypełnia je rdzeń
// (moduły Terminal i Developer), a wykonuje warstwa kanału przez port
// Uruchamiacz. Sesja nie buduje polecenia i sama go nie wykonuje.
type Polecenie struct {
	// Program — plik wykonywalny.
	Program string
	// Argumenty wiersza poleceń, bez nazwy programu.
	Argumenty []string
	// Katalog uruchomienia; pusty oznacza katalog główny okna.
	Katalog string
	// Srodowisko w postaci KLUCZ=wartość. Puste znaczy środowisko odziedziczone
	// po rdzeniu.
	Srodowisko []string
	// DziedziczSrodowisko dokłada środowisko rdzenia przed wpisami własnymi.
	// Sekrety idą przez środowisko procesu, nigdy przez bazę ani repozytorium.
	DziedziczSrodowisko bool
}
