-- Migracja 222 — oś czasu aktywności projektu.
--
-- Zdarzenie zapisuje się w chwili czynności, a nie wylicza z bytów przy
-- odczycie. Wyliczanie z bytów pokazałoby wyłącznie to, co jeszcze istnieje:
-- usunięte zadanie znikałoby również z historii, a wtedy oś czasu przestaje być
-- zapisem zdarzeń i staje się drugim widokiem stanu bieżącego.

CREATE TABLE zdarzenie_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('projectChanged','instructionsChanged',
                                                      'memoryEntryChanged','fileChanged',
                                                      'taskChanged','noteChanged',
                                                      'agentAssignmentChanged','executionRun',
                                                      'commentAdded')),
    zmiana                   TEXT    NOT NULL DEFAULT ''
                                     CHECK(zmiana IN ('','created','updated','deleted')),
    rodzaj_bytu              TEXT    NOT NULL DEFAULT ''
                                     CHECK(rodzaj_bytu IN ('','task','note','file','memoryEntry',
                                                           'instructions','comment','project')),
    byt                      TEXT    NOT NULL DEFAULT '',
    opis                     TEXT    NOT NULL,
    rodzaj_autora            TEXT    NOT NULL DEFAULT 'operator'
                                     CHECK(rodzaj_autora IN ('operator','agent')),
    autor                    TEXT    NOT NULL DEFAULT '',
    zaszlo                   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zdarzenie_projektu_czas ON zdarzenie_projektu(projekt_id, id DESC);
CREATE INDEX idx_zdarzenie_projektu_byt ON zdarzenie_projektu(projekt_id, rodzaj_bytu, byt);
