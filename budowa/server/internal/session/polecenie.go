package session

// Polecenie opisuje uruchomienie procesu okna, oddzielając wspólny kształt danych startowych od budowy wiersza poleceń wykonywanej przez warstwę kanału poprzez port Uruchamiacz.
type Polecenie struct {
	// Program — plik wykonywalny.
	Program string
	// Argumenty wiersza poleceń, bez nazwy programu.
	Argumenty []string
	// Katalog uruchomienia; pusty oznacza katalog główny okna.
	Katalog string
	// Srodowisko podaje zmienne KLUCZ=wartość; puste oznacza środowisko odziedziczone po rdzeniu.
	Srodowisko []string
	// DziedziczSrodowisko dokłada środowisko rdzenia przed wpisami własnymi; sekrety idą przez proces.
	DziedziczSrodowisko bool
}
