-- Migracja 235 zakłada tabele zestawów żetonów systemu projektowego oraz
-- żetonów wchodzących w ich skład, jako byt odrębny od motywu powłoki.

CREATE TABLE zestaw_zetonow_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    motyw                    TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeks wspiera odczyt zestawów żetonów okna w porządku od ostatnio zaktualizowanego do najstarszego.
CREATE INDEX idx_zestaw_zetonow_design_okno
    ON zestaw_zetonow_design(okno, zaktualizowano DESC, id DESC);

-- Rodzaj żetonu niesie wartość zgodną z wyliczeniem kontraktu i nie ma warunku sprawdzającego w schemacie.
CREATE TABLE zeton_design (
    zestaw_id   INTEGER NOT NULL REFERENCES zestaw_zetonow_design(id) ON DELETE CASCADE,
    nazwa       TEXT    NOT NULL,
    rodzaj      TEXT    NOT NULL,
    wartosc     TEXT    NOT NULL,
    opis        TEXT,
    odsylacz_do TEXT,
    kolejnosc   INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (zestaw_id, nazwa)
);
CREATE INDEX idx_zeton_design_zestaw ON zeton_design(zestaw_id, kolejnosc);
