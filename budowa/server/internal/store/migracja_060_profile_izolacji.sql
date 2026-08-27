-- Migracja zakłada profile izolacji, ich przypisanie do poziomu zasięgu oraz wybór
-- warstwy izolacji dla karty sesji albo okna.

CREATE TABLE profil_izolacji (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT    NOT NULL DEFAULT '',
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE przelacznik_profilu_izolacji (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id INTEGER NOT NULL REFERENCES profil_izolacji(id) ON DELETE CASCADE,
    klucz     TEXT    NOT NULL,
    wartosc   TEXT    NOT NULL,
    UNIQUE(profil_id, klucz)
);

CREATE TABLE przypisanie_profilu_izolacji (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id         INTEGER NOT NULL REFERENCES profil_izolacji(id) ON DELETE CASCADE,
    poziom_zasiegu_id INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu     TEXT    NOT NULL DEFAULT '',
    przypisano        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(poziom_zasiegu_id, klucz_zasiegu)
);

CREATE TABLE warstwa_izolacji (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom_zasiegu_id INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu     TEXT    NOT NULL DEFAULT '',
    warstwa           TEXT    NOT NULL CHECK(warstwa IN ('default','session')),
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(poziom_zasiegu_id, klucz_zasiegu)
);
