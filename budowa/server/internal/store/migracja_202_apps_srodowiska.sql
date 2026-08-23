-- Migracja 202 — moduł Apps, Deployment Panel: środowiska wdrożeniowe produktu,
-- ich zmienne, nastawy skalowania i wyniki sprawdzeń kondycji.
--
-- Środowisko ma wiersz, choć kontrakt zna trzy wartości wbudowane. `AppEnvironment`
-- niesie `code`, `name`, `order` i `domain` — cztery pola, których wyliczenie
-- `AppDeployEnvironment` nie ma i mieć nie może, bo `apps.deployment.domain.set`
-- zapisuje domenę osobno dla każdego środowiska. Kontrakt mówi wprost, że kod
-- środowiska „dla trzech wartości wbudowanych równy AppDeployEnvironment", więc
-- wiersz jest miejscem na te cztery pola, a nie drugim słownikiem środowisk.
--
-- Wpisy DNS leżą w jednej kolumnie tekstowej z surowym JSON-em żądania
-- (`dnsRecords` kontraktu jest polem `json`). Rdzeń ich nie rozkłada: nie ma po
-- czym filtrować ani sortować, a kształt wpisu należy do dostawcy domeny, nie
-- do platformy.
--
-- Zmienna środowiskowa niesie wartość jawną ALBO odwołanie do sekretu, nigdy
-- oba naraz — kontrakt mówi „wyklucza się z secretRef" i warunek CHECK niżej
-- to egzekwuje. Treści sekretu w tej tabeli nie ma i nie będzie: kolumna trzyma
-- klucz jawny, po którym warstwa sekretów wydaje wartość.
--
-- Skalowanie ma jeden wiersz na (okno, środowisko) — `apps.deployment.scale.set`
-- nadsyła nastawę w komplecie, więc zapis jest UPSERT-em po tej parze, a nie
-- dziennikiem kolejnych nastaw.
--
-- Sprawdzenie kondycji jest dziennikiem, nie stanem. `AppDeploymentHealth`
-- niesie `availabilityPercent` — udział w oknie pomiaru — którego z jednego
-- wiersza „stan bieżący" policzyć się nie da. Każde `apps.deployment.health.get`
-- dopisuje więc wynik swojego sprawdzenia, a udział liczy się z wierszy.

-- ── Środowisko wdrożeniowe produktu ──────────────────────────────────────────
CREATE TABLE srodowisko_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Kod środowiska; dla trzech wartości wbudowanych równy AppDeployEnvironment.
    kod                      TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    domena                   TEXT,
    -- Surowy JSON wpisów DNS — patrz czoło pliku.
    wpisy_dns                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, kod)
);

-- ── Zmienna środowiskowa produktu ────────────────────────────────────────────
CREATE TABLE zmienna_srodowiska_apps (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    okno              TEXT    NOT NULL,
    srodowisko        TEXT    NOT NULL CHECK(srodowisko IN ('dev','staging','production')),
    nazwa             TEXT    NOT NULL,
    wartosc           TEXT,
    -- Klucz jawny warstwy sekretów, nigdy treść sekretu — patrz czoło pliku.
    odwolanie_sekretu TEXT,
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, srodowisko, nazwa),
    CHECK (wartosc IS NULL OR odwolanie_sekretu IS NULL)
);

-- ── Nastawa skalowania usługi ────────────────────────────────────────────────
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
