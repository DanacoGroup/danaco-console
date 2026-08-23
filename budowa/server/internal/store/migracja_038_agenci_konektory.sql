-- Migracja 038 — konektory i uprawnienia eksperta.
--
-- Konektor podpina się pod istniejące mosty. Kolumna `punkt_dostepu_id` wskazuje
-- wiersz `punkt_dostepu` — ten sam katalog, z którego składacz
-- `core/most_okna.go` buduje wpisy `mcpServers` dla procesu modelu. Drugiego
-- rejestru serwerów MCP nie zakładamy: konektor rodzaju `mcp` jest
-- wskazaniem mostu, a nie jego kopią. Klucz obcy jest ON DELETE SET NULL, bo
-- skasowanie punktu odłącza konektor, lecz nie kasuje eksperta.
--
-- Rodzaj. `rodzaj` przyjmuje dosłownie wartości shared.AgentConnectorKind
-- ('mcp','plugin','api'). Kontrakt nie daje dla tego wyliczenia słownika
-- przekładu bazy, więc kolumna trzyma wartość kontraktu wprost.
--
-- Uprawnienie jest konfiguracją możliwości, nie kontrolą dostępu.
-- Stanem wyjściowym eksperta jest pełny dostęp operacyjny: przy założeniu
-- eksperta warstwa `dane` wpisuje cztery wiersze `przyznane = 1`, po jednym na
-- grupę zakresu. Brak wiersza znaczy więc „nie rozstrzygnięto”, a nie „odmowa”
-- — Permissions Center pokazuje wtedy wartość wyjściową.
--
-- Zakres szczegółowy. `zakres` jest napisem w obrębie grupy (kontrakt: pole
-- `scope`, opcjonalne). Zakres pusty oznacza całą grupę; para (grupa, zakres)
-- jest kluczem, więc uprawnienie węższe nie nadpisuje szerszego i odwrotnie.

-- ── Konektory eksperta ────────────────────────────────────────────────────────
CREATE TABLE agent_konektor (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    kod              TEXT    NOT NULL UNIQUE,
    agent_id         INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    nazwa            TEXT    NOT NULL,
    rodzaj           TEXT    NOT NULL CHECK (rodzaj IN ('mcp', 'plugin', 'api')),
    punkt_dostepu_id INTEGER REFERENCES punkt_dostepu(id) ON DELETE SET NULL,
    konfiguracja     TEXT    NOT NULL DEFAULT '{}',
    aktywny          INTEGER NOT NULL DEFAULT 1 CHECK (aktywny IN (0, 1)),
    utworzono        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX idx_agent_konektor_agent ON agent_konektor(agent_id, nazwa);

-- ── Uprawnienia eksperta w czterech grupach zakresu ───────────────────────────
CREATE TABLE agent_uprawnienie (
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    grupa     TEXT    NOT NULL CHECK (grupa IN ('files', 'network', 'processes', 'integrations')),
    zakres    TEXT    NOT NULL DEFAULT '',
    przyznane INTEGER NOT NULL DEFAULT 1 CHECK (przyznane IN (0, 1)),
    PRIMARY KEY (agent_id, grupa, zakres)
);
