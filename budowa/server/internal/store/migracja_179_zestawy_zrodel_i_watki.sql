-- Migracja 179 dopisuje do źródeł i notatek przeglądania kolumny grupowania, klasyfikacji, wątku i przypięcia oraz zakłada tabelę zestawów tematycznych.

ALTER TABLE zrodlo_przegladania ADD COLUMN grupa TEXT;
ALTER TABLE notatka_przegladania ADD COLUMN klasyfikacja TEXT;
ALTER TABLE notatka_przegladania ADD COLUMN watek TEXT;
ALTER TABLE notatka_przegladania ADD COLUMN przypieta INTEGER NOT NULL DEFAULT 0 CHECK(przypieta IN (0,1));

CREATE TABLE zestaw_zrodel_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zestaw_zrodel_okno ON zestaw_zrodel_przegladania(okno, utworzono DESC, id);

CREATE TABLE watek_notatek_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_watek_notatek_okno ON watek_notatek_przegladania(okno, utworzono DESC, id);
