-- Migracja 274 — skarbiec poświadczeń i dziennik audytu modułu Automations
-- (`automation.secret.*`, `automation.audit.list`).
--
-- Tabela poświadczeń NIE ma kolumny na wartość i mieć jej nie będzie. Wartość
-- leży w sejfie plikowym katalogu danych (`dane/sejf_poswiadczen.go`), tym
-- samym, którym jadą sekrety kont i punktów dostępu. Baza zna wyłącznie nazwę,
-- zasięg i odwołanie — dzięki temu odczyt bazy nie może wynieść sekretu, bo
-- kolumny na niego nie ma, a nie dlatego, że ktoś pamiętał o filtrze.
--
-- Zasięg idzie słownikiem kontraktu (`ConfigScope`), bo poświadczenie
-- rozstrzyga się warstwowo tak samo jak reszta modelu konfiguracji.
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

-- Dziennik audytu jest zapisem niezmiennym: wiersz raz dopisany nie jest
-- zmieniany ani kasowany przez żadną komendę modułu. Automatyka usunięta
-- zabiera swoje wpisy (ON DELETE CASCADE), lecz wpis o automatyce nieznanej
-- jest dopuszczony (kolumna pusta) — audytowana bywa czynność, która żadnej
-- automatyki nie dotyczy.
CREATE TABLE wpis_audytu_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER REFERENCES automatyka(id) ON DELETE CASCADE,
    wykonawca                TEXT    NOT NULL DEFAULT 'operator',
    czynnosc                 TEXT    NOT NULL,
    szczegoly                TEXT,
    chwila                   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- automation.audit.list czyta dziennik od najnowszego, zawężając po
-- automatyce i po zakresie dat.
CREATE INDEX idx_wpis_audytu_automatyki_czas ON wpis_audytu_automatyki(chwila DESC, id DESC);
CREATE INDEX idx_wpis_audytu_automatyki_wlasciciel
    ON wpis_audytu_automatyki(automatyka_id, chwila DESC, id DESC);
