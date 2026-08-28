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
	// TrybUprawnien jest wyliczeniem kontraktu odpowiadającym przełącznikowi
	// --permission-mode.
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
	// Wznowienie jest identyfikatorem rozmowy CLI dla przełącznika --resume.
	Wznowienie string
	// PulapKosztuUSD jest górną granicą kosztu wywołania w dolarach; zero
	// znaczy brak ograniczenia.
	PulapKosztuUSD float64
	// Srodowisko jest zestawem zmiennych dokładanych do środowiska procesu.
	Srodowisko map[string]string
	// Konto jest kodem konta wskazanego dla tego wywołania; puste znaczy
	// rotację puli.
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
	// NaStartProcesu zawiadamia o uruchomieniu procesu tury wraz z jego PID
	// systemowym.
	NaStartProcesu func(pid int)
	// NaZdarzenieZaczepu oddaje zdarzenie zaczepu odczytane ze strumienia
	// warstwie składania.
	NaZdarzenieZaczepu func(ZdarzenieZaczepu)
	// Ustawienia i Nakladka opisują to wywołanie.
	Ustawienia Ustawienia
	Nakladka   Nakladka
}
