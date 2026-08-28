-- Migracja 196 zakłada tabelę rankingu uczestników debaty, kluczowaną tożsamością kanału modelu i persony, akumulowaną między sesjami w danym zakresie.

CREATE TABLE debata_ranking (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    klucz_tozsamosci         TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    zakres                   TEXT    NOT NULL CHECK(zakres IN ('environment','project','window')),
    -- Okno puste przy zakresie środowiska i projektu.
    okno                     TEXT    NOT NULL DEFAULT '',
    algorytm                 TEXT    NOT NULL CHECK(algorytm IN ('elo','glicko','trueSkill')),
    punktacja                REAL    NOT NULL DEFAULT 1500,
    -- Odchylenie punktacji: Glicko i TrueSkill niosą w nim niepewność oszacowania, Elo go nie używa.
    odchylenie               REAL    NOT NULL DEFAULT 350,
    pojedynki                INTEGER NOT NULL DEFAULT 0,
    wygrane                  INTEGER NOT NULL DEFAULT 0,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (klucz_tozsamosci, zakres, okno, algorytm)
);
CREATE INDEX idx_debata_ranking_zakres ON debata_ranking(zakres, okno, algorytm, punktacja DESC);
