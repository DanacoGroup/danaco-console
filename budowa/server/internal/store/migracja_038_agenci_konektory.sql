-- Migracja zakłada tabelę konektorów eksperta oraz tabelę jego uprawnień operacyjnych
-- w czterech grupach zakresu.

-- Konektor podpina się pod istniejące punkty dostępu; konektor rodzaju mcp jest wskazaniem
-- mostu, nie jego kopią.
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

-- Uprawnienie jest konfiguracją możliwości, nie kontrolą dostępu; stanem wyjściowym
-- eksperta jest pełny dostęp operacyjny.
CREATE TABLE agent_uprawnienie (
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    grupa     TEXT    NOT NULL CHECK (grupa IN ('files', 'network', 'processes', 'integrations')),
    zakres    TEXT    NOT NULL DEFAULT '',
    przyznane INTEGER NOT NULL DEFAULT 1 CHECK (przyznane IN (0, 1)),
    PRIMARY KEY (agent_id, grupa, zakres)
);
