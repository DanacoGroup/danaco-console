-- Migracja 231 zakłada tabelę szablonów promptu strukturalnego modułu Design,
-- przechowywanych jako byt niezależny od wydanych promptów.

CREATE TABLE szablon_promptu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    temat                    TEXT    NOT NULL,
    styl                     TEXT,
    kompozycja               TEXT,
    oswietlenie              TEXT,
    paleta                   TEXT,
    proporcje_kadru          TEXT,
    wykluczenia              TEXT,
    ziarno                   INTEGER,
    warianty                 INTEGER,
    silnik                   TEXT,
    kreatywnosc              REAL,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeks wspiera odczyt szablonów okna w porządku od najpóźniej zmienionego do najstarszego zapisanego szablonu.
CREATE INDEX idx_szablon_promptu_design_okno
    ON szablon_promptu_design(okno, zaktualizowano DESC, id DESC);
