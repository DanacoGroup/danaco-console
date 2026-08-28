-- Migracja zakłada trwałość modułu przeglądarki: migawkę strony, źródło zebrane w toku
-- przeglądania i notatkę powiązaną ze źródłem.

-- Migawka jest tabelą historii nawigacji: każde wywołanie nawigacji lub odczytu migawki
-- wstawia nowy wiersz zrzutu stanu strony.
CREATE TABLE migawka_strony (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    tekst_odwolanie          TEXT,
    zrodlo_odwolanie         TEXT,
    zrzut_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_migawka_strony_okno ON migawka_strony(okno, utworzono DESC, id);

-- Źródło tego modułu jest odciskiem strony zebranym w toku przeglądania, zawsze
-- powiązanym z migawką tego samego okna.
CREATE TABLE zrodlo_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    migawka_zewnetrzna_id    TEXT,
    kluczowe                 INTEGER NOT NULL DEFAULT 0 CHECK(kluczowe IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zrodlo_przegladania_okno ON zrodlo_przegladania(okno, utworzono DESC, id);

-- Notatka bez wskazanego źródła jest dopuszczalna: dotyczy wtedy całej strony, nie
-- jednego zebranego odcinka treści.
CREATE TABLE notatka_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    zrodlo_zewnetrzny_id     TEXT,
    tresc                    TEXT    NOT NULL DEFAULT '',
    cytat                    TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_notatka_przegladania_okno ON notatka_przegladania(okno, utworzono DESC, id);
