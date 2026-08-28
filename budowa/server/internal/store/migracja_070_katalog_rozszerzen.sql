-- Migracja zakłada katalog rozszerzeń jako warstwę nad mostem oraz konektorem
-- i umiejętnością eksperta, nie ich kopię.

-- Pozycja katalogu rozszerzeń istnieje niezależnie od eksperta; ekspert dopiero bierze
-- pozycję, zakładając sobie konektor albo umiejętność.
CREATE TABLE rozszerzenie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Tożsamość wiersza, którą niosą polecenia konfiguracji, przełączenia i odinstalowania.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Kod pozycji stały między wydaniami, którym woła polecenie instalacji.
    kod                      TEXT    NOT NULL UNIQUE,
    -- Wartości kontraktu wprost, bez tłumaczenia na osobny słownik bazy.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('mcp', 'plugin', 'api', 'skill')),
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    wersja                   TEXT,
    zainstalowane            INTEGER NOT NULL DEFAULT 0 CHECK(zainstalowane IN (0, 1)),
    wlaczone                 INTEGER NOT NULL DEFAULT 0 CHECK(wlaczone IN (0, 1)),
    -- Wskazuje wiersz punktu dostępu tak samo jak konektor eksperta, bez powielania rejestru mostów.
    punkt_dostepu_id         INTEGER          REFERENCES punkt_dostepu(id) ON DELETE SET NULL,
    -- Napis podany w polu źródła, nie adres do pobrania; rdzeń go nie pobiera i nie uruchamia.
    zrodlo_deklarowane       TEXT    NOT NULL DEFAULT '',
    -- Konfiguracja rozszerzenia; poprawność zapisu sprawdza warstwa danych.
    konfiguracja             TEXT    NOT NULL DEFAULT '{}',
    -- Milisekundy epoki wprost, bez przekładu w obie strony.
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);

-- Indeks biegnie porządkiem zapytania wykazu pozycji, więc kolejność jest stała między
-- kolejnymi wywołaniami odczytu.
CREATE INDEX idx_rozszerzenie_wykaz ON rozszerzenie(rodzaj, nazwa, id);

-- Odpowiedź na pytanie, które pozycje katalogu wiszą na tym moście; bez indeksu odczyt
-- szedłby przeglądem całej tabeli.
CREATE INDEX idx_rozszerzenie_punkt ON rozszerzenie(punkt_dostepu_id);
