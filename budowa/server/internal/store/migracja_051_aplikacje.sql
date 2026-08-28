-- Migracja 051 tworzy trwałość modułu Apps: architekturę produktu, plik
-- warsztatu warstwy frontendu i backendu oraz dziennik przebiegów wdrożenia.

-- ── Architektura produktu — Architecture Designer ────────────────────────────
-- Wiersz trzyma definicję układu produktu w jednej wersjonowanej całości,
-- oddzielnej od komponentów i zależności zapisanych w tabelach własnych.
CREATE TABLE architektura_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    szablon                  TEXT    NOT NULL DEFAULT 'monolith'
                                     CHECK(szablon IN ('monolith','microservices','serverless')),
    -- Wersja rośnie przy każdym zapisie definicji układu, by panel mógł
    -- wskazać oglądaną wersję.
    wersja                   INTEGER NOT NULL DEFAULT 1,
    -- Zastrzeżenia walidacji układu to wynik obliczony przy zapisie, nie
    -- wejście Operatora.
    zastrzezenia_walidacji   TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_architektura_apps_okno ON architektura_apps(okno, zaktualizowano DESC);

-- ── Komponent architektury — Architecture Designer ───────────────────────────
-- Wartości kolumny `rodzaj` są wartościami kontraktu (AppComponentKind), nie
-- ich tłumaczeniem.
CREATE TABLE komponent_architektury_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    architektura_id          INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('frontend','backend','service','database','queue','external')),
    stos                     TEXT,
    opis                     TEXT,
    kontrakt_api             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (architektura_id, identyfikator_zewnetrzny)
);
CREATE INDEX idx_komponent_architektury_apps_architektura
    ON komponent_architektury_apps(architektura_id, id);

-- ── Zależność między komponentami — Architecture Designer ───────────────────
-- Więz pierwotny obejmuje parę komponentów, pętla własna odrzucona wprost, cykl
-- dłuższy rozstrzyga walidacja układu w rdzeniu.
CREATE TABLE zaleznosc_komponentu_apps (
    architektura_id INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    komponent_z     TEXT    NOT NULL,
    komponent_do    TEXT    NOT NULL,
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (architektura_id, komponent_z, komponent_do),
    CHECK (komponent_z <> komponent_do)
);
CREATE INDEX idx_zaleznosc_komponentu_apps_do
    ON zaleznosc_komponentu_apps(architektura_id, komponent_do);

-- ── Plik warsztatu — Workspace okna Apps, odrębny od modułu Workspace ───────
-- Klucz złożony z okna, warstwy i ścieżki niesie tożsamość pliku; zapis
-- działa jako nadpisanie stanu bieżącego, bez tabeli wersji obok.
CREATE TABLE plik_warsztatu_apps (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    okno           TEXT    NOT NULL,
    warstwa        TEXT    NOT NULL CHECK(warstwa IN ('frontend','backend')),
    sciezka        TEXT    NOT NULL,
    tresc          TEXT    NOT NULL,
    rozmiar        INTEGER NOT NULL DEFAULT 0,
    -- Komponent architektury dotknięty zmianą; wartość danych, nie więź
    -- obca, bo komponent mógł zniknąć.
    komponent_id   TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, warstwa, sciezka)
);
CREATE INDEX idx_plik_warsztatu_apps_okno ON plik_warsztatu_apps(okno, warstwa, sciezka);

-- ── Dziennik przebiegów wdrożenia — Deployment ───────────────────────────────
-- Tabela niesie wyłącznie ślad zlecenia wdrożenia i jego stan; sam przebieg
-- wdrożenia prowadzi rdzeń w pamięci procesu.
CREATE TABLE wdrozenie_apps (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                   TEXT    NOT NULL UNIQUE,
    okno                  TEXT    NOT NULL,
    srodowisko            TEXT    NOT NULL CHECK(srodowisko IN ('dev','staging','production')),
    strategia             TEXT    NOT NULL DEFAULT 'immediate'
                                  CHECK(strategia IN ('immediate','staged','blueGreen')),
    stan                  TEXT    NOT NULL DEFAULT 'pending'
                                  CHECK(stan IN ('pending','running','succeeded','failed','rolledBack')),
    wersja                TEXT,
    notatki_wydania       TEXT,
    adres                 TEXT,
    log_odwolanie         TEXT,
    -- Wdrożenie źródłowe, z którego nastąpiło cofnięcie; kolumna
    -- samoodwołania jednego łuku.
    cofniete_do_kodu      TEXT,
    rozpoczeto             TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono             TEXT
);
CREATE INDEX idx_wdrozenie_apps_okno ON wdrozenie_apps(okno, rozpoczeto DESC);
CREATE INDEX idx_wdrozenie_apps_stan ON wdrozenie_apps(stan, rozpoczeto DESC);
