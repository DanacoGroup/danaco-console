-- Migracja 374 zakłada tabelę wyciszeń nakładki powiadomień, obejmującą
-- wyciszenie czasowe, kontekstowe oraz wyciszenie całej klasy zdarzeń.

CREATE TABLE IF NOT EXISTS wyciszenie_nakladki (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT  NOT NULL UNIQUE,
    rodzaj                 TEXT    NOT NULL
                                   CHECK(rodzaj IN ('timed', 'contextual', 'eventClass')),
    zakres                 TEXT    NOT NULL DEFAULT ''
                                   CHECK(zakres IN ('', 'module', 'session', 'eventClass')),
    klucz_zakresu          TEXT    NOT NULL DEFAULT '',
    nazwa_zakresu          TEXT    NOT NULL DEFAULT '',
    klasa_zdarzen          TEXT    NOT NULL DEFAULT ''
                                   CHECK(klasa_zdarzen IN ('', 'executionLoopState', 'taskQueueState',
                                                           'qualityControlResult', 'moduleEvent',
                                                           'schedule', 'operatorWorkContext')),
    urzadzenie_id          TEXT    NOT NULL DEFAULT '',
    konczy_sie             TEXT    NOT NULL DEFAULT '',
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Warunek wymaga chwili końca dla wyciszenia czasowego i wskazanego bytu dla pozostałych rodzajów.
    CHECK((rodzaj = 'timed' AND zakres = '' AND klasa_zdarzen = '' AND konczy_sie <> '')
          OR (rodzaj = 'contextual' AND zakres IN ('module', 'session') AND klucz_zakresu <> '')
          OR (rodzaj = 'eventClass' AND zakres = 'eventClass' AND klasa_zdarzen <> ''))
);

-- Wskaźnik jednoznaczny obejmuje wyciszenie czasowe całej platformy oraz każdy
-- pozostały byt z osobna, tak aby powtórzone żądanie nie zakładało kolejnego wiersza.
CREATE UNIQUE INDEX IF NOT EXISTS idx_wyciszenie_nakladki_byt
    ON wyciszenie_nakladki(rodzaj, zakres, klucz_zakresu, klasa_zdarzen);
