-- Migracja 224 zakłada tabelę wersji instrukcji projektu, w której każdy zapis
-- odkłada nową wersję tworzącą pełną historię, niezależną od ustawienia
-- obowiązującego w tabeli ustawienie.

CREATE TABLE wersja_instrukcji_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    tresc                    TEXT    NOT NULL DEFAULT '',
    odcisk                   TEXT    NOT NULL DEFAULT '',
    poziom                   TEXT    NOT NULL DEFAULT 'projekt',
    klucz_zasiegu            TEXT    NOT NULL DEFAULT '',
    rodzaj_autora            TEXT    NOT NULL DEFAULT 'operator'
                                     CHECK(rodzaj_autora IN ('operator','agent')),
    autor                    TEXT    NOT NULL DEFAULT '',
    przywrocono_z            TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wersja_instrukcji_projekt
    ON wersja_instrukcji_projektu(projekt_id, poziom, klucz_zasiegu, id DESC);
