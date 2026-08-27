package session

import "io"

// UchwytProcesu jest uchwytem procesu już uruchomionego przez warstwę kanału, a pakiet session dokłada do niego objęcie drzewa potomstwa i pętlę koordynator-wykonawca.
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
	// Ubij kończy proces; sesja sięga po to, gdy nie udało się objąć drzewem, by nie zostawić sieroty.
	Ubij() error
}

// Uruchamiacz startuje proces okna według podanego polecenia, jest implementowany przez warstwę kanału modelu, a pakiet session zna wyłącznie ten interfejs.
type Uruchamiacz interface {
	UruchomProces(o Okno, p Polecenie) (UchwytProcesu, error)
}
