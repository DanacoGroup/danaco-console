-- Migracja 277 dodaje tabelę ośmiu zakresów izolacji technicznej ustawianych
-- bezpośrednio na koncie eksperta, niezależnie od reguł zasięgu.

CREATE TABLE agent_izolacja_techniczna (
    agent_id INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    -- Wartość wyliczenia IsolationTechnicalScope kontraktu; katalog zakresów
    -- sprawdza rdzeń przy zapisie.
    zakres   TEXT    NOT NULL,
    odciety  INTEGER NOT NULL DEFAULT 0 CHECK (odciety IN (0, 1)),
    zapisano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (agent_id, zakres)
);
