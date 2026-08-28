-- Migracja zakłada trwałość słownika modułu tłumaczenia: terminy glosariusza, ślady
-- importu i eksportu oraz pamięć tłumaczeń.

-- Tożsamość terminu jest własna, nie wyprowadzona z treści, bo zapis pozwala edytować
-- istniejący termin przez jego identyfikator.
CREATE TABLE termin_slownika (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zrodlo                   TEXT    NOT NULL,
    jezyk                    TEXT    NOT NULL,
    cel                      TEXT,
    nie_tlumaczyc            INTEGER NOT NULL DEFAULT 0 CHECK(nie_tlumaczyc IN (0,1)),
    uwaga                    TEXT,
    zaktualizowano           INTEGER NOT NULL
);
-- glossary.apply i glossary.occurrences szukają terminów danego języka po
-- treści źródłowej; indeks obsługuje obie ścieżki wyszukiwania.
CREATE INDEX idx_termin_slownika_jezyk_zrodlo ON termin_slownika(jezyk, zrodlo);

-- Ślad importu słownika niesie ścieżkę pliku źródłowego oraz liczbę terminów
-- zaimportowanych w tym przebiegu.
CREATE TABLE slad_importu_slownika (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sciezka                  TEXT    NOT NULL,
    liczba_zaimportowanych   INTEGER NOT NULL,
    utworzono                INTEGER NOT NULL
);
CREATE INDEX idx_slad_importu_slownika_czas ON slad_importu_slownika(utworzono DESC);

-- Bez kolumny na treść ani identyfikator pliku repozytorium, bo rdzeń nie ma magazynu
-- blobów i eksport nie wytwarza pliku.
CREATE TABLE slad_eksportu_slownika (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sciezka                  TEXT    NOT NULL,
    liczba_wyeksportowanych  INTEGER NOT NULL,
    utworzono                INTEGER NOT NULL
);
CREATE INDEX idx_slad_eksportu_slownika_czas ON slad_eksportu_slownika(utworzono DESC);

-- Wpis gromadzony z zatwierdzonej pary segmentów, nie odczyt po panelu; segment zostaje
-- kolumną tekstową wprost, bez pary odwołania.
CREATE TABLE pamiec_tlumaczen (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    panel_id                 INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    jezyk                    TEXT    NOT NULL,
    segment_zrodlowy         TEXT    NOT NULL,
    segment_docelowy         TEXT    NOT NULL,
    utworzono                INTEGER NOT NULL
);
-- memory.suggest dopasowuje przybliżenie segmentu źródłowego w obrębie języka
-- docelowego panelu — indeks po (jezyk, segment_zrodlowy) obsługuje wstępne
-- zawężenie kandydatów przed dopasowaniem przybliżonym w warstwie `dane`.
CREATE INDEX idx_pamiec_tlumaczen_jezyk_segment ON pamiec_tlumaczen(jezyk, segment_zrodlowy);
CREATE INDEX idx_pamiec_tlumaczen_panel ON pamiec_tlumaczen(panel_id);
