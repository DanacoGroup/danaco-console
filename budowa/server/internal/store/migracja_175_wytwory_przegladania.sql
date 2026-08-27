-- Migracja 175 zakłada tabele wytworów sesji przeglądania oraz zrzutów stron, z odwołaniem do bajtów w magazynie treści i pomiarami zrzutu.

CREATE TABLE wytwor_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL,
    tytul                    TEXT,
    tresc_odwolanie          TEXT    NOT NULL,
    typ_mime                 TEXT,
    rozmiar_bajtow           INTEGER,
    url_zrodla               TEXT,
    migawka_zewnetrzna_id    TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wytwor_przegladania_okno ON wytwor_przegladania(okno, utworzono DESC, id);

CREATE TABLE zrzut_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    migawka_zewnetrzna_id    TEXT,
    tresc_odwolanie          TEXT    NOT NULL,
    tryb                     TEXT    NOT NULL,
    format                   TEXT    NOT NULL,
    szerokosc                INTEGER NOT NULL DEFAULT 0,
    wysokosc                 INTEGER NOT NULL DEFAULT 0,
    rozmiar_bajtow           INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zrzut_przegladania_okno ON zrzut_przegladania(okno, utworzono DESC, id);
CREATE INDEX idx_zrzut_przegladania_odwolanie ON zrzut_przegladania(tresc_odwolanie);
