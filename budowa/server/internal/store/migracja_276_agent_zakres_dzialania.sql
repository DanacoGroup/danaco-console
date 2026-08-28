-- Migracja 276 dodaje tabelę zakresu modułów zastosowania eksperta oraz
-- kolumnę limitu podagentów sieci wieloagentowej na koncie eksperta.

CREATE TABLE agent_modul_zastosowania (
    agent_id   INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    -- Kod modułu pochodzi z rodziny module kontraktu i nie ma więzu obcego.
    kod_modulu TEXT    NOT NULL,
    zapisano   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (agent_id, kod_modulu)
);

-- Indeks obsługuje pytanie odwrotne przy przeglądzie biblioteki: wskazuje,
-- którzy eksperci mają dostęp do wybranego modułu zastosowania.
CREATE INDEX idx_agent_modul_zastosowania_modul
    ON agent_modul_zastosowania(kod_modulu, agent_id);

ALTER TABLE agent ADD COLUMN limit_podagentow INTEGER NOT NULL DEFAULT 15
    CHECK (limit_podagentow BETWEEN 0 AND 15);
