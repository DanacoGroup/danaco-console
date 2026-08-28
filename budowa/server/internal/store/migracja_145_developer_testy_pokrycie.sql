-- Migracja 145 — wyniki testów i pokrycie kodu przebiegu budowania (Developer,
-- okno Build Output).
--
-- Wynik testu powstaje z rozbioru logu przebiegu w chwili jego domknięcia,
-- a nie z ponownego odpytywania logu przy każdym żądaniu. Powód jest praktyczny:
-- w dzienniku przebiegu zostaje wyłącznie OGON logu (rdzeń przycina go do
-- kilkuset wierszy), więc wynik testu odczytany godzinę później nie miałby
-- z czego powstać. Rozbiór idzie raz, kiedy pełne wyjście jeszcze płynie.
--
-- Pokrycie ma własną tabelę, bo jest pomiarem pliku, a nie testu: jeden przebieg
-- daje setki wyników testów i dziesiątki wierszy pokrycia, i nic ich nie łączy
-- poza przebiegiem.
--
-- `wiersze_bez_pokrycia` trzymamy jako tekst z numerami rozdzielonymi
-- przecinkiem. Odbiorcą jest nakładka pokrycia w edytorze, która bierze ten
-- zbiór w całości dla jednego pliku; tabela wiersz-na-wiersz rosłaby o rząd
-- wielkości bez jednego pytania, na które odpowiadałaby lepiej.
CREATE TABLE developer_wynik_testu (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    budowanie_kod TEXT    NOT NULL,
    zestaw        TEXT,
    nazwa         TEXT    NOT NULL,
    stan          TEXT    NOT NULL CHECK(stan IN ('passed','failed','skipped')),
    czas_ms       INTEGER,
    tresc         TEXT,
    sciezka       TEXT,
    wiersz        INTEGER
);
CREATE INDEX idx_developer_wynik_testu_budowanie ON developer_wynik_testu(budowanie_kod, stan);

CREATE TABLE developer_pokrycie (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    budowanie_kod        TEXT    NOT NULL,
    sciezka              TEXT    NOT NULL,
    instrukcje           INTEGER NOT NULL DEFAULT 0,
    pokryte              INTEGER NOT NULL DEFAULT 0,
    procent              INTEGER NOT NULL DEFAULT 0,
    wiersze_bez_pokrycia TEXT,
    UNIQUE(budowanie_kod, sciezka)
);
CREATE INDEX idx_developer_pokrycie_budowanie ON developer_pokrycie(budowanie_kod);
