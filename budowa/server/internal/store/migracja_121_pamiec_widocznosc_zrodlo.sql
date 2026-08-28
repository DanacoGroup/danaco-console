-- Migracja 121 dodaje trzy braki bytów: poziomy pamięci eksperta, widoczność
-- eksperta i źródło rozszerzenia, których komendy istniejące nie miały gdzie
-- zapisać.

-- Tabela agent_pamiec_poziom przechowuje poziomy pamięci eksperta jako zbiór
-- wierszy: brak wiersza znaczy pamięć wyłączoną, nie piąty poziom.
CREATE TABLE agent_pamiec_poziom (
    agent_id INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    poziom   TEXT    NOT NULL
             CHECK (poziom IN ('global', 'project', 'session', 'environment')),

    PRIMARY KEY (agent_id, poziom)
);

-- Eksperci zastani zachowują stan wyjściowy platformy: pełny komplet poziomów.
-- Bez tego wypełnienia biblioteka zastana pokazałaby się jako „pamięć
-- wyłączona" — czyli migracja zmieniłaby zachowanie ekspertów, których nikt
-- nie tknął.
INSERT INTO agent_pamiec_poziom (agent_id, poziom)
SELECT a.id, p.poziom
  FROM agent a
  JOIN (SELECT 'global' AS poziom UNION ALL SELECT 'project'
        UNION ALL SELECT 'session' UNION ALL SELECT 'environment') p;

-- Kolumna widocznosc rozstrzyga, czy ekspert pokazuje się wszędzie, czy
-- wyłącznie tam, gdzie został przypisany do projektu.
ALTER TABLE agent ADD COLUMN widocznosc TEXT NOT NULL DEFAULT 'global'
    CHECK (widocznosc IN ('global', 'project'));

-- Kolumna zrodlo_pochodzenia niesie wartość kontraktu wprost i rozstrzyga
-- wyłącznie stan wyjściowy przy rejestracji rozszerzenia.
ALTER TABLE rozszerzenie ADD COLUMN zrodlo_pochodzenia TEXT NOT NULL DEFAULT 'personal'
    CHECK (zrodlo_pochodzenia IN ('danaco', 'personal'));
