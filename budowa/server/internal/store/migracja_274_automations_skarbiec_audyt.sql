-- Migracja 274 wprowadza skarbiec poświadczeń automatyki bez kolumny na wartość, ponieważ
-- wartość leży w sejfie plikowym katalogu danych, a baza zna wyłącznie nazwę, zasięg i odwołanie.
CREATE TABLE poswiadczenie_automatyki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    odwolanie      TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    zasieg         TEXT,
    zasieg_id      TEXT,
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(nazwa, zasieg, zasieg_id)
);
CREATE INDEX idx_poswiadczenie_automatyki_nazwa ON poswiadczenie_automatyki(nazwa, id);

-- Dziennik audytu jest zapisem niezmiennym: wiersz raz dopisany nie jest zmieniany ani kasowany
-- przez żadną komendę modułu, a wpis o automatyce nieznanej jest dopuszczony.
CREATE TABLE wpis_audytu_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER REFERENCES automatyka(id) ON DELETE CASCADE,
    wykonawca                TEXT    NOT NULL DEFAULT 'operator',
    czynnosc                 TEXT    NOT NULL,
    szczegoly                TEXT,
    chwila                   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeksy porządkują dziennik audytu od najnowszego wpisu, osobno dla całego dziennika i osobno
-- zawężony do jednej automatyki, aby oba odczyty pomijały pełne przeglądanie tabeli.
CREATE INDEX idx_wpis_audytu_automatyki_czas ON wpis_audytu_automatyki(chwila DESC, id DESC);
CREATE INDEX idx_wpis_audytu_automatyki_wlasciciel
    ON wpis_audytu_automatyki(automatyka_id, chwila DESC, id DESC);
