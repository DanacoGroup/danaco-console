-- Migracja 163 zakłada tabelę ustalenie_korekty niosącą trwałe ustalenia korekty językowej ze stanem zastosowania i propozycją poprawki.

CREATE TABLE ustalenie_korekty (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    panel_id                 INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('grammar','spelling','punctuation','style',
                                                      'register','readability','falseFriends','typography')),
    waga                     TEXT    NOT NULL CHECK(waga IN ('hint','warning','error')),
    segment                  TEXT,
    szczegol                 TEXT    NOT NULL,
    -- Propozycja poprawki; pusta oznacza usterkę bez gotowego zastąpienia, a zastosowanie wtedy odmawia.
    propozycja               TEXT,
    zastosowano              INTEGER,
    odrzucono                INTEGER,
    utworzono                INTEGER NOT NULL
);
CREATE INDEX idx_ustalenie_korekty_panel ON ustalenie_korekty(panel_id, utworzono DESC);
