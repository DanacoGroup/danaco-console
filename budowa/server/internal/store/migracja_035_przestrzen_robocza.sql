-- Migracja zakłada trwałość modułu przestrzeni roboczej: tabelę projektu, pamięci
-- projektu i przypisania eksperta do projektu.

-- Tabela projektu niesie nazwę, stan i czas ostatniej zmiany, których kolumna sesji
-- z samym kodem projektu nie niesie.
CREATE TABLE projekt (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT,
    stan           TEXT    NOT NULL DEFAULT 'active' CHECK(stan IN ('active','archived')),
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Wpis pamięci projektu niesie poziom zasięgu współdzielenia oraz pochodzenie
-- odróżniające ustalenie Operatora od propozycji modelu.
CREATE TABLE wpis_pamieci_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tresc                    TEXT    NOT NULL,
    przypiety                INTEGER NOT NULL DEFAULT 0 CHECK(przypiety IN (0,1)),
    pochodzenie              TEXT    NOT NULL DEFAULT 'operator'
                                     CHECK(pochodzenie IN ('operator','model')),
    poziom_zasiegu_id        INTEGER NOT NULL REFERENCES poziom_zasiegu(id),
    klucz_zasiegu            TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wpis_pamieci_projektu_projekt
    ON wpis_pamieci_projektu(projekt_id, przypiety DESC, zaktualizowano DESC);

-- Ekspert mieszka poza tą tabelą jako komponent własny modułu agentów; tabela zapisuje
-- wyłącznie fakt przypisania go do projektu.
CREATE TABLE przypisanie_agenta_projektu (
    projekt_id         INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    agent_kod          TEXT    NOT NULL,
    rola               TEXT,
    domyslny_wykonawca INTEGER NOT NULL DEFAULT 0 CHECK(domyslny_wykonawca IN (0,1)),
    przypisano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (projekt_id, agent_kod)
);
CREATE INDEX idx_przypisanie_agenta_projektu_domyslny
    ON przypisanie_agenta_projektu(projekt_id, domyslny_wykonawca DESC, agent_kod);
