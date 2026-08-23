-- Migracja 262 — zmienne przepływu (`automation.workflow.variables.set`).
--
-- Wartość domyślna leży jako zapis strukturalny, bo kontrakt niesie ją typem
-- `json`: zmienna bywa liczbą, tekstem i zapisem złożonym, a kolumna
-- o jednym typie skalarnym kazałaby rdzeniowi zgadywać przy odczycie.
--
-- `odwolanie_sekretu` nigdy nie niesie wartości poświadczenia — wyłącznie
-- referencję do skarbca (migracja 274). Wartość poświadczenia leży poza bazą.
CREATE TABLE zmienna_automatyki (
    automatyka_id     INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    nazwa             TEXT    NOT NULL,
    rodzaj            TEXT    NOT NULL,
    wartosc_domyslna  TEXT,
    odwolanie_sekretu TEXT,
    kolejnosc         INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (automatyka_id, nazwa)
);
