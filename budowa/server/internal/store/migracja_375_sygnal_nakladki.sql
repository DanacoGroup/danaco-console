-- Migracja 375 zakłada tabelę sygnałów klas zdarzeń wyzwalających nakładkę
-- powiadomień, niezależną od telemetrii poszczególnych modułów rdzenia.

CREATE TABLE IF NOT EXISTS sygnal_nakladki (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT  NOT NULL UNIQUE,
    klasa_zdarzen          TEXT    NOT NULL
                                   CHECK(klasa_zdarzen IN ('executionLoopState', 'taskQueueState',
                                                           'qualityControlResult', 'moduleEvent',
                                                           'schedule', 'operatorWorkContext')),
    tresc                  TEXT    NOT NULL,
    modul_kod              TEXT    NOT NULL DEFAULT '',
    sesja_kod              TEXT    NOT NULL DEFAULT '',
    liczba_wystapien       INTEGER NOT NULL DEFAULT 1 CHECK(liczba_wystapien >= 1),
    zdarzylo_sie           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE INDEX IF NOT EXISTS idx_sygnal_nakladki_klasa
    ON sygnal_nakladki(klasa_zdarzen, zdarzylo_sie DESC);
