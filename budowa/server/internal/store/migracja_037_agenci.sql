-- Migracja zakłada tabelę eksperta modułu agentów jako komponentu własnego, wraz
-- z tabelą przypisanych mu umiejętności.

-- Kod eksperta jest identyfikatorem trwałym wychodzącym na zewnątrz jako pole kontraktu,
-- niezależnym od zmiany nazwy eksperta.
CREATE TABLE agent (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                  TEXT    NOT NULL UNIQUE,
    nazwa                TEXT    NOT NULL,
    opis                 TEXT    NOT NULL DEFAULT '',
    instrukcje_systemowe TEXT    NOT NULL DEFAULT '',
    kanal_kod            TEXT,
    model                TEXT,
    -- Droga wywołania kanału; wartości transportu wprost, bo brak słownika przekładu.
    transport            TEXT,
    -- Parametry wywołania zależne od kanału; sprawdzane jako poprawny JSON przez warstwę danych.
    parametry_json       TEXT    NOT NULL DEFAULT '{}',
    wersja               INTEGER NOT NULL DEFAULT 1,
    aktywny              INTEGER NOT NULL DEFAULT 1 CHECK (aktywny IN (0, 1)),
    utworzono            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    zaktualizowano       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Wykaz biblioteki ekspertów idzie po nazwie (kontrakt: „eksperci w kolejności
-- nazw”), a filtr `enabledOnly` po kolumnie `aktywny`.
CREATE INDEX idx_agent_nazwa ON agent(aktywny, nazwa);

-- ── Umiejętności przypisane ekspertowi ────────────────────────────────────────
-- Klucz główny na parze nie pozwala przypisać tej samej umiejętności dwa razy,
-- więc powtórzone `agent.skill.add` jest bezskutkowe zamiast dublować wiersz.
CREATE TABLE agent_umiejetnosc (
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    kod       TEXT    NOT NULL,
    dodano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    PRIMARY KEY (agent_id, kod)
);
