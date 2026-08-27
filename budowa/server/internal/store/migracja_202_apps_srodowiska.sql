-- Migracja 202 zakłada tabele środowisk wdrożeniowych, zmiennych środowiskowych, nastaw skalowania oraz wyników sprawdzeń kondycji modułu Apps.

-- Zakłada tabelę srodowisko_apps niosącą kod, nazwę, domenę i wpisy DNS środowiska wdrożeniowego produktu.
CREATE TABLE srodowisko_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Kod środowiska; dla trzech wartości wbudowanych równy AppDeployEnvironment.
    kod                      TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    domena                   TEXT,
    -- Surowy JSON wpisów DNS, nierozkładany przez rdzeń.
    wpisy_dns                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, kod)
);

-- Zakłada tabelę zmienna_srodowiska_apps niosącą wartość jawną albo odwołanie do sekretu, nigdy oba naraz.
CREATE TABLE zmienna_srodowiska_apps (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    okno              TEXT    NOT NULL,
    srodowisko        TEXT    NOT NULL CHECK(srodowisko IN ('dev','staging','production')),
    nazwa             TEXT    NOT NULL,
    wartosc           TEXT,
    -- Klucz jawny warstwy sekretów, nigdy treść samego sekretu.
    odwolanie_sekretu TEXT,
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, srodowisko, nazwa),
    CHECK (wartosc IS NULL OR odwolanie_sekretu IS NULL)
);

-- Zakłada tabelę skalowanie_apps niosącą jeden wiersz nastawy na parę okna i środowiska, nadpisywany w całości.
CREATE TABLE skalowanie_apps (
    okno           TEXT    NOT NULL,
    srodowisko     TEXT    NOT NULL CHECK(srodowisko IN ('dev','staging','production')),
    instancje      INTEGER,
    min_instancji  INTEGER,
    maks_instancji INTEGER,
    -- Surowy JSON reguł wyzwalających skalowanie.
    reguly         TEXT,
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (okno, srodowisko)
);

-- ── Wynik sprawdzenia kondycji wdrożonego produktu ───────────────────────────
-- `sprawdzono` niesie milisekundy epoki wprost z kontraktu (`lastCheckAt`).
CREATE TABLE kondycja_wdrozenia_apps (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    okno       TEXT    NOT NULL,
    srodowisko TEXT    NOT NULL CHECK(srodowisko IN ('dev','staging','production')),
    dostepna   INTEGER NOT NULL CHECK(dostepna IN (0,1)),
    szczegol   TEXT,
    sprawdzono INTEGER NOT NULL
);
CREATE INDEX idx_kondycja_wdrozenia_apps ON kondycja_wdrozenia_apps(okno, srodowisko, sprawdzono DESC);
