-- Migracja zakłada rejestr komponentów własnych strefy drugiej strony głównej,
-- wskazujących byt magazynu modułowego.

-- Komponent własny jest bytem odrębnym, kaflem strefy drugiej, który wskazuje byt
-- magazynu modułowego, nie powielając go.
CREATE TABLE komponent (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Wartości kontraktu wprost, bez tłumaczenia na osobny słownik bazy.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('automations','agents','workspace','assistant')),
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    -- Identyfikator zewnętrzny bytu modułowego, opcjonalny, bez więzu obcego wobec czterech tabel.
    byt_docelowy             TEXT,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    -- Konfiguracja przekazywana adapterowi modułu; poprawność zapisu sprawdza warstwa danych.
    konfiguracja             TEXT    NOT NULL DEFAULT '{}',
    -- Przypisanie do poziomu zasięgu; brak wartości znaczy nieprzypisany, nie globalny.
    poziom_zasiegu_id        INTEGER          REFERENCES poziom_zasiegu(id) ON DELETE SET NULL,
    klucz_zasiegu            TEXT    NOT NULL DEFAULT '',
    -- Kolumny czasu niosą milisekundy epoki wprost, bez przekładu w obie strony.
    utworzono                INTEGER NOT NULL DEFAULT 0,
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);

-- Indeks biegnie porządkiem zapytania wykazu, więc kolejność jest stała między
-- wywołaniami i nie wymaga sortowania po odczycie.
CREATE INDEX idx_komponent_wykaz ON komponent(rodzaj, nazwa, id);

-- Wskazanie bytu modułowego służy odpowiedzi, czy dany byt ma już swój kafel; bez
-- indeksu odczyt szedłby przeglądem całej tabeli.
CREATE INDEX idx_komponent_byt_docelowy ON komponent(rodzaj, byt_docelowy);
