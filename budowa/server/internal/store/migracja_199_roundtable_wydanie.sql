-- Migracja 199 — szablony moderacji i artefakty wydane z debaty.
--
-- Szablon moderacji zapisuje format, liczbę tur, granicę czasu i kolejność
-- głosu (opracowanie, 2.8.5). Kolejność jest listą kodów uczestników
-- rozdzieloną znakiem nowego wiersza, a nie tabelą wiążącą: uczestnik należy do
-- okna, a szablon ma przeżyć okno, więc więz obcy do składu zabiłby szablon
-- razem z debatą, z której go zdjęto.
--
-- Artefakt debaty niesie odwołanie do bajtów w magazynie treści, nie same bajty.
-- Transkrypt, graf i nagranie idą w megabajtach, a baza rdzenia trzyma stan,
-- nie treść — tak samo jak w bibliotece i w module Design.

CREATE TABLE debata_szablon (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    format                   TEXT    NOT NULL DEFAULT 'free'
                                     CHECK(format IN ('free','structured','oxford',
                                                      'roundRobin','delphi','expertPanel')),
    granica_tur              INTEGER NOT NULL DEFAULT 0,
    granica_czasu_ms         INTEGER NOT NULL DEFAULT 0,
    kolejnosc_glosu          TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_szablon_nazwa ON debata_szablon(nazwa);

CREATE TABLE debata_artefakt (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Rodzaj mówi, co wydano: zapis debaty, graf argumentów albo odsłuch.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('transkrypt','graf','odsluch')),
    format                   TEXT    NOT NULL,
    -- Odwołanie względne magazynu treści, zawsze z ukośnikiem `/`.
    odwolanie                TEXT    NOT NULL,
    rozmiar                  INTEGER NOT NULL DEFAULT 0,
    suma_kontrolna           TEXT    NOT NULL DEFAULT '',
    -- Długość nagrania w milisekundach; zero przy artefaktach niebędących
    -- odsłuchem.
    dlugosc_ms               INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_artefakt_okno ON debata_artefakt(okno, id DESC);
