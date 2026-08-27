-- Migracja 292 dodaje tabelę historii schowka operatora wraz z indeksem
-- porządkującym wpisy przypięte i najnowsze.

CREATE TABLE wpis_schowka (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL DEFAULT 'text',
    tresc                    TEXT    NOT NULL,
    odcisk                   TEXT    NOT NULL UNIQUE,
    zajawka                  TEXT,
    rozmiar_bajtow           INTEGER NOT NULL DEFAULT 0,
    przypiety                INTEGER NOT NULL DEFAULT 0,
    wrazliwy                 INTEGER NOT NULL DEFAULT 0,
    okno_zrodlowe            TEXT,
    utworzono                INTEGER NOT NULL,
    uzyto                    INTEGER
);
-- Wykaz idzie wpisami przypiętymi na czele, a dalej od najnowszego, zgodnie
-- z porządkiem prezentacji historii schowka.
CREATE INDEX idx_wpis_schowka_wykaz ON wpis_schowka(przypiety DESC, utworzono DESC, id DESC);
CREATE INDEX idx_wpis_schowka_rodzaj ON wpis_schowka(rodzaj, utworzono DESC);
