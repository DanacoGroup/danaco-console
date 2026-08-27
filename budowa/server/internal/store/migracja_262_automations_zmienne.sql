-- Migracja 262 wprowadza zmienne przepływu automatyki: wartość domyślna leży jako zapis
-- strukturalny, a odwołanie do sekretu niesie wyłącznie referencję do skarbca, nigdy poświadczenie.
CREATE TABLE zmienna_automatyki (
    automatyka_id     INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    nazwa             TEXT    NOT NULL,
    rodzaj            TEXT    NOT NULL,
    wartosc_domyslna  TEXT,
    odwolanie_sekretu TEXT,
    kolejnosc         INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (automatyka_id, nazwa)
);
