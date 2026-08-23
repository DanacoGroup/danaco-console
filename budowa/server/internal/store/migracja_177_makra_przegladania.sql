-- Migracja 177 — nagrywarka makr przeglądania (`browser.macro.record`).
--
-- Makro jest ciągiem kroków w kolejności wykonania, a krokiem kontraktu jest
-- `AutomationStep` — ten sam kształt, którym jedzie moduł Automations. Kroki
-- idą kolumną JSON, bo kolejność i kształt kroku należą do makra: krok
-- wyjęty z makra nie jest niczym samodzielnym, nikt go nie wyszukuje osobno,
-- a tabela kroków kazałaby przy każdym doklejeniu kroku przenumerować resztę.
--
-- Kolumna `nagrywanie` jest stanem nagrywarki, nie własnością zapisu: makro
-- z ustawionym `nagrywanie = 1` przyjmuje kolejne kroki komendą
-- `browser.macro.record` z czynnością `step`. Zakończenie zapisu (`stop`)
-- zdejmuje ten stan i makro staje się zamkniętym scenariuszem.
--
-- Kolumna `automatyka_zewnetrzna_id` niesie odwołanie do przebiegu modułu
-- Automations, gdy scenariusz zostanie tam przekazany (opracowanie, rozdz.
-- 2.13). Nie jest więzem obcym: przekazanie jest konfigurowalne i makro ma
-- prawo istnieć bez niego.
CREATE TABLE makro_przegladania (
    id                        INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny  TEXT    NOT NULL UNIQUE,
    okno                      TEXT    NOT NULL,
    nazwa                     TEXT    NOT NULL,
    kroki_json                TEXT,
    nagrywanie                INTEGER NOT NULL DEFAULT 0 CHECK(nagrywanie IN (0,1)),
    automatyka_zewnetrzna_id  TEXT,
    utworzono                 TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_makro_przegladania_okno ON makro_przegladania(okno, utworzono DESC, id);
