-- Migracja 306 — przebieg wsadu modułu Studio (`studio.batch.run`).
--
-- ── Dlaczego wsad ma tabelę, a nie samą odpowiedź ────────────────────────────
-- Kontrakt oddaje `runId`. Identyfikator, za którym nic nie stoi, byłby napisem
-- bez znaczenia: Operator nie miałby czym zapytać, co się z jego dziesięcioma
-- dokumentami stało, a Execution Loop Window nie miałby czego pokazać. Przebieg
-- jest więc bytem trwałym, tak samo jak przebieg automatyki.
--
-- ── Dlaczego pozycja jest osobnym wierszem ───────────────────────────────────
-- Rozdział 4.5 opracowania żąda, żeby odmowa jednego dokumentu nie wstrzymywała
-- pozostałych. Wynik wsadu jest więc wykazem rozstrzygnięć per dokument, nie
-- jednym stanem całości: wiersz na dokument niesie stan i powód odmowy własnymi
-- słowami, a suma tych wierszy daje liczby `accepted` i `rejected`.

CREATE TABLE przebieg_wsadu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    akcja_id                 TEXT    NOT NULL,
    parametry_json           TEXT,
    przyjete                 INTEGER NOT NULL DEFAULT 0,
    odrzucone                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_przebieg_wsadu_studio_okno ON przebieg_wsadu_studio(okno, id);

CREATE TABLE pozycja_wsadu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    przebieg_id              INTEGER NOT NULL REFERENCES przebieg_wsadu_studio(id) ON DELETE CASCADE,
    dokument_kod             TEXT    NOT NULL,
    stan                     TEXT    NOT NULL DEFAULT 'przyjeta'
                                     CHECK(stan IN ('przyjeta','odrzucona')),
    powod                    TEXT,
    propozycja_kod           TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pozycja_wsadu_studio_przebieg ON pozycja_wsadu_studio(przebieg_id, id);
