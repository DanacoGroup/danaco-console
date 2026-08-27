-- Migracja 219 zakłada tabele notatek i stron wiki projektu wraz z odnośnikami treści, niosącymi hierarchię przez wskazanie strony nadrzędnej.

CREATE TABLE notatka_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    tytul                    TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL DEFAULT '',
    notatka_nadrzedna        TEXT    NOT NULL DEFAULT '',
    -- Etykiety i nagłówki rozdzielone znakiem nowego wiersza.
    etykiety                 TEXT    NOT NULL DEFAULT '',
    naglowki                 TEXT    NOT NULL DEFAULT '',
    rodzaj_autora            TEXT    NOT NULL DEFAULT 'operator'
                                     CHECK(rodzaj_autora IN ('operator','agent')),
    autor                    TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_notatka_projektu_projekt ON notatka_projektu(projekt_id, tytul);
CREATE INDEX idx_notatka_projektu_nadrzedna ON notatka_projektu(notatka_nadrzedna);

CREATE TABLE odnosnik_notatki_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    notatka_zrodlowa         TEXT    NOT NULL,
    nazwa_docelowa           TEXT    NOT NULL,
    notatka_docelowa         TEXT    NOT NULL DEFAULT '',
    rodzaj                   TEXT    NOT NULL DEFAULT 'wikilink'
                                     CHECK(rodzaj IN ('wikilink','reference','embed','attachment')),
    kontekst                 TEXT    NOT NULL DEFAULT ''
);
CREATE INDEX idx_odnosnik_notatki_zrodlo ON odnosnik_notatki_projektu(notatka_zrodlowa);
CREATE INDEX idx_odnosnik_notatki_cel ON odnosnik_notatki_projektu(projekt_id, notatka_docelowa);
CREATE INDEX idx_odnosnik_notatki_nazwa ON odnosnik_notatki_projektu(projekt_id, nazwa_docelowa);
