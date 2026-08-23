package injection

import "danacoconsole/shared"

// Ustawienia są kompletem parametrów jednego wywołania kanału. Pakiet nie zna
// domyślnego modelu, ścieżki ani katalogu — wszystko przychodzi z zewnątrz.
// Pole puste znaczy „przełącznika nie podajemy",
// nigdy „przerwij wywołanie".
type Ustawienia struct {
	// Program jest ścieżką pliku wykonywalnego `claude`.
	Program string
	// Model i ModelZapasowy trafiają do --model i --fallback-model.
	Model         string
	ModelZapasowy string
	// Naklad jest nakładem rozumowania (--effort): low, medium, high, xhigh, max.
	Naklad string
	// TrybUprawnien jest wyliczeniem kontraktu; jego wartości odpowiadają
	// dosłownie przełącznikowi --permission-mode.
	TrybUprawnien shared.PermissionMode
	// Katalogi są listą katalogów roboczych okna. Każdy jedzie
	// osobnym --add-dir.
	Katalogi []string
	// KatalogRoboczy jest katalogiem startowym procesu. Pusty znaczy: katalog
	// bieżący procesu rdzenia.
	KatalogRoboczy string
	// PlikUstawien trafia do --settings (ścieżka albo napis JSON).
	PlikUstawien string
	// KonfiguracjaMCP to wykaz plików albo napisów JSON dla --mcp-config.
	KonfiguracjaMCP []string
	// Wznowienie jest identyfikatorem rozmowy CLI dla --resume. Bierze się
	// z pola IdSesjiCLI poprzedniej tury.
	Wznowienie string
	// PulapKosztuUSD jest górną granicą kosztu tego wywołania, w dolarach —
	// trafia do --max-budget-usd (klucz `pulap_kosztu_usd`).
	//
	// Pułap jest zapobiegawczy, a nie sprawozdawczy: pole `Koszt` fragmentu
	// zamknięcia tury mówi, ile wydano, a to pole mówi, ile wydać wolno.
	// Zero (i tak samo wartość ujemna, której Operator wpisać nie powinien,
	// ale wpisać może) znaczy „przełącznika nie podajemy" — wywołanie idzie
	// bez ograniczenia.
	PulapKosztuUSD float64
	// Srodowisko jest zestawem zmiennych dokładanych do środowiska procesu.
	// CLAUDE_CONFIG_DIR ustawia pula kont i tutaj go nie podajemy.
	Srodowisko map[string]string
	// Konto jest kodem konta wskazanego dla tego wywołania: obszar
	// account konfiguracji sesji albo wiersz rejestru kanałów. Puste znaczy
	// rotację puli. Wskazanie jest rozkazem tożsamości:
	// tura nie pojedzie innym kontem niż wskazane — bez cichej podmiany.
	Konto string
}

// Zapytanie jest pojedynczym zwróceniem się do modelu w konkretnym oknie.
// Ustawienia i nakładka jadą razem z zapytaniem, bo rozstrzyga je resolver
// ośmiu poziomów zasięgu, a nie kanał.
type Zapytanie struct {
	// IdOkna i IdWiadomosci wracają w każdym fragmencie strumienia.
	IdOkna       string
	IdWiadomosci string
	// Tekst jest wypowiedzią użytkownika wysyłaną na stdin jako JSON-lines.
	Tekst string
	// NaStartProcesu zawiadamia o uruchomieniu procesu tury wraz z jego PID.
	// Warstwa sesji obejmuje ten proces uchwytem systemowym, dzięki czemu
	// zamknięcie okna kończy także jego potomstwo. Haczyk pusty nie zmienia
	// przebiegu tury — znika wyłącznie sprzątanie po niej.
	NaStartProcesu func(pid int)
	// NaZdarzenieZaczepu oddaje zdarzenie zaczepu odczytane ze strumienia.
	// Odbiorcą jest warstwa składania — dziennik zdarzeń i diagnostyka;
	// kanał zdarzenia nie interpretuje. Haczyk pusty nie zmienia przebiegu
	// tury — znika wyłącznie ślad zaczepów.
	NaZdarzenieZaczepu func(ZdarzenieZaczepu)
	// Ustawienia i Nakladka opisują to wywołanie.
	Ustawienia Ustawienia
	Nakladka   Nakladka
}
