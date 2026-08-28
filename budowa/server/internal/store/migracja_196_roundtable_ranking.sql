-- Migracja 196 — ranking uczestników akumulowany między sesjami.
--
-- Kluczem rankingu jest tożsamość, nie uczestnik. Uczestnik jest bytem okna
-- i ginie razem z debatą, a ranking ma przetrwać sesję (opracowanie, 2.5.5).
-- Tożsamością jest para „kanał modelu + nazwa persony", zapisana jako jeden
-- klucz — dzięki temu ten sam model w dwóch personach ma dwie punktacje, a ta
-- sama persona w dwóch oknach jedną.
--
-- Zakres i algorytm wchodzą do klucza jednoznaczności, bo ta sama tożsamość ma
-- odrębną punktację w rankingu środowiska i w rankingu jednego okna, a Elo
-- i Glicko liczą się inaczej i nie wolno ich sumować.

CREATE TABLE debata_ranking (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    klucz_tozsamosci         TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    zakres                   TEXT    NOT NULL CHECK(zakres IN ('environment','project','window')),
    -- Okno puste przy zakresie środowiska i projektu.
    okno                     TEXT    NOT NULL DEFAULT '',
    algorytm                 TEXT    NOT NULL CHECK(algorytm IN ('elo','glicko','trueSkill')),
    punktacja                REAL    NOT NULL DEFAULT 1500,
    -- Odchylenie punktacji: Glicko i TrueSkill trzymają w nim niepewność
    -- oszacowania. Elo go nie używa i zostawia wartość wnoszoną.
    odchylenie               REAL    NOT NULL DEFAULT 350,
    pojedynki                INTEGER NOT NULL DEFAULT 0,
    wygrane                  INTEGER NOT NULL DEFAULT 0,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (klucz_tozsamosci, zakres, okno, algorytm)
);
CREATE INDEX idx_debata_ranking_zakres ON debata_ranking(zakres, okno, algorytm, punktacja DESC);
