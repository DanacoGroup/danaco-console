# Danaco Console — przekrój pionowy: powłoka

Opracowanie opisuje powłokę natywną w postaci, w jakiej działa w bieżącej
budowie. Obejmuje punkt wejścia, prowadzenie procesu rdzenia, sekret nawiązania
oraz obszary obsługi okna i zasobnika.

## Punkt wejścia

Powłoka jest aplikacją Tauri w języku Rust; źródła leżą
w `budowa/desktop/src-tauri/src/`. Punkt wejścia `main.rs` czyta ustawienia,
składa aplikację i pracuje do jawnego zakończenia procesu. Awarię startu
obsługuje osobny tor `awaria_startu.rs`, który zamiast pustego okna pokazuje
przyczynę.

## Proces rdzenia

Powłoka startuje rdzeń obok siebie i prowadzi go przez cały czas pracy: moduł
`rdzen/` niesie uruchomienie procesu (`proces.rs`), nasłuch jego gotowości
(`nasluch.rs`) i stan (`stan.rs`). Zamknięcie powłoki kończy pracę w porządku
wyznaczonym przez `zamkniecie.rs`.

## Sekret nawiązania

Moduł `sekret.rs` wytwarza sekret nawiązania gniazda: wartość losową jednego
uruchomienia powłoki, o 32 bajtach losowości w zapisie szesnastkowym.
Powłoka podaje sekret rdzeniowi startowanemu obok siebie oraz stronie
interfejsu otwierającej gniazdo — w zapytaniu adresu gniazda. Źródłem
losowości jest system operacyjny, a jego odmowa jest awarią startu, nie powodem
do sekretu przewidywalnego. Weryfikację sekretu po stronie rdzenia opisuje
opracowanie [kanał](kanal.md).

## Okno, zasobnik i otoczenie

Pozostałe obszary powłoki: `okno.rs` i `montaz.rs` — złożenie i prowadzenie
okna aplikacji; `zasobnik.rs` i `menu_zasobnika.rs` — ikona i menu zasobnika
systemowego; `nastawy.rs` i `ustawienia.rs` — nastawy powłoki; `polecenia.rs`
— polecenia wystawione stronie interfejsu; `dialog_katalogu.rs` — okno wyboru
katalogu; `wskazanie.rs` — wskazania strony; `dziennik.rs` — dziennik;
`aktualizacja/` — aktualizacja aplikacji.
