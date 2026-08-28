-- Migracja 272 — punkty wznowienia przebiegu
-- (`automation.execution.checkpoint.list`, `.resume`).
--
-- Punkt wznowienia bez trwałości jest sprzecznością sam w sobie: służy
-- wznowieniu po awarii, a awaria zabiera pamięć procesu. Wykaz kroków
-- ukończonych leży zapisem strukturalnym, bo czyta się go w całości i tylko
-- w całości — wznowienie pyta „czego już nie powtarzać”, nie „czy krok siódmy”.
CREATE TABLE punkt_wznowienia_przebiegu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    przebieg_id              INTEGER NOT NULL REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    krok_kod                 TEXT    NOT NULL,
    kroki_ukonczone          TEXT    NOT NULL DEFAULT '[]',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_punkt_wznowienia_przebiegu_czas
    ON punkt_wznowienia_przebiegu(przebieg_id, id);
