-- Migracja 198 rozszerza stanowisko końcowe o redakcję ręczną, zakłada jego wersje, zdania odrębne uczestników oraz przekazanie do modułu docelowego.

ALTER TABLE debata_stanowisko ADD COLUMN redagowane INTEGER NOT NULL DEFAULT 0
    CHECK(redagowane IN (0,1));
ALTER TABLE debata_stanowisko ADD COLUMN zaakceptowane INTEGER NOT NULL DEFAULT 0
    CHECK(zaakceptowane IN (0,1));
ALTER TABLE debata_stanowisko ADD COLUMN kontekst TEXT;
ALTER TABLE debata_stanowisko ADD COLUMN warianty TEXT;
ALTER TABLE debata_stanowisko ADD COLUMN konsekwencje TEXT;
-- Tury objęte stanowiskiem nadane ręcznie, rozdzielone znakiem nowego wiersza.
-- Puste znaczy „tury biorą się z zapisu debaty".
ALTER TABLE debata_stanowisko ADD COLUMN tury TEXT NOT NULL DEFAULT '';

CREATE TABLE debata_stanowisko_wersja (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    stanowisko               TEXT    NOT NULL,
    wersja                   INTEGER NOT NULL,
    tresc                    TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (stanowisko, wersja)
);
CREATE INDEX idx_debata_stanowisko_wersja ON debata_stanowisko_wersja(stanowisko, wersja);

CREATE TABLE debata_zdanie_odrebne (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    stanowisko               TEXT    NOT NULL,
    uczestnik                TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Jeden uczestnik podpisuje jedno zdanie; ponowny zapis jest redakcją, nie drugim głosem.
    UNIQUE (stanowisko, uczestnik)
);
CREATE INDEX idx_debata_zdanie_odrebne ON debata_zdanie_odrebne(stanowisko, id);

CREATE TABLE debata_przekazanie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    stanowisko               TEXT    NOT NULL,
    modul                    TEXT    NOT NULL
                                     CHECK(modul IN ('studio','research','library','automations')),
    artefakt                 TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_przekazanie_okno ON debata_przekazanie(okno, id);
