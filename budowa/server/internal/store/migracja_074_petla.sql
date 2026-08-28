-- Migracja 074 dokłada do pętli automatyki obsadę uczestników, wiązanie kroku
-- z agentem portfolio oraz wybudzenie: stan biegu zawieszonego, czekającego na
-- reakcję ze świata. Pełne uzasadnienie stoi w dokumentacji architektury.

-- Wiązanie kroku z agentem idzie osobną tabelą, bo zapis kroków ją podmienia w całości.

CREATE TABLE agent_kroku_automatyki (
    automatyka_id   INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    -- Identyfikator zewnętrzny kroku, ten sam rodzaj klucza co w tabeli zaleznosc_kroku_automatyki.
    krok_zewnetrzny TEXT    NOT NULL,
    agent_id        INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok_zewnetrzny)
);

CREATE INDEX idx_agent_kroku_automatyki_agent
    ON agent_kroku_automatyki(agent_id);

-- Obsada liczy od jednego do czterech uczestników; granica jest specyfikacją, trzymaną więzem.

CREATE TABLE obsada_automatyki (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    -- Uczestnik: agent z portfolio ma pierwszeństwo przed samym modelem.
    agent_id      INTEGER          REFERENCES agent(id) ON DELETE SET NULL,
    model         TEXT,
    rola          TEXT    NOT NULL DEFAULT 'wykonawca'
                          CHECK(rola IN ('koordynator','wykonawca','samodzielne')),
    miejsce       INTEGER NOT NULL CHECK(miejsce BETWEEN 1 AND 4),
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT   NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (automatyka_id, miejsce),
    CHECK (agent_id IS NOT NULL OR model IS NOT NULL)
);

CREATE INDEX idx_obsada_automatyki_automatyka
    ON obsada_automatyki(automatyka_id, miejsce);

CREATE INDEX idx_obsada_automatyki_agent
    ON obsada_automatyki(agent_id)
    WHERE agent_id IS NOT NULL;

-- Przebudowa: SQLite nie zna polecenia zdejmującego więz CHECK.

CREATE TABLE przebieg_automatyki_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    kolejka_id               INTEGER REFERENCES kolejka(id) ON DELETE SET NULL,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','paused','oczekuje',
                                                    'succeeded','failed','stopped')),
    etap_biezacy             INTEGER NOT NULL DEFAULT 0,
    etapow                   INTEGER NOT NULL DEFAULT 0,
    proba                    INTEGER NOT NULL DEFAULT 0,
    komunikat_bledu          TEXT,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono               TEXT
);

INSERT INTO przebieg_automatyki_nowa
    (id, identyfikator_zewnetrzny, automatyka_id, kolejka_id, stan,
     etap_biezacy, etapow, proba, komunikat_bledu, rozpoczeto, zakonczono)
SELECT
     id, identyfikator_zewnetrzny, automatyka_id, kolejka_id, stan,
     etap_biezacy, etapow, proba, komunikat_bledu, rozpoczeto, zakonczono
  FROM przebieg_automatyki;

DROP TABLE przebieg_automatyki;

ALTER TABLE przebieg_automatyki_nowa RENAME TO przebieg_automatyki;

-- Indeksy migracji 040 odtworzone pod tymi samymi nazwami, wolnymi po skasowanej tabeli sprzed przebudowy schematu.
CREATE INDEX idx_przebieg_automatyki_automatyka
    ON przebieg_automatyki(automatyka_id, id DESC);

CREATE UNIQUE INDEX idx_przebieg_automatyki_kolejka
    ON przebieg_automatyki(kolejka_id)
    WHERE kolejka_id IS NOT NULL;

-- Oczekiwanie wiąże przebieg zawieszony z sygnałem i wznawia go od kroku, w którym stanął.

CREATE TABLE oczekiwanie_biegu (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    przebieg_id      INTEGER NOT NULL REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    -- Krok, w którym bieg stanął, identyfikator zewnętrzny jak w zaleznosc_kroku_automatyki.
    krok_zewnetrzny  TEXT    NOT NULL,
    -- Na co bieg czeka; kod sygnału ustala ten, kto sygnał wyśle.
    sygnal           TEXT    NOT NULL,
    -- Termin wybudzenia; wartość pusta znaczy oczekiwanie bez końca.
    termin           TEXT,
    po_terminie      TEXT    NOT NULL DEFAULT 'wznow'
                             CHECK(po_terminie IN ('wznow','ponow','przerwij')),
    zalozono         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Wartość pusta znaczy, że bieg wciąż czeka; wypełniona zamyka oczekiwanie.
    wybudzono        TEXT,
    powod_wybudzenia TEXT    CHECK(powod_wybudzenia IS NULL
                                   OR powod_wybudzenia IN ('sygnal','termin','operator')),
    -- Treść, z którą przyszedł sygnał, do wykorzystania przez wznawiany krok.
    tresc_sygnalu    TEXT,
    CHECK (wybudzono IS NULL OR powod_wybudzenia IS NOT NULL)
);

-- Jedno czynne oczekiwanie na bieg naraz; oczekiwania już zamknięte zostają w tabeli jako ślad przebiegu.
CREATE UNIQUE INDEX idx_oczekiwanie_biegu_czynne
    ON oczekiwanie_biegu(przebieg_id)
    WHERE wybudzono IS NULL;

-- Wskazuje biegi czekające właśnie na ten dany sygnał w chwili jego doręczenia rdzeniowi z zewnątrz świata.
CREATE INDEX idx_oczekiwanie_biegu_sygnal
    ON oczekiwanie_biegu(sygnal)
    WHERE wybudzono IS NULL;

-- Wskazuje oczekiwania z już minionym terminem; wiersze bez terminu zostają całkiem poza tym indeksem.
CREATE INDEX idx_oczekiwanie_biegu_termin
    ON oczekiwanie_biegu(termin)
    WHERE wybudzono IS NULL AND termin IS NOT NULL;
