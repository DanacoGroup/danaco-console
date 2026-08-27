-- Migracja 282 dodaje piątą grupę zakresu uprawnień eksperta obejmującą moduły
-- i zasoby, odtwarzając tabelę agent_uprawnienie z rozszerzonym warunkiem.

CREATE TABLE agent_uprawnienie_nowe (
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    grupa     TEXT    NOT NULL
                      CHECK (grupa IN ('files', 'network', 'processes', 'integrations', 'modules')),
    zakres    TEXT    NOT NULL DEFAULT '',
    przyznane INTEGER NOT NULL DEFAULT 1 CHECK (przyznane IN (0, 1)),
    PRIMARY KEY (agent_id, grupa, zakres)
);

INSERT INTO agent_uprawnienie_nowe (agent_id, grupa, zakres, przyznane)
SELECT agent_id, grupa, zakres, przyznane FROM agent_uprawnienie;

DROP TABLE agent_uprawnienie;
ALTER TABLE agent_uprawnienie_nowe RENAME TO agent_uprawnienie;
