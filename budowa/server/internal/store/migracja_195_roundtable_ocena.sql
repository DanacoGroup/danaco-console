-- Migracja 195 — ocena Operatora, rubryki oceny i werdykty modeli-sędziów.
--
-- Ocena Operatora i werdykt sędziego są osobnymi bytami, choć obie „oceniają".
-- Różnią się autorem i skutkiem: ocena Operatora wchodzi do rankingu jako
-- pojedynek rozstrzygnięty ręcznie, werdykt sędziego niesie punkty w kryteriach
-- rubryki i uzasadnienie wypowiedziane przez model.
--
-- Rubryka bywa wspólna dla platformy albo związana z jednym oknem. Okno puste
-- znaczy rubrykę wspólną — `roundtable.rubric.list` bez wskazania okna oddaje
-- wtedy same wspólne, a ze wskazaniem wspólne i te jednego okna.

CREATE TABLE debata_rubryka (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL DEFAULT '',
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_rubryka_okno ON debata_rubryka(okno, id);

CREATE TABLE debata_kryterium (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rubryka                  TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    waga                     REAL    NOT NULL DEFAULT 0,
    opis                     TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_debata_kryterium_rubryka ON debata_kryterium(rubryka, kolejnosc, id);

CREATE TABLE debata_ocena (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL CHECK(rodzaj IN ('star','pairwise')),
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    uczestnik                TEXT    NOT NULL DEFAULT '',
    -- Gwiazdki poza zakresem 1..5 znaczą „bez oceny w skali"; zero jest
    -- wartością domyślną przy ocenie porównawczej.
    gwiazdki                 INTEGER NOT NULL DEFAULT 0,
    wskazana                 TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_ocena_okno ON debata_ocena(okno, id);

CREATE TABLE debata_werdykt (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rubryka                  TEXT    NOT NULL,
    sedzia                   TEXT    NOT NULL,
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    uczestnik                TEXT    NOT NULL DEFAULT '',
    punkty_json              TEXT    NOT NULL DEFAULT '{}',
    wynik                    REAL    NOT NULL DEFAULT 0,
    uzasadnienie             TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_werdykt_okno ON debata_werdykt(okno, id);
