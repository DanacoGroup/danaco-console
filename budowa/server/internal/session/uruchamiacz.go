package session

import "io"

// UchwytProcesu jest uchwytem procesu już uruchomionego. Uruchamianie procesu
// ma w drzewie jedno miejsce — warstwę kanału (internal/injection), bo tylko ona
// zna wiersz poleceń, rotację kont i kształt strumienia.
//
// Pakiet session nie buduje ani nie startuje procesu — także procesu okna
// komunikacji. Bierze uchwyt procesu gotowego i dokłada to, co należy do sesji:
// objęcie całego drzewa potomstwa jednym uchwytem systemowym (przejecie.go)
// oraz pętlę koordynator–wykonawca (petla.go).
//
// Interfejs zostaje w tym pakiecie, choć sam session po niego nie sięga:
// definiuje go strona znająca kształt uruchomienia okna (Okno, Polecenie, drzewo
// potomstwa), a wypełnia warstwa kanału. Sięgają po niego moduły Terminal
// i Developer w rdzeniu — one uruchamiają procesy okna i one obejmują je
// drzewem.
type UchwytProcesu interface {
	// Pid — identyfikator systemowy uruchomionego procesu.
	Pid() int
	// Wejscie — strumień wejściowy procesu.
	Wejscie() io.WriteCloser
	// Wyjscie — strumień wyjściowy procesu.
	Wyjscie() io.Reader
	// Diagnostyka — wyjście diagnostyczne procesu.
	Diagnostyka() io.Reader
	// Czekaj czeka na zakończenie procesu i zwraca wynik zakończenia.
	Czekaj() error
	// Ubij kończy sam proces. Sesja sięga po to wyłącznie wtedy, gdy procesu nie
	// udało się objąć drzewem — uchwytu nie wolno oddać, zostawiając sierotę.
	Ubij() error
}

// Uruchamiacz startuje proces okna według podanego polecenia. Implementuje go
// warstwa kanału modelu; pakiet session zna wyłącznie ten interfejs.
type Uruchamiacz interface {
	UruchomProces(o Okno, p Polecenie) (UchwytProcesu, error)
}
