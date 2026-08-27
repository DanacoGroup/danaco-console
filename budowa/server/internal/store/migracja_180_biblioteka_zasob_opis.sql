-- Migracja 180 rozszerza plik_biblioteki o stan cyklu życia i ścieżkę wewnątrz repozytorium, zakłada opis zasobu w schemacie Dublin Core oraz definicje pól niestandardowych metadanych.

-- Rozszerza tabelę plik_biblioteki o kolumny stanu cyklu życia i ścieżki wewnątrz repozytorium wraz z indeksami wykazu czynnego i wyszukiwania duplikatów.
ALTER TABLE plik_biblioteki
    ADD COLUMN stan TEXT NOT NULL DEFAULT 'aktywny'
        CHECK(stan IN ('aktywny','zarchiwizowany'));

ALTER TABLE plik_biblioteki ADD COLUMN sciezka_repozytorium TEXT;

-- Wykaz domyślny pokazuje zasoby czynne, więc kolumna stanu wchodzi do indeksu przed porządkiem czasu aktualizacji wiersza.
CREATE INDEX idx_plik_biblioteki_stan ON plik_biblioteki(stan, zaktualizowano DESC);
CREATE INDEX idx_plik_biblioteki_sciezka_repozytorium
    ON plik_biblioteki(sciezka_repozytorium);
-- Rozpoznanie duplikatów dokładnych idzie po sumie kontrolnej — bez indeksu
-- byłoby przejściem po całym repozytorium przy każdym skanowaniu.
CREATE INDEX idx_plik_biblioteki_suma ON plik_biblioteki(suma_kontrolna);

-- Zakłada tabelę opis_zasobu_biblioteki niosącą opis zasobu w schemacie Dublin Core wraz z mapą pól niestandardowych w zapisie JSON.
CREATE TABLE opis_zasobu_biblioteki (
    plik_id              INTEGER PRIMARY KEY
                         REFERENCES plik_biblioteki(id) ON DELETE CASCADE,
    tytul                TEXT,
    tworca               TEXT,
    temat                TEXT,
    opis                 TEXT,
    wydawca              TEXT,
    wspoltworca          TEXT,
    data_zasobu          TEXT,
    rodzaj               TEXT,
    format               TEXT,
    identyfikator        TEXT,
    zrodlo               TEXT,
    jezyk                TEXT,
    powiazanie           TEXT,
    zakres               TEXT,
    prawa                TEXT,
    -- Mapa kod pola i wartość w JSON, bez kolumny na pole zmienianej migracją przy każdym polu Operatora.
    pola_niestandardowe  TEXT,
    zaktualizowano       TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Zakłada tabelę pole_schematu_biblioteki niosącą definicję pola niestandardowego metadanych wraz z jego zawężeniem stosowalności.
CREATE TABLE pole_schematu_biblioteki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    etykieta       TEXT    NOT NULL,
    rodzaj         TEXT    NOT NULL
                   CHECK(rodzaj IN ('tekst','liczba','data','logiczna','lista')),
    wymagane       INTEGER NOT NULL DEFAULT 0 CHECK(wymagane IN (0,1)),
    -- Zawężenia stosowalności pola: rodzaj treści i kolekcja; pustka jest wartością domyślną.
    mime_type      TEXT,
    kolekcja_kod   TEXT,
    -- Słownik dopuszczalnych wartości pola rodzaju `lista`, zapis JSON.
    opcje          TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pole_schematu_biblioteki_zakres
    ON pole_schematu_biblioteki(mime_type, kolekcja_kod);
