-- Migracja 273 — reguły alarmowania (`automation.alert.rule.set`, `.list`).
--
-- Kanały leżą zapisem strukturalnym, bo kontrakt niesie je wykazem tekstów
-- ustalanym w całości: reguła ma kanały takie, jakie zapisano ostatnio.
-- Rozbicie na tabelę wierszy dałoby możliwość stanu, którego kontrakt nie zna
-- — kanału dopisanego bez przepisania reguły.
CREATE TABLE regula_alarmowania_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    wyzwalacz                TEXT    NOT NULL
                                     CHECK(wyzwalacz IN ('failure','timeout','missingRun',
                                                         'successRateDrop')),
    warunek                  TEXT,
    kanaly                   TEXT    NOT NULL DEFAULT '[]',
    czynna                   INTEGER NOT NULL DEFAULT 1 CHECK(czynna IN (0,1)),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_regula_alarmowania_automatyki_wlasciciel
    ON regula_alarmowania_automatyki(automatyka_id, id);
