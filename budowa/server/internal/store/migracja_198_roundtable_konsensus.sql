-- Migracja 198 — stanowisko końcowe redagowane, jego wersje, zdania odrębne
-- i przekazanie do modułu docelowego.
--
-- Kolumna `redagowane` rozstrzyga spór dwóch autorów tej samej treści. Do dziś
-- treść stanowiska składał rdzeń z zapisu tur przy każdym odczycie
-- (`roundtable.consensus.get`). Kontrakt ma dziś `roundtable.consensus.set`,
-- czyli treść nadaną ręcznie przez Operatora — a złożenie z tur nadpisałoby ją
-- przy pierwszym otwarciu okna. Znacznik mówi rdzeniowi, że tego wiersza nie
-- składa się już z zapisu: od redakcji stanowisko należy do Operatora.
--
-- Wersje leżą osobno, bo `roundtable.consensus.version.list` ma oddać kolejne
-- redakcje do porównania. Licznik w kolumnie `wersja` mówi, ile ich było;
-- porównać da się dopiero wtedy, gdy każda została.

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
    -- Jeden uczestnik podpisuje jedno zdanie odrębne wobec jednego stanowiska;
    -- powtórny zapis jest redakcją tego samego zdania, nie drugim głosem.
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
