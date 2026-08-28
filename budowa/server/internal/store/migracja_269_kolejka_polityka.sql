-- Migracja 269 wprowadza politykę kolejki jako kolumny tabeli kolejki oraz rodzaj kolejki
-- uruchamianej z zegara, w jednym kroku przepisania tabeli.

PRAGMA foreign_keys = off;

CREATE TABLE kolejka_nowa (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                  TEXT    NOT NULL,
    rodzaj                 TEXT    NOT NULL DEFAULT 'sesyjna'
                                   CHECK(rodzaj IN ('sesyjna','multitasking','harmonogram')),
    sesja_id               INTEGER REFERENCES sesja(id) ON DELETE CASCADE,
    okno_koordynatora_id   INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    stan                   TEXT    NOT NULL DEFAULT 'bezczynna'
                                   CHECK(stan IN ('bezczynna','pracuje','wstrzymana',
                                                  'zatrzymana','wyczerpana')),
    zasieg                 TEXT    NOT NULL DEFAULT 'lokalna'
                                   CHECK(zasieg IN ('globalna','lokalna','modelu',
                                                    'eksperta','projektu')),
    zasieg_id              TEXT,
    limit_rownoleglych     INTEGER NOT NULL DEFAULT 0,
    tempo_na_minute        INTEGER NOT NULL DEFAULT 0,
    limit_prob             INTEGER NOT NULL DEFAULT 3,
    wycofanie              TEXT    NOT NULL DEFAULT 'exponential'
                                   CHECK(wycofanie IN ('fixed','exponential')),
    wycofanie_sekundy      INTEGER NOT NULL DEFAULT 0,
    rozproszenie           INTEGER NOT NULL DEFAULT 1 CHECK(rozproszenie IN (0,1)),
    zadania_martwe         INTEGER NOT NULL DEFAULT 1 CHECK(zadania_martwe IN (0,1)),
    idempotencja_zycie_sekundy INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO kolejka_nowa
    (id, nazwa, rodzaj, sesja_id, okno_koordynatora_id, stan, utworzono, zaktualizowano)
SELECT id, nazwa, rodzaj, sesja_id, okno_koordynatora_id, stan, utworzono, zaktualizowano
  FROM kolejka;

DROP TABLE kolejka;
ALTER TABLE kolejka_nowa RENAME TO kolejka;

CREATE INDEX idx_kolejka_sesja ON kolejka(sesja_id);
CREATE INDEX idx_kolejka_koordynator ON kolejka(okno_koordynatora_id);

PRAGMA foreign_keys = on;
